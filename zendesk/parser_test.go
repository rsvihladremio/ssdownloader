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

//zendesk package provides api access to the zendesk rest api
package zendesk

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestGetLinksFromComments(t *testing.T) {
	links, err := GetLinksFromComments(` {
		 	"comments": [
		 	  {
		 		"attachments": [],
				"plain_body": "Thanks for your help!example",
				"html_body": "<p>Thanks for your help!</p><a href='http://example.com'>example</a>"
			  },
			  {
				"attachments": [],
			   "plain_body": "here is some more for your help!example 2",
			   "html_body": "<p>here is some more for your help!</p><a href='http://example.com/file2'>example 2</a>"
			 }
			]
		 }`)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	expectedLinks := []CommentTextWithLink{
		{URL: "http://example.com", Body: "Thanks for your help!example"},
		{URL: "http://example.com/file2", Body: "here is some more for your help!example 2"},
	}
	if !reflect.DeepEqual(expectedLinks, links) {
		t.Errorf("expected %v but had %v", expectedLinks, links)
	}
}

func TestGetLinksFromCommentsHasInvalidJson(t *testing.T) {
	_, err := GetLinksFromComments(`{}`)
	if err == nil {
		t.Error("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(MissingJSONFieldError{}) {
		t.Errorf("expected MissingJSONFieldError but was %T", err)
	}
	expectedErr := "parsing json data '{}' failed, missing 'comments' field"
	if err.Error() != expectedErr {
		t.Errorf("expected error text '%q' but was %q", expectedErr, err.Error())
	}
}

func TestGetLinksFromCommentsIsMissingComments(t *testing.T) {
	_, err := GetLinksFromComments(``)
	if err == nil {
		t.Error("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(ParserErr{}) {
		t.Errorf("expected ParserErr but was %T", err)
	}
	expectedErr := "parsing json data '' failed, error was 'cannot parse JSON: cannot parse empty string; unparsed tail: \"\"'"
	if err.Error() != expectedErr {
		t.Errorf("expected error text '%q' but was %q", expectedErr, err.Error())
	}
}

func TestGetLinksFromCommentsHasInvalidCommentsField(t *testing.T) {
	_, err := GetLinksFromComments(`{
		"comments":{}
	}`)
	if err == nil {
		t.Error("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(ParserErr{}) {
		t.Errorf("expected ParserErr but was %T", err)
	}
	expectedErr := "parsing json data '{\n\t\t\"comments\":{}\n\t}' failed for 'comments', error was 'value doesn't contain array; it contains object'"
	if err.Error() != expectedErr {
		t.Errorf("expected error text '%q' but was %q", expectedErr, err.Error())
	}
}

func TestGetLinksFromCommentsIsMissingPlainBodyInComments(t *testing.T) {
	_, err := GetLinksFromComments(`{
		"comments": [
			{
				"html_body": "<p>hello</p>",
				"plain_body": "hello"
			},
			{ 
				"html_body": "<p>test</p>"
			}
		]
	}`)
	if err == nil {
		t.Error("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(MissingJSONFieldError{}) {
		t.Errorf("expected MissingJSONFieldError but was %T", err)
	}
	expectedErr := "missing field 'plain_body' in 'comment 1 (base index 0)'"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected error text '%q' but was %q", expectedErr, err.Error())
	}
}
func TestGetLinksFromCommentsIsMissingHTMLBodyInComments(t *testing.T) {
	_, err := GetLinksFromComments(`{
		"comments": [
			{
				"html_body": "<p>hello</p>",
				"plain_body": "hello"
			},
			{ 
				"plain_body": "test"
			}
		]
	}`)
	if err == nil {
		t.Error("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(MissingJSONFieldError{}) {
		t.Errorf("expected MissingJSONFieldError but was %T", err)
	}
	expectedErr := "missing field 'html_body' in 'comment 1 (base index 0)'"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected error text '%q' but was %q", expectedErr, err.Error())
	}
}

func TestGetLinksFromCommentsHasNoLinks(t *testing.T) {
	links, err := GetLinksFromComments(`{
		"comments": [
		  {
			"attachments": [],
			"plain_body": "Thanks for your help!",
		   "html_body": "<p>Thanks for your help!</p>"
		 },
		 {
		   "attachments": [],
		  "plain_body": "here is some more for your help!",
		  "html_body": "<p>here is some more for your help!</p>"
		}
	   ]
	}`)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if len(links) != 0 {
		t.Errorf("expected 0 but had %v", len(links))
	}
}

func TestGetAttachmentsFromCommentHaveNoID(t *testing.T) {
	_, err := GetAttachmentsFromComments(`{
		"comments": [
		  {
			"attachments": [],
		   "plain_body": "Thanks for your help!",
		   "html_body": "<p>Thanks for your help!</p>"
		 },
		 {
		   "attachments": [],
		  "plain_body": "here is some more for your help!",
		  "html_body": "<p>here is some more for your help!</p>"
		}
	   ]
	}`)
	if err == nil {
		t.Error("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(MissingJSONFieldError{}) {
		t.Errorf("expected MissingJSONFieldError but was %T", err)
	}
	expectedErr := "missing field 'id' in 'in comment 0 (base index 0)'"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected error text to container '%q' but was %q", expectedErr, err.Error())
	}
}

func TestGetAttachmentsFromCommentHaveWrongIDType(t *testing.T) {
	_, err := GetAttachmentsFromComments(`{"comments": [{"id": "1"}]}`)
	if err == nil {
		t.Error("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(ParserErr{}) {
		t.Errorf("expected ParserErr but was %T", err)
	}
	expectedErr := "parsing json data '{\"comments\": [{\"id\": \"1\"}]}' failed for ''id' field at comments index 0', error was 'value doesn't contain number; it contains string'"
	if err.Error() != expectedErr {
		t.Errorf("expected error text '%q' but was %q", expectedErr, err.Error())
	}
}

func TestGetAttachmentsFromCommentHaveNoCreatedAt(t *testing.T) {
	_, err := GetAttachmentsFromComments(`{"comments": [{"id": 1}]}`)
	if err == nil {
		t.Error("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(MissingJSONFieldError{}) {
		t.Errorf("expected MissingJSONFieldError but was %T", err)
	}
	expectedErr := "parsing json data '{\"comments\": [{\"id\": 1}]}' missing field 'created_at' in 'in comment 0 (base index 0)'"
	if err.Error() != expectedErr {
		t.Errorf("expected error text '%q' but was %q", expectedErr, err.Error())
	}
}

func TestGetAttachmentsFromCommentHaveBlankCreatedAt(t *testing.T) {
	_, err := GetAttachmentsFromComments(`{"comments": [{"id": 1,"created_at":""}]}`)
	if err == nil {
		t.Error("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(ParserErr{}) {
		t.Errorf("expected ParserErr but was %T", err)
	}
	expectedErr := "parsing json data '{\"comments\": [{\"id\": 1,\"created_at\":\"\"}]}' failed for ''created_at' field a comments index 0', error was 'parsing time \"\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"\" as \"2006\"'"
	if err.Error() != expectedErr {
		t.Errorf("expected error text '%q' but was %q", expectedErr, err.Error())
	}
}

func TestGetAttachmentsFromCommentHaveInvalidCreatedAt(t *testing.T) {
	_, err := GetAttachmentsFromComments(`{"comments": [{"id": 1,"created_at":[]}]}`)
	if err == nil {
		t.Error("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(ParserErr{}) {
		t.Errorf("expected ParserErr but was %T", err)
	}
	expectedErr := "parsing json data '{\"comments\": [{\"id\": 1,\"created_at\":[]}]}' failed for ''created_at' field a comments index 0', error was 'value doesn't contain string; it contains array'"
	if err.Error() != expectedErr {
		t.Errorf("expected error text '%q' but was %q", expectedErr, err.Error())
	}
}

func TestGetAttachmentsFromCommentsHaveNoAttachments(t *testing.T) {
	attachements, err := GetAttachmentsFromComments(`{
		"comments": [
		  {
			"id": 1,
			"created_at": "2022-01-02T15:04:05Z",
			"attachments": [],
		   "plain_body": "Thanks for your help!",
		   "html_body": "<p>Thanks for your help!</p>"
		 },
		 {
			"id": 2,
			"created_at": "2022-01-02T15:04:05Z",
		   "attachments": [],
		  "plain_body": "here is some more for your help!",
		  "html_body": "<p>here is some more for your help!</p>"
		}
	   ]
	}`)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if len(attachements) != 0 {
		t.Errorf("expected 0 but had %v", len(attachements))
	}
}

func TestGetAttachmentsFromComments(t *testing.T) {
	attachments, err := GetAttachmentsFromComments(`{
		"comments": [
		  {
			"id": 1,
			"created_at": "2022-01-02T15:04:05Z",
			"attachments": [
				{
					"file_name": "abc",
					"deleted": false,
					"content_url": "http://test.com?file='test'",
					"content_type": "application/text",
					"size": 999
				},
				{
					"file_name": "xyz",
					"deleted": true,
					"content_url": "http://test.com?file='xyz'",
					"content_type": "application/image",
					"size": 100
				}
			]
		 },
		 {
			"id": 2,
			"created_at": "2022-02-02T15:04:05Z",
		    "attachments": [
			{
				"file_name": "pop",
				"deleted": false,
				"content_url": "http://test.com?file='pop'",
				"content_type": "application/json",
				"size": 50
			}
		   ]
		}
	   ]
	}`)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	expectedLen := 3
	if len(attachments) != expectedLen {
		t.Errorf("expected %v but had %v", expectedLen, len(attachments))
	}
	first := attachments[0]
	t1, err := time.Parse(time.RFC3339, "2022-01-02T15:04:05Z")
	if err != nil {
		t.Fatal(err)
	}
	expectedFirst := Attachment{
		FileName:          "abc",
		ParentCommentDate: t1,
		ParentCommentID:   1,
		ContentURL:        "http://test.com?file='test'",
		ContentType:       "application/text",
		Size:              999,
		Deleted:           false,
	}
	if !reflect.DeepEqual(first, expectedFirst) {
		t.Errorf("expected %v but was %v", expectedFirst, first)
	}
	second := attachments[1]
	if err != nil {
		t.Fatal(err)
	}
	expectedSecond := Attachment{
		FileName:          "xyz",
		ParentCommentDate: t1,
		ParentCommentID:   1,
		ContentURL:        "http://test.com?file='xyz'",
		ContentType:       "application/image",
		Size:              100,
		Deleted:           true,
	}
	if !reflect.DeepEqual(second, expectedSecond) {
		t.Errorf("expected %v but was %v", expectedFirst, second)
	}

	third := attachments[2]
	t2, err := time.Parse(time.RFC3339, "2022-02-02T15:04:05Z")
	if err != nil {
		t.Fatal(err)
	}
	expectedThird := Attachment{
		FileName:          "pop",
		ParentCommentDate: t2,
		ParentCommentID:   2,
		ContentURL:        "http://test.com?file='pop'",
		ContentType:       "application/json",
		Size:              50,
		Deleted:           false,
	}
	if !reflect.DeepEqual(third, expectedThird) {
		t.Errorf("expected %v but was %v", expectedThird, first)
	}
}

func TestGetAttachmentsFromCommentsMissingFileName(t *testing.T) {
	_, err := GetAttachmentsFromComments(`{
		"comments": [
		  {
			"id": 1,
			"created_at": "2022-01-02T15:04:05Z",
			"attachments": [
				{
					"deleted": false,
					"content_url": "http://test.com?file='test'",
					"content_type": "application/text",
					"size": 999
				}
			]
		 }
		]
	}`)
	if err == nil {
		t.Fatal("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(MissingJSONFieldError{}) {
		t.Errorf("expected MissingJSONFieldError but was %T", err)
	}
	expectedErr := "parsing json data '{\n\t\t\"comments\": [\n\t\t  {\n\t\t\t\"id\": 1,\n\t\t\t\"created_at\": \"2022-01-02T15:04:05Z\",\n\t\t\t\"attachments\": [\n\t\t\t\t{\n\t\t\t\t\t\"deleted\": false,\n\t\t\t\t\t\"content_url\": \"http://test.com?file='test'\",\n\t\t\t\t\t\"content_type\": \"application/text\",\n\t\t\t\t\t\"size\": 999\n\t\t\t\t}\n\t\t\t]\n\t\t }\n\t\t]\n\t}' missing field 'file_name' in 'in comment 0 in attachment 0 (base index 0)'"
	if err.Error() != expectedErr {
		t.Errorf("expected error text '%q' but was %q", expectedErr, err.Error())
	}
}

func TestGetAttachmentsFromCommentsInvalidFileName(t *testing.T) {
	_, err := GetAttachmentsFromComments(`{
		"comments": [
		  {
			"id": 1,
			"created_at": "2022-01-02T15:04:05Z",
			"attachments": [
				{
					"file_name": {},
					"deleted": false,
					"content_url": "http://test.com?file='test'",
					"content_type": "application/text",
					"size": 999
				}
			]
		 }
		]
	}`)
	if err == nil {
		t.Fatal("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(ParserErr{}) {
		t.Errorf("expected ParserErr but was %T", err)
	}
	expectedErr := "parsing json data '{\n\t\t\"comments\": [\n\t\t  {\n\t\t\t\"id\": 1,\n\t\t\t\"created_at\": \"2022-01-02T15:04:05Z\",\n\t\t\t\"attachments\": [\n\t\t\t\t{\n\t\t\t\t\t\"file_name\": {},\n\t\t\t\t\t\"deleted\": false,\n\t\t\t\t\t\"content_url\": \"http://test.com?file='test'\",\n\t\t\t\t\t\"content_type\": \"application/text\",\n\t\t\t\t\t\"size\": 999\n\t\t\t\t}\n\t\t\t]\n\t\t }\n\t\t]\n\t}' failed for 'file_name field in comment 0 in attachment 0 (base index 0)', error was 'value doesn't contain string; it contains object'"
	if err.Error() != expectedErr {
		t.Errorf("expected error text '%q' but was %q", expectedErr, err.Error())
	}
}

func TestGetAttachmentsFromCommentsMissingDeleted(t *testing.T) {
	_, err := GetAttachmentsFromComments(`{
		"comments": [
		  {
			"id": 1,
			"created_at": "2022-01-02T15:04:05Z",
			"attachments": [
				{
					"file_name": "false",
					"content_url": "http://test.com?file='test'",
					"content_type": "application/text",
					"size": 999
				}
			]
		 }
		]
	}`)
	if err == nil {
		t.Fatal("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(MissingJSONFieldError{}) {
		t.Errorf("expected MissingJSONFieldError but was %T", err)
	}
	expectedErr := "parsing json data '{\n\t\t\"comments\": [\n\t\t  {\n\t\t\t\"id\": 1,\n\t\t\t\"created_at\": \"2022-01-02T15:04:05Z\",\n\t\t\t\"attachments\": [\n\t\t\t\t{\n\t\t\t\t\t\"file_name\": \"false\",\n\t\t\t\t\t\"content_url\": \"http://test.com?file='test'\",\n\t\t\t\t\t\"content_type\": \"application/text\",\n\t\t\t\t\t\"size\": 999\n\t\t\t\t}\n\t\t\t]\n\t\t }\n\t\t]\n\t}' missing field 'deleted' in 'in comment 0 in attachment 0 (base index 0)'"
	if err.Error() != expectedErr {
		t.Errorf("expected error text '%q' but was %q", expectedErr, err.Error())
	}
}

func TestGetAttachmentsFromCommentsInvalidDelete(t *testing.T) {
	_, err := GetAttachmentsFromComments(`{
		"comments": [
		  {
			"id": 1,
			"created_at": "2022-01-02T15:04:05Z",
			"attachments": [
				{
					"file_name": "f",
					"deleted": [],
					"content_url": "http://test.com?file='test'",
					"content_type": "application/text",
					"size": 999
				}
			]
		 }
		]
	}`)
	if err == nil {
		t.Fatal("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(ParserErr{}) {
		t.Errorf("expected ParserErr but was %T", err)
	}
	expectedErr := "parsing json data '{\n\t\t\"comments\": [\n\t\t  {\n\t\t\t\"id\": 1,\n\t\t\t\"created_at\": \"2022-01-02T15:04:05Z\",\n\t\t\t\"attachments\": [\n\t\t\t\t{\n\t\t\t\t\t\"file_name\": \"f\",\n\t\t\t\t\t\"deleted\": [],\n\t\t\t\t\t\"content_url\": \"http://test.com?file='test'\",\n\t\t\t\t\t\"content_type\": \"application/text\",\n\t\t\t\t\t\"size\": 999\n\t\t\t\t}\n\t\t\t]\n\t\t }\n\t\t]\n\t}' failed for 'deleted field in comment 0 in attachment 0 (base index 0)', error was 'value doesn't contain bool; it contains array'"
	if err.Error() != expectedErr {
		t.Errorf("expected error text '%q' but was %q", expectedErr, err.Error())
	}
}
func TestGetAttachmentsFromCommentsMissingContentUrl(t *testing.T) {
	_, err := GetAttachmentsFromComments(`{
		"comments": [
		  {
			"id": 1,
			"created_at": "2022-01-02T15:04:05Z",
			"attachments": [
				{
					"file_name": "false",
					"deleted": false,
					"content_type": "application/text",
					"size": 999
				}
			]
		 }
		]
	}`)
	if err == nil {
		t.Fatal("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(MissingJSONFieldError{}) {
		t.Errorf("expected MissingJSONFieldError but was %T", err)
	}
	expectedErr := "parsing json data '{\n\t\t\"comments\": [\n\t\t  {\n\t\t\t\"id\": 1,\n\t\t\t\"created_at\": \"2022-01-02T15:04:05Z\",\n\t\t\t\"attachments\": [\n\t\t\t\t{\n\t\t\t\t\t\"file_name\": \"false\",\n\t\t\t\t\t\"deleted\": false,\n\t\t\t\t\t\"content_type\": \"application/text\",\n\t\t\t\t\t\"size\": 999\n\t\t\t\t}\n\t\t\t]\n\t\t }\n\t\t]\n\t}' missing field 'content_url' in 'in comment 0 in attachment 0 (base index 0)'"
	if err.Error() != expectedErr {
		t.Errorf("expected error text '%q' but was %q", expectedErr, err.Error())
	}
}

func TestGetAttachmentsFromCommentsInvalidContentUrl(t *testing.T) {
	_, err := GetAttachmentsFromComments(`{
		"comments": [
		  {
			"id": 1,
			"created_at": "2022-01-02T15:04:05Z",
			"attachments": [
				{
					"file_name": "f",
					"deleted": false,
					"content_url": {},
					"content_type": "application/text",
					"size": 999
				}
			]
		 }
		]
	}`)
	if err == nil {
		t.Fatal("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(ParserErr{}) {
		t.Errorf("expected ParserErr but was %T", err)
	}
	expectedErr := "parsing json data '{\n\t\t\"comments\": [\n\t\t  {\n\t\t\t\"id\": 1,\n\t\t\t\"created_at\": \"2022-01-02T15:04:05Z\",\n\t\t\t\"attachments\": [\n\t\t\t\t{\n\t\t\t\t\t\"file_name\": \"f\",\n\t\t\t\t\t\"deleted\": false,\n\t\t\t\t\t\"content_url\": {},\n\t\t\t\t\t\"content_type\": \"application/text\",\n\t\t\t\t\t\"size\": 999\n\t\t\t\t}\n\t\t\t]\n\t\t }\n\t\t]\n\t}' failed for 'content_url field in comment 0 in attachment 0 (base index 0)', error was 'value doesn't contain string; it contains object'"
	if err.Error() != expectedErr {
		t.Errorf("expected error text '%q' but was %q", expectedErr, err.Error())
	}
}

func TestGetAttachmentsFromCommentsMissingContentType(t *testing.T) {
	_, err := GetAttachmentsFromComments(`{
		"comments": [
		  {
			"id": 1,
			"created_at": "2022-01-02T15:04:05Z",
			"attachments": [
				{
					"file_name": "false",
					"content_url": "http://test.com?file='test'",
					"deleted": false,
					"size": 999
				}
			]
		 }
		]
	}`)
	if err == nil {
		t.Fatal("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(MissingJSONFieldError{}) {
		t.Errorf("expected MissingJSONFieldError but was %T", err)
	}
	expectedErr := "parsing json data '{\n\t\t\"comments\": [\n\t\t  {\n\t\t\t\"id\": 1,\n\t\t\t\"created_at\": \"2022-01-02T15:04:05Z\",\n\t\t\t\"attachments\": [\n\t\t\t\t{\n\t\t\t\t\t\"file_name\": \"false\",\n\t\t\t\t\t\"content_url\": \"http://test.com?file='test'\",\n\t\t\t\t\t\"deleted\": false,\n\t\t\t\t\t\"size\": 999\n\t\t\t\t}\n\t\t\t]\n\t\t }\n\t\t]\n\t}' missing field 'content_type' in 'in comment 0 in attachment 0 (base index 0)'"
	if err.Error() != expectedErr {
		t.Errorf("expected error text '%q' but was %q", expectedErr, err.Error())
	}
}

func TestGetAttachmentsFromCommentsInvalidContentType(t *testing.T) {
	_, err := GetAttachmentsFromComments(`{
		"comments": [
		  {
			"id": 1,
			"created_at": "2022-01-02T15:04:05Z",
			"attachments": [
				{
					"file_name": "f",
					"deleted": false,
					"content_url": "http://test.com?file='test'",
					"content_type": {},
					"size": 999
				}
			]
		 }
		]
	}`)
	if err == nil {
		t.Fatal("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(ParserErr{}) {
		t.Errorf("expected ParserErr but was %T", err)
	}
	expectedErr := "parsing json data '{\n\t\t\"comments\": [\n\t\t  {\n\t\t\t\"id\": 1,\n\t\t\t\"created_at\": \"2022-01-02T15:04:05Z\",\n\t\t\t\"attachments\": [\n\t\t\t\t{\n\t\t\t\t\t\"file_name\": \"f\",\n\t\t\t\t\t\"deleted\": false,\n\t\t\t\t\t\"content_url\": \"http://test.com?file='test'\",\n\t\t\t\t\t\"content_type\": {},\n\t\t\t\t\t\"size\": 999\n\t\t\t\t}\n\t\t\t]\n\t\t }\n\t\t]\n\t}' failed for 'content_type field in comment 0 in attachment 0 (base index 0)', error was 'value doesn't contain string; it contains object'"
	if err.Error() != expectedErr {
		t.Errorf("expected error text '%q' but was %q", expectedErr, err.Error())
	}
}
func TestGetAttachmentsFromCommentsMissingSize(t *testing.T) {
	_, err := GetAttachmentsFromComments(`{
		"comments": [
		  {
			"id": 1,
			"created_at": "2022-01-02T15:04:05Z",
			"attachments": [
				{
					"file_name": "false",
					"content_url": "http://test.com?file='test'",
					"content_type": "application/text",
					"deleted": false
				}
			]
		 }
		]
	}`)
	if err == nil {
		t.Fatal("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(MissingJSONFieldError{}) {
		t.Errorf("expected MissingJSONFieldError but was %T", err)
	}
	expectedErr := "parsing json data '{\n\t\t\"comments\": [\n\t\t  {\n\t\t\t\"id\": 1,\n\t\t\t\"created_at\": \"2022-01-02T15:04:05Z\",\n\t\t\t\"attachments\": [\n\t\t\t\t{\n\t\t\t\t\t\"file_name\": \"false\",\n\t\t\t\t\t\"content_url\": \"http://test.com?file='test'\",\n\t\t\t\t\t\"content_type\": \"application/text\",\n\t\t\t\t\t\"deleted\": false\n\t\t\t\t}\n\t\t\t]\n\t\t }\n\t\t]\n\t}' missing field 'size' in 'in comment 0 in attachment 0 (base index 0)'"
	if err.Error() != expectedErr {
		t.Errorf("expected error text '%q' but was %q", expectedErr, err.Error())
	}
}

func TestGetAttachmentsFromCommentsInvalidSize(t *testing.T) {
	_, err := GetAttachmentsFromComments(`{
		"comments": [
		  {
			"id": 1,
			"created_at": "2022-01-02T15:04:05Z",
			"attachments": [
				{
					"file_name": "f",
					"deleted": false,
					"content_url": "http://test.com?file='test'",
					"content_type": "application/text",
					"size": ""
				}
			]
		 }
		]
	}`)
	if err == nil {
		t.Fatal("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(ParserErr{}) {
		t.Errorf("expected ParserErr but was %T", err)
	}
	expectedErr := "parsing json data '{\n\t\t\"comments\": [\n\t\t  {\n\t\t\t\"id\": 1,\n\t\t\t\"created_at\": \"2022-01-02T15:04:05Z\",\n\t\t\t\"attachments\": [\n\t\t\t\t{\n\t\t\t\t\t\"file_name\": \"f\",\n\t\t\t\t\t\"deleted\": false,\n\t\t\t\t\t\"content_url\": \"http://test.com?file='test'\",\n\t\t\t\t\t\"content_type\": \"application/text\",\n\t\t\t\t\t\"size\": \"\"\n\t\t\t\t}\n\t\t\t]\n\t\t }\n\t\t]\n\t}' failed for 'size field in comment 0 in attachment 0 (base index 0)', error was 'value doesn't contain number; it contains string'"
	if err.Error() != expectedErr {
		t.Errorf("expected error text '%q' but was %q", expectedErr, err.Error())
	}
}
func TestGetAttachmentsFromCommentsAreMissingAttachmentsField(t *testing.T) {
	_, err := GetAttachmentsFromComments(`{
		"comments": [
		  {
			"id": 1,
			"created_at": "2022-01-02T15:04:05Z"
		 },
		 {
			"id": 2,
			"created_at": "2022-01-02T15:04:05Z"
		}
	   ]
	}`)
	if err == nil {
		t.Error("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(MissingJSONFieldError{}) {
		t.Errorf("expected MissingJSONFieldError but was %T", err)
	}
	expectedErr := "parsing json data '{\n\t\t\"comments\": [\n\t\t  {\n\t\t\t\"id\": 1,\n\t\t\t\"created_at\": \"2022-01-02T15:04:05Z\"\n\t\t },\n\t\t {\n\t\t\t\"id\": 2,\n\t\t\t\"created_at\": \"2022-01-02T15:04:05Z\"\n\t\t}\n\t   ]\n\t}' missing field 'attachments' in 'in comment 0 (base index 0)'"
	if err.Error() != expectedErr {
		t.Errorf("expected error text '%q' but was %q", expectedErr, err.Error())
	}
}

func TestGetAttachmentsFromCommentsHaveInvalidAttachmentsField(t *testing.T) {
	_, err := GetAttachmentsFromComments(`{
		"comments": [
		  {
			"id": 1,
			"attachments": 1,
			"created_at": "2022-01-02T15:04:05Z"
		 },
		 {
			"id": 2,
			"attachments": 2,
			"created_at": "2022-01-02T15:04:05Z"
		}
	   ]
	}`)
	if err == nil {
		t.Error("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(ParserErr{}) {
		t.Errorf("expected ParserErr but was %T", err)
	}
	expectedErr := "parsing json data '{\n\t\t\"comments\": [\n\t\t  {\n\t\t\t\"id\": 1,\n\t\t\t\"attachments\": 1,\n\t\t\t\"created_at\": \"2022-01-02T15:04:05Z\"\n\t\t },\n\t\t {\n\t\t\t\"id\": 2,\n\t\t\t\"attachments\": 2,\n\t\t\t\"created_at\": \"2022-01-02T15:04:05Z\"\n\t\t}\n\t   ]\n\t}' failed for 'attachments field in comment 0 (base index 0)', error was 'value doesn't contain array; it contains number'"
	if err.Error() != expectedErr {
		t.Errorf("expected error text '%q' but was %q", expectedErr, err.Error())
	}
}
func TestGetAttachmentsFromCommentsIsMissingComments(t *testing.T) {
	_, err := GetAttachmentsFromComments(`{}`)
	if err == nil {
		t.Error("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(MissingJSONFieldError{}) {
		t.Errorf("expected MissingJSONFieldError but was %T", err)
	}
	expectedErr := "parsing json data '{}' failed, missing 'comments' field"
	if err.Error() != expectedErr {
		t.Errorf("expected error text '%q' but was %q", expectedErr, err.Error())
	}
}

func TestGetLinksFromCommentsHasInvalidJSON(t *testing.T) {
	_, err := GetAttachmentsFromComments(``)
	if err == nil {
		t.Error("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(ParserErr{}) {
		t.Errorf("expected ParserErr but was %T", err)
	}
	expectedErr := "parsing json data '' failed, error was 'cannot parse JSON: cannot parse empty string; unparsed tail: \"\"'"
	if err.Error() != expectedErr {
		t.Errorf("expected error text '%q' but was %q", expectedErr, err.Error())
	}
}

func TestGetAttachmentsFromCommentsHasInvalidCommentsField(t *testing.T) {
	_, err := GetAttachmentsFromComments(`{
		"comments":{}
	}`)
	if err == nil {
		t.Error("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(ParserErr{}) {
		t.Errorf("expected ParserErr but was %T", err)
	}
	expectedErr := "parsing json data '{\n\t\t\"comments\":{}\n\t}' failed for 'comments', error was 'value doesn't contain array; it contains object'"
	if err.Error() != expectedErr {
		t.Errorf("expected error text '%q' but was %q", expectedErr, err.Error())
	}
}
