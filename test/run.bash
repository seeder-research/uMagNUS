#! /bin/bash

set -e

mumax3cl -vet *.mx3

mumax3cl -paranoid=false -failfast -cache /tmp -f -http "" *.go *.mx3

