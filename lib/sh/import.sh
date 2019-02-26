
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

function openUpForImports() {
    export INIT_SH_IMPORT_DIRECTIVES=""
}

function importShellDir() {
    importForType shell "$@"
}

function importSubcmdDir() {
    importForType subcmd "$@"
}

function importDataDir() {
    importForType data "$@"
}


function importForType() {
    local directiveType=$1
    local importDir=$2
    local namespace=""
    if [[ $# -ge 3 ]] ; then
	namespace=$3
    fi
    importDirective $directiveType $importDir $namespace
}

function importDirective() {
    local directive=$1
    local importDir=$2
    local namespace="$3"

    local dirPrefix=$(formDirectivePrefix $directive "$namespace")
    local namePostfix=$(formNamespacePostfix "$namespace")    
    local directive=${dirPrefix}:${importDir}${namePostfix}
    export INIT_SH_IMPORT_DIRECTIVES="$INIT_SH_IMPORT_DIRECTIVES $directive"
}

function formDirectivePrefix() {
    local directive=$1
    local namespace="$2"
    if isNestedNamespace "$namespace" ; then
	echo $directive-nest
    else
	echo $directive
    fi
}

function formNamespacePostfix() {
    local namespace="$1"
    if [[ -z $namespace ]] ; then
	echo ""
    else
	echo ":$namespace"
    fi
}

function sweepImports() {
    $buildImports \
	$(retrieveProjBinView) $(retrieveProjLibView) \
	$INIT_SH_IMPORT_DIRECTIVES
    unset $INIT_SH_IMPORT_DIRECTIVES
}
