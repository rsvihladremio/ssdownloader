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
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerbosityLevel(t *testing.T) {
	defer func() {
		Verbose = false
	}()
	Verbose = true
	SetVerbosity()
	assert.Equal(t, slog.LevelDebug, programLevel.Level())
}

func TestVerbosityLevelDefault(t *testing.T) {
	defer func() {
		Verbose = false
	}()
	Verbose = false
	SetVerbosity()
	assert.Equal(t, slog.LevelInfo, programLevel.Level())
}
