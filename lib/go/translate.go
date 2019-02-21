
package initsh

import "strings"
import "path/filepath"

type AddressTranslator interface {
	translate (origin cargoAddress) cargoAddress
}

type StackedTranslators struct {
	pipeline []AddressTranslator
}

func (t StackedTranslators) translate (origin cargoAddress) cargoAddress {
	result := origin
	for _, trans := range t.pipeline {
		result = trans.translate(result)
	}
	return result
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

func (t FlattenTranslator) translate (origin cargoAddress) cargoAddress {
	fields := strings.Split(origin.slot, t.div)
	return cargoAddress{origin.bucket, fields[len(fields)-1]}
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
	slot := p.namespace + p.div + origin.slot
	return cargoAddress{p.importDir, slot}
}

type MirrorTranslator struct {
	importDir string
}

func (p MirrorTranslator) translate (origin cargoAddress) cargoAddress {
	return cargoAddress{p.importDir, origin.slot}
}
