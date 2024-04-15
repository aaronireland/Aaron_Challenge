# SED Challenge

This repository contains two demo projects:
- The `static-site` project deploys a secure static web application to AWS.
- `cc-validator` is a Go project which creates a CLI tool to validate batches of credit card numbers

## static-site

Uses Terraform to deploy an EKS cluster and Helm to deploy an NGinx web server that servers a single html file.


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
(cd cc-validator; mage build)
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
(cd cc-validator; mage test:unit)
```

or with the html coverage:

```shell
(cd cc-validator; mage test:htmlcov)
```