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
	"testing"

	"github.com/jarcoal/httpmock"
)

// This is the default happy path test, no errors
func TestRetrievePackgeById(t *testing.T) {
	// since we are using a mock http api we can use any api secret we feel like
	zdClient := NewClient("myApiKey", "mySecret", "zdsub", false)

	// pass in the resty httpy client that the SendSafelyClient uses so that
	// httpmock can replace it's transport parameter with a mock one
	// preventing remote calls from going to SendSafely
	httpmock.ActivateNonDefault(zdClient.client.GetClient())
	// make sure to reset the mock after the test
	defer httpmock.DeactivateAndReset()
	ticketID := "12314"
	// using the sample json from the sendafely website, with the files field correctly
	// filled out (the sendsafely site had the "files" field incorrectly documented)
	resp := `[{"id":"oye"}]`

	// setup a responder to the expected status code of 200 and then returning the json data setup abovee
	responder := httpmock.NewStringResponder(200, resp)

	url := URL(zdClient.subDomain, ticketID)
	// we are expecting a GET request with the exact url specified above, if that exact match happens
	// the json body setup in the responder will return instead of hitting the remote sendsafely server
	httpmock.RegisterResponder("GET", url, responder)
	comments, err := zdClient.GetTicketComentsJSON(ticketID)
	if err != nil {
		t.Fatalf("unexpected error retrieving id '%v'", err)
	}
	if comments != resp {
		t.Errorf("expected %v but received %v", resp, comments)
	}
}
