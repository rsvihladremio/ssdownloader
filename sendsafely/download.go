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
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/rsvihladremio/ssdownloader/cmd/config"
	"github.com/rsvihladremio/ssdownloader/downloader"
)

type PartRequests struct {
	StartSegment int
	EndSegment   int
}

func FileSizeMatches(fileName string, fileSize int64, verbose bool) bool {
	fi, err := os.Stat(fileName)
	if err != nil {
		if verbose {
			log.Printf("DEBUG: unable to check file size for %v due to error %v", fileName, err)
		}
		return false
	}

	if verbose {
		log.Printf("DEBUG: file comparison is %v current versus %v to download", fi.Size(), fileSize)
	}
	return fi.Size() == fileSize
}

func DownloadFilesFromPackage(d *downloader.GenericDownloader, packageId, keyCode string, c config.Config, subDirToDownload string, verbose bool) error {

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

	for _, f := range p.Files {

		fileName := f.FileName
		parts := f.Parts
		fileSize := f.FileSize
		fileId := f.FileId
		fullPath := filepath.Join(downloadDir, fileName)
		if FileSizeMatches(fullPath, fileSize, verbose) {
			log.Printf("file %v already downloaded skipping", fullPath)
			continue
		}
		if verbose {
			log.Printf("file %v has %v parts and a total file size of %v", fileName, parts, fileSize)
		}
		log.Printf("downloading %v", fullPath)

		var fileNames []string
		var errMutex sync.Mutex
		var failedFiles []string
		segmentRequestInformation := calculateExecutionCalls(parts)
		for _, segment := range segmentRequestInformation {
			start := segment.StartSegment
			end := segment.EndSegment

			urls, err := client.GetDownloadUrlsForFile(
				p,
				fileId,
				keyCode,
				start,
				end,
			)
			if err != nil {
				log.Printf("unable to download file '%v' due to error '%v' while attemping to get the download url, skipping file", fileName, err)
				continue //in threaded version is return
			}
			//var m sync.Mutex
			for i := range urls {
				index := i

				//spawning yet another go routine so adding this to the wait group
				url := urls[index]
				downloadUrl := url.Url
				filePart := url.Part
				// we add the encrypted value here to make it obvious on reading the directory what step in the download process it is at
				tmpName := fmt.Sprintf("%v.%v.encrypted", fileName, filePart)
				downloadLoc := filepath.Join(downloadDir, tmpName)
				err = d.DownloadFile(downloadLoc, downloadUrl)
				if err != nil {
					log.Printf("unable to download file %v due to error '%v'", downloadLoc, err)
					continue //return
				}
				newFileName, err := DecryptPart(downloadLoc, p.ServerSecret, keyCode)
				//m.Lock()
				fileNames = append(fileNames, newFileName)
				//m.Unlock()
				if err != nil {
					errMutex.Lock()
					failedFiles = append(failedFiles, newFileName)
					errMutex.Unlock()
					log.Printf("unable to decrypt file %v due to error '%v'", downloadLoc, err)

					continue //return
				}
				if verbose {
					log.Printf("file '%v' is decrypted", newFileName)
				}
			}
			if len(failedFiles) > 0 {
				log.Printf("there were %v failed files therefore not going to bother combining parts for file '%v'", len(failedFiles), fileName)
				continue // in threaded version is return
			}
		}
		newFile, err := CombineFiles(fileNames, verbose)
		if err != nil {
			log.Printf("unable to combine downloaded parts for fileName '%v' due to error '%v'", fileName, err)
		} else {
			log.Printf("file '%v is complete", newFile)
		}
		if !FileSizeMatches(fullPath, fileSize, verbose) {
			log.Printf("ERROR: file %v failed verification, removing", fullPath)
			if err := os.Remove(newFile); err != nil {
				log.Printf("WARN: unexpected failure removing file %v due to error %v", newFile, err)
			}
		}

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
