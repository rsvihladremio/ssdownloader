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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTicketReport(t *testing.T) {

	assert.Equal(t, `
================================
= ssdownloader summary         =
================================
= total files       : 100
= total succeeded   : 20
= total skipped     : 50
= total failed      : 30
= total bytes       : 1000 bytes
= max bytes         : 200 bytes
================================`, Report(100, 50, 30, 1000, 200))
}
