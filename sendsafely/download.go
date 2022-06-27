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
package sendsafely

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/rsvihladremio/ssdownloader/cmd/config"
)

type PartRequests struct {
	StartSegment int
	EndSegment   int
}

func DownloadFilesFromPackage(packageId, keyCode string, c config.Config, subDirToDownload string, verbose bool) error {

	client := NewSendSafelyClient(c.SsApiKey, c.SsApiSecret, verbose)
	p, err := client.RetrievePackgeById(packageId)
	if err != nil {
		return err
	}
	//making top level download directory if it does not exist
	_, err = os.Stat(c.DownloadDir)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		configDir := filepath.Dir(c.DownloadDir)
		log.Printf("making dir %v", c.DownloadDir)
		err = os.Mkdir(c.DownloadDir, 0700)
		if err != nil {
			return fmt.Errorf("unable to create download dir '%v' due to error '%v'", configDir, err)
		}
	}

	shortPackageId := p.PackageId
	//make config directory for this package code if it does not exist
	downloadDir := filepath.Join(c.DownloadDir, subDirToDownload, shortPackageId)

	_, err = os.Stat(downloadDir)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		configDir := filepath.Dir(downloadDir)
		log.Printf("making dir %v", downloadDir)
		err = os.MkdirAll(downloadDir, 0700)
		if err != nil {
			return fmt.Errorf("unable to create download dir '%v' due to error '%v'", configDir, err)
		}
	}

	//naively make a lot of requests and use primitives to wait until it's all done
	var apiCalls sync.WaitGroup

	for _, f := range p.Files {
		log.Printf("downloading %v", f.FileName)
		fileName := f.FileName
		//TODO make part of verbose flag
		//log.Printf("file %v has %v parts and a total file size of %v", f.FileName, f.Parts, f.FileSize)
		segmentRequestInformation := calculateExecutionCalls(f.Parts)
		for _, segment := range segmentRequestInformation {
			start := segment.StartSegment
			end := segment.EndSegment
			apiCalls.Add(1)
			fileId := f.FileId
			go func() {
				defer apiCalls.Done()
				urls, err := client.GetDownloadUrlsForFile(
					p,
					fileId,
					keyCode,
					start,
					end,
				)
				if err != nil {
					log.Printf("unable to download file '%v' due to error '%v' while attemping to get the download url, skipping file", fileName, err)
					return
				}
				var wgDownloadUrls sync.WaitGroup
				var m sync.Mutex
				var fileNames []string
				var errMutex sync.Mutex
				var failedFiles []string
				for i := range urls {
					index := i

					//spawning yet another go routine so adding this to the wait group
					wgDownloadUrls.Add(1)
					go func() {
						defer wgDownloadUrls.Done()
						url := urls[index]
						downloadUrl := url.Url
						filePart := url.Part
						// we add the encrypted value here to make it obvious on reading the directory what step in the download process it is at
						tmpName := fmt.Sprintf("%v.%v.encrypted", fileName, filePart)
						downloadLoc := filepath.Join(downloadDir, tmpName)
						err = downloadFile(downloadLoc, downloadUrl)
						if err != nil {
							log.Printf("unable to download file %v due to error '%v'", downloadLoc, err)
							return
						}
						newFileName, err := DecryptPart(downloadLoc, p.ServerSecret, keyCode)
						m.Lock()
						fileNames = append(fileNames, newFileName)
						m.Unlock()
						if err != nil {
							errMutex.Lock()
							failedFiles = append(failedFiles, newFileName)
							errMutex.Unlock()
							log.Printf("unable to decrypt file %v due to error '%v'", downloadLoc, err)

							return
						}
						//TODO behind verbose flag
						//log.Printf("file '%v' is decrypted", newFileName)
					}()
				}
				wgDownloadUrls.Wait()
				if len(failedFiles) > 0 {
					log.Printf("there were %v failed files therefore not going to bother combining parts for file '%v'", len(failedFiles), fileName)
					return
				}
				newFile, err := CombineFiles(fileNames, verbose)
				if err != nil {
					log.Printf("unable to combine downloaded parts for fileName '%v' due to error '%v'", fileName, err)
				} else {
					log.Printf("file '%v is complete", newFile)
				}
			}()
		}
	}
	apiCalls.Wait()
	return nil
}

func downloadFile(fileName, url string) error {
	// making sure there are no goofy file names that overwrite critical files
	cleanedFileName := filepath.Clean(fileName)
	f, err := os.Create(cleanedFileName)
	if err != nil {
		return fmt.Errorf("unable to create the destination file '%v' due to error '%v'", cleanedFileName, err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Printf("WARN: unable to close file handle for file '%v' due to error '%v'", cleanedFileName, err)
		}
	}()
	// is technically a security violation according to https://securego.io/docs/rules/g107.html
	// but in reality based on the application is unavoidable and a risk of using SendSafely
	// ignoring the rule in the ./script/audit file. Used suggestion from
	// https://stackoverflow.com/questions/70281883/golang-untaint-url-variable-to-fix-gosec-warning-g107
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("unable to retrieve url '%v' due to error '%v'", url, err)
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Printf("WARN: unable to close body handle for url '%v' due to error '%v'", url, err)
		}
	}()
	//TODO make optional with verbose flag
	//log.Printf("file %v complete with %v bytes written", fileName, bytes_written)
	//TODO make buffer size adjustable
	buf := make([]byte, 4096*1024)
	_, err = io.CopyBuffer(f, resp.Body, buf)
	if err != nil {
		return fmt.Errorf("unable to write to filename '%v' due to error '%v'", cleanedFileName, err)
	}

	return nil
}

func calculateExecutionCalls(parts int) []PartRequests {
	var requests []PartRequests
	if parts == 0 {
		return []PartRequests{}
	}

	//submit 25 parts at once
	fullRequests := parts / 25
	counter := 1
	for i := 0; i < fullRequests; i++ {
		requests = append(requests, PartRequests{
			StartSegment: counter,
			EndSegment:   counter + 24,
		})
		counter += 25
	}
	remainder := parts % 25
	if remainder > 0 {
		requests = append(requests, PartRequests{
			StartSegment: counter,
			EndSegment:   parts,
		})
	}
	return requests
}
