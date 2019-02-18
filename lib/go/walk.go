package initsh

import (
       "fmt"
	"os"
	"path/filepath"
)

func walkThru (root string, hndl PathIngester) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		hndl.ingestPath(path, info)
		return nil
	})
}
