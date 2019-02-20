
package initsh

import "errors"
import "os"
import "path/filepath"

func linkTo (tgtPath string, linkPath string) error {
	s, err := os.Stat(linkPath)
	if os.IsNotExist(err) {
		return os.Symlink(tgtPath, linkPath)
	} else if err != nil {
		return err
	} else {
		return assertLinkTo(s, linkPath, tgtPath)
	}
}

func assertLinkTo (stat os.FileInfo, linkPath string, tgtPath string) error {
	if (isSymLink(stat.Mode())) {
		current, err := filepath.EvalSymlinks(linkPath)
		if err != nil {
			return err
		}
		if current != tgtPath {
			return raiseLinkConflict(linkPath, tgtPath, current)
		}
		return nil
	} else {
		return errors.New("Cannot symLink at " + linkPath +
			". Previous non-link file exists.")
	}
}

func raiseLinkConflict (link string, tgt string, pre string) error {
	msg := "initSh Error: Name conflict when creating link=" +
		link + ". Previously exists with import target=" +
		pre + " which differs from requested import=" +
		tgt
	return errors.New(msg)
}

func mkdir (path string) error {
	s, err := os.Stat(path)
	if !(s.IsDir()) {
		return errors.New("File exists at directory targer: " + path)
	} else if os.IsNotExist(err) {
		return os.Mkdir(path, 0755)
	} else {
		return err
	}
}
