[![license](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/udhos/conbox/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/udhos/conbox)](https://goreportcard.com/report/github.com/udhos/conbox)

# conbox
[conbox](https://github.com/udhos/conbox) is a Go implementation of unix-like utilities as single static executable intended for small container images.

* [Install](#install)
* [Usage](#usage)
  * [Available applets](#available-applets)
  * [Basename usage](#basename-usage)
  * [Arg\-1 usage](#arg-1-usage)
* [Docker](#docker)
  * [Docker recipes](#docker-recipes)
* [Related work](#related-work)

Created by [gh-md-toc](https://github.com/ekalinin/github-markdown-toc.go)

# Install

```bash
git clone https://github.com/udhos/conbox
cd conbox
go install ./conbox
```

# Usage

## Available applets

List available applets:

```bash
$ conbox
conbox: version 0.0 runtime go1.12.4 GOMAXPROC=12 OS=linux ARCH=amd64
usage: conbox APPLET [ARG]... : run APPLET
       conbox -h              : show command-line help
       conbox -l              : list applets

conbox: registered applets:
cat echo ls rm 
```

See all implemented applets here:

https://github.com/udhos/conbox/tree/master/applets

## Basename usage

Create a symbolic link to 'conbox':

```bash
ln -s ~/go/bin/conbox ~/bin/cat
~/bin/cat /etc/passwd
```

## Arg-1 usage

Pass applet name as first argument to 'conbox':

```bash
conbox cat /etc/passwd
```

# Docker

Get 'conbox' as docker image `udhos/conbox:latest` from:

https://hub.docker.com/r/udhos/conbox

## Docker recipes

Build docker image:

```bash
./docker/build.sh
```

Tag image:

```bash
docker tag udhos/conbox udhos/conbox:latest
```

Push image:

```bash
docker login
docker push udhos/conbox:latest
```

Run image:

```bash
docker run --rm udhos/conbox:latest cat /etc/passwd
```

# Related work

## Go Projects

Unfortunately these projects seem inactive:

- https://github.com/surma/gobox
- https://github.com/laher/someutils
- https://github.com/shirou/toybox

## Non-Go projects

- https://www.busybox.net/
- http://landley.net/toybox/
