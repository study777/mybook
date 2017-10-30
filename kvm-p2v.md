
scp -r 11.4.13.249:/vmfs/volumes/datastore-1/win-11.4.4.45   ./



cd win-11.4.4.44/

查看esxi 磁盘信息

qemu-img  info win-11.4.4.44.vmdk

输出如下内容

```
image: win-11.4.4.44.vmdk
file format: vmdk
virtual size: 100G (107374182400 bytes)
disk size: 100G
Format specific information:
    cid: 2549577159
    parent cid: 4294967295
    create type: vmfs
    extents:
        [0]:
            virtual size: 107374182400
            filename: win-11.4.4.44-flat.vmdk
            format: VMFS

```


开始转换格式

qemu-img convert -f vmdk -O qcow2   win-11.4.4.44.vmdk    win-11.4.4.44.qcow2






