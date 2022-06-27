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
package downloader

import (
	"fmt"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestDownloadFile(t *testing.T) {

	// pass in the resty httpy client that the SendSafelyClient uses so that
	// httpmock can replace it's transport parameter with a mock one
	// preventing remote calls from going to SendSafely
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	resp := `{"response":"SUCCESS"}`

	// Exact URL match
	url := "http://example.com"
	responder := httpmock.NewStringResponder(200, resp) //yes they really log 200 when you get an error

	httpmock.RegisterResponder("GET", url, responder)
	fileName := fmt.Sprintf("%v/testFile.json", t.TempDir())
	err := DownloadFile(fileName, url)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if string(b) != resp {
		t.Errorf("expected %v but was %v", resp, string(b))
	}
}
