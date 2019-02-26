
SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

libPath = $(SELF_DIR)lib/
importBin = $(libPath)import
goCmdDir = $(SELF_DIR)src/go/initsh/cmd/

$(importBin): $(libPath) $(goCmdDir)
	goto -o $(importBin) $(goCmdDir)

$(libPath):
	mkdir -p $(libPath)

clean:
	rm -r $(libPath)
