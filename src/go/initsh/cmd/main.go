package main

import "os"
import "log"
import "fmt"
import "initsh"
import "bufio"

func main() {
	importRoot := os.Args[1]
	directs := os.Args[2:]
	treePaths := scanPaths()

	lobby := initsh.MakeLobby(importRoot, "::")
	for _, directive :=  range directs {
		path, err := initsh.ImportNeeded(directive, lobby, treePaths)
		if err != nil {
			log.Fatal("Error on Directive (\"" + directive +"\"): ", err)
		} else {
			fmt.Println(path)
		}
	}
}

func scanPaths() []string {
	var treePaths []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		treePaths = append(treePaths, scanner.Text())
	}
	return treePaths
}
