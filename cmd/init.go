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

//cmd package contains all the command line flag configuration
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rsvihladremio/ssdownloader/cmd/config"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Write and generate a configuration file with prompts or command line flags",
	Long: `Writes and generates a configuration file with prompts or command line flags

examples:

// via command line flags
ssdownloader init --ss-api-key 2ufjwid --ss-api-secret afj292 --zendesk-subdomain test --zendesk-email test@example.com --zendesk-token 3jkljf --download-dir /opt/sendsafely

// via prompts
ssdownloader init
> (sendsafely api key): 2ufjwid
> (sendsafely api secret): afj292
> (zendesk subdomain): test
> (zendesk email): test@example.com
> (zendesk token): 3jkljf
`,

	Run: func(cmd *cobra.Command, args []string) {

		if C.SsAPIKey == "" {
			fmt.Print("(sendsafely api key):")
			n, err := fmt.Scanln(&C.SsAPIKey)
			if err != nil {
				log.Fatal(err)
			}
			if n == 0 {
				fmt.Println("sendsafely api key cannot be blank")
				os.Exit(1)
			}
		}
		if C.SsAPISecret == "" {
			fmt.Print("(sendsafely api secret):")
			n, err := fmt.Scanln(&C.SsAPISecret)
			if err != nil {
				log.Fatal(err)
			}
			if n == 0 {
				fmt.Println("sendsafely api secret cannot be blank")
				os.Exit(1)
			}
		}

		if C.ZendeskDomain == "" {
			fmt.Print("(zendesk subdomain):")
			n, err := fmt.Scanln(&C.ZendeskDomain)
			if err != nil && !strings.Contains(err.Error(), "unexpected newline") {
				log.Fatal(err)
			}
			if n == 0 {
				fmt.Println("zendesk domain was blank, this means `ssdownload ticket` will not function without the --zendesk-domain flag")
			}
		}

		if C.ZendeskEmail == "" {
			fmt.Print("(zendesk email):")
			n, err := fmt.Scanln(&C.ZendeskEmail)
			if err != nil && !strings.Contains(err.Error(), "unexpected newline") {
				log.Fatal(err)
			}
			if n == 0 {
				fmt.Println("zendesk email was blank, this means `ssdownload ticket` will not function without the --zendesk-email flag")
			}
		}

		if C.ZendeskToken == "" {
			fmt.Print("(zendesk token):")
			n, err := fmt.Scanln(&C.ZendeskToken)
			if err != nil && !strings.Contains(err.Error(), "unexpected newline") {
				log.Fatal(err)
			}
			if n == 0 {
				fmt.Println("zendesk token was blank, this means `ssdownload ticket` will not function without the --zendesk-token flag")
			}
		}

		//TODO take parameter to specify configuration file location
		newConf, err := config.Save(C, "")
		if err != nil {
			fmt.Printf("unexpected error saving configuration '%v'\n", err)
			os.Exit(1)
		}
		log.Printf("config file %v written\n", newConf)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
