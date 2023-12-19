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

// sendsafely package decrypts files, combines file parts into whole files, and handles api access to the sendsafely rest api
package sendsafely

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/crypto/pbkdf2"
)

const URL = "https://app.sendsafely.com/api/v2.0"

type Client interface {
	RetrievePackageByID(packageID string) (Package, error)
	GetDownloadUrlsForFile(p Package, fileID, keyCode string, start, end int) ([]DownloadURL, error)
}

// Client uses the SendSafely REST Api to
// enable automation of SendSafely in Go
type DownloadClient struct {
	parser      *APIParser
	client      *resty.Client
	ssAPIKey    string
	ssAPISecret string
	verbose     bool
}

// NewClient is the preferred way to initialize SendSafelyClient
func NewClient(ssAPIKey, ssAPISecret string, verbose bool) Client {
	client := resty.New()

	return &DownloadClient{
		ssAPIKey:    ssAPIKey,
		ssAPISecret: ssAPISecret,
		client:      client,
		parser:      &APIParser{},
		verbose:     verbose,
	}
}

func (s *DownloadClient) RetrievePackageByID(packageID string) (Package, error) {
	now := time.Now()
	//2019-01-14T22:24:00+0000 as documented in https://sendsafely.zendesk.com/hc/en-us/articles/360027599232-SendSafely-REST-API
	ts := now.Format("2006-01-02T15:04:05-0700")
	// adding package and packageId to the base send safely URL. This is a quirk documented under URL_PATH in the sendsafely docs above
	urlPath := strings.Join([]string{"/api", "v2.0", "package", packageID}, "/")
	sig, err := s.generateRequestSignature(ts, urlPath, "")
	if err != nil {
		return Package{}, fmt.Errorf("unexpected error generating request signature '%v'", err)
	}
	// validating client is set in the first place
	if s.client == nil {
		return Package{}, errors.New("client was never initialized. Please use NewSendSafelyClient to initialize SendSafelyClient")
	}

	//this is actually usable by the rest api unlike the urlPath
	requestPath := strings.Join([]string{URL, "package", packageID}, "/")
	// add the required sendsafely headers to the request is accepted and then submit the request

	slog.Debug("retrieving package by id", "api_key", s.ssAPIKey, "request_ts_header", ts, "request_sig_header", sig, "url_path", urlPath, "request_path", requestPath)
	r, err := s.client.R().
		SetHeader("ss-api-key", s.ssAPIKey).
		SetHeader("ss-request-timestamp", ts).
		SetHeader("ss-request-signature", sig).
		Get(requestPath)
	if err != nil {
		return Package{}, fmt.Errorf("unexpected error '%v' while retrieving request '%v' error code was '%v'", err, requestPath, r.StatusCode())
	}
	rawResponseBody := r.Body()
	if s.verbose {
		var prettyJSONBuffer bytes.Buffer
		if err := json.Indent(&prettyJSONBuffer, rawResponseBody, "=", "\t"); err != nil {
			slog.Warn("unable to log debugging json for sendsafely", "package_id", packageID, "body", string(rawResponseBody), "error_msg", err)
		} else {
			slog.Debug("package response", "package_id", packageID, "http_response_body", prettyJSONBuffer.String())
		}
	}
	return s.parser.ParsePackage(packageID, string(rawResponseBody))
}

// GenerateRequestSignature is a utility method to generate the ss-request-signature header
// which is a combination of HmacSHA256(API_SECRET, API_KEY + URL_PATH + TIMESTAMP + REQUEST_BODY)
// TIMESTAMP meaning ss-request-timestamp header. The overall function is documented at the
// following link https://sendsafely.zendesk.com/hc/en-us/articles/360027599232-SendSafely-REST-API
func (s *DownloadClient) generateRequestSignature(ts string, urlPath string, requestBody string) (string, error) {

	// dump data into the hash, a combination of api_key + urlPath + timestamp + request-body
	requestData := strings.Join([]string{s.ssAPIKey, urlPath, ts, requestBody}, "")
	return s.sign(requestData)
}

// GenerateRequestSignature is a utility method to generate the checksum for download requests
// which is a combination of HmacSHA256(keycode,packageCode))
// following link https://bump.sh/doc/sendsafely-rest-api#operation-post-package-parameter-file-parameter-download-urls
func (s *DownloadClient) generateChecksum(keyCode, packageCode string) string {

	// use pbkdf2 to encrypt the keycode
	// from sendsafely docs https://sendsafely.zendesk.com/hc/en-us/articles/360027599232-SendSafely-REST-API
	// This endpoint requires you do provide a "checksum" parameter that is calculated
	// using the keycode (Client Secret) and the Package Code, which are both included in the Package Link. The checksum is generated using PBKDF2-HMAC-SHA256 with the keycode as the password, and the Package Code as the salt.
	//
	// Use the following inputs for your PBKDF2 function:
	//
	// Hashing Algorithm - SHA-256
	// Password  - Use the keycode for this value
	// Salt - Use the Package Code for this value
	// Iteration Count - 1024
	// Key Length - 32 bytes
	// later I read the code https://github.com/SendSafely/Java-Client-API/blob/ad47e899ed3bea13168b24af25d3921571285e94/SendSafelyAPI/src/com/sendsafely/utils/CryptoUtil.java#L81-L84
	// to discover this is all I needed and the other signing method used elsewhere is not needed here
	iterations := 1024
	keyLength := 32
	dk := pbkdf2.Key([]byte(keyCode), []byte(packageCode), iterations, keyLength, sha256.New)
	return hex.EncodeToString(dk)
}

func (s *DownloadClient) sign(data string) (string, error) {
	// using the api secret to setup the hmacsha256
	h := hmac.New(sha256.New, []byte(s.ssAPISecret))

	_, err := h.Write([]byte(data))
	if err != nil {
		return "", fmt.Errorf("unexpected error encoding data '%v' the following combined value was sent %v", err, data)
	}

	// Get result and encode as hexadecimal string
	return hex.EncodeToString(h.Sum(nil)), nil
}

// GetDownloadUrlsForFile retrieves up to 25 download urls per request
// read https://bump.sh/doc/sendsafely-rest-api#operation-post-package-parameter-file-parameter-download-urls
// for details
// request example
//
// # Headers
// ss-api-key: string
// ss-request-signature: string
// ss-request-timestamp: string

// # Payload
//
//	{
//	  "checksum": "20ea6189a58d6618164d6b53cbc7b5dbeef2fa9a61186b13a80cf1f620987bc1",
//	  "startSegment": 1,
//	  "endSegment": 25
//	}
func (s *DownloadClient) GetDownloadUrlsForFile(p Package, fileID, keyCode string, start, end int) ([]DownloadURL, error) {
	// validating client is set in the first place
	if s.client == nil {
		return []DownloadURL{}, errors.New("client was never initialized. Please use NewSendSafelyClient to initialize SendSafelyClient")
	}
	now := time.Now()
	//2019-01-14T22:24:00+0000 as documented in https://sendsafely.zendesk.com/hc/en-us/articles/360027599232-SendSafely-REST-API
	ts := now.Format("2006-01-02T15:04:05-0700")
	// adding package and packageId to the base send safely URL. This is a quirk documented under URL_PATH in the sendsafely docs above
	urlPath := strings.Join([]string{"/api", "v2.0", "package", p.PackageID, "file", fileID, "download-urls/"}, "/")
	//generate the check sum
	checkSum := s.generateChecksum(keyCode, p.PackageCode)
	body := fmt.Sprintf("{\"checksum\":\"%v\",\"startSegment\":%v,\"endSegment\":%v}", checkSum, start, end)

	sig, err := s.generateRequestSignature(ts, urlPath, body)
	if err != nil {
		return []DownloadURL{}, fmt.Errorf("unexpected error generating request signature '%v'", err)
	}
	//this is actually usable by the rest api unlike the urlPath
	requestPath := strings.Join([]string{URL, "package", p.PackageID, "file", fileID, "download-urls/"}, "/")
	// add the required sendsafely headers to the request is accepted and then submit the request

	r, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("ss-api-key", s.ssAPIKey).
		SetHeader("ss-request-timestamp", ts).
		SetHeader("ss-request-signature", sig).
		SetBody(body).
		Post(requestPath)
	if err != nil {
		return []DownloadURL{}, fmt.Errorf("unexpected error '%v' while retrieving request '%v'", err, requestPath)
	}
	rawResponseBody := r.Body()
	if s.verbose {
		var prettyJSONBuffer bytes.Buffer
		if err := json.Indent(&prettyJSONBuffer, rawResponseBody, "=", "\t"); err != nil {
			slog.Warn("unable to log debugging json for sendsafely urls", "package_id", p.PackageID, "http_response_body", string(rawResponseBody), "error_msg", err)
		} else {
			slog.Debug("package download urls", "package_id", p.PackageID, "http_response_body", prettyJSONBuffer.String())
		}
	}
	return s.parser.ParseDownloadUrls(string(rawResponseBody))
}
