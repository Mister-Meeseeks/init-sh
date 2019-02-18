package initsh

import (
       "fmt"
	"os"
	"path/filepath"
)

type PathHandler interface {
	ingestPath(path string, info os.FileInfo)
}

type PathPrinter struct { }

func (p PathPrinter) ingestPath (path string, info os.FileInfo) {
	fmt.Printf("Descent path %q - base %q - isDir= %q - Mode %o\n",
		path, info.Name(), info.IsDir(), info.Mode().Perm())
}

type PathExecFilter struct {
	inside PathHandler
}


func (p PathExecFilter) ingestPath (path string, info os.FileInfo) {
	if (isExecFile(info)) {
		p.inside.ingestPath(path, info)
	}
}

func walkThru (root string, hndl PathHandler) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		hndl.ingestPath(path, info)
		return nil
	})
}
