
package initsh

import "os"
import "strings"

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
	return parseFileExt(path) == "sh"
}

func hasGzipExt (path string) bool {
	return parseFileExt(path) == "gz"
}

func parseFileExt (path string) string {
	fields := strings.Split(path, ".")
	if len(fields) < 2 {
		return ""
	} else {
		return fields[len(fields)-1]
	}
}
