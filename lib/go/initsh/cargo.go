
package initsh

import "strings"

type cargoDeliverer struct {
	bucket cargoBucketer
	slotter cargoSlotter
}

type cargoAddress struct {
	bucket string
	slot string
}

func pathToCargoAddress (dest cargoAddress) string {
	return dest.bucket + "/" + dest.slot
}

type cargoSlotter interface {
	slotCargo (dest cargoAddress, originPath string) error
}

type cargoBucketer interface {
	bucketCargo (bucket string) error
}

func (d cargoDeliverer) deliver (dest cargoAddress, originPath string) error {
	err := d.bucket.bucketCargo(dest.bucket)
	if err != nil {
		return err
	}
	return d.slotter.slotCargo(dest, originPath)
}

type bucketerPathMode struct { }

func (b bucketerPathMode) bucketCargo (bucket string) error {
	db := dirBinder{}
	return bindTo(db, bucket)
}

type bucketerSubcmd struct { }

func (b bucketerSubcmd) bucketCargo (bucket string) error {
	return bindTo(makeSubcmdBinder(), bucket)
}

type binderSlotter struct {
	binder idempotentBinder
}

func (s binderSlotter) slotCargo (dest cargoAddress, originPath string) error {
	return bindTo(s.binder, pathToCargoAddress(dest))
}

type subcmdSlotter struct {
	inTreeDir cargoSlotter
}

func (s subcmdSlotter) slotCargo (dest cargoAddress, originPath string) error {
	fileDest := cargoAddress{toSubcmdTreeDir(dest.bucket), dest.slot}
	return bindCargo(makeSubcmdBinder(), fileDest)
}

func toSubcmdTreeDir (entryPath string) string {
	return strings.TrimRight(entryPath, "/") + "-subcmd/"
}


type symLinkSlotter struct { }
type dataSlotter struct { }
type gzDataSlotter struct { }

func (s symLinkSlotter) slotCargo (dest cargoAddress, originPath string) error {
	return bindCargo(linkBinder{originPath}, dest)
}

func (s dataSlotter) slotCargo (dest cargoAddress, originPath string) error {
	return bindCargo(makeReadBinder(originPath, "cat"), dest)
}

func (s gzDataSlotter) slotCargo (dest cargoAddress, originPath string) error {
	return bindCargo(makeReadBinder(originPath, "zcat"), dest)
}

func bindCargo (binder idempotentBinder, dest cargoAddress) error {
	return  bindTo(binder, pathToCargoAddress(dest))
}
