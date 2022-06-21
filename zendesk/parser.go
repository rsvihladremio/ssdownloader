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
	"strings"

	"github.com/valyala/fastjson"
	"golang.org/x/net/html"
)

// GetSendSafelyLinksFromComments is parsing out the sendsafely links from the html_
// docs are here https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_comments/#list-comments
func GetSendSafelyLinksFromComments(jsonData string) ([]string, error) {
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
					if a.Key == "href" && strings.HasPrefix(a.Val, "https://sendsafely") {
						linksFound = append(linksFound, a.Val)
						break
					}
				}

			}
		}
	}
	return linksFound, nil
}
