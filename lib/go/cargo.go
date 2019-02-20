
package initsh

type cargoDeliverer struct {
	bucket cargoBucketer
	slotter cargoSlotter
}

type cargoDest struct {
	bucket string
	slot string
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
