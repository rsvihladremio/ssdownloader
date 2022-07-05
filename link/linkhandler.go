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

//link package handles parsing of sendsafely links so that we can retrieve the identifying information in the query parameters
package link

import (
	"fmt"
	"net/url"
	"strings"
)

type Parts struct {
	Thread      string
	PackageCode string
	KeyCode     string
}

type KeyCodeIsMissingErr struct {
	InputURL string
	KeyCode  string
}

func (k KeyCodeIsMissingErr) Error() string {
	return fmt.Sprintf("expected to have fragment keyCode= in url '%v' but it is not present, the fragment detected is '%v'", k.InputURL, k.KeyCode)
}

type PackageCodeIsMissingErr struct {
	InputURL string
}

func (p PackageCodeIsMissingErr) Error() string {
	return fmt.Sprintf("expected to have packageCode in url '%v' but it is not present", p.InputURL)
}

type ThreadIsMissingErr struct {
	InputURL string
}

func (p ThreadIsMissingErr) Error() string {
	return fmt.Sprintf("expected to have thread in url '%v' but it is not present", p.InputURL)
}

type URLParseErr struct {
	URL     string
	BaseErr error
}

func (u URLParseErr) Error() string {
	return fmt.Sprintf("unable to parse url '%v' due to error '%v'", u.URL, u.BaseErr)
}

// ParseLink splits up a SendSafely package download URL into it's important parts
// This allows us to download the package
func ParseLink(inputURL string) (Parts, error) {
	// attempt to parse this as a valid url, will fail if it is malformed
	u, err := url.Parse(inputURL)
	if err != nil {
		return Parts{}, URLParseErr{
			BaseErr: err,
			URL:     inputURL,
		}
	}
	// search the query parameters for packageCode and thread
	query := u.Query()
	if !query.Has("packageCode") {
		return Parts{}, PackageCodeIsMissingErr{InputURL: inputURL}
	}
	if !query.Has("thread") {
		return Parts{}, ThreadIsMissingErr{InputURL: inputURL}
	}
	// for whatever reason keyCode is stored as a fragment, this is a bit tricker but we know what it starts with
	// however, this is the most fragile part and if the URL scheme varies a bit this will break badly
	keyCodeRaw := u.Fragment
	if !strings.HasPrefix(keyCodeRaw, "keyCode=") {
		return Parts{}, KeyCodeIsMissingErr{
			InputURL: inputURL,
			KeyCode:  keyCodeRaw,
		}
	}

	return Parts{
		Thread:      query.Get("thread"),
		PackageCode: query.Get("packageCode"),
		KeyCode:     keyCodeRaw[8:], //throwing away keyCode= and only keeping the rest of the string
	}, nil
}
