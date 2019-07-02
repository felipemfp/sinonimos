# 🇧🇷 sinonimos

[![CircleCI](https://circleci.com/gh/felipemfp/sinonimos.svg?style=shield)](https://circleci.com/gh/felipemfp/sinonimos)
[![Go Report Card](https://goreportcard.com/badge/github.com/felipemfp/sinonimmos)](https://goreportcard.com/report/github.com/felipemfp/sinonimmos)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=shield)](http://godoc.org/github.com/felipemfp/sinonimmos)
[![Release](https://img.shields.io/github/release/felipemfp/sinonimmos.svg?style=shield)](https://github.com/felipemfp/sinonimmos/releases/latest)

Find synonyms (for pt-BR) without leaving your terminal.

![Preview](sinonimos-peek.gif)

## Installing

Run this command to download the latest version of `sinonimos`:

```bash
go get github.com/felipemfp/sinonimos/...
```

Now you're ready to use `sinonimos`

```bash
sinonimos --help
```

## Hacking into `sinonimos`

These instructions will get you a copy of the project up and running on your local machine for development.

### Prerequisites

What things you need to start hacking:

- [go](https://golang.org/doc/install)
- [golang/lint](https://github.com/golang/lint#installation)

### Getting started

First [fork](https://guides.github.com/activities/forking/) and clone the project to your machine:

```
git clone https://github.com/{your-username}/sinonimos.git
```

Then install the dependencies:

```
cd sinonimos
go mod download
```

Now you're ready to go.

> For example:
>
> ```bash
> go run ./cmd/sinonimos camisa
> ```
>
> It'll run the CLI and try to find synonyms for "camisa".

## Built With

- [mow.cli](https://github.com/jawher/mow.cli)
- [scrape](https://github.com/yhat/scrape)

## Contributing

Please feel free for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/felipemfp/sinonimos-cli/tags).

## Authors

- **Felipe Pontes** - _Initial work_ - [felipemfp](https://github.com/felipemfp)

See also the list of [contributors](https://github.com/felipemfp/sinonimos-cli/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
