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
	"os"
	"path/filepath"
	"runtime/pprof"

	"github.com/spf13/cobra"

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
		if C.SsApiKey == "" {
			log.Fatalf("ss-api-key is not set and this is required")
		}
		if C.SsApiSecret == "" {
			log.Fatalf("ss-api-secret is not set and this is required")
		}
		url := args[0]
		linkParts, err := link.ParseLink(url)
		if err != nil {
			log.Fatalf("unexpected error '%v' reading url '%v'", err, url)
		}
		packageId := linkParts.PackageCode
		err = sendsafely.DownloadFilesFromPackage(packageId, linkParts.KeyCode, C, filepath.Join("packages", packageId), Verbose)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)
}
