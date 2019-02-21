
package initsh

import "os"
import "path/filepath"

type ExecFilter struct { }
type ShellLibFilter struct { }
type DataItemFilter struct { }
type GzDataFilter struct { }

func (f ExecFilter) doIngest (path string, info os.FileInfo) (bool, error) {
	return isExecFile(info), nil
}

func (f ShellLibFilter) doIngest (path string, info os.FileInfo) (bool, error) {
	return isNonExecFile(info) && hasShellExt(path), nil
}

func (f DataItemFilter) doIngest (path string, info os.FileInfo) (bool, error) {
	return isNonExecFile(info) && !hasGzipExt(path), nil
}

func (f GzDataFilter) doIngest (path string, info os.FileInfo) (bool, error) {
	return isNonExecFile(info) && hasGzipExt(path), nil
}

func hasShellExt (path string) bool {
	return filepath.Ext(path) == ".sh"
}

func hasGzipExt (path string) bool {
	return filepath.Ext(path) == ".gz"
}
