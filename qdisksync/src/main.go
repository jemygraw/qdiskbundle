package main

import (
	"disksync"
	"flag"
	"fmt"
)

func main() {
	var worker int
	var privateKey string
	var user string
	var host string
	var srcFileList string
	var destPath string
	var debugMode bool

	flag.IntVar(&worker, "worker", 1, "sync worker count")
	flag.StringVar(&privateKey, "key", "", "ssh private key with no password")
	flag.StringVar(&user, "user", "", "ssh login user")
	flag.StringVar(&host, "host", "", "ssh login host")
	flag.StringVar(&srcFileList, "file", "", "file list to sync")
	flag.StringVar(&destPath, "dest", "", "sync destination path")
	flag.BoolVar(&debugMode, "debug", false, "debug mode")

	flag.Parse()

	if privateKey == "" {
		fmt.Println("Error: no private key")
		return
	}

	if user == "" {
		fmt.Println("Error: no login user")
		return
	}

	if host == "" {
		fmt.Println("Error: no remote host")
		return
	}

	if srcFileList == "" {
		fmt.Println("Error: no src file list")
		return
	}

	if destPath == "" {
		fmt.Println("Error: no destination path")
		return
	}

	disksync.Sync(privateKey, user, host, srcFileList, destPath, worker, debugMode)

}
