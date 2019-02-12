#!/bin/bash -eu

libDir=$(dirname $(readlink -f ${BASH_SOURCE[0]}))

. $libDir/default.sh
. $libDir/discover.sh
. $libDir/import.sh
. $libDir/path.sh
. $libDir/setup.sh
. $libDir/mktemp.sh
. $libDir/context.sh

convertImportArgs=$libDir/importArgs.py
