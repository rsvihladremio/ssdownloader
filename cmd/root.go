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
package cmd

import (
	"log"
	"os"

	"github.com/rsvihladremio/ssdownloader/cmd/config"
	"github.com/spf13/cobra"
)

var cfgFile string

var C config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ssdownloader",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&C.SsApiKey, "ss-api-key", "", "the SendSafely API key")
	rootCmd.PersistentFlags().StringVar(&C.SsApiSecret, "ss-api-secret", "", "the SendSafely API secret")
	rootCmd.PersistentFlags().StringVar(&C.ZendeskDomain, "zendesk-domain", "", "the customer domain part of the zendesk url that you login against ie https://test.zendesk.com would be 'test'")
	rootCmd.PersistentFlags().StringVar(&C.ZendeskEmail, "zendesk-email", "", "zendesk email address")
	rootCmd.PersistentFlags().StringVar(&C.ZendeskToken, "zendesk-token", "", "zendesk api token")
	initConfig()
}

// initConfig reads in config file if present
func initConfig() {
	fileToLoad, err := config.ReadConfigFile(cfgFile)
	if err != nil {
		log.Fatalf("unhandled error loading configuration file '%v' with error '%v'", cfgFile, err)
	}
	_, err = os.Stat(fileToLoad)
	log.Printf("config file '%v'", fileToLoad)
	if !os.IsNotExist(err) {
		log.Printf("loading configuration file '%v'", fileToLoad)
		err := config.Load(cfgFile, &C)
		if err != nil {
			log.Printf("unable to load config file '%v' due to error '%v'", fileToLoad, err)
		}
	}
}
