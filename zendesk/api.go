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
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

func URL(subDomain, ticketID string) string {
	return fmt.Sprintf("https://%v.zendesk.com/api/v2/tickets/%v/comments.json", subDomain, ticketID)
}

// GetTicketComments returns all comments as a concatenated string so we can search them
// GET /api/v2/tickets/{ticket_id}/comments
// curl https://{subdomain}.zendesk.com/api/v2/tickets/{ticket_id}/comments.json \
//
// {
// 	"comments": [
// 	  {
// 		"attachments": [
// 		  {
// 			"content_type": "text/plain",
// 			"content_url": "https://company.zendesk.com/attachments/crash.log",
// 			"file_name": "crash.log",
// 			"id": 498483,
// 			"size": 2532,
// 			"thumbnails": []
// 		  }
// 		],
// 		"author_id": 123123,
// 		"body": "Thanks for your help!",
// 		"created_at": "2009-07-20T22:55:29Z",
// 		"id": 1274,
// 		"metadata": {
// 		  "system": {
// 			"client": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
// 			"ip_address": "1.1.1.1",
// 			"latitude": -37.000000000001,
// 			"location": "Melbourne, 07, Australia",
// 			"longitude": 144.0000000000002
// 		  },
// 		  "via": {
// 			"channel": "web",
// 			"source": {
// 			  "from": {},
// 			  "rel": "web_widget",
// 			  "to": {}
// 			}
// 		  }
// 		},
// 		"public": true,
// 		"type": "Comment"
// 	  }
// 	]
//   }
func (z *Client) GetTicketComentsJSON(ticketID string, page_url *string) (string, error) {
	url := URL(z.subDomain, ticketID)
	if page_url != nil && *page_url != "" {
		log.Printf("page_url used: %v", *page_url)
		url = *page_url
	}
	log.Printf("URL used: %v", url)
	auth := fmt.Sprintf("%v/token:%v", z.username, z.password)
	base64Auth := base64.StdEncoding.EncodeToString([]byte(auth))
	r, err := z.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Basic %v", base64Auth)).
		Get(url)
	if err != nil {
		return "", fmt.Errorf("unable to read ticket comments with error '%v'", err)
	}
	rawBody := r.Body()
	if z.verbose {
		var prettyJSONBuffer bytes.Buffer
		if err := json.Indent(&prettyJSONBuffer, rawBody, "=", "\t"); err != nil {
			log.Printf("WARN: Unable to log debugging json for ticket %v printing string '%v'", ticketID, string(rawBody))
		} else {
			log.Printf("DEBUG: Ticket %v Comments Contents '%v'", ticketID, prettyJSONBuffer.String())
		}
	}
	return string(rawBody), nil
}

type Client struct {
	client    resty.Client
	username  string
	password  string
	subDomain string
	verbose   bool
}

func NewClient(username, password, subDomain string, verbose bool) *Client {
	return &Client{
		subDomain: subDomain,
		username:  username,
		password:  password,
		client:    *resty.New(),
		verbose:   verbose,
	}
}
