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
	"runtime"

	"github.com/rsvihladremio/ssdownloader/cmd/config"
	"github.com/spf13/cobra"
)

var cfgFile string

var C config.Config
var Verbose bool
var DownloadBufferSize int
var DownloadThreads int
var MaxFileSizeGiB int

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
var GitSha = "unknown"

// Version is pulled from the branch name and set in the build and release scripts
var Version = "unknownVersion"

var platform = runtime.GOOS
var arch = runtime.GOARCH
var programLevel = new(slog.LevelVar) // Info by default

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
		slog.Error("do not have access to home directory, critical error", "error_msg", err)
		os.Exit(1)
	}
	return fmt.Sprintf("%v/.sendsafely", userDir)
}

func init() {
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(h))
	fmt.Println(PrintHeader(Version, platform, arch, GitSha))
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose logging")
	rootCmd.PersistentFlags().StringVar(&C.SsAPIKey, "ss-api-key", "", "the SendSafely API key")
	rootCmd.PersistentFlags().StringVar(&C.SsAPISecret, "ss-api-secret", "", "the SendSafely API secret")
	rootCmd.PersistentFlags().StringVar(&C.ZendeskDomain, "zendesk-subdomain", "", "the customer domain part of the zendesk url that you login against ie https://test.zendesk.com would be 'test'")
	rootCmd.PersistentFlags().StringVar(&C.ZendeskEmail, "zendesk-email", "", "zendesk email address")
	rootCmd.PersistentFlags().StringVar(&C.ZendeskToken, "zendesk-token", "", "zendesk api token")
	rootCmd.PersistentFlags().StringVar(&C.DownloadDir, "download-dir", DefaultDownloadDir(), "base directory to put downloads")
	rootCmd.PersistentFlags().IntVarP(&DownloadBufferSize, "download-buffer-size-kb", "b", 4096, "buffer size in kb to use during downloads")
	rootCmd.PersistentFlags().IntVarP(&DownloadThreads, "download-threads", "t", 8, "number of threads to use when downloading")
	rootCmd.PersistentFlags().IntVarP(&MaxFileSizeGiB, "max-file-size-gib", "m", 10, "max file size in GiB (base 1000) to download, anything over this size will be skipped")
	initConfig()
}

// made avaiable to subcommands via this method
func SetVerbosity() {
	//initialize logger with verbosity set
	if Verbose {
		programLevel.Set(slog.LevelDebug)
	} else {
		programLevel.Set(slog.LevelInfo)
	}
}

// initConfig reads in config file if present
func initConfig() {
	fileToLoad, err := config.ReadConfigFile(cfgFile)
	if err != nil {
		slog.Error("unhandled error loading configuration file", "file_name", cfgFile, "error_msg", err)
		os.Exit(1)
	}
	_, err = os.Stat(fileToLoad)
	if !os.IsNotExist(err) {
		err := config.Load(cfgFile, &C)
		if err != nil {
			slog.Error("unable to load config file", "file_name", fileToLoad, "error_msg", err)
			os.Exit(1)
		}
	}
}
