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

//config package handles the reading and writing of the app configuration file
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	SsAPIKey      string
	SsAPISecret   string
	ZendeskDomain string
	ZendeskEmail  string
	ZendeskToken  string
	DownloadDir   string
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

func Save(c Config, cfgFile string) (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("unable to convert configuration to json file due to error '%v'", err)
	}
	fileToSave, err := ReadConfigFile(cfgFile)
	if err != nil {
		return "", fmt.Errorf("trying to get the path to the configuration resulted in error '%v'", err)
	}

	// best security practice
	cleanedConfigFile := filepath.Clean(fileToSave)

	_, err = os.Stat(cleanedConfigFile)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		configDir := filepath.Dir(cleanedConfigFile)
		err = os.MkdirAll(configDir, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("unable to create configuration dir '%v' due to error '%v'", configDir, err)
		}
	}

	err = os.WriteFile(cleanedConfigFile, b, os.FileMode(0600))
	if err != nil {
		return "", fmt.Errorf("unable to write configuration file to location '%v' due to error '%v'", cleanedConfigFile, err)
	}
	return cleanedConfigFile, nil
}
