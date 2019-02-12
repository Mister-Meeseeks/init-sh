
# Use case: Projects of the form proj/lib/sh/init.sh
#                                           /[initSh scripts]
#                                        /py/[python libraries] (e.g.)
#                                        /[other library directory]/
function getCoLibDir() {
    local coLibName=$1
    echo $(discoverProjectDir)/../$coLibName/
}

function getPythonCoLibDir() {
    getCoLibDir py
}

function getHaskellCoLibDir() {
    getCoLibDir hs
}

function getRCoLibDir() {
    getCoLibDir R
}

function getParentLibDir() {
    echo $(discoverProjectDir)/../
}

# Use case: Projects of the form proj/src/sh/init.sh
#                                     etc/[data]
#                                     lib/[objs]
#                                     var/
function getProjEtcDir() {
    echo $(getParentProjectDir)/etc/
}

function getProjLibDir() {
    echo $(getParentProjectDir)/lib/
}

function getProjVarDir() {
    echo $(getParentProjectDir)/var/
}

function getParentProjectDir() {
    echo $(discoverProjectDir)/../../
}

function prepareFreshView() {
    resetViews
    setPathForProject
}

function setPathForProject() {
    if [[ $doesProjectHavePrecedence -gt 0 ]] ; then
	export PATH=$(getProjectPaths):$PATH
    else
	export PATH=$PATH:$(getProjectPaths)
    fi
}

function resetViews() {
    rm -r $(getProjDirs)
}

function getProjectPaths() {
    getProjDirs | sed 's+ +:+g'

}

function getProjDirs() {
    echo  $(retrieveProjectBinViewDir) \
	$(retrieveProjectLibViewDir) \
	$(retrieveProjectDataViewDir)
}

function pointInBinView() {
    pointLinksInView $(retrieveProjectBinViewDir)
}

function pointInLibView() {
    pointLinksInView $(retrieveProjectLibViewDir)
}

function pointInDataView() {
    pointDataExecsInView $(retrieveProjectDataViewDir)
}

function pointLinksInView() {
    local viewDir=$1
    while read linkArgs ; do
	createSymLinkInView $viewDir $linkArgs
    done
}

function pointDataExecsInView() {
    local viewDir=$1
    while read linkArgs ; do
	createDataExecutableInView $viewDir $linkArgs
    done    
}

function retrieveProjectBinViewDir() {
    retrieveInternalSubDir $viewBinSubDir
}

function retrieveProjectLibViewDir() {
    retrieveInternalSubDir $viewLibSubDir
}

function retrieveProjectDataViewDir() {
    retrieveInternalSubDir $viewDataSubDir
}

function retrieveInternalSubDir() {
    local subDir=$1
    local outDir=$INIT_SH_INSTANCE_DIR/$subDir
    mkdir -p $outDir
    echo $outDir
}

function createDataExecutableInView() {
    linkFromProjectSource writeDataExecutable $@
}

function writeDataExecutable() {
    local sourcePath=$1
    local targetPath=$2
    local absSource=$(readlink -f $sourcePath)
    (echo "#!/bin/bash -eu"; echo "cat $absSource") > $targetPath
    chmod u+x $targetPath
}

function createSymLinkInView() {
    local viewDir=$1
    local sourcePath=$2
    local targetName=$3
    linkFromProjectSource "ln -s" $@
}

function linkFromProjectSource() {
    local linkCmd="$1"
    local viewDir=$2
    local sourcePath=$3
    local targetName=$4

    local targetPath=$viewDir/$targetName
    assertNoLinkConflict $sourcePath $targetPath
    linkIdempotently "$linkCmd" $sourcePath $targetPath
}

function linkIdempotently() {
    local linkCmd="$1"
    local sourcePath=$2
    local targetPath=$3
    if [[ ! -e $targetPath ]] && [[ ! -L $targetPath ]] ; then
	$linkCmd $sourcePath $targetPath
    fi
}

function assertNoLinkConflict() {
    if doesBashCommandConflict $sourcePath $targetPath ; then
	raiseBashConflict $sourcePath $targetPath
    elif doesExistingLinkConflict $sourcePath $targetPath ; then
	raiseNameConflict $sourcePath $targetPath
    fi
}

function doesExistingLinkConflict() {
    local sourcePath=$1
    local targetPath=$2
    doesLinkExist $targetPath && areDiffTargets $sourcePath $targetPath
}

function doesLinkExist() {
    local targetPath=$1
    [[ -e $targetPath || -L $targetPath ]]
}

function doesBashCommandConflict() {
    local sourcePath=$1
    local targetPath=$2
    local cmdName=$(basename $targetPath)
    local cmdTgt=$(type -p $cmdName)
    doesOverrideCmd $sourcePath $cmdName &&
	areDiffTargets $sourcePath "$cmdTgt"
}

function doesOverrideCmd() {
    local sourcePath=$1
    local cmdName=$2
    [[ -x $sourcePath ]] && (type $cmdName >/dev/null 2>&1)
}

function areDiffTargets() {
    local targetA=$1
    local targetB=$2
    [[ $(readlink -f $targetA) != $(readlink -f $targetB) ]]
}

function raiseNameConflict() {
    local sourcePath=$1
    local targetPath=$2
    echo "init.sh Error: Name conflict at $targetPath. Attempt to link to" \
	"$sourcePath, contradicts existing link $(readlink -f $targetPath)" >&2
    exit 1
}

function raiseBashConflict() {
    local sourcePath=$1
    local targetPath=$2
    local execName=$(basename $targetPath)
    echo "init.sh Error: Name conflict at $targetPath. Conflicts with existing"\
	 "bash command: $(type $execName)" >&2
    exit 1
}
