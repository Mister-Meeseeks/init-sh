package initsh

import "os"

type ImportLobby interface {
	initDigest (digestKey string) (string, error)
	pullDigest (digestKey string) (string, error)
	pullImport (spaceName string) (string, error)
	directImport (spaceName string) (ImportDirector, error)
}

type importLobby struct {
	importDir string
	div string
}

func MakeLobby (importDir string, div string) ImportLobby {
	return importLobby{importDir, div}
}

func (l importLobby) initDigest (digestKey string) (string, error) {
	err := readyDigestDir(l.importDir)
	if (err != nil) {
		return "", err
	}
	return formDigestPath(l.importDir, digestKey), nil
}

func (l importLobby) pullDigest (digestKey string) (string, error) {
	path := formDigestPath(l.importDir, digestKey)
	_, err := os.Stat(path)
	if (err != nil) {
		return "", err
	}
	return path, nil
}

func (l importLobby) pullImport (spaceName string) (string, error) {
	path := formImportPath(l.importDir, spaceName)
	_, err := os.Stat(path)
	if (err != nil) {
		return "", err
	}
	return path, nil
}

func (l importLobby) directImport (spaceName string) (ImportDirector, error) {
	path, err := makeImportSpace(l.importDir, spaceName)
	if (err != nil) {
		return MakeImporter("", "", l.div), err
	}
	return MakeImporter(path, path, l.div), nil
}

func makeImportSpace (importDir string, spaceName string) (string, error) {
	err := prepImportParent(importDir)
	if (err != nil) {
		return "", err
	}
	return makeImportLeaf(importDir, spaceName)
}

func makeImportLeaf (importDir string, spaceName string) (string, error) {
	spaceDir := formImportPath(importDir, spaceName)
	err := mkdirIdempot(spaceDir)
	if (err != nil) {
		return "", err
	}
	return spaceDir, nil
}

func prepImportParent (importDir string) error {
	err := mkdirIdempot(importDir)
	if (err != nil) {
		return err
	}
	return mkdirIdempot(importDir + "/spaces/")
}

func formImportPath (importDir string, spaceName string) string {
	return importDir + "/spaces/" + spaceName
}

func readyDigestDir (importDir string) error {
	return mkdirIdempot(formDigestDir(importDir))
}

func formDigestDir (importDir string) string {
	return importDir + "/digests/"
}

func formDigestPath (importDir string, digestKey string) string {
	return formDigestDir(importDir) + digestKey
}

func mkdirIdempot (path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, os.ModeDir)
	} else {
		return nil
	}
}
