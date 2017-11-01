# openshift origin 3.6 部署 文档

# 系统环境：

CentOS Linux release 7.3.1611

主机配置：

内存8G  硬盘100G

# 系统安装

集群个节点 安装如下


时区 中国上海

分区
 
选择磁盘 选中 i will configure partitioning 选项  

下一界面 点击；

Click here to create them automatiaclly

删除home 分区

最终分区情况如下

/boot 1024 MiB

/     50GiB

swap  8064 MiB



最小化安装


配置本地yum

配置网卡
删除uuid mac地址等信息，以便可以方便的进行克隆

ls | grep -v  CentOS-Base.repo | xargs  rm  -rf

vi  /etc/yum.repos.d/CentOS-Base.repo

```
# CentOS-Base.repo
#
# The mirror system uses the connecting IP address of the client and the
# update status of each mirror to pick mirrors that are updated to and
# geographically close to the client.  You should use this for CentOS updates
# unless you are manually picking other mirrors.
#
# If the mirrorlist= does not work for you, as a fall back you can try the 
# remarked out baseurl= line instead.
#
#

[base]
name=CentOS-$releasever - Base
baseurl=http://172.16.2.100/base
gpgcheck=1
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-CentOS-7

#released updates 
[updates]
name=CentOS-$releasever - Updates
baseurl=http://172.16.2.100/updates
gpgcheck=1
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-CentOS-7

#additional packages that may be useful
[extras]
name=CentOS-$releasever - Extras
baseurl=http://172.16.2.100/extras
gpgcheck=1
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-CentOS-7

#additional packages that extend functionality of existing packages
[centosplus]
name=CentOS-$releasever - Plus
mirrorlist=http://mirrorlist.centos.org/?release=$releasever&arch=$basearch&repo=centosplus&infra=$infra
#baseurl=http://mirror.centos.org/centos/$releasever/centosplus/$basearch/
gpgcheck=1
enabled=0
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-CentOS-7
```

yum clean all

yum install wget git net-tools bind-utils iptables-services bridge-utils bash-completion kexec-tools sos psacct

yum update


yum install docker-1.12.6

yum -y install ansible pyOpenSSL


配置openshift yum 源

```
vi /etc/yum.repos.d/openshift.repo

[etcd]
name=etcd
baseurl=http://172.16.2.21/yum/etcd
gpgcheck=0
enable=1

[op3.6]
name=op3.6
baseurl=http://172.16.2.21/yum/op3.6
gpgcheck=0
enable=1

[commrpm]
name=commrpm
baseurl=http://172.16.2.21/yum/comm
gpgcheck=0
enable=1
```



事先安装一些依赖：


yum   install cockpit-bridge cockpit-docker cockpit-system cockpit-ws  ceph-common flannel


设置docker 镜像仓库

sed -i '/OPTIONS=.*/c\OPTIONS="--selinux-enabled --insecure-registry 172.30.0.0/16"'   /etc/sysconfig/docker


设置docker  存储

fdisk  /dev/vda
```
Command (m for help): n
Partition type:
   p   primary (2 primary, 0 extended, 2 free)
   e   extended
Select (default p): p
Partition number (3,4, default 3):回车
First sector (123486208-209715199, default 123486208): 
Using default value 123486208
Last sector, +sectors or +size{K,M,G} (123486208-209715199, default 209715199): +10G
Partition 3 of type Linux and of size 10 GiB is set

Command (m for help): w
```

partprobe


pvcreate /dev/vda3

vgcreate c2 /dev/vda3

systemctl stop docker

vim   /usr/bin/docker-storage-setup
```
VG=c2
SETUP_LVM_THIN_POOL=yes
```


lvmconf   --disable-cluster

docker-storage-setup

systemctl  start docker

systemctl  enable docker





# 为glusterfs 配置存储块设备

创建第一个块设备

fdisk  /dev/vda

Command (m for help): n

Select (default e): e

Selected partition 4

First sector (144457728-209715199, default 144457728): 回车

Last sector, +sectors or +size{K,M,G} (144457728-209715199, default 209715199): +15G

Using default value 209715199

Partition 4 of type Linux and of size 31.1 GiB is set


Command (m for help): w


partprobe

创建第二个块设备

fdisk  /dev/vda

Command (m for help): n

All primary partitions are in use

Adding logical partition 5

First sector (144459776-175915007, default 144459776): 

Using default value 144459776

Last sector, +sectors or +size{K,M,G} (144459776-175915007, default 175915007): +12G    

Partition 5 of type Linux and of size 12 GiB is set

Command (m for help): w  

The partition table has been altered!

partprobe





# 禁用ansible  自动配置yum

```
>  /root/openshift-ansible-release-3.6/roles/openshift_repos/files/origin/repos/openshift-ansible-centos-paas-sig.repo
```



# ansible 清单文件如下

cat /etc/ansible/hosts 

```
[OSEv3:children]
masters
nodes
etcd
glusterfs

[OSEv3:vars]

ansible_ssh_user=root
openshift_storage_glusterfs_namespace=glusterfs 
openshift_storage_glusterfs_name=storage 

ansible_become=true

openshift_deployment_type=origin

openshift_disable_check=memory_availability,disk_availability,docker_storage,docker_image_availability


openshift_master_identity_providers=[{'name': 'htpasswd_auth', 'login': 'true', 'challenge': 'true', 'kind': 'HTPasswdPasswordIdentityProvider', 'filename': '/etc/origin/master/htpasswd'}]

[masters]
master-160.example.com

[etcd]
master-160.example.com

[nodes]
node-161.example.com    openshift_schedulable=True  openshift_node_labels="{'region': 'infra', 'zone': 'default'}"
node-162.example.com    openshift_schedulable=True   openshift_node_labels="{'region': 'infra', 'zone': 'default'}"
node-163.example.com    openshift_schedulable=True   openshift_node_labels="{'region': 'infra', 'zone': 'default'}"

[glusterfs]
node-161.example.com   glusterfs_ip=172.16.2.161   glusterfs_devices='[ "/dev/vda4", "/dev/vda5" ]'
node-162.example.com   glusterfs_ip=172.16.2.162   glusterfs_devices='[ "/dev/vda4", "/dev/vda5" ]'
node-163.example.com   glusterfs_ip=172.16.2.163   glusterfs_devices='[ "/dev/vda4", "/dev/vda5" ]'
```

# 配置DNS 服务器

DNS  服务器地址

172.16.2.21


三台服务器 域名规划如下

master-160.example.com

node-161.example.com

node-162.example.com





# NTP  服务器

echo '/usr/sbin/ntpdate 172.16.2.21'   >>   /etc/rc.local

yum -y install ntpdate





# 配置 ssh  key

 ssh-keygen -f ~/.ssh/id_rsa -N ''
 

 for host in  master-160.example.com    node-161.example.com    node-162.example.com;  do  ssh-copy-id -i ~/.ssh/id_rsa.pub $host;  done

#  导入 基础 image

cd /root/op3.6-images/

docker load -i    hello-openshift.latest.tar

docker load -i    kubernetes.latest.tar  

docker load -i    origin-deployer3.6.0.tar  

docker load -i    origin-docker-registry.3.6.0.tar  

docker load -i    origin-haproxy-router.3.6.0.tar  

docker load -i    origin-pod.3.6.0.tar

docker load -i    origin-sti-builder.3.6.0.tar

docker load -i    gluster-centos.tar

docker load -i    heketi.tar


# 执行 ansible  开始安装


yum -y install etcd 

systemctl start etcd


yum install glusterfs-fuse


ansible-playbook -i /etc/ansible/hosts   /root/openshift-ansible-release-3.6/playbooks/byo/config.yml 



failed: [master-160.example.com] (item=master-160.example.com) => {"failed": true, "item": "master-160.example.com", "msg": {"cmd": "/usr/bin/oc label node master-160.example.com glusterfs=storage-host --overwrite", "results": {}, "returncode": 1, "stderr": "Error from server (NotFound): nodes \"master-160.example.com\" not found\n", "stdout": ""}}


手动执行命令

oc label node master-160.example.com glusterfs=storage-host --overwrite


Error from server (NotFound): nodes "master-160.example.com" not found


gluster/gluster-centos:latest


oc  policy  add-role-to-user  admin  dev   -n  default

oc  policy  add-role-to-user  admin  dev   -n  openshift

oc  policy  add-role-to-user  admin  dev   -n  glusterfs



htpasswd -b /etc/origin/master/htpasswd dev dev


oc login -u system:admin






