#!/usr/bin/env bash
set -e

cd "$(dirname $BASH_SOURCE)/src"

set -x

sudo docker build -t king011/mariabackup:10.5.9-v1.0.3 .