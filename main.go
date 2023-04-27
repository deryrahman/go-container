//go:build linux
// +build linux

package main

import (
	"fmt"
	"os"
	"syscall"
)

// go main.go run /bin/bash <args>
func main() {
	switch os.Args[1] {
	case "run":
		runCommand()
	default:
		panic("bad command")
	}
}

func runCommand() {
	fmt.Println("Current Process:", os.Getpid())

	// begin to fork the process
	childPid, err := syscall.ForkExec(os.Args[2], os.Args[2:], &syscall.ProcAttr{
		Files: []uintptr{
			os.Stdin.Fd(),
			os.Stdout.Fd(),
			os.Stderr.Fd(),
		},
	})
	must(err)
	fmt.Println("Child Process:", childPid)
	syscall.Wait4(childPid, nil, 0, nil) // wait for the termination of the container
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
