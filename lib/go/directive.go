
package initsh

type ImportDirector interface {
	importToPath (srcPath string, namer pathNamer) error
	importSubcmd (srcPath string, namer subcmdNamer) error
	importDataPath (srcPath string, namer pathNamer) error
	importDataSubcmd (srcPath string, namer subcmdNamer) error
}

type importTargetDirs struct {
	bin string
	lib string
}

type importDirector struct {
	tgt importTargetDirs
}

type pathNamer struct {
}

type subcmdNamer struct {
}
