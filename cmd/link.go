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
		if C.SsApiKey == "" {
			log.Fatalf("ss-api-key is not set and this is required")
		}
		if C.SsApiSecret == "" {
			log.Fatalf("ss-api-secret is not set and this is required")
		}
		ssClient := sendsafely.NewSendSafelyClient(C.SsApiKey, C.SsApiSecret)
		url := args[0]
		linkParts, err := link.ParseLink(url)
		if err != nil {
			log.Fatalf("unexpected error '%v' reading url '%v'", err, url)
		}
		packageId := linkParts.PackageCode
		p, err := ssClient.RetrievePackgeById(packageId)
		if err != nil {
			log.Fatalf("unexpected error '%v' retrieving packageId '%v'", err, packageId)
		}
		log.Printf("%v", p)
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)
}
