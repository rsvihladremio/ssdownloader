![build](https://github.com/rsvihladremio/ssdownloader/actions/workflows/checkin.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/rsvihladremio/ssdownloader)](https://goreportcard.com/report/github.com/rsvihladremio/ssdownloader)
[![Coverage Status](https://coveralls.io/repos/github/rsvihladremio/ssdownloader/badge.svg?branch=main&service=github)](https://coveralls.io/github/rsvihladremio/ssdownloader?branch=main&service=github)

# ssdownloader

SendSafely downloader that integrates with Zendesk

## License

[Apache 2.0 License](https://www.apache.org/licenses/LICENSE-2.0.html)

## Quickstart

On Linux, Mac or WSL do the following:

```sh
$ curl -sSfL https://raw.githubusercontent.com/rsvihladremio/ssdownloader/main/script/install | sh 
$ ssdownloader init
$ ssdownloader ticket 9999 
$ ssdownloader main-a3caa66-darwin-amd64

2022/06/23 10:35:35 making dir /Users/foo/.sendsafely
2022/06/23 10:35:35 making dir /Users/foo/.sendsafely/tickets/9999
2022/06/23 10:35:35 downloading fqfqsdfqds-fdfsd-fdqfd-fqdfq-fdqfdsffqdfq - works.zip
2022/06/23 10:35:35 downloading cbabd5ba-fdqf-fdqdf-qfd-fqdfsdfqs - problem.zip
2022/06/23 10:35:35 downloading server.log
```


## FAQ

### why go ?

Ease of deployment, easy to learn development and fast enough. Some people will say why not Python? why not Java? why not Rust?  In one way or another they will lack, I love all of them and use them regularly for other tasks, but this is neither performance sensitive, nor server based (so deployment ease matters a lot) and those that need to maintain it need a language they can spin up easily.
