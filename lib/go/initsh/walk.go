package initsh

import (
	"os"
	"path/filepath"
	"strings"
)

func WalkThru (importArg string, dir ImportDirector) error {
	ing, root, err := parseImportStr(importArg, dir)
	if (err != nil) {
		return err
	}
	return walkIngest(root, *ing)
}

func walkIngest (root string, hndl PathIngester) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if isIgnorable (path) {
			return filepath.SkipDir
		} else {
			return hndl.ingestPath(path, info)
		}
	})
}

func isIgnorable (path string) bool {
	base := filepath.Base(path)
	return isHiddenBase(base) || isScratchBase(base)
}

func isHiddenBase (path string) bool {
	return strings.HasPrefix(path, ".") &&
		path != "." && path != "./" &&
		path != ".." && path != "../"
}

func isScratchBase (path string) bool {
	return strings.HasSuffix(path, "~")
}
