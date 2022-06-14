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
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
		ssApiKey := fmt.Sprintf("%v", viper.Get("ss-api-key"))
		ssApiSecret := fmt.Sprintf("%v", viper.Get("ss-api-secret"))
		if ssApiKey == "" {
			log.Fatalf("ss-api-key is not set and this is required")
		}
		if ssApiSecret == "" {
			log.Fatalf("ss-api-secret is not set and this is required")
		}
		ssClient := sendsafely.NewSendSafelyClient(ssApiKey, ssApiSecret)
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
