#!/bin/bash

msg() {
	echo 2>&1 "$0": $@
}

hash gosimple >/dev/null && gosimple ./applets
hash golint >/dev/null && golint ./applets
hash staticcheck >/dev/null && staticcheck ./applets

gofmt -s -w ./applets
go fix ./applets/...
go vet -vettool="$(which shadow)" ./applets/...
go test ./applets/...
go install -v ./applets/...

hash gosimple >/dev/null && gosimple ./conbox
hash golint >/dev/null && golint ./conbox
hash staticcheck >/dev/null && staticcheck ./conbox

gofmt -s -w ./conbox
go fix ./conbox
go vet -vettool="$(which shadow)" ./conbox
go test ./conbox
go install -v ./conbox

