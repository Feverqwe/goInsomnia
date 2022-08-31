#!/bin/sh

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $SCRIPT_DIR/..

source ./scripts/_variables.sh

rm ./${BINARY}
go build -trimpath -o ${BINARY}
