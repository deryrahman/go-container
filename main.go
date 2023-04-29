//go:build linux
// +build linux

package main

import (
	"fmt"
	"os"
	"syscall"
)

// go main.go run /bin/sh <args>
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
			Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNS,
			UidMappings: []syscall.SysProcIDMap{{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			}},
		},
	})
	must(err)
	fmt.Println("Child Process:", childPid)
	syscall.Wait4(childPid, nil, 0, nil) // wait for the termination of the container
}

func childExec() {
	fmt.Println("Running", os.Args[2:])

	// setup all things between fork and exec
	must(syscall.Sethostname([]byte("container")))
	must(syscall.Chroot("./alpine-fs"))
	must(os.Chdir("/"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

	// exec the forked process with entrypoint
	must(syscall.Exec(os.Args[2], os.Args[2:], nil))
	must(syscall.Unmount("proc", 0))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
