# 1pif-to-csv

## Overview

`1pif-to-csv` is a tool to convert passwords exported from 1Password into a format suitable for macOS Monterey's
built-in password manager. Although 1Password can export CSV files directly, they do not contain the 2FA secrets.

## Installation

With [Homebrew](https://brew.sh) (recommended):
```shell
brew install ingmarstein/repo/1pif-to-csv
```

With Go 1.17 and later:
```shell
go install github.com/ingmarstein/1pif-to-csv@latest
```

If you are using Go 1.16 or earlier, use this command instead:
```shell
go get github.com/ingmarstein/1pif-to-csv
```

## Usage

```shell
1pif-to-csv -input export.1pif/data.1pif -output password.csv
```