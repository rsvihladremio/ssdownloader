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

// link package handles parsing of sendsafely links so that we can retrieve the identifying information in the query parameters
package link

import (
	"errors"
	"testing"
)

func TestLinkHandlerWithEncodedByGoogleUrlMissingQ(t *testing.T) {
	url := "https://www.google.com/url?wrong=https%3A%2F%2Fsendsafely.tester.com%2Freceive%2F%3Fthread%3DMYTHREAD%26packageCode%3DMYPKGCODE%23keyCode%3DMYKEYCODE&sa=D&ust=11111111&usg=JJJJJJJJJ"
	_, err := ParseLink(url)
	if err == nil {
		t.Fatal("expected error")
	}
	expectedError := QIsMissingErr{InputURL: url}
	if err.Error() != expectedError.Error() {
		t.Errorf("expected\n'%q'\n but got\n'%q'", expectedError, err)
	}
}

func TestLinkHandlerWithEncodedByGoogleUrlAfterPastedInTerminal(t *testing.T) {
	url := "https://www.google.com/url\\?q\\=https%3A%2F%2Fsendsafely.tester.com%2Freceive%2F%3Fthread%3DMYTHREAD%26packageCode%3DMYPKGCODE%23keyCode%3DMYKEYCODE\\&sa\\=D\\&ust\\=11111111\\&usg\\=JJJJJJJJJ"
	linkParts, err := ParseLink(url)
	if err != nil {
		t.Fatalf("unexected error '%v'", err)
	}
	expectedKeyCode := "MYKEYCODE"
	if linkParts.KeyCode != expectedKeyCode {
		t.Errorf("expected keycode '%v' but got '%v'", expectedKeyCode, linkParts.KeyCode)
	}

	expectedPackageCode := "MYPKGCODE"
	if linkParts.PackageCode != expectedPackageCode {
		t.Errorf("expected package code '%v' but got '%v'", expectedPackageCode, linkParts.PackageCode)
	}
}
func TestLinkHandlerWithEncodedByGoogleUrl(t *testing.T) {
	url := "https://www.google.com/url?q=https%3A%2F%2Fsendsafely.tester.com%2Freceive%2F%3Fthread%3DMYTHREAD%26packageCode%3DMYPKGCODE%23keyCode%3DMYKEYCODE&sa=D&ust=11111111&usg=JJJJJJJJJ"
	linkParts, err := ParseLink(url)
	if err != nil {
		t.Fatalf("unexected error '%v'", err)
	}
	expectedKeyCode := "MYKEYCODE"
	if linkParts.KeyCode != expectedKeyCode {
		t.Errorf("expected keycode '%v' but got '%v'", expectedKeyCode, linkParts.KeyCode)
	}

	expectedPackageCode := "MYPKGCODE"
	if linkParts.PackageCode != expectedPackageCode {
		t.Errorf("expected package code '%v' but got '%v'", expectedPackageCode, linkParts.PackageCode)
	}

}

func TestLinkHandler(t *testing.T) {
	url := "https://sendsafely.tester.com/receive/?thread=MYTHREAD&packageCode=MYPKGCODE#keyCode=MYKEYCODE"
	linkParts, err := ParseLink(url)
	if err != nil {
		t.Fatalf("unexected error '%v'", err)
	}
	expectedKeyCode := "MYKEYCODE"
	if linkParts.KeyCode != expectedKeyCode {
		t.Errorf("expected keycode '%v' but got '%v'", expectedKeyCode, linkParts.KeyCode)
	}

	expectedPackageCode := "MYPKGCODE"
	if linkParts.PackageCode != expectedPackageCode {
		t.Errorf("expected package code '%v' but got '%v'", expectedPackageCode, linkParts.PackageCode)
	}

}
func TestLinkHandlerLowerCase(t *testing.T) {
	url := "https://sendsafely.tester.com/receive/?thread=MYTHREAD&packagecode=MYPKGCODE#keycode=MYKEYCODE"
	linkParts, err := ParseLink(url)
	if err != nil {
		t.Fatalf("unexected error '%v'", err)
	}
	expectedKeyCode := "MYKEYCODE"
	if linkParts.KeyCode != expectedKeyCode {
		t.Errorf("expected keycode '%v' but got '%v'", expectedKeyCode, linkParts.KeyCode)
	}

	expectedPackageCode := "MYPKGCODE"
	if linkParts.PackageCode != expectedPackageCode {
		t.Errorf("expected package code '%v' but got '%v'", expectedPackageCode, linkParts.PackageCode)
	}

}

func TestKeyCodeMissing(t *testing.T) {
	url := "https://sendsafely.tester.com/receive/?thread=MYTHREAD&packageCode=MYPKGCODE"
	_, err := ParseLink(url)
	if err == nil {
		t.Fatalf("exected error '%v'", err)
	}
	if !errors.Is(err, KeyCodeIsMissingErr{
		InputURL: url,
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
		InputURL: url,
	}) {
		t.Errorf("expected PackageCodeIsMissingErr but got %v", err)
	} else {
		expectedError := "expected to have packageCode in url 'https://sendsafely.tester.com/receive/?thread=MYTHREAD#keyCode=MYKEYCODE' but it is not present"
		if err.Error() != expectedError {
			t.Errorf("expected error text of '%v' but was '%v'", expectedError, err.Error())
		}
	}
}

func TestInvalidUrlForGoogleURL(t *testing.T) {
	url := "https://www.google.com/url*$ù%?q=https%3A%2F%2Fsendsafely.tester.com%2Freceive%2F%3Fthread%3DMYTHREAD%26packageCode%3DMYPKGCODE%23keyCode%3DMYKEYCODE&sa=D&ust=11111111&usg=JJJJJJJJJ"
	_, err := ParseLink(url)
	if err == nil {
		t.Fatalf("exected error '%v'", err)
	}
	expectedError := URLParseErr{URL: url, BaseErr: errors.New("parse \"https://www.google.com/url*$ù%?q=https%3A%2F%2Fsendsafely.tester.com%2Freceive%2F%3Fthread%3DMYTHREAD%26packageCode%3DMYPKGCODE%23keyCode%3DMYKEYCODE&sa=D&ust=11111111&usg=JJJJJJJJJ\": invalid URL escape \"%\"")}
	if err.Error() != expectedError.Error() {
		t.Errorf("expected\n'%v'\n but got\n'%v'", expectedError, err)
	}
}

func TestInvalidURL(t *testing.T) {
	url := "*$ù%"
	_, err := ParseLink(url)
	if err == nil {
		t.Fatalf("exected error '%v'", err)
	}
	expectedError := URLParseErr{URL: url, BaseErr: errors.New("parse \"*$ù%\": invalid URL escape \"%\"")}
	if err.Error() != expectedError.Error() {
		t.Errorf("expected '%v' but got '%v'", expectedError, err)
	}
}
