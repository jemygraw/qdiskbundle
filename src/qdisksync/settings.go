package qdisksync

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"os"
)

var L *logs.BeeLogger

const (
	DEFAULT_BUFFER_SIZE  = 1 << 22
	DEFAULT_WORKER_COUNT = 1
)

type Conf struct {
	SrcVolume   string `json:"src_volume"`
	DestVolume  string `json:"dest_volume"`
	BufferSize  int64  `json:"buffer_size,omitempty"`
	WorkerCount int32  `json:"worker_count,omitempty"`
}

func InitLogs(jsonConfig string, debugMode bool) {
	L = logs.NewLogger(1)
	L.SetLevel(logs.LevelDebug)
	L.SetLogger("file", jsonConfig)
	if debugMode {
		L.SetLogger("console", jsonConfig)
	}
}

func LoadConfig(fname string) (conf *Conf, retErr error) {
	fh, err := os.Open(fname)
	if err != nil {
		L.Error("Can not open config file `%s'", fname)
		retErr = err
		return
	}
	defer fh.Close()
	cnfData, err := ioutil.ReadAll(fh)
	if err != nil {
		L.Error("Can not read from config file")
		retErr = err
		return
	}
	conf = new(Conf)
	if err := json.Unmarshal(cnfData, conf); err != nil {
		L.Error("Can not unmarshal the json config")
		retErr = err
		return
	}
	if err = checkVolumeValid(conf.SrcVolume); err != nil {
		retErr = err
		return
	}
	if err = checkVolumeValid(conf.DestVolume); err != nil {
		retErr = err
		return
	}
	if conf.BufferSize == 0 {
		conf.BufferSize = DEFAULT_BUFFER_SIZE
	}
	if conf.WorkerCount == 0 {
		conf.WorkerCount = DEFAULT_WORKER_COUNT
	}
	return
}

func checkVolumeValid(volume string) (retErr error) {
	if fi, err := os.Stat(volume); err != nil {
		L.Error("Can not stat volume `%s'", volume)
		retErr = err
	} else {
		if !fi.IsDir() {
			msg := fmt.Sprintln("Volume `%s' is not a valid path", volume)
			L.Error(msg)
			retErr = errors.New(msg)
		}
	}
	return
}
