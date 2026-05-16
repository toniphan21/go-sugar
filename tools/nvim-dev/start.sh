#!/bin/sh

set -e

export GO_SUGAR_BIN=$(pwd)/bin/go-sugar
export GO_SUGAR_LOG=$(pwd)/bin/go-sugar.lsp.log

go build -o bin/go-sugar ./cmd/go-sugar

XDG_CONFIG_HOME=$(pwd)/tools NVIM_APPNAME=nvim-dev nvim "${1:-.}"

