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
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/rsvihladremio/ssdownloader/link"
	"github.com/rsvihladremio/ssdownloader/sendsafely"
	"github.com/rsvihladremio/ssdownloader/zendesk"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var searchCmd = &cobra.Command{
	Use:   "search",
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
	//Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ss := sendsafely.NewSendSafelyClient(C.SsApiKey, C.SsApiSecret)
		zdesk := zendesk.NewZenDeskAPI(C.ZendeskEmail, C.ZendeskToken, C.ZendeskDomain)
		//var packageIds []string
		//var m sync.Mutex
		var wg sync.WaitGroup
		var counter int64
		for i := 116000; i < 117955; i++ {
			wg.Add(1)

			go func() {
				defer func() {
					wg.Done()
				}()
				results, err := zdesk.GetSendSafelyLinksFromTicket(fmt.Sprintf("%v", i))
				if err != nil {
					log.Printf("ticket %v failed with error %v", i, err)
					return
				}
				for _, url := range results {
					linkParts, err := link.ParseLink(url)
					if err != nil {
						log.Printf("unexpected error '%v' reading url '%v'", err, url)
					}
					packageId := linkParts.PackageCode
					ssPkg, err := ss.RetrievePackgeById(packageId)
					if err != nil {
						log.Println(err)
					}
					for _, f := range ssPkg.Files {
						if strings.Contains(f.FileName, "odbc") {
							log.Printf("filename: %v", f.FileName)
						}
					}

				}
				//	m.Lock()

				//	m.Unlock()
			}()
			counter++
			if i%2 == 0 {
				log.Printf("waiting after %v tickets", counter)
				wg.Wait()
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
