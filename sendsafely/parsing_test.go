package sendsafely

import (
	"testing"
)

func TestParsing(t *testing.T) {
	//taken from https://bump.sh/doc/sendsafely-rest-api#operation-getpackageinformation-200-approverlist
	p, err := ParsePackage(`{
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
			"id": "abcfile"
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

	if p.PackageId != "GVG2-MNZT" {
		t.Errorf("unexpected package id %v", p.PackageId)
	}

	if p.PackageCode != "M0AEMIrTQe9XWRgGDKiKta1pXobmpKwAVafWgXjnBsw" {
		t.Errorf("unexpected package code %v", p.PackageCode)
	}

	lenFileIds := len(p.FileIds)
	if lenFileIds != 1 {
		t.Errorf("was expected 1 element but found %v", lenFileIds)
	}

	if lenFileIds > 0 {
		if p.FileIds[0] != "abcfile" {
			t.Errorf("was expected abcfile but found %v", p.FileIds[0])
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
