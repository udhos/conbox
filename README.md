[![license](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/udhos/conbox/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/udhos/conbox)](https://goreportcard.com/report/github.com/udhos/conbox)

# conbox
Go implementation of unix-like utilities as single static executable intended for small container images.

# Install

```bash
git clone https://github.com/udhos/conbox
cd conbox
go install ./conbox
```

# Usage

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

# Related work

- https://github.com/surma/gobox
- https://github.com/laher/someutils
- https://github.com/shirou/toybox
- https://www.busybox.net/
- http://landley.net/toybox/
