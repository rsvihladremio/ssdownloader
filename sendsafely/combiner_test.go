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
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindNumberedSuffix(t *testing.T) {
	s := "abc.123"
	match, err := FindNumberedSuffix(s)
	if err != nil {
		t.Errorf("unexpected error searching for suffix for '%v' with error '%v'", s, err)
	}
	if !match {
		t.Errorf("expected match for string %v", s)
	}
	s = "abc.1"
	match, err = FindNumberedSuffix(s)
	if err != nil {
		t.Errorf("unexpected error searching for suffix for '%v' with error '%v'", s, err)
	}
	if !match {
		t.Errorf("expected match for string %v", s)
	}
	s = "abc.123.123"
	match, err = FindNumberedSuffix(s)
	if err != nil {
		t.Errorf("unexpected error searching for suffix for '%v' with error '%v'", s, err)
	}
	if !match {
		t.Errorf("expected match for string %v", s)
	}
	s = "abc.abc"
	match, err = FindNumberedSuffix(s)
	if err != nil {
		t.Errorf("unexpected error searching for suffix for '%v' with error '%v'", s, err)
	}
	if match {
		t.Errorf("expected to NOT match for string %v", s)
	}
	s = "abc"
	match, err = FindNumberedSuffix(s)
	if err != nil {
		t.Errorf("unexpected error searching for suffix for '%v' with error '%v'", s, err)
	}
	if match {
		t.Errorf("expected to NOT match for string %v", s)
	}
}
func TestRemoveAnySuffix(t *testing.T) {
	actual := RemoveAnySuffix("abc.123")
	if actual != "abc" {
		t.Errorf("expected abc but was %v", actual)
	}
	actual = RemoveAnySuffix("abc.123.123")
	if actual != "abc.123" {
		t.Errorf("expected abc but was %v", actual)
	}
	actual = RemoveAnySuffix("abc.abc")
	if actual != "abc" {
		t.Errorf("expected abc but was %v", actual)
	}
	actual = RemoveAnySuffix("abc")
	if actual != "abc" {
		t.Errorf("expected abc but was %v", actual)
	}
}

func TestCombiningMoreThanTheFirstPart(t *testing.T) {

	dirToGenerate := filepath.Join("testdata", "combining")
	err := os.MkdirAll(dirToGenerate, 0755)
	if err != nil {
		t.Fatalf("unexpected error making dir %v %v", dirToGenerate, err)
	}
	for i := 26; i < 28; i++ {
		newFile := filepath.Join(dirToGenerate, fmt.Sprintf("mylog.txt.%v", i))
		err = os.WriteFile(newFile, []byte(fmt.Sprintf("row %v\n", i)), 0644)
		if err != nil {
			t.Fatalf("unable to create file %v due to error %v", i, err)
		}
		defer func() {
			if err := os.Remove(newFile); err != nil {
				log.Printf("ignore this, but for debugging purposes we were unable to remove %v due to error %v", newFile, err)
			}
		}()
	}
	d, err := os.Open(dirToGenerate)
	if err != nil {
		t.Fatalf("unexpected error reading dir %v %v", dirToGenerate, err)
	}
	entries, err := d.ReadDir(0)
	if err != nil {
		t.Fatalf("unable to list dirs %v", err)
	}
	var files []string
	for _, e := range entries {
		files = append(files, filepath.Join(dirToGenerate, e.Name()))
	}
	f, err := CombineFiles(files)
	if err != nil {
		t.Fatalf("unexpected error combining files %v", err)
	}
	defer os.Remove(f)
	contents, err := os.ReadFile(f)
	if err != nil {
		t.Fatalf("unexpected error reading combined file with error %v", err)
	}
	l := strings.Split(string(contents), "\n")
	var lines []string
	for _, line := range l {
		if line != "" {
			lines = append(lines, line)
		}
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 but got %v for array %#v", len(lines), lines)
	}

	expectedLine := "row 26"
	if lines[0] != expectedLine {
		t.Errorf("unexpected '%v' from line 1 but was '%v'", expectedLine, lines[0])
	}
	expectedLine = "row 27"
	if lines[1] != expectedLine {
		t.Errorf("unexpected '%v' from line 2 but was '%v'", expectedLine, lines[1])
	}
}
