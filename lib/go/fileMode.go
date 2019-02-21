package initsh

import "os"

func isExecPerm (mode os.FileMode) bool {
	return mode.Perm() & 0100 > 0
}

func isSymLink (mode os.FileMode) bool {
	return mode & os.ModeSymlink > 0
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

