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
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func RemoveAnySuffix(fileName string) string {
	suffixLoc := strings.LastIndex(fileName, ".")
	if suffixLoc == -1 {
		return fileName
	}
	newFileName := fileName[:suffixLoc]
	return newFileName
}

func FindNumberedSuffix(fileName string) (bool, error) {
	match, err := regexp.MatchString("^.+\\.\\d+$", fileName) //TODO convert to struct function and compile regex, can also easy capture suffix as group and remove it from string
	if err != nil {
		return false, fmt.Errorf("error in regex %v trying to find suffix for string %v", err, fileName)
	}
	return match, nil
}

type InvalidSuffixErr struct {
	FileName string
}

func (i InvalidSuffixErr) Error() string {
	return fmt.Sprintf("expected suffix with a number but was '%v' for file '%v'", filepath.Ext(i.FileName), i.FileName)
}

type SortingErr struct {
	BaseErr error
}

func (s SortingErr) Error() string {
	return fmt.Sprintf("unable to sort due to the following error '%v'", s.BaseErr)
}

func CombineFiles(fileNames []string, verbose bool) (totalBytesWritten int64, newFileName string, err error) {
	sortErrors := make(map[string]bool)
	sort.SliceStable(fileNames, func(i, j int) bool {
		one := fileNames[i]
		suffixOne := strings.Trim(filepath.Ext(one), ".")
		intOne, err := strconv.Atoi(suffixOne)
		if err != nil {
			// using strings so I can easily create one error out of this later
			sortErrors[fmt.Sprintf("not able to parse suffix '%v' for file %v with error '%v'", suffixOne, one, err)] = true
			return false
		}
		two := fileNames[j]
		suffixTwo := strings.Trim(filepath.Ext(two), ".")
		intTwo, err := strconv.Atoi(suffixTwo)
		if err != nil {
			// using strings so I can easily create one error out of this later
			sortErrors[fmt.Sprintf("not able to parse suffix '%v' for file %v with error '%v'", suffixTwo, two, err)] = true
			return false
		}
		return intOne < intTwo
	})
	var uniqueErrors []string
	for k := range sortErrors {
		uniqueErrors = append(uniqueErrors, k)
	}
	if len(sortErrors) > 0 {
		return -1, "", SortingErr{BaseErr: errors.New(strings.Join(uniqueErrors, ","))}
	}
	if verbose {
		slog.Debug("combining files", "files_to_combine", strings.Join(fileNames, ", "))
	}
	firstFile := filepath.Clean(fileNames[0])
	match, err := FindNumberedSuffix(firstFile)
	if err != nil {
		return -1, "", fmt.Errorf("unable to find suffix for file '%v' with error '%v'", firstFile, err)
	}
	if !match {
		return -1, "", InvalidSuffixErr{
			FileName: firstFile,
		}
	}
	newFileName = RemoveAnySuffix(firstFile)
	if len(fileNames) == 1 {
		//optimize and skip the copy step
		if err := os.Rename(firstFile, newFileName); err != nil {
			return -1, "", fmt.Errorf("unable to rename file '%v' to '%v' due to error '%v'", firstFile, newFileName, err)
		}
		fileInfo, err := os.Stat(newFileName)
		if err != nil {
			return -1, "", fmt.Errorf("unable to get file information for %v due to error %v", firstFile, err)
		}
		return fileInfo.Size(), newFileName, nil
	}
	newFileHandle, err := os.Create(filepath.Clean(newFileName))
	if err != nil {
		return -1, "", fmt.Errorf("cannot create file '%v' due to error '%v'", newFileName, err)
	}
	// cleanup in case of errors so we don't leak descriptors
	defer func() {
		if err := newFileHandle.Close(); err != nil {
			slog.Debug("unable to close file, since this is a cleanup operation it is usually safe to ignore", "file_name", newFileName, "error_msg", err)
		}
	}()
	for _, f := range fileNames {
		fileHandle, err := os.Open(filepath.Clean(f))
		if err != nil {
			return -1, "", fmt.Errorf("unable to read file '%v' due to error '%v'", f, err)
		}
		// cleanup in case of errors
		defer func() {
			if err := fileHandle.Close(); err != nil {
				slog.Debug("unable to close file handle, since this is a cleanup operation it is usually safe to ignore", "file_name", f, "error_msg", err)
			}
		}()
		buf := make([]byte, 8192*1024)
		bytesWritten, err := io.CopyBuffer(newFileHandle, fileHandle, buf)
		if err != nil {
			return -1, "", fmt.Errorf("unable to copy file '%v' to file '%v' due to error '%v'", f, newFileName, err)
		}
		totalBytesWritten += bytesWritten
		if err := fileHandle.Close(); err != nil {
			return -1, "", fmt.Errorf("unable to close old file %v due to error %v", filepath.Clean(f), err)
		}

		err = os.Remove(f)
		if err != nil {
			slog.Warn("unable to remove old file after copying it's contents to the new file and it will have to be manually deleted", "file_name", f, "error_msg", err)
		}
	}
	if err := newFileHandle.Close(); err != nil {
		return -1, "", fmt.Errorf("unable to close newfile %v due to error %v, not succesfully written this means", filepath.Clean(newFileName), err)
	}

	return totalBytesWritten, newFileName, nil
}
