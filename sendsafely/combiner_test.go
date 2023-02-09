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
	for i := 1; i < 21; i++ {
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
	w, f, err := CombineFiles(files, false)
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
	if len(lines) != 20 {
		t.Fatalf("expected 20 but got %v for array %#v", len(lines), lines)
	}
	for i := 0; i < 20; i++ {
		expectedLine := fmt.Sprintf("row %v", i+1)
		if lines[i] != expectedLine {
			t.Errorf("unexpected '%v' from line %v but was '%v'", expectedLine, i, lines[i])
		}
	}
	actualFileSize := len([]byte(contents))
	if w != int64(actualFileSize) {
		t.Errorf("expected %v and %v to match", w, actualFileSize)
	}
}
func TestNoOpCombining(t *testing.T) {

	dirToGenerate := filepath.Join("testdata", "combining")
	err := os.MkdirAll(dirToGenerate, 0755)
	if err != nil {
		t.Fatalf("unexpected error making dir %v %v", dirToGenerate, err)
	}
	newFile := filepath.Join(dirToGenerate, "mylog.txt.1")
	err = os.WriteFile(newFile, []byte("row 1\n"), 0644)
	if err != nil {
		t.Fatalf("unable to create file %v due to error %v", newFile, err)
	}
	defer func() {
		if err := os.Remove(newFile); err != nil {
			log.Printf("ignore this, but for debugging purposes we were unable to remove %v due to error %v", newFile, err)
		}
	}()

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
	_, f, err := CombineFiles(files, false)
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
	if len(lines) != 1 {
		t.Fatalf("expected 1 but got %v for array %#v", len(lines), lines)
	}
	expectedLine := "row 1"
	if lines[0] != expectedLine {
		t.Errorf("unexpected '%v' from line 1 but was '%v'", expectedLine, lines[0])
	}
}

func TestCombiningOneBadSuffix(t *testing.T) {

	dirToGenerate := filepath.Join(t.TempDir(), "combining")
	err := os.MkdirAll(dirToGenerate, 0755)
	if err != nil {
		t.Fatalf("unexpected error making dir %v %v", dirToGenerate, err)
	}
	newFile := filepath.Join(dirToGenerate, "mylog.txt")
	err = os.WriteFile(newFile, []byte("row 1\n"), 0644)
	if err != nil {
		t.Fatalf("unable to create file %v due to error %v", newFile, err)
	}
	defer func() {
		if err := os.Remove(newFile); err != nil {
			log.Printf("ignore this, but for debugging purposes we were unable to remove %v due to error %v", newFile, err)
		}
	}()

	_, _, err = CombineFiles([]string{newFile}, false)
	if err == nil {
		t.Fatalf("expected error but did not have one")
	}
	if !errors.Is(err, InvalidSuffixErr{FileName: newFile}) {
		t.Errorf("expected error to be InvalidSuffixErr but was %v", err)
	}
	expectedErr := fmt.Sprintf("expected suffix with a number but was '.txt' for file '%v'", newFile)
	if err.Error() != expectedErr {
		t.Errorf("expected error '%v' but was '%v'", err.Error(), expectedErr)
	}
}

func TestCombiningSeveralBadSuffixes(t *testing.T) {

	dirToGenerate := filepath.Join(t.TempDir(), "combining")
	err := os.MkdirAll(dirToGenerate, 0755)
	if err != nil {
		t.Fatalf("unexpected error making dir %v %v", dirToGenerate, err)
	}
	newFile := filepath.Join(dirToGenerate, "mylog.txt.1")
	err = os.WriteFile(newFile, []byte("row 1\n"), 0644)
	if err != nil {
		t.Fatalf("unable to create file %v due to error %v", newFile, err)
	}
	defer func() {
		if err := os.Remove(newFile); err != nil {
			log.Printf("ignore this, but for debugging purposes we were unable to remove %v due to error %v", newFile, err)
		}
	}()
	newFile2 := filepath.Join(dirToGenerate, "mylog.txt.two")
	err = os.WriteFile(newFile2, []byte("row 2\n"), 0644)
	if err != nil {
		t.Fatalf("unable to create file %v due to error %v", newFile2, err)
	}
	defer func() {
		if err := os.Remove(newFile2); err != nil {
			log.Printf("ignore this, but for debugging purposes we were unable to remove %v due to error %v", newFile2, err)
		}
	}()

	_, _, err = CombineFiles([]string{newFile, newFile2}, false)
	if err == nil {
		t.Fatalf("expected error but did not have one")
	}
	expectedErr := SortingErr{BaseErr: fmt.Errorf("not able to parse suffix 'two' with error 'strconv.Atoi: parsing \"two\": invalid syntax'")}
	if err.Error() != expectedErr.Error() {
		t.Errorf("expected error to be\n'%v' but was \n'%v'", expectedErr, err)
	}

}

func TestCombiningSeveralBadSuffixesReverseOrder(t *testing.T) {

	dirToGenerate := filepath.Join(t.TempDir(), "combining")
	err := os.MkdirAll(dirToGenerate, 0755)
	if err != nil {
		t.Fatalf("unexpected error making dir %v %v", dirToGenerate, err)
	}
	newFile := filepath.Join(dirToGenerate, "mylog.txt.one")
	err = os.WriteFile(newFile, []byte("row 1\n"), 0644)
	if err != nil {
		t.Fatalf("unable to create file %v due to error %v", newFile, err)
	}
	defer func() {
		if err := os.Remove(newFile); err != nil {
			log.Printf("ignore this, but for debugging purposes we were unable to remove %v due to error %v", newFile, err)
		}
	}()
	newFile2 := filepath.Join(dirToGenerate, "mylog.txt.2")
	err = os.WriteFile(newFile2, []byte("row 2\n"), 0644)
	if err != nil {
		t.Fatalf("unable to create file %v due to error %v", newFile2, err)
	}
	defer func() {
		if err := os.Remove(newFile2); err != nil {
			log.Printf("ignore this, but for debugging purposes we were unable to remove %v due to error %v", newFile2, err)
		}
	}()

	_, _, err = CombineFiles([]string{newFile, newFile2}, false)
	if err == nil {
		t.Fatalf("expected error but did not have one")
	}
	expectedErr := SortingErr{BaseErr: fmt.Errorf("not able to parse suffix 'one' with error 'strconv.Atoi: parsing \"one\": invalid syntax'")}
	if err.Error() != expectedErr.Error() {
		t.Errorf("expected error to be\n'%v' but was \n'%v'", expectedErr, err)
	}

}
