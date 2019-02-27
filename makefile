
PREFIX = /usr/local/

SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

libPath = $(SELF_DIR)lib/
importBin = $(libPath)import
shPath = $(SELF_DIR)src/sh/
goCmdDir = $(SELF_DIR)src/go/initsh/cmd/
binPath = $(SELF_DIR)/bin/

repoBins = $(wildcard $(binPath)*)
repoShells = $(wildcare $(shPath)*)

sysBinDir = $(PREFIX)bin/
shebangs = $(patsubst $(binPath)%,$(sysBinDir)%,$(repoBins))

.PHONY: clean install

$(importBin): $(libPath) $(goCmdDir)
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
