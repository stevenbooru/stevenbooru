#!/bin/sh

set -e
set -x

source ./env.sh

mkdir bin ||:
cdbin

go build stevenbooru.cf/cmd/...
