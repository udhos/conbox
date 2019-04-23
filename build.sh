#!/bin/bash

msg() {
	echo 2>&1 "$0": $@
}

hash gosimple 2>/dev/null && gosimple ./applets
hash golint 2>/dev/null && golint ./applets
hash staticcheck 2>/dev/null && staticcheck ./applets

gofmt -s -w ./applets
go fix ./applets/...
go vet -vettool="$(which shadow)" ./applets/...
go test ./applets/...
go install -v ./applets/...

hash gosimple 2>/dev/null && gosimple ./conbox
hash golint 2>/dev/null && golint ./conbox
hash staticcheck 2>/dev/null && staticcheck ./conbox

gofmt -s -w ./conbox
go fix ./conbox
go vet -vettool="$(which shadow)" ./conbox
go test ./conbox
go install -v ./conbox

