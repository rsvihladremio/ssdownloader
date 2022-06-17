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

func DownloadFilesFromPackage(packageId, keyCode string, c config.Config) error {

	client := NewSendSafelyClient(c.SsApiKey, c.SsApiSecret)
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

	//make config directory for this package code if it does not exist
	downloadDir := filepath.Join(c.DownloadDir, p.PackageCode)

	_, err = os.Stat(downloadDir)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		configDir := filepath.Dir(downloadDir)
		log.Printf("making dir %v", downloadDir)
		err = os.Mkdir(downloadDir, 0700)
		if err != nil {
			return fmt.Errorf("unable to create download dir '%v' due to error '%v'", configDir, err)
		}
	}

	//naively make a lot of requests and use primitives to wait until it's all done
	var wgDownloadUrls sync.WaitGroup

	for _, f := range p.Files {
		log.Printf("downloading %v", f.FileName)
		fileName := f.FileName
		segmentRequestInformation := calculateExecutionCalls(f.Parts)
		for _, segment := range segmentRequestInformation {
			start := segment.StartSegment
			end := segment.EndSegment
			wgDownloadUrls.Add(1)
			go func() {
				defer wgDownloadUrls.Done()
				urls, err := client.GetDownloadUrlsForFile(
					p,
					f.FileId,
					keyCode,
					start,
					end,
				)
				if err != nil {
					log.Printf("unable to download file '%v' due to error '%v' while attemping to get the download url, skipping file", f.FileName, err)
					return
				}

				for _, url := range urls {
					downloadUrl := url.Url
					filePart := url.Part
					//spawning yet another go routine so adding this to the wait group
					wgDownloadUrls.Add(1)
					go func() {
						defer wgDownloadUrls.Done()
						// we add the encrypted value here to make it obvious on reading the directory what step in the download process it is at
						tmpName := fmt.Sprintf("%v.%v.encrypted", fileName, filePart)
						downloadLoc := filepath.Join(downloadDir, tmpName)
						err = downloadFile(downloadLoc, downloadUrl)
						if err != nil {
							log.Printf("unable to download file %v due to error '%v'", downloadLoc, err)
							return
						}
						newFileName, err := DecryptPart(downloadLoc, p.ServerSecret, keyCode)
						if err != nil {
							log.Printf("unable to decrype file %v due to error '%v'", downloadLoc, err)
							return
						}
						log.Printf("file '%v' is decrypted", newFileName)
					}()
				}
			}()
		}
	}
	wgDownloadUrls.Wait()
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
	bytes_written, err := io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("unable to write to filename '%v' due to error '%v'", cleanedFileName, err)
	}

	log.Printf("file %v complete with %v bytes written", fileName, bytes_written)
	//TODO make buffer size adjustable
	// buf := make([]byte, 4096)
	// bytes_read := 0
	// bytes_written := 0
	// for {
	// 	n, err := resp.Body.Read(buf)
	// 	if errors.Is(err, io.EOF) {
	// 		log.Printf("file %v complete with %v bytes written", fileName, bytes_written)
	// 		break
	// 	}
	// 	if err != nil {
	// 		return fmt.Errorf("unable to read body into buffer due to error '%v' while downloading url '%v' and writing to file '%v', already read %v bytes", err, url, cleanedFileName, bytes_read)
	// 	}
	// 	bytes_read += n
	// 	if n == 0 {
	// 		break
	// 	}
	// 	if nw, err := f.Write(buf[:n]); err != nil {
	// 		return fmt.Errorf("unable to write to filename '%v' due to error '%v'", cleanedFileName, err)
	// 	} else {
	// 		bytes_written += nw
	// 	}
	// }
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
