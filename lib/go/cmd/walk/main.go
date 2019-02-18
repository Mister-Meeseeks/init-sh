package main

import "os"
import "path/filepath"
import "log"
import "fmt"

func main() {
	paths := os.Args[1:]
	for _, path :=  range paths {
		canon, err := filepath.EvalSymlinks(path)
		if err != nil {
			log.Fatal("hello", err)
		}
		fmt.Printf("Resolve symlink %q -> %q\n", path, canon)
	}
}

