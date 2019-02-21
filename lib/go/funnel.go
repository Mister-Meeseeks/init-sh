package initsh

import "os"
import "errors"
import "strings"
import "path/filepath"

type ImportFunnel struct {
	translator AddressTranslator
	shipper cargoDeliverer
	tester ImportFilter
	importRoot string
}

func (p ImportFunnel) ingestPath (importable string, info os.FileInfo) error {
	isGood, err := p.tester.canImport(importable, info)
	if err != nil {
		return err
	} else if (isGood) {
		return p.importPath(importable)
	} else {
		return nil
	}
}

func (p ImportFunnel) importPath (importable string) error {
	origin, err := splitToCargo(importable, p.importRoot)
	if (err != nil) {
		return err
	}
	dest := p.translator.translate(origin)
	return p.shipper.deliver(dest, importable)
}

func splitToCargo (importPath string, importRoot string) (cargoAddress, error) {
	canonRoot := filepath.Clean(importRoot)
	canonPath := filepath.Clean(importPath)
	if (!(strings.HasPrefix(canonPath, canonRoot))) {
		return raiseImportRoot(importRoot, importPath)
	} else {
		slot := strings.TrimPrefix(canonPath, canonRoot)
		return cargoAddress{importRoot, slot}, nil
	}
}

func raiseImportRoot (root string, path string) (cargoAddress, error) {
	empty := cargoAddress{"", ""}
	return empty, errors.New("Internal initSh error: importRoot=" +
		root + " not found in importPath=" + path)
}
