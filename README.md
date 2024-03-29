![build](https://github.com/rsvihladremio/ssdownloader/actions/workflows/checkin.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/rsvihladremio/ssdownloader)](https://goreportcard.com/report/github.com/rsvihladremio/ssdownloader)
<a href='https://coveralls.io/github/rsvihladremio/ssdownloader?branch=main'><img src='https://coveralls.io/repos/github/rsvihladremio/ssdownloader/badge.svg?branch=main&service=test' alt='Coverage Status' /></a>


# ssdownloader

ssdownloader has implemented the zendesk and sendsafely rest APIs 
to provide support for search for sendsafely links in tickets and downloading
all files found.

Other features include: 

 * support for downloading zendesk attachments
 * ability to download sendsafely links with no zendesk information
 * storage of api credentials
 * download of all content into well known directory structures
 * support for verbose logging
 * multithreaded with support for adjusting the number of threads for you performance needs

## License

[Apache 2.0 License](https://www.apache.org/licenses/LICENSE-2.0.html)

## Quickstart

First there are binaries [here](https://github.com/rsvihladremio/ssdownloader/releases) so you can always get those
if any of the following methods do not work for you.

On mac I suggest homebrew:

```sh
brew tap rsvihladremio/tools
brew install ssdownloader
```

```sh
ssdownloader init
ssdownloader ticket 9999 
ssdownloader main-a3caa66-darwin-amd64
```
and then you should see something like the following

```sh
2022/06/23 10:35:35 making dir /Users/foo/.sendsafely
2022/06/23 10:35:35 making dir /Users/foo/.sendsafely/tickets/9999
2022/06/23 10:35:35 downloading fqfqsdfqds-fdfsd-fdqfd-fqdfq-fdqfdsffqdfq - works.zip
2022/06/23 10:35:35 downloading cbabd5ba-fdqf-fdqdf-qfd-fqdfsdfqs - problem.zip
2022/06/23 10:35:35 downloading server.log
```

On Linux or WSL do the following:

```sh
curl -sSfL https://raw.githubusercontent.com/rsvihladremio/ssdownloader/main/script/install | sh 
```

then you can run

```sh
ssdownloader init
ssdownloader ticket 9999 
ssdownloader main-a3caa66-darwin-amd64
```

and then you should see something like the following

```sh
2022/06/23 10:35:35 making dir /Users/foo/.sendsafely
2022/06/23 10:35:35 making dir /Users/foo/.sendsafely/tickets/9999
2022/06/23 10:35:35 downloading fqfqsdfqds-fdfsd-fdqfd-fqdfq-fdqfdsffqdfq - works.zip
2022/06/23 10:35:35 downloading cbabd5ba-fdqf-fdqdf-qfd-fqdfsdfqs - problem.zip
2022/06/23 10:35:35 downloading server.log
```

On Windows do the following:

```sh
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser # Optional: Needed to run a remote script the first time
irm https://raw.githubusercontent.com/rsvihladremio/ssdownloader/main/script/install.ps1  | iex 
```

then you can run

```sh
.\ssdownloader.exe init
.\ssdownloader.exe ticket 9999 
.\ssdownloader.exe main-a3caa66-darwin-amd64
```
and you should see something like the following output

```sh
2022/06/23 10:35:35 making dir /Users/foo/.sendsafely
2022/06/23 10:35:35 making dir /Users/foo/.sendsafely/tickets/9999
2022/06/23 10:35:35 downloading fqfqsdfqds-fdfsd-fdqfd-fqdfq-fdqfdsffqdfq - works.zip
2022/06/23 10:35:35 downloading cbabd5ba-fdqf-fdqdf-qfd-fqdfsdfqs - problem.zip
2022/06/23 10:35:35 downloading server.log
```

## Developing

On Linux, Mac, and WSL there are some shell scripts modeled off the [GitHub ones](https://github.com/github/scripts-to-rule-them-all)

to get started run

```sh
./script/bootstrap
```

after a pull it is a good idea to run

```sh
./script/update
```

tests

```sh
./script/test
```

before checkin run

```sh
./script/cibuild
```

to cut a release do the following

```sh
#dont forget to update changelog.md with the release notes
git tag v0.1.1
./script/release v0.1.1
gh repo view -w
# review the draft and when done set it to publish
```
### Windows
Similarly on Windows there are powershell scripts of the same design

to get started run

```powershell
.\script\bootstrap.ps1
```

after a pull it is a good idea to run

```powershell
.\script\update.ps1
```

tests

```powershell
.\script\test.ps1
```

before checkin run

```powershell
.\script\cibuild.ps1
```

to cut a release do the following

```powershell
#dont forget to update changelog.md with the release notes
git tag v0.1.1
.\script\release.ps1 v0.1.1
gh repo view -w
# review the draft and when done set it to publish
```

## FAQ

### why go ?

Ease of deployment, easy to learn development and fast enough. Some people will say why not Python? why not Java? why not Rust?  In one way or another they will lack, I love all of them and use them regularly for other tasks, but this is neither performance sensitive, nor server based (so deployment ease matters a lot) and those that need to maintain it need a language they can spin up easily.
