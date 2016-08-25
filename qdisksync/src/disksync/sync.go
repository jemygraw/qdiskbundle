package disksync

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var initSyncOnce sync.Once
var syncTasks chan func()
var syncWaitGroup sync.WaitGroup

func doSync(tasks chan func()) {
	for {
		task := <-tasks
		task()
	}
}

func Sync(privateKey, user, host, srcFileList, destPath string, worker int, debugMode bool) {
	srcFp, openErr := os.Open(srcFileList)
	if openErr != nil {
		fmt.Println("Error: open src file list error", openErr)
		return
	}

	defer srcFp.Close()

	//init task
	initSyncOnce.Do(func() {
		syncTasks = make(chan func(), worker)
		for i := 0; i < worker; i++ {
			go doSync(syncTasks)
		}
	})

	syncWaitGroup = sync.WaitGroup{}

	var totalSyncCount int64
	var totalFailedCount int64
	var totalSuccessCount int64
	var syncStartTime = time.Now()
	var totalSyncFSize int64
	var enableSize bool

	//iter each file to sync
	bScanner := bufio.NewScanner(srcFp)
	for bScanner.Scan() {
		line := bScanner.Text()
		items := strings.Split(line, "\t")
		if len(items) != 3 {
			fmt.Println("Error: please use disklist to get src file list")
			return
		}

		fileSrcPath := items[0]
		fileDestKey := items[1]
		fileSize, _ := strconv.ParseInt(items[2], 10, 64)

		relDir := filepath.Dir(fileDestKey)

		if relDir != "" && relDir != "." {
			mkDir := filepath.Join(destPath, relDir)
			if runtime.GOOS == "windows" {
				mkDir = strings.Replace(mkDir, "\\", "/", -1)
			}

			mkDir = escapeShell(mkDir)
			mkDirCmdStr := fmt.Sprintf("`mkdir -p %s`", mkDir)

			sshParams := []string{
				"-i", privateKey,
				fmt.Sprintf("%s@%s", user, host),
				mkDirCmdStr,
			}
			mkdirCmd := exec.Command("ssh", sshParams...)
			if debugMode {
				fmt.Println("Debug:", mkdirCmd)
			}
			runErr := mkdirCmd.Run()
			if runErr != nil {
				fmt.Println("Error: mkdir failed", mkDir, runErr)
			}
		}

		// start to copy files
		syncWaitGroup.Add(1)
		atomic.AddInt64(&totalSyncCount, 1)
		syncTasks <- func() {
			defer syncWaitGroup.Done()

			//sync logic
			fileDestPath := filepath.Join(destPath, fileDestKey)

			//check for windows, use cgydrive

			if runtime.GOOS == "windows" {
				//check file src path
				driveIndex := strings.Index(fileSrcPath, ":")
				if driveIndex == -1 {
					fmt.Println("Error: invalid line", line)
					return
				}
				driveName := fileSrcPath[:driveIndex]
				fileSrcPath = strings.Replace(fileSrcPath, "\\", "/", -1)
				driveRight := fileSrcPath[driveIndex+1:]
				fileSrcPath = fmt.Sprintf("/cygdrive/%s%s", driveName, driveRight)

				//check file dest path
				fileDestPath = strings.Replace(fileDestPath, "\\", "/", -1)
			}

			//escape shell chars
			fileSrcPath = escapeShell(fileSrcPath)
			fileDestPath = escapeShell(fileDestPath)

			var scpDest string
			if runtime.GOOS == "windows" {
				scpDest = fmt.Sprintf("%s@%s:'%s'", user, host, fileDestPath)
			} else {
				scpDest = fmt.Sprintf("%s@%s:%s", user, host, fileDestPath)
			}

			//form scp cmd
			scpParams := []string{
				"-i", privateKey,
				fileSrcPath,
				scpDest,
			}
			syncCmd := exec.Command("scp", scpParams...)
			if debugMode {
				fmt.Println("Debug:", syncCmd)
			}
			runErr := syncCmd.Run()
			if runErr != nil {
				atomic.AddInt64(&totalFailedCount, 1)
				fmt.Println("Info: sync", fileSrcPath, "=>", fileDestPath, "failed,", runErr)
			} else {
				atomic.AddInt64(&totalSuccessCount, 1)
				atomic.AddInt64(&totalSyncFSize, fileSize)
				fmt.Println("Info: sync", fileSrcPath, "=>", fileDestPath, "success")
			}
		}
	}

	syncWaitGroup.Wait()

	//sum
	fmt.Println("----------------------")
	fmt.Println("Total:", totalSyncCount)
	fmt.Println("Success:", totalSuccessCount)
	fmt.Println("Failed:", totalFailedCount)
	if enableSize {
		fmt.Println("Amount:", totalSyncFSize, "Bytes")
	}
	fmt.Println("Duration", time.Since(syncStartTime))
	fmt.Println("----------------------")
}

func escapeShell(from string) (to string) {
	//escape some characters
	to = from
	escapeChars := []string{"(", ")"}
	for _, c := range escapeChars {
		to = strings.Replace(to, c, "\\"+c, -1)
	}
	return
}
