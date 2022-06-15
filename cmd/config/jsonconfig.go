/*
Copyright © 2022 Ryan SVIHLA

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	SsApiKey      string
	SsApiSecret   string
	ZendeskDomain string
	ZendeskEmail  string
	ZendeskToken  string
}

func ReadConfigFile(cfgFile string) (string, error) {
	if cfgFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("unable to read home dir due to error'%v'", err)
		}
		return strings.Join([]string{home, ".config", "ssdownloader", "creds.json"}, string(os.PathSeparator)), nil
	}
	return cfgFile, nil
}

func Load(cfgFile string, c *Config) error {
	fileToLoad, err := ReadConfigFile(cfgFile)
	if err != nil {
		return fmt.Errorf("trying to read configuration resulted in error '%v'", err)
	}
	// best security practice
	cleanedConfigFile := filepath.Clean(fileToLoad)
	b, err := os.ReadFile(cleanedConfigFile)
	if err != nil {
		return fmt.Errorf("unable to read file '%v' due to error '%v'", cleanedConfigFile, err)
	}
	err = json.Unmarshal(b, &c)
	if err != nil {
		return fmt.Errorf("unable to process the file '%v' this may indicate the file format is incorrect, the error was '%v'", cleanedConfigFile, err)
	}
	return nil
}

func Save(c Config, cfgFile string) error {
	b, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("unable to convert configuration to json file due to error '%v'", err)
	}
	fileToSave, err := ReadConfigFile(cfgFile)
	if err != nil {
		return fmt.Errorf("trying to get the path to the configuration resulted in error '%v'", err)
	}
	// best security practice
	cleanedConfigFile := filepath.Clean(fileToSave)
	err = os.WriteFile(cleanedConfigFile, b, os.FileMode(0600))
	if err != nil {
		return fmt.Errorf("unable to write configuration file to location '%v' due to error '%v'", cleanedConfigFile, err)
	}
	return nil
}
