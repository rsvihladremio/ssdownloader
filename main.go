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

//ssdownloader has implemented the zendesk and sendsafely rest APIs
// to provide support for search for sendsafely links in tickets and downloading
// all files found.
// other features include
// * support for downloading zendesk attachments
// * ability to download sendsafely links with no zendesk information
// * storage of api credentials
// * download of all content into well known directory structures
// * support for verbose logging
// * multithreaded with support for adjusting the number of threads for you performance needs
package main

import (
	"github.com/rsvihladremio/ssdownloader/cmd"
)

func main() {
	cmd.Execute()
}
