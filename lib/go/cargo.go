
package initsh

type cargoDeliverer struct {
	bucket cargoBucketer
	slotter cargoSlotter
}

type cargoDest struct {
	bucket string
	slot string
}

func pathToCargoDest (dest cargoDest) string {
	return dest.bucket + "/" + dest.slot
}

type cargoSlotter interface {
	slotCargo (dest cargoDest, originPath string) error
}

type cargoBucketer interface {
	bucketCargo (bucket string) error
}

func (d *cargoDeliverer) deliver (dest cargoDest, originPath string) error {
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

func (s binderSlotter) slotCargo (dest cargoDest, originPath string) error {
	return bindTo(s.binder, pathToCargoDest(dest))
}

func makeSymLinkSlotter() cargoSlotter {
	return binderSlotter{linkBinder{}}
}

func makeDataSlotter (path string) cargoSlotter {
	return binderSlotter{makeReadBinder(path, "cat")}
}

func makeGzDataSlotter (path string) cargoSlotter {
	return binderSlotter{makeReadBinder(path, "zcat")}
}
