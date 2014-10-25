#qdisksync 文件同步工具

###用途  
本工具主要设计用来进行磁盘数据同步或备份，文件夹数据同步或备份等。

###用法
```
QDiskSync

Usage:
	Sync the data between the volumes

Commands:
	qdisksync cache cacheFile - Make a new snapshot of the tree of the volume
	qdisksync sync cacheFile - Start to sync the data by the snapshot

Build:
	v1.0.0
```

###配置
本工具的使用需要一个名称为`qdisksync.conf`的配置文件。配置文件内容类似下面：

```
{
	"src_volume" : "/Users/jemy/Projects/qdisksync/",
	"dest_volume" : "/Users/jemy/Test",
	"buffer_size" : 4096,
	"worker_count" : 10
}
```

该配置文件的内容为`json`格式。详细说明如下：

|  参数名称 | 参数说明 |   可选   |
|----------|---------|---------|
| src_volume| 同步源路径，可以为文件夹路径，也可以为磁盘路径，即同步整个磁盘。该路径为全路径。|必填|
| dest_volume| 同步目标路径，可以为文件夹路径，也可以为磁盘路径。该路径为全路径。|必填|
| buffer_size| 拷贝文件的缓冲区大小，如果不填，默认为4M，该参数单位为字节。|选填|
| worker_count| 并发拷贝文件的线程数量，如果不填，默认为1。|选填|

**注意⚠：如果你不需要填某一个参数，请从配置文件里面删除，不要留值为空。**


###原理
本工具首先需要使用`cache`命令对需要同步的源路径下面的文件夹和文件进行一次快照，快照的内容主要包括文件夹和文件的相对路径（即不含设置的源路径）以及文件的大小和权限。快照的内容保存在运行`cache`命令时所指定的快照文件中，比如`qdisksync.cache`文件中。

```
$./qdisksync cache qdisksync.cache
```

`qdisksync.cache`内容类似于：

```
LICENSE	10253	420  
README.md	374	420  
src/main/build.sh	105	493  
src/main/fsize.py	165	420  
src/main/qdisksync	5111724	493  
src/main/qdisksync.cache	4096	420  
src/main/qdisksync.cache.old	461990	420  
src/main/qdisksync.conf	152	420  
src/main/qdisksync.go	1181	420  
src/main/qdisksync.log	1443484	420  
src/qdisksync/scanner.go	1668	420  
src/qdisksync/settings.go	1697	420  
src/qdisksync/sync.go	4094	420  
```

**注意⚠：每次生成新的快照的同时，都会把`老的同名快照`加后缀`.old`重命名。比如这里的`qdisksync.cache`如果在下次快照时仍然使用这个文件名，那么老的快照文件就会被重命名为`qdisksync.cache.old`，如果你多次进行快照操作，那么旧有的一些快照内容不会被自动保留，最多保留最近两次的快照。**

然后我们使用`sync`命令来根据快照文件内容，比如这里的`qdisksync.cache`的内容来单向地将数据从源路径拷贝到目标路径。

**注意⚠：`sync`命令完全是根据所指定的快照内容，譬如这里的`qdisksync.cache`的内容来拷贝数据的。**

###增量同步
增量同步包括增量更新和增量删除，目前尚未实现。
