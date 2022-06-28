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
	"errors"
	"fmt"
	"io"
	"log"
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

func CombineFiles(fileNames []string, verbose bool) (string, error) {
	var sortErrors []string
	sort.SliceStable(fileNames, func(i, j int) bool {
		one := fileNames[i]
		two := fileNames[j]
		suffixOne := strings.Trim(filepath.Ext(one), ".")
		suffixTwo := strings.Trim(filepath.Ext(two), ".")
		intOne, err := strconv.Atoi(suffixOne)
		if err != nil {
			// using strings so I can easily create one error out of this later
			sortErrors = append(sortErrors, fmt.Sprintf("not able to parse suffix '%v' with error '%v'", suffixOne, err))
			return false
		}
		intTwo, err := strconv.Atoi(suffixTwo)
		if err != nil {
			// using strings so I can easily create one error out of this later
			sortErrors = append(sortErrors, fmt.Sprintf("not able to parse suffix '%v' with error '%v'", suffixTwo, err))
			return false
		}
		return intOne < intTwo
	})
	if len(sortErrors) > 0 {
		return "", SortingErr{BaseErr: errors.New(strings.Join(sortErrors, ","))}
	}
	if verbose {
		log.Printf("DEBUG: combining the following files: %v\n-", strings.Join(fileNames, "\n-"))
	}
	firstFile := filepath.Clean(fileNames[0])
	match, err := FindNumberedSuffix(firstFile)
	if err != nil {
		return "", fmt.Errorf("unable to find suffix for file '%v' with error '%v'", firstFile, err)
	}
	if !match {
		return "", InvalidSuffixErr{
			FileName: firstFile,
		}
	}
	newFileName := RemoveAnySuffix(firstFile)
	if len(fileNames) == 1 {
		//optimize and skip the copy step
		if err := os.Rename(firstFile, newFileName); err != nil {
			return "", fmt.Errorf("unable to rename file '%v' to '%v' due to error '%v'", firstFile, newFileName, err)
		}
		return newFileName, nil
	}
	newFileHandle, err := os.Create(filepath.Clean(newFileName))
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
