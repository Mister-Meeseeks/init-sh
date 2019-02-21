
package initsh

import "errors"
import "os"
import "path/filepath"
import "strings"
import "io/ioutil"

type idempotentBinder interface {
	makeFresh (path string) error
	assertMatch (path string, stat os.FileInfo) error
}

func bindTo (b idempotentBinder, path string) error {
	s, err := os.Stat(path)
	if os.IsNotExist(err) {
		return b.makeFresh(path)
	} else if err != nil {
		return err
	} else {
		return b.assertMatch(path, s)
	}
}

type linkBinder struct {
	tgtPath string
}

func (b linkBinder) makeFresh (path string) error {
	return os.Symlink(path, b.tgtPath)
}

func (b linkBinder) assertMatch (path string, stat os.FileInfo) error {
	if (isSymLink(stat.Mode())) {
		prePath, err := filepath.EvalSymlinks(path)
		if err != nil {
			return err
		}
		if prePath != b.tgtPath {
			return raiseLinkConflict(path, b.tgtPath, prePath)
		}
		return nil
	} else {
		return errors.New("Cannot symLink at " + path +
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

type readBinder struct {
	tgtPath string
	readCmd string
}

func (b readBinder) makeFresh (path string) error {
	err := b.assertTargetSanit(path)
	if (err != nil) {
		return err
	}
	return ioutil.WriteFile(path, b.fmtContent(), 0755)
}

func (b readBinder) assertTargetSanit (path string) error {
	if (strings.ContainsAny(b.tgtPath, "'")) {
		return errors.New("Cannot read bind a target path with " +
			"quotations: " + b.tgtPath)
	}
	return nil
}

func (b readBinder) fmtContent() []byte {
	content := b.fmtStrContent()
	return []byte(content)
}

func (b readBinder) fmtStrContent() string {
	lines := []string{"#!/bin/bash -eu", b.readCmd + " '" + b.tgtPath + "'"}
	return strings.Join(lines, "\n")
}

func (b readBinder)  assertMatch (path string, stat os.FileInfo) error {
	if (stat.IsDir()) {
		return errors.New("Cannot create a wrapper at " + path +
			". Previous directory exists at path")
	} else if (isSymLink(stat.Mode())) {
		return errors.New("Cannot create a wrapper at " + path +
			". Previous symLink exists at path")
	} else {
		return b.assertFileMatch(path)
	}
}

func (b readBinder) assertFileMatch (path string) error {
	byteContent, err := ioutil.ReadFile(path)
	if (err != nil) {
		return err
	}
	content := string(byteContent)
	if (content != b.fmtStrContent()) {
		return errors.New("Read wrapper previosuly exists at " + path +
			"with different target than " + b.tgtPath)
	}
	return nil
}

type dirBinder struct { }


func (b dirBinder) makeFresh (path string) error {
	return os.Mkdir(path, 0755)
}

func (b dirBinder) assertFresh (path string, stat os.FileInfo) error {
	if !(stat.IsDir()) {
		return errors.New("File exists at directory targer: " + path)
	}
	return nil
}
