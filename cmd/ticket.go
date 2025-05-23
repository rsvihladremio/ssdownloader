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

// cmd package contains all the command line flag configuration
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/panjf2000/ants/v2"
	"github.com/rsvihladremio/ssdownloader/downloader"
	"github.com/rsvihladremio/ssdownloader/futils"
	"github.com/rsvihladremio/ssdownloader/link"
	"github.com/rsvihladremio/ssdownloader/reporting"
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
	Run: func(_ *cobra.Command, args []string) {
		SetVerbosity()
		d := downloader.NewGenericDownloader(DownloadBufferSize)
		password := ""
		if useZendeskPassword {
			fmt.Println("enter password:")
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				slog.Error("unexpected error reading password", "error_msg", err)
				os.Exit(1)
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
		nextPage := &emptyString

		for nextPage != nil {
			var commentResults []zendesk.CommentTextWithLink
			results, err := zendeskAPI.GetTicketComentsJSON(ticketID, nextPage)
			if err != nil {
				slog.Error("unexpected error getting ticket comments", "error_msg", err)
				os.Exit(1)
			}
			commentResults, nextPage, err = zendesk.GetLinksFromComments(results)
			if err != nil {
				slog.Error("unable to parse ticket comments", "error_msg", err)
				os.Exit(1)
			}
			// Append to array for comments (short hand with "...")
			commentLinkTuples = append(commentLinkTuples, commentResults...)

			attResults, _, err := zendesk.GetAttachmentsFromComments(results)
			if err != nil {
				slog.Error("unable parse attachments", "error_msg", err)
				os.Exit(1)
			}
			// Append to array for attachments
			attachments = append(attachments, attResults...)
		}

		p, err := ants.NewPool(DownloadThreads)
		if err != nil {
			slog.Error("cannot initialize thread pool", "error_msg", err)
			os.Exit(1)
		}

		defer p.Release()
		var m sync.Mutex
		var allInvalidFiles []string
		var wg sync.WaitGroup
		client := sendsafely.NewClient(C.SsAPIKey, C.SsAPISecret, Verbose)
		for _, c := range commentLinkTuples {
			url := c.URL
			if strings.HasPrefix(url, "https://sendsafely") {
				linkParts, err := link.ParseLink(url)
				if err != nil {
					slog.Error("unexpected error reading url", "error_msg", err, "url", url)
				}
				packageID := linkParts.PackageCode
				wg.Add(1)
				err = p.Submit(func() {
					keyCode := linkParts.KeyCode
					a := sendsafely.DownloadArgs{
						PackageID:        packageID,
						KeyCode:          keyCode,
						SubDirToDownload: filepath.Join("tickets", ticketID),
						DownloadDir:      C.DownloadDir,
						MaxFileSizeByte:  int64(MaxFileSizeGiB) * 1000000000,
						Verbose:          Verbose,
						SkipList:         []string{},
					}
					outDir, invalidFiles, err := sendsafely.DownloadFilesFromPackage(client, d, a)
					if err != nil {
						slog.Error("error downloading files from package", "error_msg", err, "package_id", packageID)
					} else {
						m.Lock()
						allInvalidFiles = append(allInvalidFiles, invalidFiles...)
						m.Unlock()
						outputFile := filepath.Join(outDir, "comment.txt")
						err = os.WriteFile(outputFile, []byte(c.Body), 0600)
						if err != nil {
							slog.Error("error writing comment text", "error_msg", err, "comment_url", c.URL, "output_file", outputFile)
						}
					}
					wg.Done()
				})
				if err != nil {
					slog.Error("cannot initialize sendsafely download", "error_msg", err)
					os.Exit(1)
				}
			}
		}
		if !onlySendSafelyLinks {
			for _, a := range attachments {
				wg.Add(1)
				err = p.Submit(func() {
					if invalidFiles, err := DownloadNonSendSafelyLink(d, a, ticketID); err != nil {
						slog.Warn("error processing attachment; skipping", "error_msg", err, "attachement", a.FileName)
					} else {
						m.Lock()
						allInvalidFiles = append(allInvalidFiles, invalidFiles...)
						m.Unlock()
					}
					wg.Done()
				})
				if err != nil {
					slog.Error("cannot initialize sendsafely download", "error_msg", err)
					os.Exit(1)
				}
			}
		}
		wg.Wait()
		if result := InvalidFilesReport(allInvalidFiles); result != "" {
			fmt.Println(result)
		}
		fmt.Println(Report(reporting.GetTotalFiles(), reporting.GetTotalSkipped(), reporting.GetTotalFailed(), reporting.GetTotalBytes(), reporting.GetMaxFileSizeBytes()))
	},
}

func DownloadNonSendSafelyLink(d downloader.GenericDownloader, a zendesk.Attachment, ticketID string) (invalidFiles []string, err error) {
	reporting.AddFile()
	if a.Deleted {
		reporting.AddFailed()
		return invalidFiles, fmt.Errorf("attachment '%v' from comment %v created on %v is marked as deleted, skipping", a.FileName, a.ParentCommentID, a.ParentCommentDate)
	}
	commentDir := fmt.Sprintf("%v_%v", a.ParentCommentDate.Format("2006-01-02T150405Z0700"), a.ParentCommentID)
	downloadDir := filepath.Join(C.DownloadDir, "tickets", ticketID, "attachments", commentDir)
	newFileName := filepath.Join(downloadDir, a.FileName)
	exists, err := futils.FileExists(newFileName)
	if err != nil {
		reporting.AddFailed()
		return invalidFiles, fmt.Errorf("unable to see if there is an existing file named %v due to error %v, skipping download", newFileName, err)
	}
	if exists {
		if err := sendsafely.FileSizeCheck(newFileName, a.Size); err != nil {
			reporting.AddFailed()
			return invalidFiles, fmt.Errorf("already downloaded file %v failed verification and is not readable %v", newFileName, a.Size)
		}
		reporting.AddSkip()
		slog.Debug("file already downloaded skipping", "file_name", newFileName)
		return invalidFiles, nil
	}
	fmt.Print(".")
	slog.Debug("downloading attachment", "file_name", newFileName)
	if _, err := os.Stat(downloadDir); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(downloadDir, 0700)
			if err != nil {
				reporting.AddFailed()
				return invalidFiles, fmt.Errorf("unable to make dir %v due to error '%v", downloadDir, err)
			}
		} else {
			reporting.AddFailed()
			return invalidFiles, fmt.Errorf("unable to read dir %v due to error '%v", downloadDir, err)
		}
	}

	if err := d.DownloadFile(newFileName, a.ContentURL); err != nil {
		reporting.AddFailed()
		return invalidFiles, fmt.Errorf("cannot download %v due to error %v, skipping", newFileName, err)
	}
	if err := sendsafely.FileSizeCheck(newFileName, a.Size); err != nil {
		invalidFiles = append(invalidFiles, newFileName)
		return invalidFiles, fmt.Errorf("newly downloaded file %v failed verification and is not readable %v", newFileName, a.Size)
	}
	fmt.Print(".")
	slog.Debug("attachement download complete", "file_name", newFileName)
	reporting.AddBytes(a.Size)

	return invalidFiles, nil
}

func Report(totalFiles int, totalSkipped int, totalFailed int, totalBytes int64, maxBytes int64) string {
	return fmt.Sprintf(`
================================
= ssdownloader summary         =
================================
= total files       : %v
= total succeeded   : %v
= total skipped     : %v
= total failed      : %v
= total bytes       : %v
= max bytes         : %v
================================`, totalFiles, totalFiles-(totalSkipped+totalFailed), totalSkipped, totalFailed, sendsafely.Human(totalBytes), sendsafely.Human(maxBytes))
}

func init() {
	rootCmd.AddCommand(ticketCmd)
	ticketCmd.Flags().BoolVarP(&useZendeskPassword, "zendesk-password", "p", false, "Use a password instead of an api key to authenticate against zendesk")
	ticketCmd.Flags().BoolVar(&onlySendSafelyLinks, "sendsafely-only", false, "when true only sendsafely links will be downloaded")
}
