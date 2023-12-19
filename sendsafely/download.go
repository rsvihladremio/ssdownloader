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

// sendsafely package decrypts files, combines file parts into whole files, and handles api access to the sendsafely rest api
package sendsafely

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/rsvihladremio/ssdownloader/cmd/config"
	"github.com/rsvihladremio/ssdownloader/downloader"
	"github.com/rsvihladremio/ssdownloader/futils"
	"github.com/rsvihladremio/ssdownloader/reporting"
)

type PartRequests struct {
	StartSegment int
	EndSegment   int
}

// FileSizeCheck this is to help us be aware of future patterns in the sendsafely api, often it will return a different size
// than it actually encountered in reality, these files are so far always usable so we are just going to log the difference
func FileSizeCheck(fileName string, fileSize int64) error {
	fi, err := os.Stat(fileName)
	if err != nil {
		return fmt.Errorf("unable to check file size %v", fileName)
	}

	if fi.Size() != fileSize {
		slog.Warn("file size in sendsafely and reality do not match, this is often not a real issue", "current_file_size_bytes", fi.Size(), "file_to_download_bytes", fileSize)
	}
	return nil
}

func DownloadFilesFromPackage(d *downloader.GenericDownloader, packageID, keyCode string, c config.Config, subDirToDownload string, verbose bool) (outDir string, invalidFiles []string, err error) {
	client := NewClient(c.SsAPIKey, c.SsAPISecret, verbose)
	p, err := client.RetrievePackageByID(packageID)
	if err != nil {
		return "", []string{}, err
	}
	//making top level download directory if it does not exist
	_, err = os.Stat(c.DownloadDir)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		configDir := filepath.Dir(c.DownloadDir)
		slog.Debug("making download dir since it is not present", "dir", c.DownloadDir)
		err = os.Mkdir(c.DownloadDir, 0700)
		if err != nil {
			return "", []string{}, fmt.Errorf("unable to create download dir '%v' due to error '%v'", configDir, err)
		}
	}

	shortPackageID := p.PackageID
	//Add timestamp for sorting
	fullPackageName := fmt.Sprintf("%v_%v", p.PackageTimestamp.Format("20060102T150405"), shortPackageID)
	//make config directory for this package code if it does not exist
	outDir = filepath.Join(c.DownloadDir, subDirToDownload, fullPackageName)

	_, err = os.Stat(outDir)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		configDir := filepath.Dir(outDir)
		slog.Debug("making package directory since it does not exist", "dir", outDir)
		err = os.MkdirAll(outDir, 0700)
		if err != nil {
			return "", []string{}, fmt.Errorf("unable to create download dir '%v' due to error '%v'", configDir, err)
		}
	}

	for _, f := range p.Files {

		reporting.AddFile()
		fileName := f.FileName
		parts := f.Parts
		fileSize := f.FileSize
		fileID := f.FileID
		fullPath := filepath.Join(outDir, fileName)
		exists, err := futils.FileExists(fullPath)
		if err != nil {
			reporting.AddFailed()
			slog.Error("unable to check if file exists. Skipping file to prevent overwriting existing one.", "file_name", fullPath, "error_msg", err)
			continue
		}
		if exists {
			if err := FileSizeCheck(fullPath, fileSize); err != nil {
				reporting.AddFailed()
				invalidFiles = append(invalidFiles, fullPath)
				slog.Error("unable to validate new file", "file_name", fullPath, "err", err)
				continue
			}
			reporting.AddSkip()
			slog.Debug("file already downloaded skipping", "file_name", fullPath)
			continue
		}
		slog.Debug("reported file information from sendsafely", "file_name", fileName, "number_parts", parts, "file_size_in_bytes", fileSize)
		fmt.Print(".")
		slog.Debug("downloading", "file_name", fullPath)

		var fileNames []string
		var failedFiles []string
		segmentRequestInformation := calculateExecutionCalls(parts)
		for _, s := range segmentRequestInformation {
			segment := s
			start := segment.StartSegment
			end := segment.EndSegment

			urls, err := client.GetDownloadUrlsForFile(
				p,
				fileID,
				keyCode,
				start,
				end,
			)
			if err != nil {
				reporting.AddFailed()
				slog.Error("while attempting to get the download url we encountered an error, skipping file", "file_name", fileName, "error_msg", err)
				continue
			}
			for i := range urls {
				index := i

				url := urls[index]
				downloadURL := url.URL
				filePart := url.Part
				// we add the encrypted value here to make it obvious on reading the directory what step in the download process it is at
				tmpName := fmt.Sprintf("%v.%v.encrypted", fileName, filePart)
				downloadLoc := filepath.Join(outDir, tmpName)
				err = d.DownloadFile(downloadLoc, downloadURL)
				if err != nil {
					reporting.AddFailed()
					slog.Debug("unable to download file", "file_name", downloadLoc, "error_msg", err)
					continue
				}
				newFileName, err := DecryptPart(downloadLoc, p.ServerSecret, keyCode)
				fileNames = append(fileNames, newFileName)
				if err != nil {
					reporting.AddFailed()
					failedFiles = append(failedFiles, newFileName)
					slog.Debug("unable to decrypt file", "file_name", downloadLoc, "error_msg", err)
					continue
				}
				slog.Debug("file decrypted", "file_name", newFileName)
			}
			if len(failedFiles) > 0 {
				reporting.AddFailed()
				slog.Error("there were failed downloads of parts of the file skipping", "failed_file_parts_count", len(failedFiles), "file_name", fileName)
				continue
			}
		}
		written, newFile, err := CombineFiles(fileNames, verbose)
		if err != nil {
			reporting.AddFailed()
			return "", invalidFiles, fmt.Errorf("unable to combine downloaded parts for file %v: %v", fileName, err)
		}
		if err := FileSizeCheck(fullPath, fileSize); err != nil {
			reporting.AddFailed()
			return "", invalidFiles, fmt.Errorf("unable to validate new file: %v: %v", fileName, err)
		}
		fmt.Print(".")
		slog.Debug("file is complete", "file_name", newFile, "file_size", Human(written), "file_size_in_bytes", written)
		reporting.AddBytes(fileSize)

	}
	return outDir, invalidFiles, nil
}

func Human(bytes int64) string {
	if bytes > 1024*1024*1024 {
		return fmt.Sprintf("%.2f gb", float64(bytes)/(1024.0*1024.0*1024.0))
	} else if bytes > 1024*1024 {
		return fmt.Sprintf("%.2f mb", float64(bytes)/(1024.0*1024.0))
	} else if bytes > 1024 {
		return fmt.Sprintf("%.2f kb", float64(bytes)/1024.0)
	}
	return fmt.Sprintf("%v bytes", bytes)
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
