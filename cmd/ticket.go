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
	"path/filepath"

	"github.com/rsvihladremio/ssdownloader/link"
	"github.com/rsvihladremio/ssdownloader/sendsafely"
	"github.com/rsvihladremio/ssdownloader/zendesk"
	"github.com/spf13/cobra"
)

// ticketCmd represents the ticket command
var ticketCmd = &cobra.Command{
	Use:   "ticket <zendesk ticket id>",
	Short: "downloads all files for a given ticket id",
	Long: `download all sendsafely files for a given ticket id example:

		//ticket id is 11111
		ssdownloader ticket 111111
		`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		zendeskApi := zendesk.NewZenDeskAPI(C.ZendeskEmail, C.ZendeskToken, C.ZendeskDomain)
		ticketId := args[0]
		results, err := zendeskApi.GetTicketComents(ticketId)
		if err != nil {
			log.Fatal(err)
		}
		for _, url := range results {
			linkParts, err := link.ParseLink(url)
			if err != nil {
				log.Fatalf("unexpected error '%v' reading url '%v'", err, url)
			}
			packageId := linkParts.PackageCode
			err = sendsafely.DownloadFilesFromPackage(packageId, linkParts.KeyCode, C, filepath.Join("tickets", ticketId))
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(ticketCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ticketCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ticketCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
