# Reddit Frontpage Analyzer
[![Travis](https://travis-ci.org/swordbeta/reddit-frontpage-analyzer.svg?branch=master)](https://travis-ci.org/swordbeta/reddit-frontpage-analyzer)
[![Go Report Card](https://goreportcard.com/badge/github.com/swordbeta/reddit-frontpage-analyzer)](https://goreportcard.com/report/github.com/swordbeta/reddit-frontpage-analyzer)
[![codecov](https://codecov.io/gh/swordbeta/reddit-frontpage-analyzer/branch/master/graph/badge.svg)](https://codecov.io/gh/swordbeta/reddit-frontpage-analyzer)

This go application does the following:

- Fetches the top 25 posts from [/r/all][1]
- Analyzes (tags) images with the [Microsoft Computer Vision API][3]
- Stores all data in PostgreSQL

For the NodeJS AWS Lambda version using AWS DynamoDB see [here][2]. (Does not work as good as I hoped.)

### Installation / Running

Clone the repository
```
λ git clone git@github.com:swordbeta/reddit-frontpage-analyzer.git && cd reddit-frontpage-analyzer
```

Copy the config file and edit
```
cp config.yaml.default config.yaml && vim config.yaml
```

Build docker image
```
λ docker build -t reddit-frontpage-analyzer .
```

Run docker container
```
λ docker run --rm --name reddit-frontpage-analyzer reddit-frontpage-analyzer
```

[1]: https://reddit.com/r/all
[2]: https://github.com/swordbeta/reddit-frontpage-analyzer-nodejs
[3]: https://www.microsoft.com/cognitive-services/en-us/computer-vision-api
