#!/bin/sh

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -v cmd/clifront.go

scp -i ~/key-store/id_ed25519 clifront pojos@10.19.143.115:
