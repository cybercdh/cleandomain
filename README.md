## cleandomain
Takes a list of domains piped from stdin, then aims to clean them of special chars for parsing into other programs

## Recommended Usage
$ cat domains | cleandomain | something else...

## Install
You need to have the latest version (1.19+) of Go installed and configured (i.e. with $GOPATH/bin in your $PATH):

go install github.com/cybercdh/cleandomain@latest