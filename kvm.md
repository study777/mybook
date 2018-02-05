virt-install --virt-type kvm  --arch=x86_64 --name win2008r2 --ram  8192 --vcpus=4
--cdrom  /kvm/iso/win2008r2.iso
--noautoconsole --os-type windows --os-variant win2k8 
--disk path=/home/cloud03/virtio-win-0.1-30.iso,device=cdrom,perms=ro 
--disk path=/var/lib/libvirt/images/win2k8.img,size=50 
--graphics vnc,password=foobar,listen=0.0.0.0,keymap=en-us



success

virt-install --virt-type kvm --name CentOS-7-x86_64 --ram 2048 --cdrom=/local_kvm/CentOS-7-x86_64-DVD-1611.iso  --disk path=/local_kvm/test.qcow2  --network network=br0  --graphics vnc,keymap=en-us,listen=0.0.0.0
