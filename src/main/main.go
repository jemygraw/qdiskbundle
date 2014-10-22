package main

import (
	"fmt"
	"os"
	qds "qdisksync"
)

var conf *qds.Conf

func initConfig() {
	var err error
	qds.InitLogs(`{"filename":"qdisksync.log"}`, true)
	confFile := "qdisksync.conf"
	conf, err = qds.LoadConfig(confFile)
	if err != nil {
		qds.L.Error("Failed to load config file `%s'", confFile)
		os.Exit(2)
	}

	qds.L.Informational("Src Volume: `%s'", conf.SrcVolume)
	qds.L.Informational("Dest Volume: `%s'", conf.DestVolume)
	qds.L.Informational("Buffer Size: `%d' bytes", conf.BufferSize)
	qds.L.Informational("Worker Count: `%d'", conf.WorkerCount)
}

func help() {
	var helpDoc = `QDiskSync
	
Usage:	
	Sync the data between the volumes

Commands:
	qdisksync cache - make a new snapshot of the tree of the volume
	qdisksync sync - start to sync the data by the snapshot

Build:
	v1.0.0
`
	fmt.Println(helpDoc)
}
func main() {
	cmdArgs := os.Args
	if len(cmdArgs) != 2 {
		help()
		return
	}
	initConfig()
	cmdName := cmdArgs[1]
	switch cmdName {
	case "cache":
		//cache volume tree
		qds.CacheVolumeTree(conf.SrcVolume)
	case "sync":
		//sync data by snapshot
		qds.SyncVolumeData(conf.SrcVolume, conf.DestVolume, conf.BufferSize, conf.WorkerCount)
	default:
		qds.L.Error("Unknow command `%s'", cmdName)
	}
	qds.L.Close()
}
