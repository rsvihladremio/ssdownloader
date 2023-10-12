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

import (
	"sync"
	"testing"


	"github.com/stretchr/testify/assert"
)

func TestNonThreaded(t *testing.T) {
	defer func() {
		totalFiles = 0
		totalFailed = 0
		totalBytes = 0
		totalSkipped = 0
		maxFileSize = 0
	}()
	AddFile()
	assert.Equal(t, 1, GetTotalFiles())
	AddSkip()
	assert.Equal(t, 1, GetTotalSkipped())
	AddFailed()
	assert.Equal(t, 1, GetTotalFailed())
	AddBytes(99)
	AddBytes(101)
	assert.Equal(t, int64(200), GetTotalBytes())
	assert.Equal(t, int64(101), GetMaxFileSizeBytes())
}

func TestThreadSafeCounts(t *testing.T) {
	defer func() {
		totalFiles = 0
		totalFailed = 0
		totalBytes = 0
		totalSkipped = 0
		maxFileSize = 0
	}()
	total := 10000
	var wg sync.WaitGroup
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func() {
			AddFile()
			wg.Done()
		}()
	}
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func() {
			AddSkip()
			wg.Done()
		}()
	}
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func() {
			AddFailed()
			wg.Done()
		}()
	}
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func() {
			AddBytes(10)
			wg.Done()
		}()
	}
	wg.Wait()
	assert.Equal(t, total, GetTotalFiles())
	assert.Equal(t, total, GetTotalSkipped())
	assert.Equal(t, total, GetTotalFailed())
	assert.Equal(t, int64(10*total), GetTotalBytes())
	assert.Equal(t, int64(10), GetMaxFileSizeBytes())
}
