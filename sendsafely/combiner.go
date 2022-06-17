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
package sendsafely

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func CombineFiles(fileNames []string) (string, error) {

	sort.Strings(fileNames)
	firstFile := fileNames[0]
	if !strings.HasSuffix(firstFile, ".1") {
		return "", fmt.Errorf("expected suffix of .1 but was %v", filepath.Ext(firstFile))
	}
	newFileName := filepath.Clean(strings.TrimSuffix(firstFile, ".1"))
	if len(fileNames) == 1 {
		//optimize and skip the copy step
		if err := os.Rename(firstFile, newFileName); err != nil {
			return "", fmt.Errorf("unable to rename file '%v' to '%v' due to error '%v'", firstFile, newFileName, err)
		}
		return newFileName, nil
	}
	newFileHandle, err := os.Create(newFileName)
	if err != nil {
		return "", fmt.Errorf("cannot create file '%v' due to error '%v'", newFileName, err)
	}
	defer func() {
		err = newFileHandle.Close()
		if err != nil {
			log.Printf("WARN: unable to close file handle for file '%v' due to error '%v'", newFileName, err)
		}
	}()
	for _, f := range fileNames {
		fileHandle, err := os.Open(filepath.Clean(f))
		if err != nil {
			return "", fmt.Errorf("unable to read file '%v' due to error '%v'", f, err)
		}
		close := func() {
			err = fileHandle.Close()
			if err != nil {
				log.Printf("WARN: unable to close file handle for file '%v' due to error '%v'", f, err)
			}
		}
		buf := make([]byte, 8192*1024)
		_, err = io.CopyBuffer(newFileHandle, fileHandle, buf)
		if err != nil {
			close()
			return "", fmt.Errorf("unable to copy file '%v' to file '%v' due to error '%v'", f, newFileName, err)
		}
		close()
		err = os.Remove(f)
		if err != nil {
			log.Printf("WARN unable to remove old file '%v' after copying it's contents to the new file due to error '%v' and it will have to be manually deleted", f, err)
		}
	}
	return newFileName, nil
}
