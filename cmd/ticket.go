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
	"sync"
	"syscall"

	"github.com/panjf2000/ants/v2"
	"github.com/rsvihladremio/ssdownloader/downloader"
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
		zendeskApi := zendesk.NewZenDeskAPI(C.ZendeskEmail, password, C.ZendeskDomain, Verbose)
		ticketId := args[0]

		results, err := zendeskApi.GetTicketComentsJSON(ticketId)
		if err != nil {
			log.Fatal(err)
		}
		urls, err := zendesk.GetLinksFromComments(results)
		if err != nil {
			log.Fatalf("unable parse comments with error '%v'", err)
		}
		p, err := ants.NewPool(DownloadThreads)
		if err != nil {
			log.Fatalf("cannot initialize thread pool due to error %v", err)
		}

		defer p.Release()
		var wg sync.WaitGroup
		for _, url := range urls {
			if strings.HasPrefix(url, "https://sendsafely") {
				linkParts, err := link.ParseLink(url)
				if err != nil {
					log.Fatalf("unexpected error '%v' reading url '%v'", err, url)
				}
				packageId := linkParts.PackageCode
				wg.Add(1)
				err = p.Submit(func() {
					err := sendsafely.DownloadFilesFromPackage(d, packageId, linkParts.KeyCode, C, filepath.Join("tickets", ticketId), Verbose)
					if err != nil {
						log.Printf("error downloading %v", err)
					}
					wg.Done()
				})
				if err != nil {
					log.Fatalf("cannot itialize sendsafely download due to error %v", err)
				}
			}
		}
		if !onlySendSafelyLinks {
			attachments, err := zendesk.GetAttachementsFromComments(results)
			if err != nil {
				log.Fatal(err)
			}
			for _, a := range attachments {
				wg.Add(1)
				err = p.Submit(func() {
					if err := DownloadNonSendSafelyLink(d, a, ticketId); err != nil {
						log.Printf("WARN: error '%v' processing attachement %v skipping", err, a.FileName)
					}
					wg.Done()
				})
				if err != nil {
					log.Fatalf("cannot itialize sendsafely download due to error %v", err)
				}
			}
		}
		wg.Wait()
	},
}

func DownloadNonSendSafelyLink(d *downloader.GenericDownloader, a zendesk.Attachment, ticketId string) error {
	if a.Deleted {
		return fmt.Errorf("attachment '%v' from comment %v created on %v is marked as deleted, skipping", a.FileName, a.ParentCommentID, a.ParentCommentDate)
	}
	commentDir := fmt.Sprintf("%v_%v", a.ParentCommentDate.Format("2006-01-02T150405Z0700"), a.ParentCommentID)
	downloadDir := filepath.Join(C.DownloadDir, "tickets", ticketId, "attachments", commentDir)
	newFileName := filepath.Join(downloadDir, a.FileName)
	if sendsafely.FileSizeMatches(newFileName, a.Size, Verbose) {
		log.Printf("file '%v' already downloaded skipping", newFileName)
		return nil
	}
	log.Printf("downloading attachment '%v'", newFileName)
	if _, err := os.Stat(downloadDir); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(downloadDir, 0700)
			if err != nil {
				return fmt.Errorf("unable to make dir %v due to error '%v", downloadDir, err)
			}
		} else {
			return fmt.Errorf("unable to read dir %v due to error '%v", downloadDir, err)
		}
	}

	if err := d.DownloadFile(newFileName, a.ContentURL); err != nil {
		return fmt.Errorf("cannot download %v due to error %v, skipping", newFileName, err)
	}
	if !sendsafely.FileSizeMatches(newFileName, a.Size, Verbose) {
		defer func() {
			if err := os.Remove(newFileName); err != nil {
				log.Printf("WARN: unexpected error '%v' trying to remove file '%v', remove manually", err, newFileName)
			}
		}()
		return fmt.Errorf("file %v failed verification and did not meet expected size %v", newFileName, a.Size)
	}
	log.Printf("attachement '%v' is complete", newFileName)

	return nil
}

func init() {
	rootCmd.AddCommand(ticketCmd)
	ticketCmd.Flags().BoolVarP(&useZendeskPassword, "zendesk-password", "p", false, "Use a password instead of an api key to authenticate against zendesk")
	ticketCmd.Flags().BoolVar(&onlySendSafelyLinks, "sendsafely-only", false, "when true only sendsafely links will be downloaded")
}
