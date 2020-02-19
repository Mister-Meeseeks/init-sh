
function discoverProjectInit() {
    local searchDir=$(getRunScriptDir)
    while [[ $searchDir != $rootDir ]] ; do
	if [[ -e $searchDir/$initFileName ]] ; then
	    echo $searchDir/$initFileName
	    return 0
	fi
	searchDir=$(dirname "$searchDir")
    done
    raiseNoInitDiscovered
}

function getRunScriptDir() {
    local canonPath=$(readlink -f $runScript)
    if [[ -z $canonPath ]] ; then
	pwd
    elif [[ -d $canonPath ]] ; then
	echo $canonPath
    else
	dirname "$canonPath"
    fi
}

function discoverProjectDir() {
    dirname "$(discoverProjectInit)"
}

function raiseNoInitDiscovered() {
    echo "init.sh Error: init.sh not found in any parent directory" \
	"for $runScript. Cannot determine project location" >&2
    exit 1
}

function blueprintProject() {
    openUpForImports
    addLocalProject
    sourceLocalProjectInit
}

function addLocalProject() {
    importShellDir $(discoverProjectDir) $localProjectNamespace
}

function sourceLocalProjectInit() {
    . $(discoverProjectInit) 
}

function randomizeInstanceDir() {
    local instParent=$(retrieveRuntimeDir)/instances/
    local randTemplate=tmpInst.XXXXXXXX
    mkdir -p $instParent
    mktemp -d $instParent/$randTemplate
}

function hashImportDir() {
    local importParent=$(retrieveRuntimeDir)/imports/
    local projHash=$(hashProjectPath)
    mkdir -p $importParent/$projHash
    echo $importParent/$projHash
}

function hashProjectPath() {
    path=$(readlink -f $(discoverProjectDir))
    # 128 bits of entropy oughta be enough for anybody...
    echo $path | sha256sum | cut -b 1-32
}

function retrieveRuntimeDir() {
    local instParent=$(discoverProjectDir)/$initShSubDir/
    mkdir -p $instParent
    echo $instParent
}
