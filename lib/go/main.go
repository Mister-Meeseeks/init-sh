package main

import "os"
import "log"
import "initsh"

func main() {
	binDir := os.Args[1]
	libDir := os.Args[2]
	directs := os.Args[3:]

	importer := initsh.importDirector{binDir, libDir, "::"}
	for _, directive :=  range directs {
		err := initsh.walkThru(directive, importer)
		if err != nil {
			log.Fatal("Error on parse=" + directive, err)
		}
	}
}
