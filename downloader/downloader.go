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

// downloader package provides http download capability
package downloader

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

type IllegalBufferSize struct {
	BufferSizeKB int
}

func (e IllegalBufferSize) Error() string {
	return fmt.Sprintf("buffer size kb was %v and cannot be smaller than 1 please initialize the downloader with the NewGenericDownloader() function to guard against this happending", e.BufferSizeKB)
}

func NewGenericDownloader(bufferSizeKB int) GenericDownloader {
	if bufferSizeKB < 1 {
		slog.Debug("buffer size cannot be smaller than 1 setting to default of 4096")
		bufferSizeKB = 4096
	}
	return &HTTPGenericDownloader{bufferSizeKB: bufferSizeKB}
}

type GenericDownloader interface {
	DownloadFile(fileName, url string) error
}

type HTTPGenericDownloader struct {
	bufferSizeKB int
}

func (d *HTTPGenericDownloader) DownloadFile(fileName, url string) error {
	if d.bufferSizeKB < 1 {
		return IllegalBufferSize{BufferSizeKB: d.bufferSizeKB}
	}

	// making sure there are no goofy file names that overwrite critical files
	cleanedFileName := filepath.Clean(fileName)
	f, err := os.Create(cleanedFileName)
	if err != nil {
		return fmt.Errorf("unable to create the destination file '%v' due to error '%v'", cleanedFileName, err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			slog.Debug("unable to close file handle for file on cleanup. This is safe to ignore usually", "file_name", cleanedFileName, "error_msg", err)
		}
	}()
	// is technically a security violation according to https://securego.io/docs/rules/g107.html
	// but in reality based on the application is unavoidable and a risk of using SendSafely
	// ignoring the rule in the ./script/audit file. Used suggestion from
	// https://stackoverflow.com/questions/70281883/golang-untaint-url-†variable-to-fix-gosec-warning-g107
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("unable to retrieve url '%v' due to error '%v'", url, err)
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			slog.Warn("unable to close body handle for url", "url", url, "error_msg", err)
		}
	}()

	buf := make([]byte, d.bufferSizeKB*1024)
	_, err = io.CopyBuffer(f, resp.Body, buf)
	if err != nil {
		return fmt.Errorf("unable to write to filename '%v' due to error '%v'", cleanedFileName, err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("unable to close file %v due to error %v", cleanedFileName, err)
	}
	return nil
}
