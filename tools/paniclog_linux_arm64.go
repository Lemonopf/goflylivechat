// Log the panic under unix to the log file

//go:build linux && arm64
// +build linux,arm64

package tools

import (
	"log"
	"os"
	"syscall"
)

// redirectStderr to the file passed in
// arm64 的 Linux 没有 Dup2 系统调用，使用 Dup3 等价替代
func RedirectStderr(f *os.File) {
	err := syscall.Dup3(int(f.Fd()), int(os.Stderr.Fd()), 0)
	if err != nil {
		log.Printf("Failed to redirect stderr to file: %v", err)
	}
}
