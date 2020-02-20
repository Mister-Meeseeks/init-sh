package initsh

import "github.com/google/uuid"

func ImportNeeded (directive string, lobby ImportLobby, pathUniv []string) (string, error) {
	paths, err := FilterImportDirective(directive, pathUniv)
	if (err != nil) {
		return "", err
	}
	return overPaths(directive, lobby, paths)
}

func overPaths (directive string, lobby ImportLobby, paths []string) (string, error) {
	digest := DigestTarget(directive, paths, lobby)
	if (digest.isMatch) {
		return lobby.pullImport(digest.space)
	} else {
		return cacheImport(directive, lobby, paths)
	}
}

func cacheImport (directive string, lobby ImportLobby, paths []string) (string, error) {
	loader, err := spawnImporter(lobby)
	if (err != nil) {
		return "", err
	}
	
	err = WalkPreScanned(directive, (*loader).importer, paths)
	if (err != nil) {
		return "", err
	}

	image := DigestPreImage{directive, (*loader).space, paths}
	spaceDir, err2 := finalizeDigest(image, (*loader).space, lobby)
	if (err2 != nil) {
		return "", err2
	}
	return spaceDir, nil
}

type loadView struct {
	space string
	importer ImportDirector
}

func spawnImporter (lobby ImportLobby) (*loadView, error) {
	space, err := randomizeSpaceName()
	if (err != nil) {
		return nil, err
	}
	
	importer, err2 := lobby.directImport(space)
	if (err2 != nil) {
		return nil, err2
	}
	
	return &loadView{space, importer}, nil
}

func randomizeSpaceName() (string, error) {
	space, err := uuid.NewUUID()
	if (err != nil) {
		return "", err
	}
	
	return space.String(), nil
}

func finalizeDigest (digest DigestPreImage, space string, lobby ImportLobby) (string, error) {
	err := WriteDigest(digest, lobby)
	if (err != nil) {
		return "", err
	}
	return lobby.pullImport(space)
}
