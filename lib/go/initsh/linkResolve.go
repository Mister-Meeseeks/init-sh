package initsh

import "fmt"
import "os"
import "path/filepath"

func collapseLinks (path string, info os.FileInfo) (os.FileInfo, error) {
	if (isSymLink(info.Mode())) {
		return collapseLinked(path)
	} else {
		return info, nil
	}
}

func collapseLinked (path string) (os.FileInfo, error) {
	info, err := statThruLinks(path)
	if err == nil {
		return info, err
	} else if (info.IsDir()) {
		return info, DirLinkImportError{path}
	} else {
		return info, nil
	}
}

func statThruLinks (path string) (os.FileInfo, error) {
	canon, err := filepath.EvalSymlinks(path)
	if err != nil {
		return nil, err
	}
	return os.Stat(canon)
}

type DirLinkImportError struct {
	linkPath string
}

func (e DirLinkImportError) Error() string {
	return fmt.Sprintf("initSh Does not support importing symLinks " +
		"to directories - %s", e.linkPath)
}
