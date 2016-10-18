# Reddit Frontpage Analyzer (go)

This go application does the following:

- Fetches the top 25 posts from [/r/all][1]
- Stores all data in PostgreSQL
- Analyzes (tags) images with the [Microsoft Computer Vision API][3]

For the NodeJS serverless version see [here][2].

### Installation / Running

TODO

### Roadmap

- [X] Fetch reddit frontpage
- [X] Save unique posts to PostgreSQL
- [X] Tag images with [Microsoft Computer Vision API][3]

[1]: https://reddit.com/r/all
[2]: https://github.com/swordbeta/reddit-frontpage-analyzer-nodejs
[3]: https://www.microsoft.com/cognitive-services/en-us/computer-vision-api