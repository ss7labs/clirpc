#!/bin/sh

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -v cmd/clifront.go

