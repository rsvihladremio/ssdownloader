package sendsafely

import (
	"fmt"
	"time"

	"github.com/valyala/fastjson"
)

// SendSafelyPackage is the struct we need that maps to the fields here:
// https://bump.sh/doc/sendsafely-rest-api#operation-getpackageinformation
// this is intentionally not complete as we do nto need all the fields
type SendSafelyPackage struct {
	PackageId        string
	PackageCode      string
	FileIds          []string
	DirectoryIds     []string
	State            string
	PackageTimestamp time.Time
	Response         string
}

func missingFieldError(fieldName, jsonBody string) error {
	return fmt.Errorf("unable to get %v from '%v'", fieldName, jsonBody)
}

// ParsePackage reads a json from here https://bump.sh/doc/sendsafely-rest-api#operation-getpackageinformation
// which looks like the following
// {
//  "packageId": "GVG2-MNZT",
//  "packageCode": "M0AEMIrTQe9XWRgGDKiKta1pXobmpKwAVafWgXjnBsw",
//  "serverSecret": "ACbuj9NKTkvjZ71Gc0t5zuU1xvba9XAouA",
//  "recipients": [
//    {
//      "recipientId": "5d504769-78c4-4c0a-b982-945845ea2075",
//      "email": "recip1@example.com",
//      "fullName": "External User",
//      "needsApproval": false,
//      "recipientCode": "YN0P1G0xbS9mBSwohP9xPJSqwgKXMq4bCI5uTcx1KKM",
//      "confirmations": {
//        "ipAddress": "127.0.0.1",
//        "timestamp": "Dec 12, 2018 2:24:38 PM",
//        "timeStampStr": "Dec 12, 2018 at 14:24",
//        "isMessage": true
//      },
//      "isPackageOwner": false,
//      "checkForPublicKeys": false,
//      "roleName": "VIEWER"
//    }
//  ],
//  "contactGroups": [
//    {
//      "id": "string"
//    }
//  ],
//  "files": [
//    {
//      "id": "string"
//    }
//  ],
//  "directories": [
//    {
//      "id": "string"
//    }
//  ],
//  "approverList": [
//    {}
//  ],
//  "needsApproval": false,
//  "state": "PACKAGE_STATE_IN_PROGRESS",
//  "passwordRequired": false,
//  "life": 10,
//  "isVDR": false,
//  "isArchived": false,
//  "packageSender": "user@companyabc.com",
//  "packageTimestamp": "Feb 1, 2019 2:07:28 PM",
//  "rootDirectoryId": "8c3c2184-e73e-4137-be92-e9c5b5661258",
//  "response": "SUCCESS"
//}
func ParsePackage(packageJson string) (SendSafelyPackage, error) {
	var ssp SendSafelyPackage
	// if we were parsing lots of these we want to reuse the jsonParser to minimize allocations
	var jsonParser fastjson.Parser
	v, err := jsonParser.Parse(packageJson)
	if err != nil {
		return SendSafelyPackage{}, fmt.Errorf("unexpected error parsing package json string '%v' with error '%v'", packageJson, err)
	}

	packageId := v.Get("packageId")
	if !packageId.Exists() {
		return SendSafelyPackage{}, missingFieldError("packageId", packageJson)
	}
	ssp.PackageId = string(packageId.GetStringBytes())

	packageCode := v.Get("packageCode")
	if !packageCode.Exists() {
		return SendSafelyPackage{}, missingFieldError("packageCode", packageJson)
	}
	ssp.PackageCode = string(packageCode.GetStringBytes())

	// looping through the id values for files
	var fileIds []string
	filesArray := v.GetArray("files")
	for i, e := range filesArray {
		fileElement := e.Get("id")
		if !fileElement.Exists() {
			return SendSafelyPackage{}, fmt.Errorf("missing id in the %v element of the files array (indexed at 1)", i+1)
		}
		fileIds = append(fileIds, string(fileElement.GetStringBytes()))
	}
	ssp.FileIds = fileIds

	// looping through the id values for directories
	var directoryIds []string
	directoriesArray := v.GetArray("directories")
	for i, e := range directoriesArray {
		// this is the only value we are interested in
		directoryElement := e.Get("id")
		if !directoryElement.Exists() {
			return SendSafelyPackage{}, fmt.Errorf("missing id in the %v element of the directories array (indexed at 1)", i+1)
		}
		directoryIds = append(directoryIds, string(directoryElement.GetStringBytes()))
	}
	ssp.DirectoryIds = directoryIds

	// this is the package state, we may or may not need this, at minimum it should be useful for logging
	state := v.Get("state")
	if !packageCode.Exists() {
		return SendSafelyPackage{}, missingFieldError("state", packageJson)
	}
	ssp.State = string(state.GetStringBytes())

	// this is the packageTimestamp also primarily intended for logging
	packageTimestamp := v.Get("packageTimestamp")
	if !packageTimestamp.Exists() {
		return SendSafelyPackage{}, missingFieldError("packageTimestamp", packageJson)
	}
	rawTimestamp := packageTimestamp.GetStringBytes()
	// the format is rather curious but this is what sendsafely is providing, I can find no standard that matches this
	// example "Feb 1, 2019 2:07:28 PM"
	ts, err := time.Parse("Jan 2, 2006 3:04:05 PM", string(rawTimestamp))
	if err != nil {
		return SendSafelyPackage{}, fmt.Errorf("unparseable packageTimestamp '%v'", err)
	}
	ssp.PackageTimestamp = ts

	// response success or failure, also primarily useful logging and longer term I can see this being used for
	response := v.Get("response")
	if !response.Exists() {
		return SendSafelyPackage{}, missingFieldError("response", packageJson)
	}
	ssp.Response = string(response.GetStringBytes())
	return ssp, nil
}