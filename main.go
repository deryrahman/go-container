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
	fmt.Println("Current Process:", os.Getpid())
	switch os.Args[1] {
	case "run":
		runCommand()
	case "child":
		childExec()
	default:
		panic("bad command")
	}
}

func runCommand() {
	args := []string{os.Args[0], "child"}
	args = append(args, os.Args[2:]...)

	// begin to fork the process
	childPid, err := syscall.ForkExec(args[0], args, &syscall.ProcAttr{
		Files: []uintptr{
			os.Stdin.Fd(),
			os.Stdout.Fd(),
			os.Stderr.Fd(),
		},
		Sys: &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWUSER | syscall.CLONE_NEWUTS,
			UidMappings: []syscall.SysProcIDMap{{
				ContainerID: 0,
				HostID:      1000,
				Size:        1,
			}},
		},
	})
	must(err)
	fmt.Println("Child Process:", childPid)
	syscall.Wait4(childPid, nil, 0, nil) // wait for the termination of the container
}

func childExec() {
	// exec the forked process with entrypoint
	fmt.Println("Running", os.Args[2:])
	must(syscall.Sethostname([]byte("container")))
	must(syscall.Exec(os.Args[2], os.Args[2:], nil))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
