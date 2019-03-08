package initsh

import "os"
import "path/filepath"

func isExecPerm (mode os.FileMode) bool {
	return mode.Perm() & 0100 > 0
}

func isSymLink (mode os.FileMode) bool {
	return mode & os.ModeSymlink > 0
}

func doesPathExist (path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func isBrokenLink (path string) bool {
	_, err := filepath.EvalSymlinks(path)
	return err != nil
}

func isExecFile (info os.FileInfo) bool {
	return !info.IsDir() &&
		!isSymLink(info.Mode()) &&
		isExecPerm(info.Mode())
}

func isNonExecFile (info os.FileInfo) bool {
	return !info.IsDir() &&
		!isSymLink(info.Mode()) &&
		!isExecPerm(info.Mode())	
}

