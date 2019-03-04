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
	return filterWalkErrs(walkIngest(root, *ing))
}

func filterWalkErrs (err error) error {
	if (err == filepath.SkipDir) {
		return nil
	} else {
		return err
	}
}

func walkIngest (root string, hndl PathIngester) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if (err != nil) {
			return err
		} else if isIgnorable (path) {
			return filepath.SkipDir
		} else {
			return ingestWalk(path, hndl)
		}
	})
}

func ingestWalk (path string, hndl PathIngester) error {
	linkInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	return hndl.ingestPath(path, linkInfo)
}

func isIgnorable (path string) bool {
	return isBrokenLink(path) || isIgnorableName(path)
}

func isIgnorableName (path string) bool {
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
