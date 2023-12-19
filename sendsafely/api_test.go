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
	"strings"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

// This is the default happy path test, no errors
func TestRetrievePackgeById(t *testing.T) {
	// since we are using a mock http api we can use any api secret we feel like
	ssClient := NewClient("myApiKey", "mySecret", false).(*DownloadClient)

	// pass in the resty httpy client that the SendSafelyClient uses so that
	// httpmock can replace it's transport parameter with a mock one
	// preventing remote calls from going to SendSafely
	httpmock.ActivateNonDefault(ssClient.client.GetClient())
	// make sure to reset the mock after the test
	defer httpmock.DeactivateAndReset()
	packageID := "ABDC-DDFAF"
	// using the sample json from the sendafely website, with the files field correctly
	// filled out (the sendsafely site had the "files" field incorrectly documented)
	resp := `{
		  "packageId": "ABDC-DDFAF",
		  "packageCode": "M0AEMIrTQe9XWRgGDKiKta1pXobmpKwAVafWgXjnBsw",
		  "serverSecret": "ACbuj9NKTkvjZ71Gc0t5zuU1xvba9XAouA",
		  "recipients": [
		    {
		      "recipientId": "5d504769-78c4-4c0a-b982-945845ea2075",
		      "email": "recip1@example.com",
		      "fullName": "External User",
		      "needsApproval": false,
		      "recipientCode": "YN0P1G0xbS9mBSwohP9xPJSqwgKXMq4bCI5uTcx1KKM",
		      "confirmations": {
		        "ipAddress": "127.0.0.1",
		        "timestamp": "Dec 12, 2018 2:24:38 PM",
		        "timeStampStr": "Dec 12, 2018 at 14:24",
		        "isMessage": true
		      },
		      "isPackageOwner": false,
		      "checkForPublicKeys": false,
		      "roleName": "VIEWER"
		    }
		  ],
		  "contactGroups": [
		    {
		      "id": "string"
		    }
		  ],
		  "files": [
		    {
		      "fileId": "string",
			  "fileName": "testName",
			  "fileSize": "12344",
			  "parts": 1,
			  "fileUploaded": "Feb 9, 2022 11:14:24 PM",
			  "fileUploadedStr": "Feb 9, 2022 11:14:24 PM",
			  "createdByEmail": "test@tester.com",
			  "fileVersion": "12"
		    }
		  ],
		  "directories": [
		    {
		      "id": "string"
		    }
		  ],
		  "approverList": [
		    {}
		  ],
		  "needsApproval": false,
		  "state": "PACKAGE_STATE_IN_PROGRESS",
		  "passwordRequired": false,
		  "life": 10,
		  "isVDR": false,
		  "isArchived": false,
		  "packageSender": "user@companyabc.com",
		  "packageTimestamp": "Feb 1, 2019 2:07:28 PM",
		  "rootDirectoryId": "8c3c2184-e73e-4137-be92-e9c5b5661258",
		  "response": "SUCCESS"
		}`

	// setup a responder to the expected status code of 200 and then returning the json data setup abovee
	responder := httpmock.NewStringResponder(200, resp)

	url := strings.Join([]string{URL, "package", packageID}, "/")
	// we are expecting a GET request with the exact url specified above, if that exact match happens
	// the json body setup in the responder will return instead of hitting the remote sendsafely server
	httpmock.RegisterResponder("GET", url, responder)
	pkg, err := ssClient.RetrievePackageByID(packageID)
	if err != nil {
		t.Fatalf("unexpected error retrieving id '%v'", err)
	}
	//choosing not to retest the entire json parsing for this as it is already covered in other tests in
	// github.com/rsvihladremio/ssdownloader/sendsafely under the parsing_test.go file
	if pkg.PackageID != packageID {
		t.Errorf("expected packageId '%v' but was '%v'", packageID, pkg.PackageID)
	}
}

// the missing package case, this mimics the actual production api as of 2022-07-12
func TestRetrievePackageIsMissing(t *testing.T) {
	// since we are using a mock http api we can use any api secret we feel like
	ssClient := NewClient("myApiKey", "mySecret", false).(*DownloadClient)

	// pass in the resty httpy client that the SendSafelyClient uses so that
	// httpmock can replace it's transport parameter with a mock one
	// preventing remote calls from going to SendSafely
	httpmock.ActivateNonDefault(ssClient.client.GetClient())
	defer httpmock.DeactivateAndReset()
	packageID := "ABDC-DDFAF"
	resp := `{"needsApproval":false,"passwordRequired":false,"life":0,"isVDR":false,"isArchived":false,"totalDirectories":0,"totalFiles":0,"allowReplyAll":false,"packageContainsMessage":false,"response":"UNKNOWN_PACKAGE","message":"Package ID does not exist"}`

	// Exact URL match
	url := strings.Join([]string{URL, "package", packageID}, "/")
	responder := httpmock.NewStringResponder(200, resp) //yes they really log 200 when you get an error

	httpmock.RegisterResponder("GET", url, responder)
	_, err := ssClient.RetrievePackageByID(packageID)
	if err == nil {
		t.Fatal("expected error retrieving id")
	}
	expectedError := "unable to find package ABDC-DDFAF as it is likely expired"
	if err.Error() != expectedError {
		t.Errorf("expected error '%v' but was '%v'", expectedError, err)
	}
}

// the bad auth case, this mimics the actual production api as of 2022-06-20
func TestRetrievePackageHasBadAuth(t *testing.T) {
	// since we are using a mock http api we can use any api secret we feel like
	ssClient := NewClient("myApiKey", "mySecret", false).(*DownloadClient)

	// pass in the resty httpy client that the SendSafelyClient uses so that
	// httpmock can replace it's transport parameter with a mock one
	// preventing remote calls from going to SendSafely
	httpmock.ActivateNonDefault(ssClient.client.GetClient())
	defer httpmock.DeactivateAndReset()
	packageID := "ABDC-DDFAF"
	resp := `{"response":"AUTHENTICATION_FAILED","message":"Invalid API Key"}`

	// Exact URL match
	url := strings.Join([]string{URL, "package", packageID}, "/")
	responder := httpmock.NewStringResponder(200, resp) //yes they really log 200 when you get an error

	httpmock.RegisterResponder("GET", url, responder)
	_, err := ssClient.RetrievePackageByID(packageID)
	if err == nil {
		t.Fatal("expected error retrieving id")
	}
	expectedError := "failed authentication for package ABDC-DDFAF due to 'Invalid API Key'"
	if err.Error() != expectedError {
		t.Errorf("expected error '%v' but was '%v'", expectedError, err)
	}
}

// this tests the generated hash for the header ss-request-signature,
// the main purpose of this test is not to explain the function but to lock in time the behavior
// so that if there is a breaking change we will catch it
func TestGenerateSignature(t *testing.T) {
	ssClient := NewClient("", "", false).(*DownloadClient)
	ts, err := time.Parse(time.RFC3339, "2022-05-31T18:11:21Z")
	if err != nil {
		t.Fatalf("bad test setup since we were not able to use our datetime due to error '%v'", err)
	}
	sig, err := ssClient.generateRequestSignature(ts.String(), "/api/v2/packages", `{ "data" : 1}`)
	if err != nil {
		t.Fatalf("unexected error '%v'", err)
	}
	// calculated this, not very meaningful to read, but this will lock the tested behavior and guard against
	// any future regressions
	expectedSig := "aa08337f9fe994bad07a441b447abb18a0c2f43e11394dfc2f530eb898920328"
	if expectedSig != sig {
		t.Errorf("expected '%v' but was '%v'", expectedSig, sig)
	}
	//validate it is not always the same value
	sig2, err := ssClient.generateRequestSignature(time.Now().String(), "/api/v2/packages", `{ "data" : 1}`)
	if err != nil {
		t.Fatalf("unexected error '%v'", err)
	}
	if sig == sig2 {
		t.Error("signature did not change with new time	")
	}
}

// this tests the generated hash for the checksum field in the body
// the main purpose of this test is not to explain the function but to lock in time the behavior
// so that if there is a breaking change we will catch it
func TestGenerateCheckSum(t *testing.T) {
	ssClient := NewClient("", "", false).(*DownloadClient)
	checkSum := ssClient.generateChecksum("abc", "def")

	// calculated this, not very meaningful to read, but this will lock the tested behavior and guard against
	// any future regressions
	expectedChkSum := "34ffd317a709b2ac3a848328888b8e8ed56658ac5d3f32512071579576beaa02"
	if expectedChkSum != checkSum {
		t.Errorf("expected '%v' but was '%v'", expectedChkSum, checkSum)
	}
	//validate it is always the same value
	checkSum2 := ssClient.generateChecksum("abc", "def")

	if checkSum != checkSum2 {
		t.Error("signature changed and is not deterministic")
	}
}
