![build](https://github.com/rsvihladremio/ssdownloader/actions/workflows/checkin.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/rsvihladremio/ssdownloader)](https://goreportcard.com/report/github.com/rsvihladremio/ssdownloader)

# ssdownloader

SendSafely downloader that integrates with Zendesk

## License

[Apache 2.0 License](https://www.apache.org/licenses/LICENSE-2.0.html)

## Install

On Linux, Mac or WSL do the following:

```sh
curl -sSfL https://raw.githubusercontent.com/rsvihladremio/ssdownloader/main/script/install | sh 
```

## FAQ

### why go ?

Ease of deployment, easy to learn development and fast enough. Some people will say why not Python? why not Java? why not Rust?  In one way or another they will lack, I love all of them and use them regularly for other tasks, but this is neither performance sensitive, nor server based (so deployment ease matters a lot) and those that need to maintain it need a language they can spin up easily.