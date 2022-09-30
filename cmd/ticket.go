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
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/panjf2000/ants/v2"
	"github.com/rsvihladremio/ssdownloader/downloader"
	"github.com/rsvihladremio/ssdownloader/futils"
	"github.com/rsvihladremio/ssdownloader/link"
	"github.com/rsvihladremio/ssdownloader/sendsafely"
	"github.com/rsvihladremio/ssdownloader/zendesk"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var useZendeskPassword bool
var onlySendSafelyLinks bool

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
		d := downloader.NewGenericDownloader(DownloadBufferSize)
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
		zendeskAPI := zendesk.NewClient(C.ZendeskEmail, password, C.ZendeskDomain, Verbose)
		ticketID := args[0]

		// Handle paging when ticket comments > 100
		var commentLinkTuples []zendesk.CommentTextWithLink
		var attachments []zendesk.Attachment
		emptyString := ""
		var nextPage *string = &emptyString

		for nextPage != nil {
			var commentResults []zendesk.CommentTextWithLink
			results, err := zendeskAPI.GetTicketComentsJSON(ticketID, nextPage)
			if err != nil {
				log.Fatal(err)
			}
			commentResults, nextPage, err = zendesk.GetLinksFromComments(results)
			if err != nil {
				log.Fatalf("unable parse comments with error '%v'", err)
			}
			// Append to array for comments (short hand with "...")
			commentLinkTuples = append(commentLinkTuples, commentResults...)

			attResults, _, err := zendesk.GetAttachmentsFromComments(results)
			if err != nil {
				log.Fatalf("unable parse attachments with error '%v'", err)
			}
			// Append to array for attachments
			attachments = append(attachments, attResults...)
		}

		p, err := ants.NewPool(DownloadThreads)
		if err != nil {
			log.Fatalf("cannot initialize thread pool due to error %v", err)
		}

		defer p.Release()
		var m sync.Mutex
		var allInvalidFiles []string
		var wg sync.WaitGroup
		for _, c := range commentLinkTuples {
			url := c.URL
			if strings.HasPrefix(url, "https://sendsafely") {
				linkParts, err := link.ParseLink(url)
				if err != nil {
					log.Printf("unexpected error '%v' reading url '%v'", err, url)
				}
				packageID := linkParts.PackageCode
				wg.Add(1)
				err = p.Submit(func() {
					outDir, invalidFiles, err := sendsafely.DownloadFilesFromPackage(d, packageID, linkParts.KeyCode, C, filepath.Join("tickets", ticketID), Verbose)
					if err != nil {
						log.Printf("error downloading %v", err)
					} else {
						m.Lock()
						allInvalidFiles = append(allInvalidFiles, invalidFiles...)
						m.Unlock()
						err = os.WriteFile(filepath.Join(outDir, "comment.txt"), []byte(c.Body), 0600)
						if err != nil {
							log.Printf("error writing comments %v", err)
						}
					}
					wg.Done()
				})
				if err != nil {
					log.Fatalf("cannot itialize sendsafely download due to error %v", err)
				}
			}
		}
		if !onlySendSafelyLinks {
			for _, a := range attachments {
				wg.Add(1)
				err = p.Submit(func() {
					if invalidFiles, err := DownloadNonSendSafelyLink(d, a, ticketID); err != nil {
						log.Printf("WARN: error '%v' processing attachment %v skipping", err, a.FileName)
					} else {
						m.Lock()
						allInvalidFiles = append(allInvalidFiles, invalidFiles...)
						m.Unlock()
					}
					wg.Done()
				})
				if err != nil {
					log.Fatalf("cannot itialize sendsafely download due to error %v", err)
				}
			}
		}
		wg.Wait()
		if result := InvalidFilesReport(allInvalidFiles); result != "" {
			fmt.Println(result)
		}
	},
}

func DownloadNonSendSafelyLink(d *downloader.GenericDownloader, a zendesk.Attachment, ticketID string) (invalidFiles []string, err error) {
	if a.Deleted {
		return invalidFiles, fmt.Errorf("attachment '%v' from comment %v created on %v is marked as deleted, skipping", a.FileName, a.ParentCommentID, a.ParentCommentDate)
	}
	commentDir := fmt.Sprintf("%v_%v", a.ParentCommentDate.Format("2006-01-02T150405Z0700"), a.ParentCommentID)
	downloadDir := filepath.Join(C.DownloadDir, "tickets", ticketID, "attachments", commentDir)
	newFileName := filepath.Join(downloadDir, a.FileName)
	exists, err := futils.FileExists(newFileName)
	if err != nil {
		return invalidFiles, fmt.Errorf("unable to see if there is an existing file named %v due to error %v, skipping download", newFileName, err)
	}
	if exists {
		if !sendsafely.FileSizeMatches(newFileName, a.Size, Verbose) {
			return invalidFiles, fmt.Errorf("already downloaded file %v failed verification and did not meet expected size %v", newFileName, a.Size)
		}
		log.Printf("file '%v' already downloaded skipping", newFileName)
		return invalidFiles, nil
	}

	log.Printf("downloading attachment '%v'", newFileName)
	if _, err := os.Stat(downloadDir); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(downloadDir, 0700)
			if err != nil {
				return invalidFiles, fmt.Errorf("unable to make dir %v due to error '%v", downloadDir, err)
			}
		} else {
			return invalidFiles, fmt.Errorf("unable to read dir %v due to error '%v", downloadDir, err)
		}
	}

	if err := d.DownloadFile(newFileName, a.ContentURL); err != nil {
		return invalidFiles, fmt.Errorf("cannot download %v due to error %v, skipping", newFileName, err)
	}
	if !sendsafely.FileSizeMatches(newFileName, a.Size, Verbose) {
		invalidFiles = append(invalidFiles, newFileName)
		return invalidFiles, fmt.Errorf("newly downloaded file %v failed verification and did not meet expected size %v", newFileName, a.Size)
	}
	log.Printf("attachement '%v' is complete", newFileName)

	return invalidFiles, nil
}

func init() {
	rootCmd.AddCommand(ticketCmd)
	ticketCmd.Flags().BoolVarP(&useZendeskPassword, "zendesk-password", "p", false, "Use a password instead of an api key to authenticate against zendesk")
	ticketCmd.Flags().BoolVar(&onlySendSafelyLinks, "sendsafely-only", false, "when true only sendsafely links will be downloaded")
}
