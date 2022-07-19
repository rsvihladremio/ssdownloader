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

//cmd package contains all the command line flag configuration
package cmd

import (
	"fmt"
	"strings"
)

func InvalidFilesReport(invalidFiles []string) string {
	str := ""
	if len(invalidFiles) > 0 {
		str := `
the following files failed validation
-------------------------------------
`
		rows := []string{}
		for _, f := range invalidFiles {
			rows = append(rows, fmt.Sprintf("* %v\n", f))
		}
		return str + strings.Join(rows, "")
	}
	return str
}
