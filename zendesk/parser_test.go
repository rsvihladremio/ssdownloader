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
	"testing"
)

func TestGetLinksFromComments(t *testing.T) {
	links, err := GetLinksFromComments(` {
		 	"comments": [
		 	  {
		 		"attachments": [],
				"html_body": "<p>Thanks for your help!</p><a href='http://example.com'>example</a>"
			  },
			  {
				"attachments": [],
			   "html_body": "<p>here is some more for your help!</p><a href='http://example.com/file2'>example 2</a>"
			 }
			]
		 }`)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	expectedLinks := []string{"http://example.com", "http://example.com/file2"}
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

func TestGetLinksFromCommentsIsMissingHTMLBodyInComments(t *testing.T) {
	_, err := GetLinksFromComments(`{
		"comments": [
			{
				"html_body": "<p>hello</p>"
			},
			{}
		]
	}`)
	if err == nil {
		t.Error("expected error but was nil")
	}
	if reflect.TypeOf(err) != reflect.TypeOf(MissingJSONFieldError{}) {
		t.Errorf("expected ParserErr but was %T", err)
	}
	expectedErr := "parsing json data '{\n\t\t\"comments\": [\n\t\t\t{\n\t\t\t\t\"html_body\": \"<p>hello</p>\"\n\t\t\t},\n\t\t\t{}\n\t\t]\n\t}' missing field 'html_body' in 'comment 1 (base index 0)'"
	if err.Error() != expectedErr {
		t.Errorf("expected error text '%q' but was %q", expectedErr, err.Error())
	}
}

func TestGetLinksFromCommentsHasNoLinks(t *testing.T) {
	links, err := GetLinksFromComments(`{
		"comments": [
		  {
			"attachments": [],
		   "html_body": "<p>Thanks for your help!</p>"
		 },
		 {
		   "attachments": [],
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
