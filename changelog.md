v0.3.0
------
* now downloads attachments by default
* can skip attachmments with --sendsafely-only flag
* lots more testing of more parts of the code base, this is now beta

v0.2.5
------
* verifies file size after download matches
* will not redownload a file already downloaded now

v0.2.4
------
* package id is shortened to match what zendesk shows
* a trivial sort bug was breaking larger files

v0.2.3
------
* use package name for prefix
* was unable to combine especially large files due to error in logic
* more testing!!!

v0.2.2
------
* had init check for prompt backwords now works
* updated docs and help to show subdmain

v0.2.1
------
* require init to have certain parameters
* prompts work for ssdownloader init, one does not have to pass all the flags

v0.2.0
------
* Support for zendesk api via api key
* more automated testing still alpha though

v0.1.0
------
* does not support the ticket functionality yet
* can download via links
* use at your own risk