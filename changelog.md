## [v0.3.3] - 2022-07-06
### Added
- added a decent number of tests
- support for sendsafely links wrapped by gmail url's wrappers
### Changed
- updated changelog format to one based on https://keepachangelog.com/en/1.0.0/
### Fixed
- fixed some minor reporting errors
### Removed
- CPU and Memory profiling options as they were rarely used

## [v0.3.3] - 2022-07-04
### Fixed
- forgot to add wait group counter for attachment download leading to race condition

## [v0.3.1] - 2022-07-01
### Fixed
- date formatting bug in sendsafely api left requests expiring

## [v0.3.0] - 2022-06-29
### Added
- now downloads attachments by default
- can skip attachmments with --sendsafely-only flag
### Changed
- lots more testing of more parts of the code base, this is now beta

## [v0.2.5] - 2022-06-27
### Added
- verifies file size after download matches
### Changed
- will not redownload a file already downloaded now

## [v0.2.4] - 2022-06-27
### Changed
- package id is shortened to match what zendesk shows
### Fixed
- a trivial sort bug was breaking larger files

## [v0.2.3] - 2022-06-24
### Added
- more testing!!!
### Changed
- use package name for prefix
### Fixed
- was unable to combine especially large files due to error in logic

## [v0.2.2] - 2022-06-23
### Changed
 updated docs and help to show subdmain
### Fixed
- had init check for prompt backwords now works

## [v0.2.1] - 2022-06-22
### Changed
- require init to have certain parameters
### Fixed
- prompts work for ssdownloader init, one does not have to pass all the flags

## [v0.2.0] - 2022-06-21
### Added
- Support for zendesk api via api key
- more automated testing still alpha though

## [v0.1.0] - 2022-06-17
### Added
- can download via links, but does not support the ticket functionality yet. Use at your own risk!
