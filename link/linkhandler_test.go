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

import "testing"

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
