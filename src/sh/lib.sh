#!/bin/bash -eu

shLibDir=$(dirname $(readlink -f ${BASH_SOURCE[0]}))

. $shLibDir/default.sh
. $shLibDir/discover.sh
. $shLibDir/import.sh
. $shLibDir/path.sh
. $shLibDir/setup.sh
. $shLibDir/mktemp.sh
. $shLibDir/context.sh

buildImports=$shLibDir/../../lib/import
