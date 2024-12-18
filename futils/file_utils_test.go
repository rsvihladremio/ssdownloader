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

// package futils provides file utilities for very common ops
package futils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIfNotExists(t *testing.T) {
	fileName := filepath.Join(t.TempDir(), "test.txt")
	err := os.WriteFile(fileName, []byte("test file"), 0640)
	if err != nil {
		t.Fatal(err)
	}
	exists, err := FileExists(fileName)
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Errorf("file %v does not exist", fileName)
	}
	incorrectFileName := "notFindable"
	exists, err = FileExists(incorrectFileName)
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Errorf("file %v does exist and should not", incorrectFileName)
	}
}
