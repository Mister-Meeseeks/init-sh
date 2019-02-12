
function addPythonPath() {
    local pyPath=$1
    : ${PYTHONPATH:=""}
    export PYTHONPATH=$PYTHONPATH:$pyPath
}

function addHaskellPath() {
    local hsPath=$1
    : ${INIT_SH_HASKELL_PATH:=""}
    export INIT_SH_HASKELL_PATH=$INIT_SH_HASKELL_PATH:$hsPath
}

function addRSrcPath() {
    local rPath=$1
    : ${R_SRC_PATH:=""}
    export R_SRC_PATH=$R_SRC_PATH:$rPath
}

function addPythonCoLib() {
    addPythonPath $(getPythonCoLibDir)
}

function addHaskellCoLib() {
    addHaskellPath $(getHaskellCoLibDir)
}

function addRCoLib() {
    addRSrcPath $(getRCoLibDir)
}

function addProjEtcData() {
    importDataDir $(getProjEtcDir) etc::
}

function importShellMaybe() {
    local importDir=$1
    shift; local addArgs=$@
    if [[ -d $importDir ]] ; then
	importShellDir $importDir $addArgs
    fi
}

function importShellDir() {
    set -o pipefail
    local importDir=$1
    [[ $# -ge 2 ]] && local namespace=$2 || local namespace=""
    importExecutables $importDir "$namespace"
    importShellLibs $importDir "$namespace"
}

function importExecutables() {
    local importDir=$1
    local namespace=$2
    findExecutables $importDir | importBinsToProject "$namespace"
}

function importShellLibs() {
    local importDir=$1
    local namespace=$2
    findShellLibs $importDir | importLibsToProject "$namespace"
}

function importDataDir() {
    local importDir=$1
    local namespace=$2
    findDataFiles $importDir | importDataToProject "$namespace"
}

function findExecutables() {
    local findDir=$1
    findFilesUnderDir $findDir "-executable" ""
}

function findShellLibs() {
    local findDir=$1
    findFilesUnderDir $findDir "! -executable" "$shellLibPattern"
}

function findDataFiles() {
    local findDir=$1
    findFilesUnderDir $findDir "! -executable" ""
}

function findFilesUnderDir() {
    local findDir=$1
    local findArgs=$2
    local matchPattern=$3
    canonDir=$(expandFindDirectoryPath $findDir)
    find $canonDir -path '*/\.*' -prune -o \
	\( -type f -or -type l \) -and \( $findArgs \) -printf "%H %P\n" \
	| grepMatchFiles "$matchPattern"
}

function expandFindDirectoryPath() {
    if [[ ! -e $findDir ]] ; then
	echo "initSh Error: Import target $findDir doesn't exist" >&2
	exit 1
    fi
    readlink -f $findDir
}

function grepMatchFiles() {
    local matchPattern="$1"
    (egrep -v "$ignoreFilePattern" \
	    | egrep -v "$keywordFilePattern" \
	    | egrep "$matchPattern") \
	|| true
}

function importBinsToProject() {
    local namespace=$1
    $convertImportArgs "$namespace" 1 \
	| pointInBinView
}

function importLibsToProject() {
    local namespace=$1
    $convertImportArgs "$namespace" 0 \
	| pointInLibView
}

function importDataToProject() {
    local namespace=$1
    $convertImportArgs "$namespace" 0 \
	| pointInDataView    
}
