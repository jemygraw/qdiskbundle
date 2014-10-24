package qdisksync

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CacheVolumeTree(volume string) {
	cacheFile := "qdisksync.cache"
	if _, err := os.Stat(cacheFile); err != nil {
		L.Informational("No cache file `%s' found, will create one", cacheFile)
	} else {
		if rErr := os.Rename(cacheFile, cacheFile+".old"); rErr != nil {
			L.Error("Unable to rename cache file, plz manually delete `%s' and `%s.old'",
				cacheFile, cacheFile)
			return
		}
	}
	cacheFileH, err := os.OpenFile(cacheFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		L.Error("Failed to open cache file `%s'", cacheFile)
		return
	}
	defer cacheFileH.Close()
	bWriter := bufio.NewWriter(cacheFileH)
	walkStart := time.Now()
	L.Informational("Walk `%s' start from `%s'", volume, walkStart.String())
	filepath.Walk(volume, func(path string, fi os.FileInfo, err error) error {
		var retErr error
		L.Debug("Walking through `%s'", volume)
		if !fi.IsDir() {
			relPath := strings.TrimPrefix(strings.TrimPrefix(path, volume), "/")
			fsize := fi.Size()
			fperm := fi.Mode().Perm()
			L.Debug("Hit file `%s' size: `%d' perm: `%s'", relPath, fsize, fperm)
			fmeta := fmt.Sprintln(fmt.Sprintf("%s\t%d\t%d", relPath, fsize, fperm))
			if _, err := bWriter.WriteString(fmeta); err != nil {
				L.Error("Failed to write data `%s' to cache file", fmeta)
				retErr = err
			}
		}
		return retErr
	})
	if err := bWriter.Flush(); err != nil {
		L.Error("Failed to flush to cache file `%s'", cacheFile)
	}

	walkEnd := time.Now()
	L.Informational("Walk `%s' end at `%s'", volume, walkEnd.String())
	L.Informational("Walk `%s' last for `%s'", volume, time.Since(walkStart))
}
