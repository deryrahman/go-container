#!/bin/bash
# this script is intended to create a cgroup in current machine which will be used for the container cgroup

cgname="demo"
pidmax="20"
memmax="50000000" # 50MB
mkdir -p /sys/fs/cgroup/pids/${cgname}
mkdir -p /sys/fs/cgroup/memory/${cgname}
chmod 666 /sys/fs/cgroup/pids/${cgname}/cgroup.procs
chmod 666 /sys/fs/cgroup/memory/${cgname}/cgroup.procs
echo $pidmax > /sys/fs/cgroup/pids/${cgname}/pids.max
echo $memmax > /sys/fs/cgroup/memory/${cgname}/memory.limit_in_bytes