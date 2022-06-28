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
package zendesk

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/valyala/fastjson"
	"golang.org/x/net/html"
)

// GetLinksFromComments is parsing out the links from the html_
// docs are here https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_comments/#list-comments
func GetLinksFromComments(jsonData string) ([]string, error) {
	jsonParser := fastjson.Parser{}
	result, err := jsonParser.Parse(jsonData)
	if err != nil {
		return []string{}, fmt.Errorf("parsing json data '%v' failed, error was '%v'", jsonData, err)
	}
	commentsValue := result.Get("comments")
	if !commentsValue.Exists() {
		return []string{}, fmt.Errorf("parsing json data '%v' failed, missing comments field", jsonData)
	}
	comments, err := commentsValue.Array()
	if err != nil {
		return []string{}, fmt.Errorf("parsing comments for jsonData '%v' failed, error was '%v'", jsonData, err)
	}
	var linksFound []string
	for i, comment := range comments {
		htmlBodyValue := comment.Get("html_body")
		if !htmlBodyValue.Exists() {
			return []string{}, fmt.Errorf("parsing json data '%v' failed, missing html_body field for the '%v' comment (base index 0)", jsonData, i)
		}
		htmlBody := htmlBodyValue.GetStringBytes()
		z := html.NewTokenizer(bytes.NewBuffer(htmlBody))

		for {
			if z.Next() == html.ErrorToken {
				// Returning io.EOF indicates success.
				err = z.Err()
				if errors.Is(err, io.EOF) {
					break
				}
				return []string{}, fmt.Errorf("parsing html '%v' failed, with error %v for the '%v' comment (base index 0)", jsonData, z.Err(), i)
			}
			token := z.Token()
			if token.Data == "a" {
				for _, a := range token.Attr {
					if a.Key == "href" {
						linksFound = append(linksFound, a.Val)
						break
					}
				}
			}
		}
	}
	return linksFound, nil
}

// Attachment maps to
// 				{
//					"url": "https://dremio.zendesk.com/api/v2/attachments/1.json",
//					"id": 1,
//					"file_name": "test.txt",
//					"content_url": "https://tester.zendesk.com/attachments/token/abc/?name=test.txt",
//					"mapped_content_url": "https://test.tester.com/attachments/token/abc/?name=test.txt",
//					"content_type": "text/plain",
//					"size": 1111,
//					"width": null,
//					"height": null,
//					"inline": false,
//					"deleted": false,
//					"thumbnails": []
//				}
// with additional data from parent comment
type Attachment struct {
	ParentCommentDate time.Time // "created_at": "2000-01-01T11:11:07Z",
	ParentCommentID   int64
	FileName          string
	ContentURL        string
	ContentType       string
	Size              int64
	Deleted           bool
}

// GetLinksFromComments is parsing out the links from the html_
// docs are here https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_comments/#list-comments
func GetAttachementsFromComments(jsonData string) ([]Attachment, error) {

	jsonParser := fastjson.Parser{}
	result, err := jsonParser.Parse(jsonData)
	if err != nil {
		return []Attachment{}, fmt.Errorf("parsing json data '%v' failed, error was '%v'", jsonData, err)
	}
	commentsValue := result.Get("comments")
	if !commentsValue.Exists() {
		return []Attachment{}, fmt.Errorf("parsing json data '%v' failed, missing comments field", jsonData)
	}
	comments, err := commentsValue.Array()
	if err != nil {
		return []Attachment{}, fmt.Errorf("parsing comments for jsonData '%v' failed, error was '%v'", jsonData, err)
	}
	var attachments []Attachment

	for i, comment := range comments {
		parentIdValue := comment.Get("id")
		parentId, err := parentIdValue.Int64()
		if err != nil {
			return []Attachment{}, fmt.Errorf("parsing attachments for json data '%v' failed due to error '%v' for 'id' field a comments index %v", jsonData, err, i)
		}
		parentCreatedAtValue := comment.Get("created_at")
		parentCreatedAtRaw, err := parentCreatedAtValue.StringBytes()
		if err != nil {
			return []Attachment{}, fmt.Errorf("parsing attachments for json data '%v' failed due to error '%v' for  'created_at' field a comments index %v", jsonData, err, i)
		}
		createdAt, err := time.Parse(time.RFC3339, string(parentCreatedAtRaw))
		if err != nil {
			return []Attachment{}, fmt.Errorf(" parsing datetime for '%v' failed' due to error '%v' field a comments index %v", jsonData, err, i)
		}

		attachmentsValues := comment.Get("attachments")
		if !attachmentsValues.Exists() {
			return []Attachment{}, fmt.Errorf("parsing comments for jsonData '%v' failed, missing attachments field in comment index %v", jsonData, i)
		}
		attachmentsFromJson := attachmentsValues.GetArray()
		for ai, a := range attachmentsFromJson {
			fileNameValue := a.Get("file_name")
			fileNameBytes, err := fileNameValue.StringBytes()
			if err != nil {
				return []Attachment{}, fmt.Errorf("parsing attachments for json data '%v' failed due to error '%v' for 'file_name' field at attachments index %v and comments index %v", jsonData, err, ai, i)
			}
			fileName := string(fileNameBytes)
			boolValue := a.Get("deleted")
			isDeleted, err := boolValue.Bool()
			if err != nil {
				return []Attachment{}, fmt.Errorf("parsing attachments for json data '%v' failed due to error '%v' for 'deleted' field at attachments index %v and comments index %v", jsonData, err, ai, i)
			}

			contentUrlValue := a.Get("content_url")
			contentUrlBytes, err := contentUrlValue.StringBytes()
			if err != nil {
				return []Attachment{}, fmt.Errorf("parsing attachments for json data '%v' failed due to error '%v' for 'content_url' field at attachments index %v and comments index %v", jsonData, err, ai, i)
			}

			contentUrl := string(contentUrlBytes)
			contentTypeValue := a.Get("content_type")
			contentTypeBytes, err := contentTypeValue.StringBytes()
			if err != nil {
				return []Attachment{}, fmt.Errorf("parsing attachments for json data '%v' failed due to error '%v' for 'content_type' field at attachments index %v and comments index %v", jsonData, err, ai, i)
			}
			contentType := string(contentTypeBytes)

			sizeValue := a.Get("size")
			size, err := sizeValue.Int64()
			if err != nil {
				return []Attachment{}, fmt.Errorf("parsing attachments for json data '%v' failed due to error '%v' for 'size' field at attachments index %v and comments index %v", jsonData, err, ai, i)
			}
			attachments = append(attachments, Attachment{
				ParentCommentID:   parentId,
				ParentCommentDate: createdAt,
				FileName:          fileName,
				Deleted:           isDeleted,
				ContentURL:        contentUrl,
				ContentType:       contentType,
				Size:              size,
			})
		}
	}
	return attachments, nil

}
