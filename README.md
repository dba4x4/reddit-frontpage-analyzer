# Reddit Frontpage Analyzer (go)
[![Travis](https://travis-ci.org/swordbeta/reddit-frontpage-analyzer-go.svg?branch=master)](https://travis-ci.org/swordbeta/reddit-frontpage-analyzer-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/swordbeta/reddit-frontpage-analyzer-go)](https://goreportcard.com/report/github.com/swordbeta/reddit-frontpage-analyzer-go)
[![codecov](https://codecov.io/gh/swordbeta/reddit-frontpage-analyzer-go/branch/master/graph/badge.svg)](https://codecov.io/gh/swordbeta/reddit-frontpage-analyzer-go)
[![codebeat badge](https://codebeat.co/badges/e587155a-c69e-406a-a1b5-c219513b2400)](https://codebeat.co/projects/github-com-swordbeta-reddit-frontpage-analyzer-go)

This go application does the following:

- Fetches the top 25 posts from [/r/all][1]
- Stores all data in PostgreSQL
- Analyzes (tags) images with the [Microsoft Computer Vision API][3]

For the NodeJS serverless version see [here][2].

### Installation / Running

Clone the repository
```
λ git clone git@github.com:swordbeta/reddit-frontpage-analyzer-go.git && cd reddit-frontpage-analyzer-go
```

Copy the config file and edit
```
cp config.yaml.default config.yaml && vim config.yaml
```

Build docker image
```
λ docker build -t reddit-frontpage-analyzer-go .
```

Run docker container
```
λ docker run --rm --name reddit-frontpage-analyzer-go reddit-frontpage-analyzer-go
```

### Roadmap

- [X] Fetch reddit frontpage
- [X] Save unique posts to PostgreSQL
- [X] Tag images with [Microsoft Computer Vision API][3]
- [X] Add instructions for running in README
- [X] Gracefully exit current run when hitting rate limits
- [X] Add tests
- [X] Add Travis CI support
- [X] Add code coverage and other badges

[1]: https://reddit.com/r/all
[2]: https://github.com/swordbeta/reddit-frontpage-analyzer-nodejs
[3]: https://www.microsoft.com/cognitive-services/en-us/computer-vision-api
