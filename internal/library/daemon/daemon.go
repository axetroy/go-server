package daemon

import (
	"fmt"
	"github.com/axetroy/go-fs"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"syscall"
)

type Action func() error

func getPidFilePath() (string, error) {
	cwd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	executableFilePath, err := os.Executable()

	executableName := path.Base(executableFilePath)

	if err != nil {
		return "", err
	}

	return path.Join(cwd, executableName) + ".pid", nil
}

func Start(action Action, shouldRunInDaemon bool) error {
	if shouldRunInDaemon && os.Getppid() != 1 {
		// 将命令行参数中执行文件路径转换成可用路径
		filePath, _ := filepath.Abs(os.Args[0])
		cmd := exec.Command(filePath, os.Args[1:]...)
		// 将其他命令传入生成出的进程
		// cmd.Stdin = os.Stdin // 给新进程设置文件描述符，可以重定向到文件中
		//cmd.Stdout = ioutil.Discard
		//cmd.Stderr = ioutil.Discard
		err := cmd.Start() // 开始执行新进程，不等待新进程退出
		return err
	} else {
		pidFilePath, err := getPidFilePath()

		if err != nil {
			return err
		}

		if err := fs.WriteFile(pidFilePath, []byte(fmt.Sprintf("%d", os.Getpid()))); err != nil {
			return err
		}

		return action()
	}
}

func Stop() error {
	pidFilePath, err1 := getPidFilePath()

	if err1 != nil {
		return err1
	}

	if !fs.PathExists(pidFilePath) {
		return nil
	}

	b, err2 := fs.ReadFile(pidFilePath)

	if err2 != nil {
		return nil
	}

	pidStr := string(b)

	pid, err3 := strconv.Atoi(pidStr)

	if err3 != nil {
		return err3
	}

	ps, err4 := os.FindProcess(pid)

	if err4 != nil {
		return err4
	}

	if err5 := ps.Signal(syscall.SIGKILL); err5 != nil {
		return err5
	}

	psState, err6 := ps.Wait()

	if err6 != nil {
		return err6
	}

	haveBeenKill := psState.Exited()

	if haveBeenKill {
		log.Printf("进程 %d 已结束.\n", psState.Pid())

		_ = fs.Remove(pidFilePath)
	} else {
		log.Printf("进程 %d 结束失败.\n", psState.Pid())
	}

	return nil
}
