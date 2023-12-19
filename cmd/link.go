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

	"github.com/spf13/cobra"

	"github.com/rsvihladremio/ssdownloader/downloader"
	"github.com/rsvihladremio/ssdownloader/link"
	"github.com/rsvihladremio/ssdownloader/sendsafely"
)

// linkCmd represents the link command
var linkCmd = &cobra.Command{
	Use:   "link <sendsafely package link>",
	Short: "link command downloads all files in a package provided by a url from sendsafely",
	Long: `link command downloads all files in a package provided by a url from sendsafely. Example below:

	ssdownloader link "https://sendsafely.tester.com/receive/?thread=MYTHREAD&packageCode=MYPKGCODE#keyCode=MYKEYCODE"
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetVerbosity()
		if C.SsAPIKey == "" {
			slog.Error("ss-api-key is not set and this is required")
			os.Exit(1)
		}
		if C.SsAPISecret == "" {
			slog.Error("ss-api-secret is not set and this is required")
			os.Exit(1)
		}
		url := args[0]
		linkParts, err := link.ParseLink(url)
		if err != nil {
			slog.Error("unexpected error reading url", "url", url, "error_msg", err)
			os.Exit(1)
		}
		packageID := linkParts.PackageCode
		d := downloader.NewGenericDownloader(DownloadBufferSize)
		client := sendsafely.NewClient(C.SsAPIKey, C.SsAPISecret, Verbose)
		a := sendsafely.DownloadArgs{
			DownloadDir:      C.DownloadDir,
			MaxFileSizeByte:  int64(MaxFileSizeGiB) * 1000000000,
			KeyCode:          linkParts.KeyCode,
			PackageID:        packageID,
			SkipList:         []string{},
			SubDirToDownload: "packages",
			Verbose:          Verbose,
		}
		_, invalidFiles, err := sendsafely.DownloadFilesFromPackage(client, d, a)
		if err != nil {
			slog.Error("unexpected error downloading files", "error_msg", err)
			os.Exit(1)
		}

		if result := InvalidFilesReport(invalidFiles); result != "" {
			fmt.Println(result)
		}
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)
}
