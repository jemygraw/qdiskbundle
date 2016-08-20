# 简介

这是一个利用 `scp` 工具来进行数据同步的工具，其主要特点是对 `ssh` 和 `scp` 工具进行封装以支持在同步的目标远程机器上面创建和源机器上文件相同的相对路径。

基本原理如下：

假设源机器目录上面有路径 /mnt/src 这个目录下面有如下的文件结构：

```
.
├── apple
│   ├── hello.txt
│   └── world.txt
├── hello.txt
└── world.txt
```

我们现在需要将这个目录下的文件全部拷贝到远程机器 `/mnt/dest`  下面，当然如果是整个目录都进行拷贝，我们很容易利用命令:

```
scp -i ~/.ssh/nopass_rsa -r /mnt/src root@remote-host-ip:/mnt/dest
```

这样就可以把 `/mnt/src` 整个目录到拷贝到目的机器的路径 `/mnt/dest` 下面。

```
└── src
    ├── apple
    │   ├── hello.txt
    │   └── world.txt
    ├── hello.txt
    ├── list.txt
    └── world.txt
```

不过，在有些场景之下，这个方法就没有用了，而我们所说的场景就是这个工具所要解决的问题。

# 场景

我们需要把源机器上面的数据（100TB级别以上）拷贝到目标机器上去，但是有一个限制条件是目标机器上挂载的磁盘是彼此独立的，不能做Raid处理，每块盘最大容量为标准的 4TB（实际容量可能稍微小一点）。


# 使用

```
Usage of ./qdisksync:
  -file string
      file list to sync
  -dest string
    	sync destination path
  -key string
    	ssh private key with no password
  -user string
    	ssh login user
  -host string
		ssh login host
  -worker int
    	sync worker count (default 1)
```