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
package link

import (
	"errors"
	"testing"
)

func TestLinkHandler(t *testing.T) {
	url := "https://sendsafely.tester.com/receive/?thread=MYTHREAD&packageCode=MYPKGCODE#keyCode=MYKEYCODE"
	linkParks, err := ParseLink(url)
	if err != nil {
		t.Fatalf("unexected error '%v'", err)
	}
	expectedKeyCode := "MYKEYCODE"
	if linkParks.KeyCode != expectedKeyCode {
		t.Errorf("expected keycode '%v' but got '%v'", expectedKeyCode, linkParks.KeyCode)
	}

	expectedPackageCode := "MYPKGCODE"
	if linkParks.PackageCode != expectedPackageCode {
		t.Errorf("expected package code '%v' but got '%v'", expectedPackageCode, linkParks.PackageCode)
	}

	expectedThread := "MYTHREAD"
	if linkParks.Thread != expectedThread {
		t.Errorf("expected thread '%v' but got '%v'", expectedThread, linkParks.Thread)
	}
}

func TestKeyCodeMissing(t *testing.T) {
	url := "https://sendsafely.tester.com/receive/?thread=MYTHREAD&packageCode=MYPKGCODE"
	_, err := ParseLink(url)
	if err == nil {
		t.Fatalf("exected error '%v'", err)
	}
	if !errors.Is(err, KeyCodeIsMissingErr{
		InputUrl: url,
		KeyCode:  "",
	}) {
		t.Errorf("expected KeyCodeIsMissingErr but got %v", err)
	} else {
		expectedError := "expected to have fragment keyCode= in url 'https://sendsafely.tester.com/receive/?thread=MYTHREAD&packageCode=MYPKGCODE' but it is not present, the fragment detected is ''"
		if err.Error() != expectedError {
			t.Errorf("expected error text of '%v' but was '%v'", expectedError, err.Error())
		}
	}
}

func TestPackageCodeMissing(t *testing.T) {
	url := "https://sendsafely.tester.com/receive/?thread=MYTHREAD#keyCode=MYKEYCODE"
	_, err := ParseLink(url)
	if err == nil {
		t.Fatalf("exected error '%v'", err)
	}
	if !errors.Is(err, PackageCodeIsMissingErr{
		InputUrl: url,
	}) {
		t.Errorf("expected PackageCodeIsMissingErr but got %v", err)
	} else {
		expectedError := "expected to have packageCode in url 'https://sendsafely.tester.com/receive/?thread=MYTHREAD#keyCode=MYKEYCODE' but it is not present"
		if err.Error() != expectedError {
			t.Errorf("expected error text of '%v' but was '%v'", expectedError, err.Error())
		}
	}
}

func TestThreadMissing(t *testing.T) {
	url := "https://sendsafely.tester.com/receive/?packageCode=MYPKGCODE#keyCode=MYKEYCODE"
	_, err := ParseLink(url)
	if err == nil {
		t.Fatalf("exected error '%v'", err)
	}
	if !errors.Is(err, ThreadIsMissingErr{
		InputUrl: url,
	}) {
		t.Errorf("expected ThreadIsMissingErr but got %v", err)
	} else {
		expectedError := "expected to have thread in url 'https://sendsafely.tester.com/receive/?packageCode=MYPKGCODE#keyCode=MYKEYCODE' but it is not present"
		if err.Error() != expectedError {
			t.Errorf("expected error text of '%v' but was '%v'", expectedError, err.Error())
		}
	}
}

func TestInvalidUrl(t *testing.T) {
	url := "*$ù%"
	_, err := ParseLink(url)
	if err == nil {
		t.Fatalf("exected error '%v'", err)
	}
	expectedError := UrlParseErr{Url: url, BaseErr: errors.New("parse \"*$ù%\": invalid URL escape \"%\"")}
	if err.Error() != expectedError.Error() {
		t.Errorf("expected '%v' but got '%v'", expectedError, err)
	}
}
