
package initsh

import "errors"
import "os"
import "path/filepath"
import "strings"
import "bytes"
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

func (b *linkBinder) makeFresh (path string) error {
	return os.Symlink(path, b.tgtPath)
}

func (b *linkBinder) assertMatch (path string, stat os.FileInfo) error {
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

type fileContentBinder struct {
	contentLines []string
}

func (b fileContentBinder) makeFresh (path string) error {
	err := prepareDir(path)
	if err != nil {
		return err
	}
	return writeContent(path)
}

func (b fileContentBinder) prepareDir (path string) error {
	dirPath := filepath.Dirname(path)
	bindTo(dirBinder{}, dirPath)
}

func (b fileContentBinder) writeContent (path string) error {
	return ioutil.WriteFile(path, b.fmtContent(), 0755)
}

func (b fileContentBinder) fmtContent() []byte {
	content := strings.Join(b.contentLines, "\n")
	return []byte(content)
}

func (b fileContentBinder) assertMatch (path string, stat os.FileInfo) error {
	if (stat.IsDir()) {
		return errors.New("Cannot create a file at " + path +
			". Previous directory exists at path")
	} else if (isSymLink(stat.Mode())) {
		return errors.New("Cannot create a file at " + path +
			". Previous symLink exists at path")
	} else {
		return b.assertFileMatch(path)
	}
}

func (b fileContentBinder) assertFileMatch (path string) error {
	byteContent, err := ioutil.ReadFile(path)
	if (err != nil) {
		return err
	}
	if (!(bytes.Equal(byteContent, b.fmtContent()))) {
		return errors.New("Read wrapper previosuly exists at " + path +
			"with different content")
	}
	return nil
}

func makeReadBinder (tgtPath string, readCmd string) fileContentBinder {
	lines := []string{"#!/bin/bash -eu", readCmd + " '" + tgtPath + "'"}
	return fileContentBinder{lines}
}

func makeSubcmdBinder() fileContentBinder {
	lines := []string{"#!/usr/bin/env subcmd"}
	return fileContentBinder{lines}
}

type dirBinder struct { }

func (b dirBinder) makeFresh (path string) error {
	return os.Mkdir(path, 0755)
}

func (b dirBinder) assertMatch (path string, stat os.FileInfo) error {
	if !(stat.IsDir()) {
		return errors.New("File exists at directory targer: " + path)
	}
	return nil
}

