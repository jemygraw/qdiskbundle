package qdisksync

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

//read from snapshot and copy
func SyncVolumeData(srcVolume string, destVolume string, bufferSize int64, workerCount int32) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cacheFile := "qdisksync.cache"
	cacheFileH, err := os.Open(cacheFile)
	if err != nil {
		L.Error("Load cache file `%s' failed", cacheFile)
		return
	}
	defer cacheFileH.Close()

	bReader := bufio.NewScanner(cacheFileH)
	bReader.Split(bufio.ScanLines)
	//init channel
	var allWorkers int32 = 0
	syncStart := time.Now()
	L.Informational("Sync `%s' -> `'%s start from `%s'", srcVolume, destVolume, syncStart.String())
	for bReader.Scan() {
		line := bReader.Text()
		//split to name and size
		items := strings.Split(line, "\t")
		if len(items) != 3 {
			L.Error("Line data `%s'' error of cache file `%s'", line, cacheFile)
			continue
		}
		fname := items[0]
		fsize, err := strconv.ParseInt(items[1], 10, 64)
		if err != nil {
			L.Error("File length error `%s' for line `%s'", items[1], line)
			continue
		}
		fperm, err := strconv.ParseInt(items[2], 10, 64)
		if err != nil {
			L.Error("File perm error `%s' for line `%s'", items[2], line)
			continue
		}
		//join the path
		srcFullPath := filepath.Join(srcVolume, fname)
		destFullPath := filepath.Join(destVolume, fname)
		//check src and dest file
		srcFileH, srcErr := os.Open(srcFullPath)
		if srcErr != nil {
			L.Error("Open src file `%s' error `%s'", srcFullPath, srcErr.Error())
			continue
		}
		//create path if necessary
		lastSlashIndex := strings.LastIndex(destFullPath, "/")
		destFullPathBase := destFullPath[:lastSlashIndex]
		if err := os.MkdirAll(destFullPathBase, 0775); err != nil {
			L.Error("Failed to create dir `%s' due to error `%s'", destFullPathBase, err.Error())
			continue
		}
		destFileH, destErr := os.OpenFile(destFullPath, os.O_CREATE|os.O_RDWR, os.FileMode(fperm))
		if destErr != nil {
			L.Error("Open dest file `%s' error `%s'", destFullPath, destErr.Error())
			continue
		}

		//check whether it's time to run copy
		for {
			curWorkers := atomic.LoadInt32(&allWorkers)
			L.Debug("Current Workers: `%d'", curWorkers)
			if curWorkers < workerCount {
				atomic.AddInt32(&allWorkers, 1)
				go copy(srcFileH, destFileH, fsize, bufferSize, srcFullPath, destFullPath, &allWorkers)
				break
			} else {
				//wait a time to avoid infinite cycle
				<-time.After(time.Microsecond * 1)
			}
		}
	}
	syncEnd := time.Now()
	L.Informational("Sync `%s' -> `%s' end at `%s'", srcVolume, destVolume, syncEnd.String())
	L.Informational("Sync `%s' -> `%s' end at `%s'", srcVolume, destVolume, time.Since(syncStart))
}

func copy(srcFileH, destFileH *os.File, fsize int64, bufferSize int64, srcFullPath, destFullPath string, allWorkers *int32) {
	defer func() {
		atomic.AddInt32(allWorkers, -1)
		runtime.Gosched()
	}()
	L.Debug("Copying from `%s' to `%s'", srcFullPath, destFullPath)
	buffer := make([]byte, bufferSize)
	var cpErr error
	var cpNum int64
	for {
		numRead, errRead := srcFileH.Read(buffer)
		if errRead == io.EOF {
			break
		} else {
			if errRead != nil {
				L.Error("Read from `%s' error: `%s'", srcFullPath, errRead.Error())
				cpErr = errRead
				break
			} else {
				numWrite, errWrite := destFileH.Write(buffer[:numRead])
				if errWrite != nil {
					L.Error("Write to `%s' error: `%s'", destFullPath, errWrite.Error())
					cpErr = errWrite
					break
				}
				cpNum += int64(numWrite)
			}
		}
	}
	defer srcFileH.Close()
	defer destFileH.Close()
	if cpErr != nil || cpNum != fsize {
		L.Error("Copy from `%s' to `%s' failed, error: `%s'", srcFullPath, destFullPath, cpErr.Error())
	} else {
		L.Debug("Copy from `%s' to `%s' succcess", srcFullPath, destFullPath)
	}

}
