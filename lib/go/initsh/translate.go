
package initsh

import "strings"
import "path/filepath"
import "fmt"

type AddressTranslator interface {
	translate (origin cargoAddress) cargoAddress
}

type StackedTranslator struct {
	pipeline []AddressTranslator
}

func (t StackedTranslator) translate (origin cargoAddress) cargoAddress {
	result := origin
	for _, trans := range t.pipeline {
		result = trans.translate(result)
	}
	return result
}

func stackTrans (x AddressTranslator, y AddressTranslator) AddressTranslator {
	transs := []AddressTranslator{x, y}
	return StackedTranslator{transs}
}

type DropExtTranslator struct { }

func (t DropExtTranslator) translate (origin cargoAddress) cargoAddress {
	return cargoAddress{origin.bucket, dropPathsExtension(origin.slot)}
}

func dropPathsExtension (path string) string {
	return strings.TrimSuffix(path, filepath.Ext(path))
}

type FlattenTranslator struct {
	div string
}

type NestedTranslator struct {
	origDiv string
	div string
}

func (t FlattenTranslator) translate (origin cargoAddress) cargoAddress {
	fields := strings.Split(origin.slot, t.div)
	return cargoAddress{origin.bucket, fields[len(fields)-1]}
}

func (t NestedTranslator) translate (origin cargoAddress) cargoAddress {
	fields := strings.Split(origin.slot, t.origDiv)
	slot := strings.Join(fields, t.div)
	cleaned := strings.TrimPrefix(slot, t.div)
	return cargoAddress{origin.bucket, cleaned}
}

type SubcmdTranslator struct {
	namespace string
	importDir string
}

func (t SubcmdTranslator) translate (origin cargoAddress) cargoAddress {
	bucket := t.importDir + "/" + origin.bucket
	slot := origin.slot
	return cargoAddress{bucket, slot}
}

type NamespaceTranslator struct {
	importDir string
	namespace string
	div string
}

func (p NamespaceTranslator) translate (origin cargoAddress) cargoAddress {
	fmt.Println("Hello")
	slot := p.namespace + p.div + origin.slot
	return cargoAddress{p.importDir, slot}
}

type MirrorTranslator struct {
	importDir string
}

func (p MirrorTranslator) translate (origin cargoAddress) cargoAddress {
	return cargoAddress{p.importDir, origin.slot}
}
