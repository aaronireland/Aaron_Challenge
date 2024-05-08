# SED Challenge
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Coverage](https://github.com/aaronireland/Aaron_Challenge/wiki/coverage.svg)](https://raw.githack.com/wiki/aaronireland/Aaron_Challenge/coverage.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/aaronireland/Aaron_Challenge?)](https://goreportcard.com/report/github.com/aaronireland/Aaron_Challenge)

This repository contains two demo projects:
- The `static-site` project deploys a secure static web application to AWS.
- `cc-validator` is a Go project which creates a CLI tool to validate batches of credit card numbers

## static-site

Uses Terraform to deploy an AWS CloudFront distribution that serves a single html file from an S3 bucket.

#### Demo

1. HTTPS: <a href="https://sed-challenge.aaronireland.net" target="_blank">https://sed-challenge.aaronireland.net</a>
2. HTTP(redirect): <a href="http://sed-challenge.aaronireland.net" target="_blank">http://sed-challenge.aaronireland.net</a>

#### Requirements

1. [Terraform](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli)
2. [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-quickstart.html)
3. [Mage](https://go.dev/doc/install)

#### Set up

TODO

## cc-validator

Solution for the [HackerRank - Validating Credit Card Numbers challenge](https://www.hackerrank.com/challenges/validating-credit-card-number/problem)

Uses [Cobra](https://github.com/spf13/cobra) to generate a basic CLI binary with a single command which validates batches of credit card numbers for a fictitious bank, ABCD Bank.
It accepts a batch as text data from a file path flag or piped directly from the shell and outputs a sequence of `VALID` of `INVALID` for each card line-delimited. The validator checks
used are fairly extensible but mostly rely on regex patterns.

#### Requirements

1. [Go](https://go.dev/doc/install)
2. [Mage](https://go.dev/doc/install)


```shell
go version
```

```shell
mage -version
```


### Build

This project was built with Go version 1.22 on a darwin/amd64 machine. The binary builds with Mage but is simple enough to 
build directly with `go` commands and/or can be run directly with `go run`

```shell
mage ccvalidator:build
```

The binary will be build to the `bin/` directory at the project root and can be run like this:

```shell
./bin/cc-validator
```

It should output an error indicating no input text file was provided along with usage info for the commmand

```
Error: the file provided does not exist
Usage:
  validate [flags]

Flags:
  -f, --file string   path to a file containing the data to validate
  -h, --help          help for validate
  -v, --verbose       output full validation result with input text and the error messages
```

### Usage

Text data can be piped directly like this:

```shell
echo -e "1\n4141-1234-1211-0001" | ./bin/cc-validator
```

or a file can be passed with the `-f` flag:

```shell
./bin/cc-validator -f ./testdata/hackerrank.txt
```

If run with the `-v` flag, the card numbers and validation errors will be shown.


### Testing

To run tests with mage: 

```shell
mage ccvalidator:test
```

or with the html coverage:

```shell
mage ccvalidator:coverage
```
