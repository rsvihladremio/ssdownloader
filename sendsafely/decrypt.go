/*
   Copyright 2022 Ryan SVIHLA

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

//sendsafely package decrypts files, combines file parts into whole files, and handles api access to the sendsafely rest api
package sendsafely

import (
	"bytes"
	"crypto"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	pgpErrors "github.com/ProtonMail/go-crypto/openpgp/errors"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

// DecryptPart decrypts one of the file parts using the gpp protocol with the following logic from the sendsafely docs
// https://sendsafely.zendesk.com/hc/en-us/articles/360027599232-SendSafely-REST-API under b. Download and Decrypt File Parts
//
// Each file part will need to be individually downloaded and decrypted using PGP. You will need to use the "Server Secret" (included in the Package Information response from Step 1) and the keycode (Client Secret) in order to compute the required decryption key.
//
// When decrypting each file part, make sure you use the following PGP options:
//
// * Symmetric-Key Algorithm should be 9 (AES-256)
// * Compression Algorithm should be 0 (Uncompressed)
// * Hash Algorithm should be 8 (SHA-256)
// * Passphrase:  Server Secret concatenated with a random 256-bit Client Secret
// * S2k-count: 65535
// * Mode: b (62)
func DecryptPart(filePart, serverSecret, keyCode string) (string, error) {
	//super super quirky docs and implementation here
	password := serverSecret + keyCode
	firstTimeCalled := true

	var prompt = func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		if firstTimeCalled {
			firstTimeCalled = false
			return []byte(password), nil
		}
		// Re-prompt still occurs if SKESK pasrsing fails (i.e. when decrypted cipher algo is invalid).
		// For most (but not all) cases, inputting a wrong passwords is expected to trigger this error.
		return nil, errors.New("gopenpgp: wrong password in symmetric decryption")
	}

	config := &packet.Config{
		DefaultCipher:     packet.CipherAES256,
		DefaultHash:       crypto.Hash(crypto.SHA256),
		Time:              getTimeGenerator(),
		S2KCount:          65535,
		CompressionConfig: &packet.CompressionConfig{Level: 0},
	}

	var emptyKeyRing openpgp.EntityList
	cleanedFilePart := filepath.Clean(filePart)
	encryptedIO, err := os.Open(cleanedFilePart)
	if err != nil {
		return "", fmt.Errorf("unable to read %v due to error %v", cleanedFilePart, err)
	}
	md, err := openpgp.ReadMessage(encryptedIO, emptyKeyRing, prompt, config)
	if err != nil {
		// Parsing errors when reading the message are most likely caused by incorrect password, but we cannot know for sure
		return "", fmt.Errorf("gopenpgp: error in reading password protected message: wrong password or malformed message %v", err)
	}

	defer func() {
		err := encryptedIO.Close()
		if err != nil {
			log.Printf("WARN encrypted io handler for file '%v' failed to close due to '%v'", cleanedFilePart, err)
		}

	}()

	messageBuf := bytes.NewBuffer(nil)
	buf := make([]byte, 4096*1024)
	_, err = io.CopyBuffer(messageBuf, md.UnverifiedBody, buf)
	if errors.Is(err, pgpErrors.ErrMDCHashMismatch) {
		// This MDC error may also be triggered if the password is correct, but the encrypted data was corrupted.
		// To avoid confusion, we do not inform the user about the second possibility.
		return "", errors.New("gopenpgp: wrong password in symmetric decryption")
	} else if err != nil {
		// Parsing errors after decryption, triggered before parsing the MDC packet, are also usually the result of wrong password
		return "", fmt.Errorf("gopenpgp: error in reading password protected message: wrong password or malformed message this is happening during parsing after decryption. %v", err)
	}

	//remove the "encrypted" suffix
	newFileName := strings.TrimSuffix(cleanedFilePart, ".encrypted")

	//TODO behind verbose logs
	//log.Printf("new file name is %v", newFileName)
	err = os.WriteFile(newFileName, messageBuf.Bytes(), 0600)
	if err != nil {
		return "", fmt.Errorf("unable to write file '%v' due to error '%v'", newFileName, err)
	}
	if err := os.Remove(cleanedFilePart); err != nil {
		log.Printf("WARN unable to delete '%v' due to error '%v' so you will need to manually clean this file up", filePart, err)
	}
	return newFileName, nil
}

// getNow returns the latest server time.
func getNow() time.Time {
	pgp.lock.RLock()
	defer pgp.lock.RUnlock()

	if pgp.latestServerTime == 0 {
		return time.Now()
	}

	return time.Unix(pgp.latestServerTime, 0)
}

// getTimeGenerator Returns a time generator function.
func getTimeGenerator() func() time.Time {
	return getNow
}

var pgp = GopenPGP{
	latestServerTime: 0,
	generationOffset: 0,
	lock:             &sync.RWMutex{},
}

// GopenPGP is used as a "namespace" for many of the functions in this package.
// It is a struct that keeps track of time skew between server and client.
type GopenPGP struct {
	latestServerTime int64
	generationOffset int64
	lock             *sync.RWMutex
}
