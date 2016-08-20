package disksync

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

func Sync(privateKey, user, host, srcFileList, destPath string, worker int) {
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
			sshParams := []string{
				"-t",
				"-i", privateKey,
				fmt.Sprintf("%s@%s", user, host),
				fmt.Sprintf("`mkdir -p %s`", mkDir),
			}
			mkdirCmd := exec.Command("ssh", sshParams...)
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
			fileDestPath := fmt.Sprintf("%s@%s:%s", user, host, filepath.Join(destPath, fileDestKey))
			scpParams := []string{
				"-i", privateKey,
				fileSrcPath,
				fileDestPath,
			}
			syncCmd := exec.Command("scp", scpParams...)
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
