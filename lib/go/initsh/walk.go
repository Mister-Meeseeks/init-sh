package initsh

import (
	"os"
	"path/filepath"
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
		}
		return hndl.ingestPath(path, info)
	})
}
