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
package downloader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadFile(fileName, url string) error {
	// making sure there are no goofy file names that overwrite critical files
	cleanedFileName := filepath.Clean(fileName)
	f, err := os.Create(cleanedFileName)
	if err != nil {
		return fmt.Errorf("unable to create the destination file '%v' due to error '%v'", cleanedFileName, err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Printf("WARN: unable to close file handle for file '%v' due to error '%v'", cleanedFileName, err)
		}
	}()
	// is technically a security violation according to https://securego.io/docs/rules/g107.html
	// but in reality based on the application is unavoidable and a risk of using SendSafely
	// ignoring the rule in the ./script/audit file. Used suggestion from
	// https://stackoverflow.com/questions/70281883/golang-untaint-url-variable-to-fix-gosec-warning-g107
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("unable to retrieve url '%v' due to error '%v'", url, err)
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Printf("WARN: unable to close body handle for url '%v' due to error '%v'", url, err)
		}
	}()
	//TODO make optional with verbose flag
	//log.Printf("file %v complete with %v bytes written", fileName, bytes_written)
	//TODO make buffer size adjustable
	buf := make([]byte, 4096*1024)
	_, err = io.CopyBuffer(f, resp.Body, buf)
	if err != nil {
		return fmt.Errorf("unable to write to filename '%v' due to error '%v'", cleanedFileName, err)
	}

	return nil
}
