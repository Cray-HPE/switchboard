#!/bin/sh

export GO111MODULE=on
go get
go test -v ./... -cover
