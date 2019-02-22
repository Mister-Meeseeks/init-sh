
package initsh

type ImportDirector interface {
	importShell (path string, namespace *string) PathIngester
	importNested (path string, namespace *string) PathIngester
	importSubcmd (path string, namespace string) PathIngester
	importNestSubcmd (path string, namespace string) PathIngester
	importData (path string, namespace *string) PathIngester
	importNestData (path string, namespace *string) PathIngester
}

type importDirector struct {
	binPath string
	libPath string
	div string
}

func (d importDirector) importShell (path string, namespace *string) PathIngester {
	ex := ImportFunnel{ExecFilter{}, d.flatBinTrans(namespace),
		d.linkPathShipper(), path}
	lib := ImportFunnel{ShellLibFilter{}, d.flatLibTrans(namespace),
		d.linkPathShipper(), path}
	return mergeIngesters(ex, lib)
}

func (d importDirector) importNested (path string, namespace *string) PathIngester {
	ex := ImportFunnel{ExecFilter{}, d.flatBinTrans(namespace),
		d.linkPathShipper(), path}
	lib := ImportFunnel{ShellLibFilter{}, d.baseLibTrans(namespace),
		d.linkPathShipper(), path}
	return mergeIngesters(ex, lib)	
}

func (d importDirector) importData (path string, namespace *string) PathIngester {
	ex := ImportFunnel{DataItemFilter{}, d.flatUndropBinTrans(namespace),
		d.dataShipper(), path}
	lib := ImportFunnel{GzDataFilter{}, d.flatBinTrans(namespace),
		d.gzDataShipper(), path}
	return mergeIngesters(ex, lib)
}

func (d importDirector) importNestData (path string, namespace *string) PathIngester {
	ex := ImportFunnel{DataItemFilter{}, d.baseBinTrans(namespace),
		d.dataShipper(), path}
	lib := ImportFunnel{GzDataFilter{}, d.dropBinTrans(namespace),
		d.gzDataShipper(), path}
	return mergeIngesters(ex, lib)
}

func (d importDirector) importSubcmd (path string, namespace string) PathIngester {
	ex := ImportFunnel{ExecFilter{}, d.subcmdFlatBinTrans(namespace),
		d.subcmdShipper(), path}
	lib := ImportFunnel{ShellLibFilter{}, d.subcmdFlatLibTrans(namespace),
		d.subcmdShipper(), path}
	return mergeIngesters(ex, lib)
}

func (d importDirector) importNestSubcmd (path string, namespace string) PathIngester {
	ex := ImportFunnel{ExecFilter{}, SubcmdTranslator{namespace, d.binPath},
		d.subcmdShipper(), path}
	lib := ImportFunnel{ShellLibFilter{}, SubcmdTranslator{namespace, d.libPath},
		d.subcmdShipper(), path}
	return mergeIngesters(ex, lib)
}

func (d importDirector) subcmdFlatBinTrans (namespace string) AddressTranslator {
	return stackTrans(FlattenTranslator{"/"},
		SubcmdTranslator{namespace, d.binPath})
}

func (d importDirector) subcmdFlatLibTrans (namespace string) AddressTranslator {
	return stackTrans(FlattenTranslator{"/"},
		SubcmdTranslator{namespace, d.libPath})
}

func (d importDirector) subcmdBinTrans (namespace string) AddressTranslator {
	return SubcmdTranslator{namespace, d.binPath}
}

func (d importDirector) flatBinTrans (namespace *string) AddressTranslator {
	return stackTrans(FlattenTranslator{"/"}, d.dropBinTrans(namespace))
}

func (d importDirector) flatUndropBinTrans (namespace *string) AddressTranslator {
	return stackTrans(FlattenTranslator{"/"}, d.baseBinTrans(namespace))
}

func (d importDirector) flatLibTrans (namespace *string) AddressTranslator {
	return stackTrans(FlattenTranslator{"/"}, d.baseLibTrans(namespace))
}

func (d importDirector) dropBinTrans (namespace *string) AddressTranslator {
	return stackTrans(DropExtTranslator{}, d.baseBinTrans(namespace))
}

func (d importDirector) baseBinTrans (namespace *string) AddressTranslator {
	return d.baseTrans(d.binPath, namespace)
}

func (d importDirector) baseLibTrans (namespace *string) AddressTranslator {
	return d.baseTrans(d.binPath, namespace)
}

func (d importDirector) baseTrans (dest string, namespace *string) AddressTranslator {
	if (namespace == nil) {
		return MirrorTranslator{dest}
	} else {
		return NamespaceTranslator{dest, *namespace, d.div}
	}
}

func (d importDirector) linkPathShipper() cargoDeliverer {
	return cargoDeliverer{bucketerPathMode{}, symLinkSlotter{}}
}

func (d importDirector) subcmdShipper() cargoDeliverer {
	return cargoDeliverer{bucketerSubcmd{}, subcmdSlotter{}}
}

func (d importDirector) dataShipper() cargoDeliverer {
	return cargoDeliverer{bucketerPathMode{}, dataSlotter{}}
}

func (d importDirector) gzDataShipper() cargoDeliverer {
	return cargoDeliverer{bucketerPathMode{}, gzDataSlotter{}}
}
