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

//sendsafely package decrypts files, combines file parts into whole files, and handles api access to the sendsafely rest api
package sendsafely

import (
	"testing"
)

func TestParsingDownloadUrls(t *testing.T) {
	parser := APIParser{}
	d, err := parser.ParseDownloadUrls(`{
		"downloadUrls": [
		  {
			"part": 1,
			"url": "https://myserver/commercial/part?AWSAccessKeyId=ExampleKey&Expires=1554862678&Signature=MySignature"
		  }
		],
		"response": "SUCCESS"
	  }`)
	if err != nil {
		t.Fatalf("unable to parse download url %v", err)
	}
	if len(d) != 1 {
		t.Fatalf("expected %v download url but was %v", 1, len(d))
	}
	part := d[0].Part
	expectedPart := 1
	if part != expectedPart {
		t.Errorf("Part expected to be %v but was %v", expectedPart, part)
	}
	url := d[0].URL
	expectedURL := "https://myserver/commercial/part?AWSAccessKeyId=ExampleKey&Expires=1554862678&Signature=MySignature"
	if url != expectedURL {
		t.Errorf("Url expected to be %v but was %v", expectedURL, url)
	}
}

func TestParsingPackage(t *testing.T) {
	parser := APIParser{}
	//taken from https://bump.sh/doc/sendsafely-rest-api#operation-getpackageinformation-200-approverlist
	p, err := parser.ParsePackage("GVG2-MNZT", `{
		"packageId": "GVG2-MNZT",
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
			"fileId": "abcfile",
			"fileName": "fname",
			"fileSize": "12",
			"parts": 1,
			"createdByEmail": "test@test.com",
			"fileUploaded": "Jun 9, 2022 1:32:34 PM",
			"fileUploadedStr": "Jun 9, 2022 1:32:34 PM",
			"fileVersion" : "1"
		  }
		],
		"directories": [
		  {
			"id": "abcdir"
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
	  }`)
	if err != nil {
		t.Fatalf("unable to parse with error %v", err)
	}

	if p.PackageID != "GVG2-MNZT" {
		t.Errorf("unexpected package id %v", p.PackageID)
	}

	if p.PackageCode != "M0AEMIrTQe9XWRgGDKiKta1pXobmpKwAVafWgXjnBsw" {
		t.Errorf("unexpected package code %v", p.PackageCode)
	}

	lenFileIds := len(p.Files)
	if lenFileIds != 1 {
		t.Errorf("was expected 1 element but found %v", lenFileIds)
	}

	if lenFileIds > 0 {
		if p.Files[0].FileID != "abcfile" {
			t.Errorf("was expected abcfile but found %v", p.Files[0].FileID)
		}
	}

	lenDirIds := len(p.DirectoryIds)
	if lenDirIds != 1 {
		t.Errorf("was expected 1 element but found %v", lenDirIds)
	}

	if lenDirIds > 0 {
		if p.DirectoryIds[0] != "abcdir" {
			t.Errorf("was expected abcfile but found %v", p.DirectoryIds[0])
		}
	}

	ts := p.PackageTimestamp.String()
	//original format Feb 1, 2019 2:07:28 PM",
	if ts != "2019-02-01 14:07:28 +0000 UTC" {
		t.Errorf("unexpected packageTimestamp %v", ts)
	}

	response := p.Response
	if response != "SUCCESS" {
		t.Errorf("unexpeced response expected SUCCESS but got %v", response)
	}

	state := p.State
	if state != "PACKAGE_STATE_IN_PROGRESS" {
		t.Errorf("unexpected state, expected PACKAGE_STATE_IN_PROGRESS but got %v", state)
	}
}
