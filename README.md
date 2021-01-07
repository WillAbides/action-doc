# action-doc

[![godoc](https://godoc.org/github.com/willabides/action-doc?status.svg)](https://godoc.org/github.com/willabides/action-doc)
[![ci](https://github.com/WillAbides/action-doc/workflows/ci/badge.svg?branch=main&event=push)](https://github.com/WillAbides/action-doc/actions?query=workflow%3Aci+branch%3Amaster+event%3Apush)

action-doc takes your action's action.yml as input and outputs some markdown suitable for a README.md

## Usage

```shell
cat action.yml | action-doc > README.md
```

## Install

### go get

`go get -u github.com/willabides/action-doc/cmd/action-doc`

### bindown

Add a [bindown](https://github.com/willabides/bindown) dependency:

``` shell
$ bindown template-source add action-doc https://raw.githubusercontent.com/WillAbides/action-doc/main/bindown.yml
$ bindown dependency add action-doc action-doc#action-doc
Please enter a value for required variable "version":	<latest version>
```
