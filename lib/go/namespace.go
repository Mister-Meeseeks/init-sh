package initsh

import "strings"
import "path/filepath"

type ImportNamer interface {
	nameImportEntry (path string) string
}

type RootImportNamer struct {
	rootPath string
	wrapped ImportNamer
}

func (r RootImportNamer) nameImportEntry (path string) string {
	cleanPath := filepath.Clean(path)
	subPath := strings.TrimPrefix(cleanPath, filepath.Clean(r.rootPath))
	return r.wrapped.nameImportEntry(subPath)
}

type BaseImportNamer struct { }

func (b BaseImportNamer) nameImportEntry (path string) string {
	return filepath.Base(path)
}

type MinusExtensionNamer struct { }

func (n MinusExtensionNamer) nameImportEntry (path string) string {
	baseName := filepath.Base(path)
	fields := strings.Split(baseName, ".")
	if (len(fields) == 1) {
		return baseName
	} else {
		return strings.Join(fields[:len(fields)-1], ".")
	}
}
