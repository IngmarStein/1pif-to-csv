# 1pif-to-csv

## Overview

`1pif-to-csv` is a tool to convert passwords exported from 1Password into a format suitable for macOS Monterey's
built-in password manager. Although 1Password can export CSV files directly, they do not contain the 2FA secrets.

## Installation

```shell
go get github.com/IngmarStein/1pif-to-csv
```

## Usage

```shell
1pif-to-csv -input export.1pif/data.1pif -output password.csv
```