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

//zendesk package provides api access to the zendesk rest apis
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

//ParserErr provides location, raw json data parsed and nested error
type ParserErr struct {
	Err      error
	JSONData string
	Location string
}

//Error provides location if one is present otherwise it will omit that text
func (p ParserErr) Error() string {
	if p.Location == "" {
		// with no location we can return a shorter cleaner message
		return fmt.Sprintf("parsing json data '%v' failed, error was '%v'", p.JSONData, p.Err)
	}
	return fmt.Sprintf("parsing json data '%v' failed for '%v', error was '%v'", p.JSONData, p.Location, p.Err)
}

//MissingJSONFieldError provides location, field name and raw json data parsed
type MissingJSONFieldError struct {
	FieldName string
	JSONData  string
	Location  string
}

//Error provides location if one is present otherwise it will omit that text
func (m MissingJSONFieldError) Error() string {
	if m.Location == "" {
		// with no location we can return a shorter cleaner message
		return fmt.Sprintf("parsing json data '%v' failed, missing '%v' field", m.JSONData, m.FieldName)
	}
	return fmt.Sprintf("parsing json data '%v' missing field '%v' in '%v'", m.JSONData, m.FieldName, m.Location)
}

// GetLinksFromComments is parsing out the links from the html_
// docs are here https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_comments/#list-comments
func GetLinksFromComments(jsonData string) ([]string, error) {
	// using fastjson instead of the default golang json encoding libraries, fastjson can be 15 faster qnd
	jsonParser := fastjson.Parser{}
	result, err := jsonParser.Parse(jsonData)
	if err != nil {
		// this usually means the json is not to spec and is invalid, return the json back to the client for analysis
		return []string{}, ParserErr{
			Err:      err,
			JSONData: jsonData,
		}
	}
	// read comments field
	commentsValue := result.Get("comments")
	if !commentsValue.Exists() {
		return []string{}, MissingJSONFieldError{
			JSONData:  jsonData,
			FieldName: "comments",
		}
	}
	// try and convert comments into an array (which is expected)
	comments, err := commentsValue.Array()
	if err != nil {
		// if the comments value is somehow not an array return an error back to the client
		return []string{}, ParserErr{
			Err:      err,
			JSONData: jsonData,
			Location: "comments",
		}
	}
	var linksFound []string
	//search the html_body of all the comments
	for i, comment := range comments {
		htmlBodyValue := comment.Get("html_body")
		// if we get no html_body then this is failed parse and we are missing some data
		if !htmlBodyValue.Exists() {
			return []string{}, MissingJSONFieldError{
				JSONData:  jsonData,
				FieldName: "html_body",
				Location:  fmt.Sprintf("comment %v (base index 0)", i),
			}
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
				// return error with location of error so the client can diagnosis the issue
				return []string{}, ParserErr{
					Err:      err,
					JSONData: jsonData,
					Location: fmt.Sprintf("html_body field for the comment %v (base index 0)", i),
				}
			}
			token := z.Token()
			// if it is a link search for the href
			if token.Data == "a" {
				for _, a := range token.Attr {
					if a.Key == "href" {
						// if we find ANY href go ahead and add it to result set
						// this is to allow future searching of different kinds of links in the text
						// filtering hqppens later for sendsafely links
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
func GetAttachmentsFromComments(jsonData string) ([]Attachment, error) {

	jsonParser := fastjson.Parser{}
	result, err := jsonParser.Parse(jsonData)
	if err != nil {
		return []Attachment{}, ParserErr{
			Err:      err,
			JSONData: jsonData,
		}
	}
	commentsValue := result.Get("comments")
	if !commentsValue.Exists() {
		return []Attachment{}, MissingJSONFieldError{
			JSONData:  jsonData,
			FieldName: "comments",
		}
	}
	comments, err := commentsValue.Array()
	if err != nil {
		return []Attachment{}, ParserErr{
			Err:      err,
			JSONData: jsonData,
			Location: "comments",
		}
	}
	var attachments []Attachment

	for i, comment := range comments {
		parentIDValue := comment.Get("id")
		parentID, err := parentIDValue.Int64()
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
		attachmentsFromJSON := attachmentsValues.GetArray()
		for ai, a := range attachmentsFromJSON {
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

			contentURLValue := a.Get("content_url")
			contentURLBytes, err := contentURLValue.StringBytes()
			if err != nil {
				return []Attachment{}, fmt.Errorf("parsing attachments for json data '%v' failed due to error '%v' for 'content_url' field at attachments index %v and comments index %v", jsonData, err, ai, i)
			}

			contentURL := string(contentURLBytes)
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
				ParentCommentID:   parentID,
				ParentCommentDate: createdAt,
				FileName:          fileName,
				Deleted:           isDeleted,
				ContentURL:        contentURL,
				ContentType:       contentType,
				Size:              size,
			})
		}
	}
	return attachments, nil

}
