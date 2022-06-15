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
package link

import (
	"fmt"
	"net/url"
	"strings"
)

type LinkParts struct {
	Thread      string
	PackageCode string
	KeyCode     string
}

// ParseLink splits up a SendSafely package download URL into it's important parts
// This allows us to download the package
func ParseLink(inputUrl string) (LinkParts, error) {
	// attempt to parse this as a valid url, will fail if it is malformed
	u, err := url.Parse(inputUrl)
	if err != nil {
		return LinkParts{}, fmt.Errorf("unable to parse url '%v'", inputUrl)
	}
	// search the query parameters for packageCode and thread
	query := u.Query()
	if !query.Has("packageCode") {
		return LinkParts{}, fmt.Errorf("expected to have packageCode in url '%v' but it is not present", inputUrl)
	}
	if !query.Has("thread") {
		return LinkParts{}, fmt.Errorf("expected to have thread in url '%v' but it is not present", inputUrl)
	}
	// for whatever reason keyCode is stored as a fragment, this is a bit tricker but we know what it starts with
	// however, this is the most fragile part and if the URL scheme varies a bit this will break badly
	keyCodeRaw := u.Fragment
	if !strings.HasPrefix(keyCodeRaw, "keyCode=") {
		return LinkParts{}, fmt.Errorf("expected to have fragment keyCode= in url '%v' but it is not present, the fragment detected is '%v'", inputUrl, keyCodeRaw)
	}

	return LinkParts{
		Thread:      query.Get("thread"),
		PackageCode: query.Get("packageCode"),
		KeyCode:     keyCodeRaw[8:], //throwing away keyCode= and only keeping the rest of the string
	}, nil
}
