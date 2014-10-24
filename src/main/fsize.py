fp=open("qdisksync.cache")
fsize_all=0
for line in fp:
	items=line.split("\t")
	fsize_all+=int(items[1])
fp.close()
print("Size:"+str(fsize_all/1024.0/1024.0)+"MB")
