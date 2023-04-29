# go-container

Simple container implementation in Go. All of this explanation is not guarantee to be correct.

## Motivation
As we see many, containers are defacto for shipping the services on heterogenous complex system. Understand the flow of container creation would definitely help us to understand what happen under the hood. In consequences, debugging any container related issue would be easier.

## Goals
We will mimicking the docker command for creating the container:

`docker run <image_name> <command> <args>`

Omitting the `<image_name>` to simplify the process, the final command we have is equivalent to this:

`go-container run <command> <args>`

For this exercise I use alpine mini root file system as an image, download [here](https://www.alpinelinux.org/downloads/) and extract to folder alpine-fs

## Basic
Container is just a process which live in some restriction. The capability of container is somewhat similar with what OS can do but without provisioning entire machine virtually.

There're 3 basic minimal things to achive this:
- Namespaces: what things the current process can see
- Chroot: operation to change the root folder for existing process, limiting the current process working directory
- Cgroups: limiting the resources for the existing process

Important syscalls to create the process for initialization of the container:
- `fork()`: creating a copy of process from current one
- `exec()`: replacing the process's program with new program

Namespaces we use for this exercise:
- `CLONE_NEWPID`: to create new process id on new namespace, different from the parent namespace (since Linux 2.6.24, linux provides multiple IDs assigned to same process, it would be useful for namespacing the process) (pid_namespaces)
- `CLONE_NEWUTS`: to provide the isolation of two system identifier (hostname & domainname) (uts_namespaces)
- `CLONE_NEWUSER`: to make the new process living in the different user namespace, it will associated the process with the namespaced user/group (user_namespaces)
- `CLONE_NEWNS`: provide the isolation of the list of mounts which process can see (mount_namespaces)

## Quick Run
1. Provision the linux on virtual machine
2. Clone this repo
3. Download [alpine mini root fs](https://www.alpinelinux.org/downloads/) and extract to folder alpine-fs, so working directory tree would be like this:
```sh
❯ tree -L 1
.
├── README.md
├── alpine-fs
├── alpine-minirootfs-3.17.3-x86_64.tar.gz
├── cgroup.sh
└── main.go

2 directories, 4 files
```
4. 
```sh
sudo ./cgroup.sh # to initialize the cgroup
go run main.go run /bin/sh # run the container with /bin/sh as an entrypoint command
```

**To test out the cgroup is working**
Run command to create process > 10
```sh
for i in `seq 1 100`
do
  sleep 600 &
done
```
Run command to store memory > 50mb (TBD)
```sh
#random="$(dd if=/dev/urandom bs=1M count=50)"
```

## Recap

```md
+---+             +---+
| P | --fork()--> | P |
+---+     |       +---|
          |
          |    +---+             +---+
          +--> | C | --exec()--> | C'|
               +---+             +---+
```

- forking the process, bring the current file descriptor for stdin, stdout, stderr on the child process, and set the clone flags to setup the namespaces
- before executing the entrypoint `<command>`, some syscall operation should be happen (we can call this stage as init)
<TBD>

## Reference
https://man7.org/linux/man-pages/man2/clone.2.html
https://www.youtube.com/watch?v=8fi7uSYlOdc
https://github.com/nathanagez/c_container/blob/master/container/src/main.c
https://www.toptal.com/linux/separation-anxiety-isolating-your-system-with-linux-namespaces
https://ericchiang.github.io/post/containers-from-scratch/
https://theboreddev.com/understanding-linux-namespaces/