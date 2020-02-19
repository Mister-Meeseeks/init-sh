package main

import "os"
import "log"
import "initsh"
import "bufio"

func main() {
	binDir := os.Args[1]
	libDir := os.Args[2]
	directs := os.Args[3:]

	var treePaths []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		treePaths = append(treePaths, scanner.Text())
	}

	importer := initsh.MakeImporter(binDir, libDir, "::")
	for _, directive :=  range directs {
		err := initsh.WalkPreScanned(directive, importer, treePaths)
		if err != nil {
			log.Fatal("Error on Directive (\"" + directive +"\"): ", err)
		}
	}
}
