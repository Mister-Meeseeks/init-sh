
PREFIX = /usr/local/

SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

libPath = $(SELF_DIR)lib/
importBin = $(libPath)import
shPath = $(SELF_DIR)src/sh/
goSrcDir = $(SELF_DIR)src/go/initsh/
goCmdDir = $(goSrcDir)cmd/
binPath = $(SELF_DIR)/bin/

repoGoSrcs = $(wildcard $(goSrcDir)*.go $(goCmdDir)*.go)
repoBins = $(wildcard $(binPath)*)
repoShells = $(wildcard $(shPath)*)

sysBinDir = $(PREFIX)bin/
shebangs = $(patsubst $(binPath)%,$(sysBinDir)%,$(repoBins))

.PHONY: clean install

$(importBin): $(libPath) $(goCmdDir) $(repoGoSrcs)
	goto -o $(importBin) $(goCmdDir)

$(libPath):
	mkdir -p $(libPath)

install: $(shebangs)

$(sysBinDir)%: $(binPath)%
	test -e $@ && rm $@ || true
	test -L $@ && unlink $@ || true
	ln -s ${CURDIR}/$< $@

clean:
	rm -r $(libPath)
