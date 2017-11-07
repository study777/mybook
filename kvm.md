virt-install --virt-type kvm  --arch=x86_64 --name win2008r2 --ram  8192 --vcpus=4
--cdrom  /kvm/iso/win2008r2.iso
--noautoconsole --os-type windows --os-variant win2k8 
--disk path=/home/cloud03/virtio-win-0.1-30.iso,device=cdrom,perms=ro 
--disk path=/var/lib/libvirt/images/win2k8.img,size=50 
--graphics vnc,password=foobar,listen=0.0.0.0,keymap=en-us
