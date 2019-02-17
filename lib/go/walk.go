package main

import (
       "fmt"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("On Unix:")
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		fmt.Printf("Descent path %q - base %q - isDir= %q - Mode %d\n", path, info.Name(), info.IsDir(), info.Mode())
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %v\n", err)
		return
	}
}
