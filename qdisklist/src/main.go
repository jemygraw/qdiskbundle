package main

import (
	"disklist"
	"flag"
	"fmt"
	"path/filepath"
)

func main() {
	var listDir string
	var prefix string
	var result string

	flag.StringVar(&listDir, "dir", "", "list dir")
	flag.StringVar(&prefix, "prefix", "", "key prefix")
	flag.StringVar(&result, "result", "", "list result file")
	flag.Parse()

	listAbsDir, cErr := filepath.Abs(listDir)
	if cErr != nil {
		fmt.Println("Error: invalid list dir")
		return
	}

	disklist.ListDir(listAbsDir, prefix, result)

}
