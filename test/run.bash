#! /bin/bash

set -e

uMagNUS -vet *.mx3

uMagNUS -paranoid=false -failfast -cache /tmp -f -http "" *.go *.mx3

