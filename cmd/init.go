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
ssdownloader init --ss-api-key 2ufjwid --ss-api-secret afj292 --zendesk-domain test --zendesk-email test@example.com --zendesk-token 3jkljf

// via prompts
ssdownloader init
> (sendsafely api key): 2ufjwid
> (sendsafely api secret): afj292
> (zendesk domain): test
> (zendesk email): test@example.com
> (zendesk token): 3jkljf`,

	Run: func(cmd *cobra.Command, args []string) {
		//TODO take parameter to specify configuration file location
		if err := config.Save(C, ""); err != nil {
			log.Fatalf("unexpected error saving configuration '%v'", err)
		}

		log.Println("config file written")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
