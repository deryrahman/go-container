# go-container

Simple container implementation in Go.

## Motivation
As we see many, containers are defacto for shipping the services on heterogenous complex system. Understand the flow of container creation would definitely help us to understand what happen under the hood. So, debugging on any container related issue would be easier.

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
- Chroot: operation to change the root folder for existing process
- Cgroups: limiting the resources for the existing process

Important syscalls to create the process for initialization of the container:
- `fork()`: creating a copy of process from current one
- `exec()`: replacing the process's program with new program

Namespaces we use for this exercise:
- `CLONE_NEWPID`: to create new process id on the forked process
- `CLONE_NEWUTS`: to change the hostname
- `CLONE_NEWUSER`: to create new user which can be acted as root(?)
- `CLONE_NEWNS`:

## Reference
https://www.youtube.com/watch?v=8fi7uSYlOdc
https://github.com/nathanagez/c_container/blob/master/container/src/main.c
https://www.toptal.com/linux/separation-anxiety-isolating-your-system-with-linux-namespaces