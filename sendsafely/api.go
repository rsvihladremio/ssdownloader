package sendsafely

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const SS_URL = "https://app.sendsafely.com/api/v2.0"

type SendSafelyClient struct {
	client      *resty.Client
	ssApiKey    string
	ssApiSecret string
}

func NewSendSafelyClient(ssApiKey, ssApiSecret string) *SendSafelyClient {
	client := resty.New()

	return &SendSafelyClient{
		ssApiKey:    ssApiKey,
		ssApiSecret: ssApiSecret,
		client:      client,
	}
}

func (s *SendSafelyClient) RetrievePackgeById(packageId string) (SendSafelyPackage, error) {
	now := time.Now()
	//2019-01-14T22:24:00+0000 as documented in https://sendsafely.zendesk.com/hc/en-us/articles/360027599232-SendSafely-REST-API
	ts := now.Format("2006-02-03T15:04:05-0700")
	// adding package and packageId to the base send safely URL. This is a quirk documented under URL_PATH in the sendsafely docs above
	urlPath := strings.Join([]string{"/api", "v2.0", "package", packageId}, "/")
	sig, err := s.generateRequestSignature(ts, urlPath, "")
	if err != nil {
		return SendSafelyPackage{}, fmt.Errorf("unexpected error generating request signature '%v'", err)
	}
	// validating client is set in the first place
	if s.client == nil {
		return SendSafelyPackage{}, errors.New("client was never initialized. Please use NewSendSafelyClient to initialize SendSafelyClient")
	}

	//this is actually usable by the rest api unlike the urlPath
	requestPath := strings.Join([]string{SS_URL, "package", packageId}, "/")
	// add the required sendsafely headers to the request is accepted and then submit the request
	r, err := s.client.R().
		SetHeader("ss-api-key", s.ssApiKey).
		SetHeader("ss-request-timestamp", ts).
		SetHeader("ss-request-signature", sig).
		Get(requestPath)
	if err != nil {
		return SendSafelyPackage{}, fmt.Errorf("unexpeced error '%v' while retrieving request '%v'", err, requestPath)
	}
	return ParsePackage(string(r.Body()))
}

// GenerateRequestSignature is a utility method to generate the ss-request-signature header
// which is a combination of HmacSHA256(API_SECRET, API_KEY + URL_PATH + TIMESTAMP + REQUEST_BODY)
// TIMESTAMP meaning ss-request-timestamp header. The overal function is documented at the
// following link https://sendsafely.zendesk.com/hc/en-us/articles/360027599232-SendSafely-REST-API
func (s *SendSafelyClient) generateRequestSignature(ts string, urlPath string, requestBody string) (string, error) {
	// using the api secret to setup the hmacsha256
	h := hmac.New(sha256.New, []byte(s.ssApiSecret))

	// dump data into the hash, a combination of api_key + urlPath + timestamp + request-body
	requestData := strings.Join([]string{s.ssApiKey, urlPath, ts, requestBody}, "")
	_, err := h.Write([]byte(requestData))
	if err != nil {
		return "", fmt.Errorf("unexpected error encoding data '%v' the following combined value was sent %v", err, requestData)
	}

	// Get result and encode as hexadecimal string
	return hex.EncodeToString(h.Sum(nil)), nil
}
