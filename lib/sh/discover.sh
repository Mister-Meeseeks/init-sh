
function discoverProjectInit() {
    local searchDir=$(getRunScriptDir)
    while [[ $searchDir != $rootDir ]] ; do
	if [[ -e $searchDir/$initFileName ]] ; then
	    echo $searchDir/$initFileName
	    return 0
	fi
	searchDir=$(dirname $searchDir)
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
	dirname $canonPath
    fi
}

function discoverProjectDir() {
    dirname $(discoverProjectInit)
}

function raiseNoInitDiscovered() {
    echo "init.sh Error: init.sh not found in any parent directory" \
	"for $runScript. Cannot determine project location" >&2
    exit 1
}

function blueprintProject() {
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
    local instParentDir=$(retrieveInstanceParent)
    local randTemplate=tmpInst.XXXXXXXX
    mktemp -d $instParentDir/$randTemplate
}

function retrieveInstanceParent() {
    local instParentDir=$(discoverProjectDir)/$initShSubDir/
    mkdir -p $instParentDir
    echo $instParentDir
}
