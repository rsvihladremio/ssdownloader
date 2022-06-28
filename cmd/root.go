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
	"runtime"

	"github.com/rsvihladremio/ssdownloader/cmd/config"
	"github.com/spf13/cobra"
)

var cfgFile string

var C config.Config
var Verbose bool
var DownloadBufferSize int
var DownloadThreads int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ssdownloader",
	Short: "ssdownloader downloads file from sendsafely either via ticket number in zendesk or via sendsafely link",
	Long: `ssdownloader downloads file from sendsafely either via ticket number in zendesk or via sendsafely link see
	the following examples:

//by link
ssdownloader link "https://sendsafely.tester.com/receive/?thread=MYTHREAD&packageCode=MYPKGCODE#keyCode=MYKEYCODE"

//by zendesk ticket
ssdownloader ticket 111111
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func PrintHeader(version, platform, arch, gitSha string) string {
	return fmt.Sprintf("ssdownloader %v-%v-%v-%v\n", version, gitSha, platform, arch)
}

// GitSha is added from the build and release scripts
var GitSha string = "unknown"

// Version is pulled from the branch name and set in the build and release scripts
var Version string = "unknownVersion"

var platform string = runtime.GOOS
var arch string = runtime.GOARCH
var CpuProfile string
var MemProfile string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func DefaultDownloadDir() string {
	userDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("do not have access to home directory, critical error '%v'", err)
	}
	return fmt.Sprintf("%v/.sendsafely", userDir)
}

func init() {
	fmt.Println(PrintHeader(Version, platform, arch, GitSha))
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose logging")
	rootCmd.PersistentFlags().StringVar(&C.SsApiKey, "ss-api-key", "", "the SendSafely API key")
	rootCmd.PersistentFlags().StringVar(&C.SsApiSecret, "ss-api-secret", "", "the SendSafely API secret")
	rootCmd.PersistentFlags().StringVar(&C.ZendeskDomain, "zendesk-subdomain", "", "the customer domain part of the zendesk url that you login against ie https://test.zendesk.com would be 'test'")
	rootCmd.PersistentFlags().StringVar(&C.ZendeskEmail, "zendesk-email", "", "zendesk email address")
	rootCmd.PersistentFlags().StringVar(&C.ZendeskToken, "zendesk-token", "", "zendesk api token")
	rootCmd.PersistentFlags().StringVar(&C.DownloadDir, "download-dir", DefaultDownloadDir(), "base directory to put downloads")
	rootCmd.PersistentFlags().StringVar(&CpuProfile, "cpu-profile", "", "where to generate a cpu profile for diagnosing performance issues")
	rootCmd.PersistentFlags().StringVar(&MemProfile, "mem-profile", "", "where to generate a mem profile for diagnosing performance issues")
	rootCmd.PersistentFlags().IntVarP(&DownloadBufferSize, "download-buffer-size-kb", "b", 4096, "buffer size in kb to use during downloads")
	rootCmd.PersistentFlags().IntVarP(&DownloadThreads, "download-threads", "t", 8, "number of threads to use when downloading")
	initConfig()
}

// initConfig reads in config file if present
func initConfig() {
	fileToLoad, err := config.ReadConfigFile(cfgFile)
	if err != nil {
		log.Fatalf("unhandled error loading configuration file '%v' with error '%v'", cfgFile, err)
	}
	_, err = os.Stat(fileToLoad)
	if !os.IsNotExist(err) {
		err := config.Load(cfgFile, &C)
		if err != nil {
			log.Printf("unable to load config file '%v' due to error '%v'", fileToLoad, err)
		}
	}
}
