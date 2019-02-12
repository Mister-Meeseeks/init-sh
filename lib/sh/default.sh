
doesProjectHavePrecedence=1

rootDir=/
initFileName=init.sh

initShSubDir=.initSh/
viewSubDir=pathViews/
viewBinSubDir=$viewSubDir/bin/
viewLibSubDir=$viewSubDir/lib/
viewDataSubDir=$viewSubDir/data/

ignoreFilePattern='(~$|^#)'
keywordFilePattern="(^|[/])($initFileName)\$"
shellLibPattern='(^|[/])(..*)([.]sh)$'

localProjectNamespace=""
shellCmd="bash -eu"
