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
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"
	"syscall"

	"github.com/rsvihladremio/ssdownloader/link"
	"github.com/rsvihladremio/ssdownloader/sendsafely"
	"github.com/rsvihladremio/ssdownloader/zendesk"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var useZendeskPassword bool

// ticketCmd represents the ticket command
var ticketCmd = &cobra.Command{
	Use:   "ticket <zendesk ticket id>",
	Short: "downloads all files for a given ticket id",
	Long: `download all sendsafely files for a given ticket id example:

		//ticket id is 1111
		ssdownloader ticket 1111

		//ticket id is 1111 and use password instead of zendesk api key (not best practice security)
		ssdownloader ticket 1111 -p 

		`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if CpuProfile != "" {
			f, err := os.Create(CpuProfile)
			if err != nil {
				log.Fatalf("unable to create cpu profile '%v' due to error '%v'", CpuProfile, err)
			}
			if err := pprof.StartCPUProfile(f); err != nil {
				log.Fatalf("unable to start cpu profiling due to error %v", err)
			}
			defer pprof.StopCPUProfile()
		}
		if MemProfile != "" {
			f, err := os.Create(MemProfile)
			if err != nil {
				log.Fatalf("unable to create mem profile at '%v' due to error '%v'", MemProfile, err)
			}
			defer func() {
				if err := f.Close(); err != nil {
					log.Printf("WARN unable to close file handle for file '%v' due to error '%v'", MemProfile, err)
				}
			}()
			if err := pprof.WriteHeapProfile(f); err != nil {
				log.Fatalf("unable to start mem profiling due to error '%v'", err)
			}

		}
		password := ""
		if useZendeskPassword {
			fmt.Println("enter password:")
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				log.Fatalf("unexpected error reading password '%v'", err)
			}
			password = strings.TrimSpace(string(bytePassword))
		} else {
			password = C.ZendeskToken
		}
		zendeskApi := zendesk.NewZenDeskAPI(C.ZendeskEmail, password, C.ZendeskDomain, Verbose)
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
			err = sendsafely.DownloadFilesFromPackage(packageId, linkParts.KeyCode, C, filepath.Join("tickets", ticketId, packageId), Verbose)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(ticketCmd)
	ticketCmd.Flags().BoolVarP(&useZendeskPassword, "zendesk-password", "p", false, "Use a password instead of an api key to authenticate against zendesk")
}
