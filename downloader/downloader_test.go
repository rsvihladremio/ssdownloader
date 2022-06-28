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
	"errors"
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
	d := NewGenericDownloader(4096)
	err := d.DownloadFile(fileName, url)
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

func TestInvalidBufferSizeResultsInDefault(t *testing.T) {
	d := NewGenericDownloader(-1)
	if d.bufferSizeKB != 4096 {
		t.Errorf("expected 4096 but was %v", d.bufferSizeKB)
	}

	d = NewGenericDownloader(0)
	if d.bufferSizeKB != 4096 {
		t.Errorf("expected 4096 but was %v", d.bufferSizeKB)
	}
}

func TestUsingDefaultBufferSizeResultsInError(t *testing.T) {
	d := GenericDownloader{}
	err := d.DownloadFile("", "")
	if err == nil {
		t.Error("expected an error but there was none")
	}
	if !errors.Is(err, IllegalBufferSize{}) {
		t.Errorf("expected error of %v but was %v", IllegalBufferSize{}, err)
	}
}
