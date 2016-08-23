package disklist

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func ListDir(listAbsDir, prefix, listResultFile string) {
	listResultFileH, err := os.Create(listResultFile)
	if err != nil {
		fmt.Println("Error: failed to open cache file", listResultFile)
		return
	}
	defer listResultFileH.Close()
	bWriter := bufio.NewWriter(listResultFileH)

	filepath.Walk(listAbsDir, func(path string, fi os.FileInfo, err error) error {
		var retErr error

		if err != nil {
			retErr = err
		} else {
			if !fi.IsDir() {
				relPath := strings.TrimPrefix(strings.TrimPrefix(path, listAbsDir), string(os.PathSeparator))
				fileKey := fmt.Sprintf("%s%s", prefix, relPath)
				if runtime.GOOS == "windows" {
					fileKey = strings.Replace(fileKey, "\\", "/", -1)
				}
				fsize := fi.Size()
				fmeta := fmt.Sprintln(fmt.Sprintf("%s\t%s\t%d", path, fileKey, fsize))

				if _, err := bWriter.WriteString(fmeta); err != nil {
					fmt.Println("Error: Error: failed to write data to result file", err)
				}
			}
		}
		return retErr
	})

	if err := bWriter.Flush(); err != nil {
		fmt.Print("Error: failed to flush to cache file", listResultFile)
	}

	return
}
