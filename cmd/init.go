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
