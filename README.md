[![license](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/udhos/conbox/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/udhos/conbox)](https://goreportcard.com/report/github.com/udhos/conbox)

# conbox
[conbox](https://github.com/udhos/conbox) is a Go implementation of unix-like utilities as single static executable intended for small container images.

* [Install](#install)
* [Usage](#usage)
  * [Available applets](#available-applets)
  * [Basename usage](#basename-usage)
  * [Subcommand usage](#subcommand-usage)
  * [Shell usage](#shell-usage)
* [Adding new applet](#adding-new-applet)
* [Docker](#docker)
  * [Run in docker](#run-in-docker)
  * [Docker recipes](#docker-recipes)
* [Related work](#related-work)
  * [Go Projects](#go-projects)
  * [Non\-Go projects](#non-go-projects)

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
conbox: version 0.0 runtime go1.12.4 GOMAXPROC=1 OS=linux ARCH=amd64

usage: conbox APPLET [ARG]... : run APPLET
       conbox -h              : show command-line help
       conbox -l              : list applets

conbox: registered applets:
cat echo ls mkdir printenv pwd rm rmdir shell which
```

See all implemented applets here:

https://github.com/udhos/conbox/tree/master/applets

## Basename usage

You can create a symbolic link for a supported applet pointing to 'conbox':

```bash
ln -s ~/go/bin/conbox ~/bin/cat
~/bin/cat /etc/passwd
```

## Subcommand usage

Pass applet name as subcommand to 'conbox':

```bash
conbox cat /etc/passwd
```

## Shell usage

All applets are also directly available from within conbox shell:

```bash
$ conbox shell
conbox: version 0.0 runtime go1.12.4 GOMAXPROC=1 OS=linux ARCH=amd64

welcome to conbox shell.
this tiny shell is very limited in features.
however you can run external programs normally.
some hints:
       - use 'conbox' to see all applets available as shell commands.
       - use 'help' to list shell built-in commands.
       - 'exit' terminates the shell.

shell built-in commands:
builtin cd exit help

conbox shell$
```

# Adding new applet

1. Create a new package for the applet under directory 'applets'. The package must export the function Run() as show in example below.

```
$ more applets/myapp/run.go
package myapp // put applet myapp in package myapp

import (
        "fmt"

        "github.com/udhos/conbox/common"
)

// Run executes the applet.
func Run(tab map[string]common.AppletFunc, args []string) int {

        fmt.Println("myapp: hello")

        return 0 // exit status
}
```

2. In file 'conbox/applets.go', import the applet package and include its Run() function in the applet table: 

```
$ more conbox/applets.go
package main

import (
	// (...)
        "github.com/udhos/conbox/applets/myapp" // <-- import the applet package
	// (...)
)

func loadApplets() map[string]common.AppletFunc {
        tab := map[string]common.AppletFunc{
		// (...)
                "myapp": myapp.Run, // <-- point applet name to its Run() function
		// (...)
        }
        return tab
}
```

3. Rebuild conbox and test the new applet:

```
$ go install ./conbox
$ conbox myapp
myapp: hello
```

# Docker

Get 'conbox' as docker image `udhos/conbox:latest` from:

https://hub.docker.com/r/udhos/conbox

## Run in docker

Run applet:

```bash
docker run --rm udhos/conbox:latest cat /etc/passwd
```

Run interactive shell:

```bash
docker run --rm -ti udhos/conbox:latest shell
```

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

# Related work

## Go Projects

- https://github.com/u-root/u-root

Unfortunately these projects seem inactive:

- https://github.com/surma/gobox
- https://github.com/laher/someutils
- https://github.com/shirou/toybox
- https://github.com/u35s/busybox

## Non-Go projects

- https://www.busybox.net/
- http://landley.net/toybox/
