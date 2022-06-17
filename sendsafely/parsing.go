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
package sendsafely

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/valyala/fastjson"
)

// SendSafelyApiParser stores the jsonParser so it can be shared between
// operations, this mainly benefits the parsing of the file parts
type SendSafelyApiParser struct {
	jsonParser fastjson.Parser
}

// SendSafelyFile are curious as near as they do not match the SendSafely documentation for what a package returns
// This was discovered by returning the values
type SendSafelyFile struct {
	FileId          string
	FileName        string
	FileSize        int64
	Parts           int
	FileUploaded    time.Time
	FileUploadedStr string
	FileVersion     string
	CreatedByEmail  string
}

// SendSafelyPackage is the struct we need that maps to the fields here:
// https://bump.sh/doc/sendsafely-rest-api#operation-getpackageinformation
// this is intentionally not complete as we do nto need all the fields
type SendSafelyPackage struct {
	PackageId        string
	PackageCode      string
	Files            []SendSafelyFile
	DirectoryIds     []string
	State            string
	PackageTimestamp time.Time
	Response         string
	ServerSecret     string
}

// SendSafelyDownloadUrl provides the part id and the actual url to get the file
// the Part field tells you the order of the parts so you can reconstruct the file
// after downloading it
// https://bump.sh/doc/sendsafely-rest-api#operation-post-package-parameter-file-parameter-download-urls
type SendSafelyDownloadUrl struct {
	Part int
	Url  string
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
//    { //NOTE THIS IS NOTHING LIKE WHAT IT RETURNS see SendSafelyFile for the accurate field names and types
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
func (s *SendSafelyApiParser) ParsePackage(packageJson string) (SendSafelyPackage, error) {
	var ssp SendSafelyPackage

	// if we were parsing lots of these we want to reuse the jsonParser to minimize allocations
	v, err := s.jsonParser.Parse(packageJson)
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

	serverSecret := v.Get("serverSecret")
	if !serverSecret.Exists() {
		return SendSafelyPackage{}, missingFieldError("serverSecret", packageJson)
	}
	ssp.ServerSecret = string(serverSecret.GetStringBytes())

	// looping through the id values for files
	var fileIds []SendSafelyFile
	filesArray := v.GetArray("files")
	for i, e := range filesArray {
		fileElement := e.Get("fileId")
		if !fileElement.Exists() {
			return SendSafelyPackage{}, fmt.Errorf("missing id in the %v element of the files array (indexed at 1). Array was '%v'", i+1, filesArray)
		}

		fileName := e.Get("fileName")
		if !fileName.Exists() {
			return SendSafelyPackage{}, fmt.Errorf("missing fileName in the %v element of the files array (indexed at 1). Array was '%v'", i+1, filesArray)
		}

		fileSize := e.Get("fileSize")
		if !fileSize.Exists() {
			return SendSafelyPackage{}, fmt.Errorf("missing fileSize in the %v element of the files array (indexed at 1). Array was '%v'", i+1, filesArray)
		}
		log.Printf("file size is %v", e.Get("fileSize"))

		parts := e.Get("parts")
		if !parts.Exists() {
			return SendSafelyPackage{}, fmt.Errorf("missing parts in the %v element of the files array (indexed at 1). Array was '%v'", i+1, filesArray)
		}

		createdByEmail := e.Get("createdByEmail")
		if !createdByEmail.Exists() {
			return SendSafelyPackage{}, fmt.Errorf("missing createdByEmail in the %v element of the files array (indexed at 1). Array was '%v'", i+1, filesArray)
		}

		fileUploadedRaw := e.Get("fileUploaded")
		if !fileUploadedRaw.Exists() {
			return SendSafelyPackage{}, fmt.Errorf("missing fileUploaded in the %v element of the files array (indexed at 1). Array was '%v'", i+1, filesArray)
		}

		// comes back in this format Jun 9, 2022 1:32:34 PM
		fileUploaded, err := time.Parse(DATE_FMT, string(fileUploadedRaw.GetStringBytes()))
		if err != nil {
			return SendSafelyPackage{}, fmt.Errorf("fileUploaded has the incorrect format and caused error '%v' in the %v element of the files array (indexed at 1). Array was '%v' and raw string was '%v'", err, i+1, filesArray, fileUploadedRaw)
		}

		fileUploadedStr := e.Get("fileUploadedStr")
		if !fileUploadedStr.Exists() {
			return SendSafelyPackage{}, fmt.Errorf("missing fileUploadedStr in the %v element of the files array (indexed at 1). Array was '%v'", i+1, filesArray)
		}

		fileVersion := e.Get("fileVersion")
		if !fileVersion.Exists() {
			return SendSafelyPackage{}, fmt.Errorf("missing fileVersion in the %v element of the files array (indexed at 1). Array was '%v'", i+1, filesArray)
		}
		fileSizeInt, err := strconv.ParseInt(string(fileSize.GetStringBytes()), 10, 64)
		if err != nil {
			return SendSafelyPackage{}, fmt.Errorf("unable to convert fileSize field with value '%v' into int due to error '%v'", string(fileSize.GetStringBytes()), err)
		}
		fileIds = append(fileIds, SendSafelyFile{
			FileId:          string(fileElement.GetStringBytes()),
			FileName:        string(fileName.GetStringBytes()),
			FileSize:        fileSizeInt,
			Parts:           int(parts.GetInt64()),
			CreatedByEmail:  string(createdByEmail.GetStringBytes()),
			FileUploaded:    fileUploaded,
			FileUploadedStr: string(fileUploadedStr.GetStringBytes()),
			FileVersion:     string(fileVersion.GetStringBytes()),
		})
	}
	ssp.Files = fileIds

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
	ts, err := time.Parse(DATE_FMT, string(rawTimestamp))
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

const DATE_FMT = "Jan 2, 2006 3:04:05 PM"

// ParseDownloadUrls reads the json response provided here https://bump.sh/doc/sendsafely-rest-api#operation-post-package-parameter-file-parameter-download-urls
// here is an example
// {
//   "downloadUrls": [
//     {
//       "part": 1,
//       "url": "https://sendsafely-dual-region-us.s3-accelerate.amazonaws.com/commercial/AAAAAAAA-BBBB-CCCC-DDDD-EEEEEEEEEEEE/11111111-2222-3333-4444-555555555555-1?AWSAccessKeyId=AKIAIOSFODNN7EXAMPLE&Expires=1554862678&Signature=OTP5Z0DIutXKbRRT4NwmxQG9jFk%3D"
//     }
//   ],
//   "response": "SUCCESS"
// }
func (s *SendSafelyApiParser) ParseDownloadUrls(downloadJson string) ([]SendSafelyDownloadUrl, error) {
	var response []SendSafelyDownloadUrl
	v, err := s.jsonParser.Parse(downloadJson)
	if err != nil {
		return []SendSafelyDownloadUrl{}, fmt.Errorf("unexpected error parsing downloadUrls json string '%v' with error '%v'", downloadJson, err)
	}
	responseStatus := v.GetStringBytes("response")
	if string(responseStatus) != "SUCCESS" {
		log.Println(downloadJson)
		message := v.Get("message")
		if !message.Exists() {
			return []SendSafelyDownloadUrl{}, fmt.Errorf("unexpected response from json with response status '%v', full json was '%v'", string(responseStatus), downloadJson)
		}
		return []SendSafelyDownloadUrl{}, fmt.Errorf("failed download due to %v %v", string(responseStatus), message)
	}

	downloadUrls := v.GetArray("downloadUrls")
	for _, e := range downloadUrls {
		if e.Exists() {
			part := e.GetInt("part")
			url := string(e.GetStringBytes("url"))
			response = append(response, SendSafelyDownloadUrl{
				Part: part,
				Url:  url,
			})
		}
	}
	return response, nil
}
