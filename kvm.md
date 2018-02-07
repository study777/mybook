virt-install --virt-type kvm  --arch=x86_64 --name win2008r2 --ram  8192 --vcpus=4
--cdrom  /kvm/iso/win2008r2.iso
--noautoconsole --os-type windows --os-variant win2k8 
--disk path=/home/cloud03/virtio-win-0.1-30.iso,device=cdrom,perms=ro 
--disk path=/var/lib/libvirt/images/win2k8.img,size=50 
--graphics vnc,password=foobar,listen=0.0.0.0,keymap=en-us



success

export LANG=en_US.UTF-8

qemu-img create -f qcow2  /local_kvm/test.qcow2  80G


virt-install --virt-type kvm --name CentOS-7-x86_64 --ram 2048 --cdrom=/local_kvm/CentOS-7-x86_64-DVD-1611.iso  --disk path=/local_kvm/test.qcow2  --network network=br0  --graphics vnc,keymap=en-us,listen=0.0.0.0


example1

virt-install --virt-type kvm --name CentOS-7-x86_64 --ram 2048 --cdrom=/home/kvm/iso/CentOS-6.5-x86_64-bin-DVD1.iso  --disk path=/home/kvm/image/test.qcow2  --network network=br0 --graphics vnc,keymap=en-us,port=5901,listen=0.0.0.0


example2

qemu-img create -f qcow2  /home/kvm/image/10.10.10.167.qcow2 150G

virt-install --virt-type kvm --hvm --name 10.10.10.167  --ram 4096 --vcpus=4 --os-type=linux  --cdrom=/home/kvm/iso/CentOS-6.5-x86_64-bin-DVD1.iso --disk path=/home/kvm/image/10.10.10.167.qcow2,bus=virtio --network network=br0 --graphics vnc,keymap=en-us,port=5901,listen=0.0.0.0  --accelerate

直接从装好系统的磁盘引导，跳过系统安装

 cp centos6.4-80G.qcow2  test.qcow2

 virt-install --virt-type kvm --hvm --name test --ram 4096 --vcpus=4 --os-type=linux   --import    --disk path=/kvm/image/test.qcow2,bus=virtio --network network=br0 --graphics vnc,keymap=en-us,port=5901,listen=0.0.0.0 --accelerate



拷贝本地文件到磁盘镜像

cp centos6.4-80G.qcow2  test1.qcow2

virt-copy-in  -a test1.qcow2    /tmp/test  /tmp/

安装之前指定虚拟机ip

cp centos6.4-80G.qcow2  test2.qcow2

virt-copy-in -a  /kvm/image/test2.qcow2 /tmp/ifcfg-eth0   /etc/sysconfig/network-scripts/

virt-install --virt-type kvm --hvm --name test2 --ram 4096 --vcpus=4 --os-type=linux --import --disk path=/kvm/image/test2.qcow2,bus=virtio --network network=br0 --graphics vnc,keymap=en-us,port=5903,listen=0.0.0.0 --accelerate
