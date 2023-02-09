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

// sendsafely package decrypts files, combines file parts into whole files, and handles api access to the sendsafely rest api
package sendsafely

import (
	"crypto"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

// verify the pgp server time function we wrapped here
func TestGetNow(t *testing.T) {
	pgp.latestServerTime = 1000
	now := getNow()
	if now.Unix() != pgp.latestServerTime {
		t.Errorf("unexpected time of '%v', expected was '%v'", now.Unix(), 1000)
	}
}

// verify that time.Now is used even when the pgpServerTime is 0
func TestGetNowWithNoServerTime(t *testing.T) {
	//get now and we will test the time is close to now
	now := time.Now()
	//set pgp server time to 0 to get time.Now() functionality in method
	pgp.latestServerTime = 0
	newNow := getNow()
	// tolerance of 2 seconds
	tolerance := int64(2)
	maxAcceptable := now.Unix() + tolerance
	minAcceptable := now.Unix() - tolerance
	if newNow.Unix() > maxAcceptable {
		t.Errorf("time '%v' is too old and expected max time was '%v'", newNow.Unix(), maxAcceptable)
	}
	if newNow.Unix() < minAcceptable {
		t.Errorf("time '%v' is to new and expected min time was '%v'", newNow.Unix(), minAcceptable)
	}
}

func TestDecrypt(t *testing.T) {
	originalFile := "testdata/data-for-sendsafely-1-source.csv"
	password := "serverSecretkeyCode"
	sourceFile, err := os.ReadFile(originalFile)
	if err != nil {
		t.Fatalf("unable to copy file %v", err)
	}

	startingFile := t.TempDir() + "/data-for-sendsafely-1.csv"
	err = os.WriteFile(startingFile, sourceFile, 0755)
	if err != nil {
		t.Fatalf("unable to create file to be copied %v", err)
	}
	defer func() {
		if err := os.Remove(startingFile); err != nil {
			log.Printf("safe to ignore this if it fails. unable to remove '%v' due to error '%v'", startingFile, err)
		}
	}()
	newFile, err := EncryptFile(startingFile, password)
	if err != nil {
		t.Fatalf("error setting up test file due to '%v'", err)
	}
	unecryptedFile, err := DecryptPart(newFile, "serverSecret", "keyCode")
	if err != nil {
		t.Fatalf("unable to decrypt newly encrypted file '%v' error was '%v'", newFile, err)
	}
	defer func() {
		if err := os.Remove(unecryptedFile); err != nil {
			log.Printf("safe to ignore if this fails. unable to remove '%v' due to error '%v'", unecryptedFile, err)
		}
	}()
	unecryptedBytes, err := os.ReadFile(unecryptedFile)
	if err != nil {
		t.Fatalf("unable to read newly encrypted file '%v' error was '%v'", newFile, err)
	}
	startingFileBytes, err := os.ReadFile(startingFile)
	if err != nil {
		t.Fatalf("unable to read newly original file '%v' error was '%v'", startingFile, err)
	}
	startingFileStr := string(startingFileBytes)
	unecryptedStr := string(unecryptedBytes)
	if len(startingFileStr) == 0 {
		t.Fatalf("file '%v' was empty", originalFile)
	}
	if len(unecryptedStr) == 0 {
		t.Fatalf("file '%v' was empty", unecryptedFile)
	}
	log.Printf("unecrypted file length is %v and original file length is %v", len(unecryptedBytes), len(startingFileBytes))
	if startingFileStr != unecryptedStr {
		t.Errorf("original file and unecrypted file are not equal:\n diff is\n%v", Diff(startingFile, unecryptedStr))
	}
}

func Diff(s1, s2 string) string {
	linesDiff := []string{}
	shortestLength := 0
	longestLength := 0
	longestStr := ""
	if len(s1) > len(s2) {
		shortestLength = len(s2)
		longestLength = len(s1)
		longestStr = s1
	} else {
		shortestLength = len(s1)
		longestLength = len(s2)
		longestStr = s2
	}
	s1Lines := strings.Split(s1, "\n")
	s2Lines := strings.Split(s2, "\n")
	for i := 0; i < shortestLength; i++ {
		line1 := s1Lines[0]
		line2 := s2Lines[1]
		if line1 != line2 {
			linesDiff = append(linesDiff, line1+"-----"+line2)
		}
	}
	longestLines := strings.Split(longestStr, "\n")
	for i := shortestLength; i < longestLength; i++ {
		linesDiff = append(linesDiff, longestLines[i])
	}
	return strings.Join(linesDiff, "\n")
}

func EncryptFile(fileName, password string) (string, error) {
	config := &packet.Config{
		DefaultCipher:     packet.CipherAES256,
		DefaultHash:       crypto.Hash(crypto.SHA256),
		Time:              getTimeGenerator(),
		S2KCount:          65535,
		CompressionConfig: &packet.CompressionConfig{Level: 0},
	}

	cleanedFilePart := filepath.Clean(fileName)
	fileToEncrypt, err := os.Open(cleanedFilePart)
	if err != nil {
		return "", fmt.Errorf("unable to read %v due to error %v", cleanedFilePart, err)
	}
	defer func() {
		err := fileToEncrypt.Close()
		if err != nil {
			log.Printf("WARN encrypted io handler for file '%v' failed to close due to '%v'", cleanedFilePart, err)
		}

	}()
	newFileName := cleanedFilePart + ".encrypted"
	newFile, err := os.Create(newFileName)
	if err != nil {
		return "", fmt.Errorf("unable to create %v due to error %v", newFileName, err)
	}
	defer func() {
		err := newFile.Close()
		if err != nil {
			log.Printf("WARN file '%v' failed to close due to '%v'", newFileName, err)
		}

	}()
	md, err := openpgp.SymmetricallyEncrypt(newFile, []byte(password), &openpgp.FileHints{IsBinary: true}, config)
	if err != nil {
		// Parsing errors when reading the message are most likely caused by incorrect password, but we cannot know for sure
		return "", fmt.Errorf("gopenpgp: error in reading password protected message: wrong password or malformed message %v", err)
	}
	defer func() {
		if err := md.Close(); err != nil {
			log.Printf("WARN unable to close encrypted file handler with error '%v'", err)
		}
	}()
	if _, err := io.Copy(md, fileToEncrypt); err != nil {
		return "", fmt.Errorf("unable to write to new file from encrypted file due to error '%v'", err)
	}
	return newFileName, nil
}
