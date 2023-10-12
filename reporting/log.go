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

package reporting

import "sync"

var totalFiles int
var totalFilesLock sync.Mutex
var totalFailed int
var totalFailedLock sync.Mutex
var totalSkipped int
var totalSkippedLock sync.Mutex
var totalBytes int64
var totalBytesLock sync.Mutex
var maxFileSize int64

func AddFile() {
	totalFilesLock.Lock()
	totalFiles++
	totalFilesLock.Unlock()
}

func GetTotalFiles() int {
	totalFilesLock.Lock()
	defer totalFilesLock.Unlock()
	return totalFiles
}

func AddSkip() {
	totalSkippedLock.Lock()
	totalSkipped++
	totalSkippedLock.Unlock()
}

func GetTotalSkipped() int {
	totalSkippedLock.Lock()
	defer totalSkippedLock.Unlock()
	return totalSkipped
}

func AddFailed() {
	totalFailedLock.Lock()
	totalFailed++
	totalFailedLock.Unlock()
}

func GetTotalFailed() int {
	totalFailedLock.Lock()
	defer totalFailedLock.Unlock()
	return totalFailed
}

func AddBytes(i int64) {
	totalBytesLock.Lock()
	if i > maxFileSize {
		maxFileSize = i
	}
	totalBytes += i
	totalBytesLock.Unlock()
}

func GetTotalBytes() int64 {
	totalBytesLock.Lock()
	defer totalBytesLock.Unlock()
	return totalBytes
}

func GetMaxFileSizeBytes() int64 {
	totalBytesLock.Lock()
	defer totalBytesLock.Unlock()
	return maxFileSize
}
