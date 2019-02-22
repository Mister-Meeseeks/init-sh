package initsh

import (
       "fmt"
	"os"
	"path/filepath"
)

func walkThru (root string, hndl PathIngester) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return hndl.ingestPath(path, info)
	})
}
