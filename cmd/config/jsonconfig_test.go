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
package config

import "testing"

func TestLoadConfig(t *testing.T) {
	var c Config
	err := Load("testdata/creds.json", &c)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	expectedDir := "mydir"
	if c.DownloadDir != expectedDir {
		t.Errorf("expected %v but was %v", expectedDir, c.DownloadDir)
	}
	expectedZDToken := "zdtoken"
	if c.ZendeskToken != expectedZDToken {
		t.Errorf("expected %v but was %v", expectedZDToken, c.ZendeskToken)
	}

	expectedZDEmail := "test@example.com"
	if c.ZendeskEmail != expectedZDEmail {
		t.Errorf("expected %v but was %v", expectedZDEmail, c.ZendeskEmail)
	}

	expectedZDDomain := "tester"
	if c.ZendeskDomain != expectedZDDomain {
		t.Errorf("expected %v but was %v", expectedZDDomain, c.ZendeskDomain)
	}

	expectedSSKey := "ssapikey"
	if c.SsApiKey != expectedSSKey {
		t.Errorf("expected %v but was %v", expectedSSKey, c.SsApiKey)
	}

	expectedSSSecret := "ssapisecret"
	if c.SsApiSecret != expectedSSSecret {
		t.Errorf("expected %v but was %v", expectedSSSecret, c.SsApiSecret)
	}
}
