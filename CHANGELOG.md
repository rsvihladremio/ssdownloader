# Changelog

## [0.4.12] - 2025-03-13

## Fixed

* golang.org/x/net from 0.33.0 to 0.36.0 security fix

## [0.4.11] - 2024-04-23

### Fixed

* golang.org/x/net from 0.19.0 to 0.23.0 (#44)

## [0.4.10] - 2023-12-20

### Added 

- more logging and more simple skip list check

## [0.4.9] - 2023-12-19

### Added

- support for filtering out files over a certain size
- API support for downloading only a single file by file ID

## [0.4.8] - 2023-12-19

### Fixed
- security fixes
- some improved copying of variables

## [0.4.7] - 2023-12-12

### Fixed
- combining 0 files no longer leads to a panic

### Changed
- we now just log when files do not match file size between sendsafely and the file system, this was nearly always a valid file

## [0.4.6] - 2023-10-12
### Added
- added new report and made the logs less chatty by default
### Fixed
- version bump for CVEs

### Changed
- removed posting of changelog to release notes

## [0.4.5] - 2023-05-16

### Fixed 

- version bump for CVEs

## [0.4.4] - 2023-02-09

### Changed

- handling of URLs that have lowercase parameters

## [0.4.3] - 2023-02-09
### Fixed
- tests fail on some platforms due to inconsistent behavior of restly, worked around by checking for status code of response in zendesk api
- no longer parsing threadid as all urls do not have it, and we do not need it

## [0.4.2] - 2022-09-30
### Added
- support for paging to the zedesk API, tickets longer than 100 comments will now download all attachments
### Fixed
- tickets with sendsafely links in comments that did not have a package ID threw a fatal error. Now will log and continue

## [0.4.1] - 2022-07-19
### Added
- files that fail validation are reported at the end of the output
### Changed
- if the file exists at all but is not valid, this is logged instead of clobbered by a new download
- files that are invalid are not deleted after download

## [0.4.0] - 2022-07-12
### Added
- now download the ticket text that is associated with a given package
### Changed
- directory structure of downloads now includes a timestamp in front of the package name, this makes it easy to scan to the latest package
- now have bytes written in logs
### Fixed
- better log message when package not found

## [0.3.3] - 2022-07-06
### Added
- added a decent number of tests
- support for sendsafely links wrapped by gmail url's wrappers
### Changed
- updated changelog format to one based on https://keepachangelog.com/en/1.0.0/
### Fixed
- fixed some minor reporting errors
### Removed
- CPU and Memory profiling options as they were rarely used

## [0.3.2] - 2022-07-04
### Fixed
- forgot to add wait group counter for attachment download leading to race condition

## [0.3.1] - 2022-07-01
### Fixed
- date formatting bug in sendsafely api left requests expiring

## [0.3.0] - 2022-06-29
### Added
- now downloads attachments by default
- can skip attachmments with --sendsafely-only flag
### Changed
- lots more testing of more parts of the code base, this is now beta

## [0.2.5] - 2022-06-27
### Added
- verifies file size after download matches
### Changed
- will not redownload a file already downloaded now

## [0.2.4] - 2022-06-27
### Changed
- package id is shortened to match what zendesk shows
### Fixed
- a trivial sort bug was breaking larger files

## [0.2.3] - 2022-06-24
### Added
- more testing!!!
### Changed
- use package name for prefix
### Fixed
- was unable to combine especially large files due to error in logic

## [0.2.2] - 2022-06-23
### Changed
 updated docs and help to show subdmain
### Fixed
- had init check for prompt backwords now works

## [0.2.1] - 2022-06-22
### Changed
- require init to have certain parameters
### Fixed
- prompts work for ssdownloader init, one does not have to pass all the flags

## [0.2.0] - 2022-06-21
### Added
- Support for zendesk api via api key
- more automated testing still alpha though

## [0.1.0] - 2022-06-17
### Added
- can download via links, but does not support the ticket functionality yet. Use at your own risk!
  
[0.4.12]: https://github.com/rsvihladremio/ssdownloader/compare/v0.4.11...v0.4.12
[0.4.11]: https://github.com/rsvihladremio/ssdownloader/compare/v0.4.10...v0.4.11
[0.4.10]: https://github.com/rsvihladremio/ssdownloader/compare/v0.4.9...v0.4.10
[0.4.9]: https://github.com/rsvihladremio/ssdownloader/compare/v0.4.8...v0.4.9
[0.4.8]: https://github.com/rsvihladremio/ssdownloader/compare/v0.4.7...v0.4.8
[0.4.7]: https://github.com/rsvihladremio/ssdownloader/compare/v0.4.6...v0.4.7
[0.4.6]: https://github.com/rsvihladremio/ssdownloader/compare/v0.4.5...v0.4.6
[0.4.5]: https://github.com/rsvihladremio/ssdownloader/compare/v0.4.4...v0.4.5
[0.4.4]: https://github.com/rsvihladremio/ssdownloader/compare/v0.4.3...v0.4.4
[0.4.3]: https://github.com/rsvihladremio/ssdownloader/compare/v0.4.2...v0.4.3
[0.4.2]: https://github.com/rsvihladremio/ssdownloader/compare/v0.4.1...v0.4.2
[0.4.1]: https://github.com/rsvihladremio/ssdownloader/compare/v0.4.0...v0.4.1
[0.4.0]: https://github.com/rsvihladremio/ssdownloader/compare/v0.3.3...v0.4.0
[0.3.3]: https://github.com/rsvihladremio/ssdownloader/compare/v0.3.2...v0.3.3
[0.3.2]: https://github.com/rsvihladremio/ssdownloader/compare/v0.3.1...v0.3.2
[0.3.1]: https://github.com/rsvihladremio/ssdownloader/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/rsvihladremio/ssdownloader/compare/v0.2.5...v0.3.0
[0.2.5]: https://github.com/rsvihladremio/ssdownloader/compare/v0.2.4...v0.2.5
[0.2.4]: https://github.com/rsvihladremio/ssdownloader/compare/v0.2.3...v0.2.4
[0.2.3]: https://github.com/rsvihladremio/ssdownloader/compare/v0.2.2...v0.2.3
[0.2.2]: https://github.com/rsvihladremio/ssdownloader/compare/v0.2.1...v0.2.2
[0.2.1]: https://github.com/rsvihladremio/ssdownloader/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/rsvihladremio/ssdownloader/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/rsvihladremio/ssdownloader/releases/tag/v0.1.0
