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

设置主机名

```
hostnamectl  set-hostname  master-160.example.com

hostnamectl  set-hostname  node-161.example.com

hostnamectl  set-hostname  node-162.example.com

hostnamectl  set-hostname  node-163.example.com

hostnamectl  set-hostname  node-164.example.com
```

配置本地yum

配置网卡  
删除uuid mac地址等信息，以便可以方便的进行克隆

ls \| grep -v  CentOS-Base.repo \| xargs  rm  -rf

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

sed -i '/OPTIONS=.\*/c\OPTIONS="--selinux-enabled --insecure-registry 172.30.0.0/16"'   /etc/sysconfig/docker

设置docker  存储

事先创建VG
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

直接使用裸设备
systemctl  stop docker

rm -rf /var/lib/docker/*


cat /etc/sysconfig/docker-storage-setup 
DEVS=/dev/sdb
VG=docker-vg
SETUP_LVM_THIN_POOL=yes



lvmconf --disable-cluster

docker-storage-setup

systemctl  start docker



docker info


Storage Driver: devicemapper
 Pool Name: docker--vg-docker--pool
 Pool Blocksize: 524.3 kB
 Base Device Size: 10.74 GB
 Backing Filesystem: xfs
 Data file: 
 Metadata file: 
 Data Space Used: 20.45 MB
 Data Space Total: 8.535 GB
 Data Space Available: 8.515 GB
 Metadata Space Used: 40.96 kB
 Metadata Space Total: 25.17 MB
 Metadata Space Available: 25.12 MB
 Thin Pool Minimum Free Space: 853.5 MB



# 为glusterfs 配置存储块设备

创建第一个块设备

fdisk  /dev/vda

Command \(m for help\): n

Select \(default e\): e

Selected partition 4

First sector \(144457728-209715199, default 144457728\): 回车

Last sector, +sectors or +size{K,M,G} \(144457728-209715199, default 209715199\): +15G

Using default value 209715199

Partition 4 of type Linux and of size 31.1 GiB is set

Command \(m for help\): w

partprobe

创建第二个块设备

fdisk  /dev/vda

Command \(m for help\): n

All primary partitions are in use

Adding logical partition 5

First sector \(144459776-175915007, default 144459776\):

Using default value 144459776

Last sector, +sectors or +size{K,M,G} \(144459776-175915007, default 175915007\): +12G

Partition 5 of type Linux and of size 12 GiB is set

Command \(m for help\): w

The partition table has been altered!

partprobe

# 禁用ansible  自动配置yum

```
cd /root

wget https://github.com/openshift/openshift-ansible/archive/release-3.6.zip

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
node-161.example.com    
node-162.example.com   
node-163.example.com    
node-164.example.com    openshift_schedulable=True   openshift_node_labels="{'region': 'infra', 'zone': 'default'}"

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

echo '/usr/sbin/ntpdate 172.16.2.21'   &gt;&gt;   /etc/rc.local

yum -y install ntpdate

# 配置 ssh  key

```
 ssh-keygen -f ~/.ssh/id_rsa -N ''


 for host in  master-160.example.com    node-161.example.com    node-162.example.com node-163.example.com node-164.example.com;  do  ssh-copy-id -i ~/.ssh/id_rsa.pub $host;  done
```

# 导入 基础 image

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

docker load -i heketi-dev.tar

docker load -i origin-service-catalog.tar

# 执行 ansible  开始安装

yum -y install etcd

systemctl start etcd

yum install glusterfs-fuse

ansible-playbook -i /etc/ansible/hosts   /root/openshift-ansible-release-3.6/playbooks/byo/config.yml

failed: \[master-160.example.com\] \(item=master-160.example.com\) =&gt; {"failed": true, "item": "master-160.example.com", "msg": {"cmd": "/usr/bin/oc label node master-160.example.com glusterfs=storage-host --overwrite", "results": {}, "returncode": 1, "stderr": "Error from server \(NotFound\): nodes \"master-160.example.com\" not found\n", "stdout": ""}}

手动执行命令

oc label node master-160.example.com glusterfs=storage-host --overwrite

Error from server \(NotFound\): nodes "master-160.example.com" not found

gluster/gluster-centos:latest

oc  policy  add-role-to-user  admin  dev   -n  default

oc  policy  add-role-to-user  admin  dev   -n  openshift

oc  policy  add-role-to-user  admin  dev   -n  glusterfs

htpasswd -b /etc/origin/master/htpasswd dev dev

oc login -u system:admin

git clone [https://github.com/openshift/openshift-ansible.git](https://github.com/openshift/openshift-ansible.git)

> /root/openshift-ansible/roles/openshift\_repos/templates/CentOS-OpenShift-Origin36.repo.j2

ansible-playbook -i /etc/ansible/hosts   /root/openshift-ansible/playbooks/byo/config.yml

卸载：

ansible-playbook -i /etc/ansible/hosts   /root/openshift-ansible-release-3.6/playbooks/adhoc/uninstall.yml

# success playbook

```
[OSEv3:children]
masters
nodes
etcd
glusterfs

[OSEv3:vars]

ansible_ssh_user=root


ansible_become=true

openshift_deployment_type=origin
openshift_version=3.6.0
openshift_storage_glusterfs_namespace=glusterfs 
openshift_storage_glusterfs_name=storage
openshift_disable_check=memory_availability,disk_availability,docker_storage,docker_image_availability


openshift_master_identity_providers=[{'name': 'htpasswd_auth', 'login': 'true', 'challenge': 'true', 'kind': 'HTPasswdPasswordIdentityProvider', 'filename': '/etc/origin/master/htpasswd'}]

[masters]
master-161.example.com

[etcd]
master-161.example.com

[nodes]
node-165.example.com  openshift_schedulable=True  openshift_node_labels="{'region': 'infra', 'zone': 'default'}"
node-162.example.com
node-163.example.com
node-164.example.com
[glusterfs]
node-162.example.com     glusterfs_ip=172.16.10.162   glusterfs_devices='[ "/dev/sda5" ]'
node-163.example.com     glusterfs_ip=172.16.10.163   glusterfs_devices='[ "/dev/sda5" ]'
node-164.example.com     glusterfs_ip=172.16.10.164   glusterfs_devices='[ "/dev/sda5" ]'
```

# GlusterFS

oc project glusterfs

查看 Gluster Endpoints

oc get endpoints

heketi-db-storage-endpoints   172.16.10.162:1,172.16.10.163:1,172.16.10.164:1   1d

oc get endpoints  heketi-db-storage-endpoints  -o yaml

apiVersion: v1

kind: Endpoints

metadata:

creationTimestamp: 2017-11-03T16:39:18Z  
  name: heketi-db-storage-endpoints  
  namespace: glusterfs  
  resourceVersion: "2264"  
  selfLink: /api/v1/namespaces/glusterfs/endpoints/heketi-db-storage-endpoints  
  uid: 8964a176-c0b5-11e7-8a0a-000c298f426c  
subsets:

* addresses:
  * ip: 172.16.10.162
  * ip: 172.16.10.163
  * ip: 172.16.10.164
    ports:
  * port: 1
    protocol: TCP

创建pv

cat pv.yaml

```
apiVersion: v1
kind: PersistentVolume
metadata:
  name: gluster-default-volume 
spec:
  capacity:
    storage: 2Gi 
  accessModes: 
    - ReadWriteMany
  glusterfs: 
    endpoints: heketi-db-storage-endpoints 
    path: myVol1 
    readOnly: false
  persistentVolumeReclaimPolicy: Retain
```

创建pvc

```
cat gluster-claim.yaml 

apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: gluster-claim
spec:
  accessModes:
  - ReadWriteMany 
  resources:
     requests:
       storage: 1Gi
```

确认pv和 pvc  状态都是bound

oc get pv

NAME                     CAPACITY   ACCESSMODES   RECLAIMPOLICY   STATUS    CLAIM                     STORAGECLASS

REASON    AGE

gluster-default-volume   2Gi        RWX           Retain          Bound     glusterfs/gluster-claim                            12m

oc get pvc

NAME            STATUS    VOLUME                   CAPACITY   ACCESSMODES   STORAGECLASS   AGE

gluster-claim   Bound     gluster-default-volume   2Gi        RWX                          11m

## 持久化镜像仓库

```
oc project default

oc get pod

docker-registry-1-m8tkm

oc get dc


oc volumes dc/docker-registry --all

deploymentconfigs/docker-registry

  empty directory as registry-storage

    mounted at /registry

  secret/registry-certificates as registry-certificates

    mounted at /etc/secrets
```

查看当前挂载的本地目录使用大小情况

oc rsh  docker-registry-1-m8tkm  'du' '-sh'  '/registry'

0    /registry

当前 并未使用任何空间

如果已经存在数据 可以通过以下方式进行备份

mkdir  /root/backup

cd /root/backup/

oc rsync  docker-registry-1-m8tkm:/registry .

创建 pv

cat registry\_pv.yaml

apiVersion: v1  
kind: PersistentVolume  
metadata:  
  name: registry-volume   
spec:  
  capacity:  
    storage: 5Gi   
  accessModes:

* ReadWriteMany
  glusterfs: 
  endpoints: heketi-db-storage-endpoints 
  path: registry 
  readOnly: false
  persistentVolumeReclaimPolicy: Retain 

oc create -f registry\_pv.yaml

查看pv  状态   Available  即为可用

oc get pv

NAME                     CAPACITY   ACCESSMODES   RECLAIMPOLICY   STATUS

registry-volume          5Gi        RWX           Retain          Available

创建pvc

cat registry\_pvc.yaml

apiVersion: v1  
kind: PersistentVolumeClaim  
metadata:  
  name: docker-registry-claim  
spec:  
  accessModes:

* ReadWriteMany 
  resources:
   requests:
     storage: 5Gi 

oc create  -f registry\_pvc.yaml

oc get pvc   状态为 Bound

NAME                    STATUS    VOLUME            CAPACITY   ACCESSMODES   STORAGECLASS   AGE

docker-registry-claim   Bound     registry-volume   5Gi        RWX                          6s

关联持久化请求

为registry 的容器添加持久化卷请求 docker-registry-claim

并与挂载点 registry-storage  关联

oc volumes  dc/docker-registry  --add --name=registry-storage -t pvc --claim-name=docker-registry-claim --overwrite deploymentconfigs/docker-registry

deploymentconfig "docker-registry" updated

deploymentconfig "docker-registry" updated

再次查看registry 的数据卷信息

oc volumes dc/docker-registry --all  
deploymentconfigs/docker-registry  
  pvc/docker-registry-claim \(allocated 5GiB\) as registry-storage  
    mounted at /registry  
  secret/registry-certificates as registry-certificates  
    mounted at /etc/secrets

报错信息

```
Unable to mount volumes for pod "docker-registry-2-hg798_default(d47a229f-c22b-11e7-a4fa-000c298f426c)": timeout expired waiting for volumes to attach/mount for pod "default"/"docker-registry-2-hg798". list of unattached/unmounted volumes=[registry-storage]
```

oc delete pod docker-registry-1-m8tkm

# 对 registry 存储的操作有误如何修复

oc edit  dc docker-registry

修改为如下

volumes:

* emptyDir: {}
  name: registry-storage
* name: registry-certificates
  secret:
    defaultMode: 420
    secretName: registry-certificates

# 为 mysql 配置 持久存储

使用mysql 持续存储模板创建应用

```
oc get pvc -o yaml
apiVersion: v1
items:
- apiVersion: v1
  kind: PersistentVolumeClaim
  metadata:
    creationTimestamp: 2017-11-08T05:35:03Z
    labels:
      app: mysql-persistent
      template: mysql-persistent-template
    name: mysql
    namespace: mysql
    resourceVersion: "93187"
    selfLink: /api/v1/namespaces/mysql/persistentvolumeclaims/mysql
    uid: 92448c96-c446-11e7-972c-52540011feca
  spec:
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
  status:
    phase: Pending
kind: List
metadata: {}
resourceVersion: ""
selfLink: ""
```

报错信息

```
no persistent volumes available for this claim and no storage class is set
```

修改  mysql-persistent

oc edit template mysql-persistent -n openshift

```
kind: PersistentVolumeClaim
  metadata:
    annotations:
      olume.beta.kubernetes.io/storage-provisioner: kubernetes.io/glusterfs
      volume.beta.kubernetes.io/storage-class: glusterfs-storage
    name: ${DATABASE_SERVICE_NAME}
```

```
报错信息

MountVolume.SetUp failed for volume "kubernetes.io/secret/d6499aa0-c452-11e7-972c-52540011feca-deployer-token-5z6gr" (spec.Name: "deployer-token-5z6gr") pod "d6499aa0-c452-11e7-972c-52540011feca" (UID: "d6499aa0-c452-11e7-972c-52540011feca") with: secret "test"/"deployer-token-5z6gr" not registered
```

```
报错信息

Readiness probe failed: sh: cannot set terminal process group (-1): Inappropriate ioctl for device sh: no job control in this shell ERROR 2003 (HY000): Can't connect to MySQL server on '127.0.0.1' (111)
```

Username: root Password: root Database Name: sampledb Connection URL: mysql://mysql:3306/

![install mysql glusterfs](Gluster-mysql.png)

测试 mysql

查看数据卷大小

```
oc volumes dc/mysql --all -n mysql-t
deploymentconfigs/mysql
  pvc/mysql (allocated 1GiB) as mysql-data
    mounted at /var/lib/mysql/data


oc rsh mysql-1-33nc1  'du' '-sh' '/var/lib/mysql/data'
189M    /var/lib/mysql/data


oc rsh mysql-1-33nc1  'df' '-h'
```

登陆mysql  写入数据

```
oc rsh  mysql-1-33nc1 


普通用户登录
mysql -uuser -h mysql.mysql-t.svc -puser

root  用户登录

mysql -uroot -h mysql.mysql-t.svc -proot

show databases;

use sampledb; 


create table tutorials_tbl(
   tutorial_id INT NOT NULL AUTO_INCREMENT,
   tutorial_title VARCHAR(100) NOT NULL,
   tutorial_author VARCHAR(40) NOT NULL,
   submission_date DATE,
   PRIMARY KEY ( tutorial_id )
);
```

#### 制作一个镜像 用于mysql的客户端 以便测试 mysql 服务器

```
设置使docker 容器可以连接外部网络

echo "net.ipv4.ip_forward=1" >> /etc/sysctl.conf

sysctl  -p


mkdir test

cd test

cat Dockerfile 
FROM docker.io/centos:latest

RUN touch /tmp/test && \

yum -y install mysql

CMD tail -f /tmp/test


Docker build -t mytest .

为新镜像打上标签 并推送至私有registry

docker tag mytest 172.16.2.31:5000/mytest:latest

docker push 172.16.2.31:5000/mytest:latest

openshift  master 将镜像 导入 为 is


oc import-image  172.16.2.31:5000/mytest:latest   -n openshift --confirm --insecure



web console  部署 这个 is

进入容器内部 测试连接 mysql 服务器

oc rsh  mytest-1-pxcmd 


mysql -uroot -h mysql.mysql-t.svc -proot


show databases;
```

#### 测试删除mysql  pod 数据是否都在

```
oc delete pod mysql-1-33nc1

之后会自动生成一个pod
oc get pod | grep mysql
mysql-1-klvk3         1/1       Running   0          33s

在上一步的 客户端容器内查看

oc project test-mysql

oc rsh  mytest-1-pxcmd 

mysql -uroot -h mysql.mysql-t.svc -proot

show databases;

数据都在
```

# 在web console 使用持久存储  应用于 mysql

创建一个 pvc

Persistent Volume Claims

新建project  test1  web 页面点击 storage

点击 create storage

选中 之前创建的 gluster-storage

Name  自定义

mysql

Access Mode  Single User

Size

1 G

Create

volumes:

* name: mysql-data
  persistentVolumeClaim:
    claimName: mysql

```
Unable to mount volumes for pod "mysql-1-cfbdq_test1(edc9376a-c239-11e7-a4fa-000c298f426c)": timeout expired waiting for volumes to attach/mount for pod "test1"/"mysql-1-cfbdq". list of unattached/unmounted volumes=[mysql-data]
2 times in the last
```

```
SchedulerPredicates failed due to PersistentVolumeClaim is not bound: "mysql", which is unexpected. (repeated 4 times)
```

```
no persistent volumes available for this claim and no storage class is set
```

```
Failed to provision volume with StorageClass "glusterfs-storage": glusterfs: create volume err: error creating volume Post http://heketi-storage-glusterfs.router.default.svc.cluster.local/volumes: dial tcp: lookup heketi-storage-glusterfs.router.default.svc.cluster.local: no such host.
```

```
oc get dc/mysql -o yaml

volumes:
      - name: mysql-data
        persistentVolumeClaim:
          claimName: mysql
```

```
报错信息

TASK [openshift_service_catalog : Label master-160.example.com for APIServer and controller deployment] *****************************************
fatal: [master-160.example.com]: FAILED! => {"changed": false, "failed": true, "msg": {"cmd": "/usr/bin/oc label node master-160.example.com openshift-infra=apiserver --overwrite", "results": {}, "returncode": 1, "stderr": "Error from server (NotFound): nodes \"master-160.example.com\" not found\n", "stdout": ""}}
    to retry, use: --limit @/root/openshift-ansible/playbooks/byo/config.retry


oc label node master-160.example.com  openshift-infra=apiserver --overwrite
```

# Success install a cluster  and  don't Manual start etcd

```
[OSEv3:children]
masters
nodes
etcd

[OSEv3:vars]

ansible_ssh_user=root


ansible_become=true

openshift_deployment_type=origin
openshift_version=3.6.0
openshift_disable_check=memory_availability,disk_availability,docker_storage,docker_image_availability


openshift_master_identity_providers=[{'name': 'htpasswd_auth', 'login': 'true', 'challenge': 'true', 'kind': 'HTPasswdPasswordIdentityProvider', 'filename': '/etc/origin/master/htpasswd'}]

[masters]
master-160.example.com

[etcd]
master-160.example.com

[nodes]
node-161.example.com  openshift_schedulable=True  openshift_node_labels="{'region': 'infra', 'zone': 'default'}"
node-162.example.com  openshift_schedulable=True  openshift_node_labels="{'region': 'infra', 'zone': 'default'}"
```

ansible-playbook -i /etc/ansible/hosts /root/openshift-ansible-release-3.6/playbooks/byo/config.yml

# TEST

```
[OSEv3:children]
masters
nodes
etcd
glusterfs
[OSEv3:vars]
ansible_ssh_user=root
ansible_become=true
openshift_deployment_type=origin
openshift_version=3.6.0
openshift_storage_glusterfs_namespace=glusterfs 
openshift_storage_glusterfs_name=storage
openshift_disable_check=memory_availability,disk_availability,docker_storage,docker_image_availability
openshift_master_identity_providers=[{'name': 'htpasswd_auth', 'login': 'true', 'challenge': 'true', 'kind': 'HTPasswdPasswordIdentityProvider', 'filename': '/etc/origin/master/htpasswd'}]

[masters]
master-161.example.com
[etcd]
master-161.example.com
[nodes]
master-161.example.com
node-165.example.com  openshift_schedulable=True  openshift_node_labels="{'region': 'infra', 'zone': 'default'}"
node-162.example.com
node-163.example.com
node-164.example.com
[glusterfs]
node-162.example.com     glusterfs_ip=172.16.10.162   glusterfs_devices='[ "/dev/sda5" ]'
node-163.example.com     glusterfs_ip=172.16.10.163   glusterfs_devices='[ "/dev/sda5" ]'
node-164.example.com     glusterfs_ip=172.16.10.164   glusterfs_devices='[ "/dev/sda5" ]'
```

docker load -i heketi-dev.tar

docker load -i origin-service-catalog.tar

docker  pull docker.io/openshift/origin-service-catalog

git clone [https://github.com/openshift/openshift-ansible.git](https://github.com/openshift/openshift-ansible.git)

cd /root/openshift-ansible/roles/openshift\_repos/templates/

&gt; CentOS-OpenShift-Origin14.repo.j2

&gt; CentOS-OpenShift-Origin15.repo.j2

&gt; CentOS-OpenShift-Origin36.repo.j2

&gt; CentOS-OpenShift-Origin.repo.j2

ansible-playbook -i /etc/ansible/hosts /root/openshift-ansible/playbooks/byo/config.yml

output

```
PLAY RECAP *******************************************************************************************************************************************************************************************************
localhost                  : ok=13   changed=0    unreachable=0    failed=0   
master-161.example.com     : ok=587  changed=120  unreachable=0    failed=0   
node-162.example.com       : ok=174  changed=14   unreachable=0    failed=0   
node-163.example.com       : ok=174  changed=14   unreachable=0    failed=0   
node-164.example.com       : ok=174  changed=14   unreachable=0    failed=0   
node-165.example.com       : ok=170  changed=13   unreachable=0    failed=0   


INSTALLER STATUS *************************************************************************************************************************************************************************************************
Initialization             : Complete
Health Check               : Complete
etcd Install               : Complete
Master Install             : Complete
Master Additional Install  : Complete
Node Install               : Complete
GlusterFS Install          : Complete
Hosted Install             : Complete
Service Catalog Install    : Complete
```

# success but  systemctl restart  origin-master   output  Failed to restart origin-master.service: Unit is masked.

```
[OSEv3:children]
masters
nodes
etcd
glusterfs
[OSEv3:vars]
openshift_enable_service_catalog=false
ansible_ssh_user=root
ansible_become=true
openshift_deployment_type=origin
openshift_version=3.6.0
openshift_storage_glusterfs_namespace=glusterfs 
openshift_storage_glusterfs_name=storage
openshift_disable_check=memory_availability,disk_availability,docker_storage,docker_image_availability
openshift_master_identity_providers=[{'name': 'htpasswd_auth', 'login': 'true', 'challenge': 'true', 'kind': 'HTPasswdPasswordIdentityProvider', 'filename': '/etc/origin/master/htpasswd'}]

[masters]
master-160.example.com
[etcd]
master-160.example.com
[nodes]
master-160.example.com
node-161.example.com  openshift_schedulable=True  openshift_node_labels="{'region': 'infra', 'zone': 'default'}"
node-162.example.com
node-163.example.com
node-164.example.com
[glusterfs]
node-162.example.com     glusterfs_ip=172.16.2.162   glusterfs_devices='[ "/dev/vda5" ]'
node-163.example.com     glusterfs_ip=172.16.2.163   glusterfs_devices='[ "/dev/vda5" ]'
node-164.example.com     glusterfs_ip=172.16.2.164   glusterfs_devices='[ "/dev/vda5" ]'
```

hostname -f

### 动态扩容卷

```
自己试验内容

先看下库的内容

MySQL [sampledb]> select * from sampledb.testtb;
+----+------+------+
| id | name | age  |
+----+------+------+
|  1 | aa   |   12 |
+----+------+------+


oc edit pv pvc-c526a5d4-c45a-11e7-972c-52540011feca 

更改大小


capacity:
    storage: 5Gi



oc replace | oc get pv -o yaml

oc replace | oc get pv  pvc-c526a5d4-c45a-11e7-972c-52540011feca  -o yaml


再次查库 最好重新登录

MySQL [sampledb]> select * from sampledb.testtb;
+----+------+------+
| id | name | age  |
+----+------+------+
|  1 | aa   |   12 |
+----+------+------+






查看 mysql 容器 挂载卷的大小

oc rsh mysql-1-33nc1  'df' '-h'


oc rsh  mysql-1-klvk3   'df' '-h'
Filesystem                                                                                          Size  Used Avail Use% Mounted on
/dev/mapper/docker-253:0-67836054-dc018f37f49bffbbef322d837921b152ddb2452c68ead0cd37784842b1f59752   10G  437M  9.6G   5% /
tmpfs                                                                                               3.9G     0  3.9G   0% /dev
tmpfs                                                                                               3.9G     0  3.9G   0% /sys/fs/cgroup
/dev/mapper/cl-root                                                                                  50G  8.3G   42G  17% /etc/hosts
shm                                                                                                  64M     0   64M   0% /dev/shm
172.16.2.162:vol_90c0d59eeffb2dac420fe630334fcac1                                                  1016M  223M  793M  22% /var/lib/mysql/data
tmpfs                                                                                               3.9G   16K  3.9G   1% /run/secrets/kubernetes.io/serviceaccount





显示还是1G


存储后台  要扩容 

文档链接

https://blog.openshift.com/container-native-storage-for-the-openshift-masses/

172.16.2.162:vol_90c0d59eeffb2dac420fe630334fcac1 

heketi-cli volume expand --volume=0e8a8adc936cd40c2df3698b2f06bba9 --expand-size=2


oc rsh heketi-storage-1-7h38q 



heketi-cli volume  list
Error: Unknown user



后台报错
[negroni] Started GET /volumes
[negroni] Completed 401 Unauthorized in 128.042µs



http://redhatstorage.redhat.com/
```

### 登录一个glusterfs node

```
oc rsh glusterfs-storage-249pb

gluster volume list
heketidbstorage
vol_90c0d59eeffb2dac420fe630334fcac1


gluster volume info

Volume Name: heketidbstorage
Type: Replicate
Volume ID: 7948e5d8-9681-4c4e-b653-f76e94255264
Status: Started
Snapshot Count: 0
Number of Bricks: 1 x 3 = 3
Transport-type: tcp
Bricks:
Brick1: 172.16.2.164:/var/lib/heketi/mounts/vg_538ad1b922315855356110956015a128/brick_a36eb4501204d67e8ce0cf9df5241d98/brick
Brick2: 172.16.2.163:/var/lib/heketi/mounts/vg_43ca2753e822fa045739ad774c5ca8da/brick_1e70cec7cda84a9a4c4531ec3bc37214/brick
Brick3: 172.16.2.162:/var/lib/heketi/mounts/vg_bca58b5f894155d2a4a30f9811943f53/brick_14f7ccd5fab6bceec987bdf4665bd570/brick
Options Reconfigured:
transport.address-family: inet
nfs.disable: on

Volume Name: vol_90c0d59eeffb2dac420fe630334fcac1
Type: Replicate
Volume ID: b02f0f62-e794-4995-9d06-3445cab846d9
Status: Started
Snapshot Count: 0
Number of Bricks: 1 x 3 = 3
Transport-type: tcp
Bricks:
Brick1: 172.16.2.164:/var/lib/heketi/mounts/vg_538ad1b922315855356110956015a128/brick_0f3068b010dd897d1ec922fa89997faa/brick
Brick2: 172.16.2.163:/var/lib/heketi/mounts/vg_43ca2753e822fa045739ad774c5ca8da/brick_78f76ad4ec9611333613a71b352bfaf7/brick
Brick3: 172.16.2.162:/var/lib/heketi/mounts/vg_bca58b5f894155d2a4a30f9811943f53/brick_c1e5c24ca347c954eb304d63c1c4c0c3/brick
Options Reconfigured:
transport.address-family: inet
nfs.disable: on
```

```
gluster volume  info vol_90c0d59eeffb2dac420fe630334fcac1

Volume Name: vol_90c0d59eeffb2dac420fe630334fcac1
Type: Replicate
Volume ID: b02f0f62-e794-4995-9d06-3445cab846d9
Status: Started
Snapshot Count: 0
Number of Bricks: 1 x 3 = 3
Transport-type: tcp
Bricks:
Brick1: 172.16.2.164:/var/lib/heketi/mounts/vg_538ad1b922315855356110956015a128/brick_0f3068b010dd897d1ec922fa89997faa/brick
Brick2: 172.16.2.163:/var/lib/heketi/mounts/vg_43ca2753e822fa045739ad774c5ca8da/brick_78f76ad4ec9611333613a71b352bfaf7/brick
Brick3: 172.16.2.162:/var/lib/heketi/mounts/vg_bca58b5f894155d2a4a30f9811943f53/brick_c1e5c24ca347c954eb304d63c1c4c0c3/brick
Options Reconfigured:
transport.address-family: inet
nfs.disable: on
```

[http://note.youdao.com/share/?id=cf9b1b32f9e6b7c368b9edad5663f7ec&type=note\#/](http://note.youdao.com/share/?id=cf9b1b32f9e6b7c368b9edad5663f7ec&type=note#/)

[https://blog.openshift.com/3-6-gluster-storage-containers/](https://blog.openshift.com/3-6-gluster-storage-containers/)

### success hosts-v1

```
[OSEv3:children]
masters
nodes
etcd
[OSEv3:vars]
ansible_ssh_user=root
openshift_deployment_type=origin
openshift_release=3.6.0
openshift_master_default_subdomain=pass.com
openshift_master_identity_providers=[{'name': 'htpasswd_auth', 'login': 'true', 'challenge': 'true', 'kind': 'HTPasswdPasswordIdentityProvider', 'filename': '/etc/origin/master/htpasswd'}]
openshift_disable_check=memory_availability,disk_availability,docker_storage
oreg_url=registry.example.com/openshift/origin-${component}:${version}
openshift_examples_modify_imagestreams=true
openshift_docker_additional_registries=registry.example.com
[masters]
master01.pass.com
[etcd]
master01.pass.com
[nodes]
master01.pass.com
node01.pass.com openshift_node_labels="{'region': 'infra', 'zone': 'default'}"
node02.pass.com openshift_node_labels="{'region': 'primary', 'zone': 'east'}"
```

host origin example

    # This is an example of a bring your own (byo) host inventory

    # Create an OSEv3 group that contains the masters and nodes groups
    [OSEv3:children]
    masters
    nodes
    etcd
    lb
    # nfs

    # Set variables common for all OSEv3 hosts
    [OSEv3:vars]
    # SSH user, this user should allow ssh based auth without requiring a
    # password. If using ssh key based auth, then the key should be managed by an
    # ssh agent.
    ansible_ssh_user=root

    # If ansible_ssh_user is not root, ansible_become must be set to true and the
    # user must be configured for passwordless sudo
    #ansible_become=yes

    # Debug level for all OpenShift components (Defaults to 2)
    debug_level=2

    # deployment type valid values are origin, online, atomic-enterprise and openshift-enterprise
    deployment_type=origin

    # Specify the generic release of OpenShift to install. This is used mainly just during installation, after which we
    # rely on the version running on the first master. Works best for containerized installs where we can usually
    # use this to lookup the latest exact version of the container images, which is the tag actually used to configure
    # the cluster. For RPM installations we just verify the version detected in your configured repos matches this
    # release.
    # openshift_release=v1.4

    # Specify an exact container image tag to install or configure.
    # WARNING: This value will be used for all hosts in containerized environments, even those that have another version installed.
    # This could potentially trigger an upgrade and downtime, so be careful with modifying this value after the cluster is set up.
    openshift_image_tag=v1.4.1
    containerized=true

    # Specify an exact rpm version to install or configure.
    # WARNING: This value will be used for all hosts in RPM based environments, even those that have another version installed.
    # This could potentially trigger an upgrade and downtime, so be careful with modifying this value after the cluster is set up.
    #openshift_pkg_version=-1.2.0

    # Install the openshift examples
    openshift_install_examples=true

    # Configure logoutURL in the master config for console customization
    # See: https://docs.openshift.org/latest/install_config/web_console_customization.html#changing-the-logout-url
    #openshift_master_logout_url=http://example.com

    # Configure extensionScripts in the master config for console customization
    # See: https://docs.openshift.org/latest/install_config/web_console_customization.html#loading-custom-scripts-and-stylesheets
    #openshift_master_extension_scripts=['/path/to/script1.js','/path/to/script2.js']

    # Configure extensionStylesheets in the master config for console customization
    # See: https://docs.openshift.org/latest/install_config/web_console_customization.html#loading-custom-scripts-and-stylesheets
    #openshift_master_extension_stylesheets=['/path/to/stylesheet1.css','/path/to/stylesheet2.css']

    # Configure extensions in the master config for console customization
    # See: https://docs.openshift.org/latest/install_config/web_console_customization.html#serving-static-files
    #openshift_master_extensions=[{'name': 'images', 'sourceDirectory': '/path/to/my_images'}]

    # Configure extensions in the master config for console customization
    # See: https://docs.openshift.org/latest/install_config/web_console_customization.html#serving-static-files
    #ophenshift_master_oauth_template=/path/to/login-template.html

    # Configure imagePolicyConfig in the master config
    # See: https://godoc.org/github.com/openshift/origin/pkg/cmd/server/api#ImagePolicyConfig
    #openshift_master_image_policy_config={"maxImagesBulkImportedPerRepository": 3, "disableScheduledImport": true}

    # Docker Configuration
    # Add additional, insecure, and blocked registries to global docker configuration
    # For enterprise deployment types we ensure that registry.access.redhat.com is
    # included if you do not include it
    # 自己的私有镜像仓，加速部署而已
    openshift_docker_additional_registries=hub.xxp.cn
    #openshift_docker_insecure_registries=registry.example.com
    #openshift_docker_blocked_registries=registry.hacker.com
    # Disable pushing to dockerhub
    #openshift_docker_disable_push_dockerhub=True
    # Items added, as is, to end of /etc/sysconfig/docker OPTIONS
    # Default value: "--log-driver=json-file --log-opt max-size=50m"
    openshift_docker_options="-l warn --ipv6=false -s overlay --selinux-enabled=false"

    # Specify exact version of Docker to configure or upgrade to.
    # Downgrades are not supported and will error out. Be careful when upgrading docker from < 1.10 to > 1.10.
    # docker_version="1.10.3"

    # Skip upgrading Docker during an OpenShift upgrade, leaves the current Docker version alone.
    # docker_upgrade=False

    # Alternate image format string, useful if you've got your own registry mirror
    #oreg_url=example.com/openshift3/ose-${component}:${version}
    # If oreg_url points to a registry other than registry.access.redhat.com we can
    # modify image streams to point at that registry by setting the following to true
    #openshift_examples_modify_imagestreams=true

    # Origin copr repo
    #openshift_additional_repos=[{'id': 'openshift-origin-copr', 'name': 'OpenShift Origin COPR', 'baseurl': 'https://copr-be.cloud.fedoraproject.org/results/maxamillion/origin-next/epel-7-$basearch/', 'enabled': 1, 'gpgcheck': 1, 'gpgkey': 'https://copr-be.cloud.fedoraproject.org/results/maxamillion/origin-next/pubkey.gpg'}]

    # htpasswd auth
    openshift_master_identity_providers=[{'name': 'htpasswd_auth', 'login': 'true', 'challenge': 'true', 'kind': 'HTPasswdPasswordIdentityProvider', 'filename': '/etc/origin/master/htpasswd'}]
    # Defining htpasswd users
    #openshift_master_htpasswd_users={'user1': '<pre-hashed password>', 'user2': '<pre-hashed password>'}
    # or
    #openshift_master_htpasswd_file=<path to local pre-generated htpasswd file>

    # Allow all auth
    #openshift_master_identity_providers=[{'name': 'allow_all', 'login': 'true', 'challenge': 'true', 'kind': 'AllowAllPasswordIdentityProvider'}]

    # LDAP auth
    #openshift_master_identity_providers=[{'name': 'my_ldap_provider', 'challenge': 'true', 'login': 'true', 'kind': 'LDAPPasswordIdentityProvider', 'attributes': {'id': ['dn'], 'email': ['mail'], 'name': ['cn'], 'preferredUsername': ['uid']}, 'bindDN': '', 'bindPassword': '', 'ca': 'my-ldap-ca.crt', 'insecure': 'false', 'url': 'ldap://ldap.example.com:389/ou=users,dc=example,dc=com?uid'}]
    #
    # Configure LDAP CA certificate
    # Specify either the ASCII contents of the certificate or the path to
    # the local file that will be copied to the remote host. CA
    # certificate contents will be copied to master systems and saved
    # within /etc/origin/master/ with a filename matching the "ca" key set
    # within the LDAPPasswordIdentityProvider.
    #
    #openshift_master_ldap_ca=<ca text>
    # or
    #openshift_master_ldap_ca_file=<path to local ca file to use>

    # OpenID auth
    #openshift_master_identity_providers=[{"name": "openid_auth", "login": "true", "challenge": "false", "kind": "OpenIDIdentityProvider", "client_id": "my_client_id", "client_secret": "my_client_secret", "claims": {"id": ["sub"], "preferredUsername": ["preferred_username"], "name": ["name"], "email": ["email"]}, "urls": {"authorize": "https://myidp.example.com/oauth2/authorize", "token": "https://myidp.example.com/oauth2/token"}, "ca": "my-openid-ca-bundle.crt"}]
    #
    # Configure OpenID CA certificate
    # Specify either the ASCII contents of the certificate or the path to
    # the local file that will be copied to the remote host. CA
    # certificate contents will be copied to master systems and saved
    # within /etc/origin/master/ with a filename matching the "ca" key set
    # within the OpenIDIdentityProvider.
    #
    #openshift_master_openid_ca=<ca text>
    # or
    #openshift_master_openid_ca_file=<path to local ca file to use>

    # Request header auth
    #openshift_master_identity_providers=[{"name": "my_request_header_provider", "challenge": "true", "login": "true", "kind": "RequestHeaderIdentityProvider", "challengeURL": "https://www.example.com/challenging-proxy/oauth/authorize?${query}", "loginURL": "https://www.example.com/login-proxy/oauth/authorize?${query}", "clientCA": "my-request-header-ca.crt", "clientCommonNames": ["my-auth-proxy"], "headers": ["X-Remote-User", "SSO-User"], "emailHeaders": ["X-Remote-User-Email"], "nameHeaders": ["X-Remote-User-Display-Name"], "preferredUsernameHeaders": ["X-Remote-User-Login"]}]
    #
    # Configure request header CA certificate
    # Specify either the ASCII contents of the certificate or the path to
    # the local file that will be copied to the remote host. CA
    # certificate contents will be copied to master systems and saved
    # within /etc/origin/master/ with a filename matching the "clientCA"
    # key set within the RequestHeaderIdentityProvider.
    #
    #openshift_master_request_header_ca=<ca text>
    # or
    #openshift_master_request_header_ca_file=<path to local ca file to use>

    # Cloud Provider Configuration
    #
    # Note: You may make use of environment variables rather than store
    # sensitive configuration within the ansible inventory.
    # For example:
    #openshift_cloudprovider_aws_access_key="{{ lookup('env','AWS_ACCESS_KEY_ID') }}"
    #openshift_cloudprovider_aws_secret_key="{{ lookup('env','AWS_SECRET_ACCESS_KEY') }}"
    #
    # AWS
    #openshift_cloudprovider_kind=aws
    # Note: IAM profiles may be used instead of storing API credentials on disk.
    #openshift_cloudprovider_aws_access_key=aws_access_key_id
    #openshift_cloudprovider_aws_secret_key=aws_secret_access_key
    #
    # Openstack
    #openshift_cloudprovider_kind=openstack
    #openshift_cloudprovider_openstack_auth_url=http://openstack.example.com:35357/v2.0/
    #openshift_cloudprovider_openstack_username=username
    #openshift_cloudprovider_openstack_password=password
    #openshift_cloudprovider_openstack_domain_id=domain_id
    #openshift_cloudprovider_openstack_domain_name=domain_name
    #openshift_cloudprovider_openstack_tenant_id=tenant_id
    #openshift_cloudprovider_openstack_tenant_name=tenant_name
    #openshift_cloudprovider_openstack_region=region
    #openshift_cloudprovider_openstack_lb_subnet_id=subnet_id
    #
    # GCE
    #openshift_cloudprovider_kind=gce

    # Project Configuration
    #osm_project_request_message=''
    #osm_project_request_template=''
    #osm_mcs_allocator_range='s0:/2'
    #osm_mcs_labels_per_project=5
    #osm_uid_allocator_range='1000000000-1999999999/10000'

    # Configure additional projects
    #openshift_additional_projects={'my-project': {'default_node_selector': 'label=value'}}

    # Enable cockpit
    #osm_use_cockpit=true
    #
    # Set cockpit plugins
    #osm_cockpit_plugins=['cockpit-kubernetes']

    # Native high availability cluster method with optional load balancer.
    # If no lb group is defined, the installer assumes that a load balancer has
    # been preconfigured. For installation the value of
    # openshift_master_cluster_hostname must resolve to the load balancer
    # or to one or all of the masters defined in the inventory if no load
    # balancer is present.
    openshift_master_cluster_method=native
    openshift_master_cluster_hostname=192.168.56.109
    openshift_master_cluster_public_hostname=192.168.31.116

    # Pacemaker high availability cluster method.
    # Pacemaker HA environment must be able to self provision the
    # configured VIP. For installation openshift_master_cluster_hostname
    # must resolve to the configured VIP.
    #openshift_master_cluster_method=pacemaker
    #openshift_master_cluster_password=openshift_cluster
    #openshift_master_cluster_vip=192.168.133.25
    #openshift_master_cluster_public_vip=192.168.133.25
    #openshift_master_cluster_hostname=openshift-ansible.test.example.com
    #openshift_master_cluster_public_hostname=openshift-ansible.test.example.com

    # Override the default controller lease ttl
    #osm_controller_lease_ttl=30

    # Configure controller arguments
    #osm_controller_args={'resource-quota-sync-period': ['10s']}

    # Configure api server arguments
    #osm_api_server_args={'max-requests-inflight': ['400']}

    # default subdomain to use for exposed routes
    #openshift_master_default_subdomain=apps.test.example.com

    # additional cors origins
    #osm_custom_cors_origins=['foo.example.com', 'bar.example.com']

    # default project node selector
    #osm_default_node_selector='region=primary'

    # Override the default pod eviction timeout
    #openshift_master_pod_eviction_timeout=5m

    # Override the default oauth tokenConfig settings:
    # openshift_master_access_token_max_seconds=86400
    # openshift_master_auth_token_max_seconds=500

    # Override master servingInfo.maxRequestsInFlight
    #openshift_master_max_requests_inflight=500

    # default storage plugin dependencies to install, by default the ceph and
    # glusterfs plugin dependencies will be installed, if available.
    #osn_storage_plugin_deps=['ceph','glusterfs','iscsi']

    # OpenShift Router Options
    #
    # An OpenShift router will be created during install if there are
    # nodes present with labels matching the default router selector,
    # "region=infra". Set openshift_node_labels per node as needed in
    # order to label nodes.
    #
    # Example:
    # [nodes]
    # node.example.com openshift_node_labels="{'region': 'infra'}"
    #
    # Router selector (optional)
    # Router will only be created if nodes matching this label are present.
    # Default value: 'region=infra'
    #openshift_hosted_router_selector='region=infra'
    #
    # Router replicas (optional)
    # Unless specified, openshift-ansible will calculate the replica count
    # based on the number of nodes matching the openshift router selector.
    #openshift_hosted_router_replicas=2
    #
    # Router force subdomain (optional)
    # A router path format to force on all routes used by this router
    # (will ignore the route host value)
    #openshift_hosted_router_force_subdomain='${name}-${namespace}.apps.example.com'
    #
    # Router certificate (optional)
    # Provide local certificate paths which will be configured as the
    # router's default certificate.
    #openshift_hosted_router_certificate={"certfile": "/path/to/router.crt", "keyfile": "/path/to/router.key", "cafile": "/path/to/router-ca.crt"}
    #
    # Disable management of the OpenShift Router
    #openshift_hosted_manage_router=false

    # Openshift Registry Options
    #
    # An OpenShift registry will be created during install if there are
    # nodes present with labels matching the default registry selector,
    # "region=infra". Set openshift_node_labels per node as needed in
    # order to label nodes.
    #
    # Example:
    # [nodes]
    # node.example.com openshift_node_labels="{'region': 'infra'}"
    #
    # Registry selector (optional)
    # Registry will only be created if nodes matching this label are present.
    # Default value: 'region=infra'
    #openshift_hosted_registry_selector='region=infra'
    #
    # Registry replicas (optional)
    # Unless specified, openshift-ansible will calculate the replica count
    # based on the number of nodes matching the openshift registry selector.
    #openshift_hosted_registry_replicas=2
    #
    # Disable management of the OpenShift Registry
    #openshift_hosted_manage_registry=false

    # Registry Storage Options
    #
    # NFS Host Group
    # An NFS volume will be created with path "nfs_directory/volume_name"
    # on the host within the [nfs] host group.  For example, the volume
    # path using these options would be "/exports/registry"
    #openshift_hosted_registry_storage_kind=nfs
    #openshift_hosted_registry_storage_access_modes=['ReadWriteMany']
    #openshift_hosted_registry_storage_nfs_directory=/exports
    #openshift_hosted_registry_storage_nfs_options='*(rw,root_squash)'
    #openshift_hosted_registry_storage_volume_name=registry
    #openshift_hosted_registry_storage_volume_size=10Gi
    #
    # External NFS Host
    # NFS volume must already exist with path "nfs_directory/_volume_name" on
    # the storage_host. For example, the remote volume path using these
    # options would be "nfs.example.com:/exports/registry"
    #openshift_hosted_registry_storage_kind=nfs
    #openshift_hosted_registry_storage_access_modes=['ReadWriteMany']
    #openshift_hosted_registry_storage_host=nfs.example.com
    #openshift_hosted_registry_storage_nfs_directory=/exports
    #openshift_hosted_registry_storage_volume_name=registry
    #openshift_hosted_registry_storage_volume_size=10Gi
    #
    # Openstack
    # Volume must already exist.
    #openshift_hosted_registry_storage_kind=openstack
    #openshift_hosted_registry_storage_access_modes=['ReadWriteOnce']
    #openshift_hosted_registry_storage_openstack_filesystem=ext4
    #openshift_hosted_registry_storage_openstack_volumeID=3a650b4f-c8c5-4e0a-8ca5-eaee11f16c57
    #openshift_hosted_registry_storage_volume_size=10Gi
    #
    # AWS S3
    # S3 bucket must already exist.
    #openshift_hosted_registry_storage_kind=object
    #openshift_hosted_registry_storage_provider=s3
    #openshift_hosted_registry_storage_s3_accesskey=aws_access_key_id
    #openshift_hosted_registry_storage_s3_secretkey=aws_secret_access_key
    #openshift_hosted_registry_storage_s3_bucket=bucket_name
    #openshift_hosted_registry_storage_s3_region=bucket_region
    #openshift_hosted_registry_storage_s3_chunksize=26214400
    #openshift_hosted_registry_storage_s3_rootdirectory=/registry
    #openshift_hosted_registry_pullthrough=true
    #openshift_hosted_registry_acceptschema2=true
    #openshift_hosted_registry_enforcequota=true
    #
    # Any S3 service (Minio, ExoScale, ...): Basically the same as above
    # but with regionendpoint configured
    # S3 bucket must already exist.
    #openshift_hosted_registry_storage_kind=object
    #openshift_hosted_registry_storage_provider=s3
    #openshift_hosted_registry_storage_s3_accesskey=access_key_id
    #openshift_hosted_registry_storage_s3_secretkey=secret_access_key
    #openshift_hosted_registry_storage_s3_regionendpoint=https://myendpoint.example.com/
    #openshift_hosted_registry_storage_s3_bucket=bucket_name
    #openshift_hosted_registry_storage_s3_region=bucket_region
    #openshift_hosted_registry_storage_s3_chunksize=26214400
    #openshift_hosted_registry_storage_s3_rootdirectory=/registry
    #openshift_hosted_registry_pullthrough=true
    #openshift_hosted_registry_acceptschema2=true
    #openshift_hosted_registry_enforcequota=true

    # Metrics deployment
    # See: https://docs.openshift.com/enterprise/latest/install_config/cluster_metrics.html
    #
    # By default metrics are not automatically deployed, set this to enable them
    # openshift_hosted_metrics_deploy=true
    #
    # Storage Options
    # If openshift_hosted_metrics_storage_kind is unset then metrics will be stored
    # in an EmptyDir volume and will be deleted when the cassandra pod terminates.
    # Storage options A & B currently support only one cassandra pod which is
    # generally enough for up to 1000 pods. Additional volumes can be created
    # manually after the fact and metrics scaled per the docs.
    #
    # Option A - NFS Host Group
    # An NFS volume will be created with path "nfs_directory/volume_name"
    # on the host within the [nfs] host group.  For example, the volume
    # path using these options would be "/exports/metrics"
    #openshift_hosted_metrics_storage_kind=nfs
    #openshift_hosted_metrics_storage_access_modes=['ReadWriteOnce']
    #openshift_hosted_metrics_storage_nfs_directory=/exports
    #openshift_hosted_metrics_storage_nfs_options='*(rw,root_squash)'
    #openshift_hosted_metrics_storage_volume_name=metrics
    #openshift_hosted_metrics_storage_volume_size=10Gi
    #
    # Option B - External NFS Host
    # NFS volume must already exist with path "nfs_directory/_volume_name" on
    # the storage_host. For example, the remote volume path using these
    # options would be "nfs.example.com:/exports/metrics"
    #openshift_hosted_metrics_storage_kind=nfs
    #openshift_hosted_metrics_storage_access_modes=['ReadWriteOnce']
    #openshift_hosted_metrics_storage_host=nfs.example.com
    #openshift_hosted_metrics_storage_nfs_directory=/exports
    #openshift_hosted_metrics_storage_volume_name=metrics
    #openshift_hosted_metrics_storage_volume_size=10Gi
    #
    # Option C - Dynamic -- If openshift supports dynamic volume provisioning for
    # your cloud platform use this.
    #openshift_hosted_metrics_storage_kind=dynamic
    #
    # Override metricsPublicURL in the master config for cluster metrics
    # Defaults to https://hawkular-metrics.{{openshift_master_default_subdomain}}/hawkular/metrics
    # Currently, you may only alter the hostname portion of the url, alterting the
    # `/hawkular/metrics` path will break installation of metrics.
    #openshift_hosted_metrics_public_url=https://hawkular-metrics.example.com/hawkular/metrics

    # Logging deployment
    #
    # Currently logging deployment is disabled by default, enable it by setting this
    #openshift_hosted_logging_deploy=true
    #
    # Logging storage config
    # Option A - NFS Host Group
    # An NFS volume will be created with path "nfs_directory/volume_name"
    # on the host within the [nfs] host group.  For example, the volume
    # path using these options would be "/exports/logging"
    #openshift_hosted_logging_storage_kind=nfs
    #openshift_hosted_logging_storage_access_modes=['ReadWriteOnce']
    #openshift_hosted_logging_storage_nfs_directory=/exports
    #openshift_hosted_logging_storage_nfs_options='*(rw,root_squash)'
    #openshift_hosted_logging_storage_volume_name=logging
    #openshift_hosted_logging_storage_volume_size=10Gi
    #
    # Option B - External NFS Host
    # NFS volume must already exist with path "nfs_directory/_volume_name" on
    # the storage_host. For example, the remote volume path using these
    # options would be "nfs.example.com:/exports/logging"
    #openshift_hosted_logging_storage_kind=nfs
    #openshift_hosted_logging_storage_access_modes=['ReadWriteOnce']
    #openshift_hosted_logging_storage_host=nfs.example.com
    #openshift_hosted_logging_storage_nfs_directory=/exports
    #openshift_hosted_logging_storage_volume_name=logging
    #openshift_hosted_logging_storage_volume_size=10Gi
    #
    # Option C - Dynamic -- If openshift supports dynamic volume provisioning for
    # your cloud platform use this.
    #openshift_hosted_logging_storage_kind=dynamic
    #
    # Option D - none -- Logging will use emptydir volumes which are destroyed when
    # pods are deleted
    #
    # Other Logging Options -- Common items you may wish to reconfigure, for the complete
    # list of options please see roles/openshift_hosted_logging/README.md
    #
    # Configure loggingPublicURL in the master config for aggregate logging, defaults
    # to https://kibana.{{ openshift_master_default_subdomain }}
    #openshift_master_logging_public_url=https://kibana.example.com
    # Configure the number of elastic search nodes, unless you're using dynamic provisioning
    # this value must be 1
    #openshift_hosted_logging_elasticsearch_cluster_size=1
    #openshift_hosted_logging_hostname=logging.apps.example.com
    # Configure the prefix and version for the deployer image
    #openshift_hosted_logging_deployer_prefix=registry.example.com:8888/openshift3/
    #openshift_hosted_logging_deployer_version=3.3.0

    # Configure the multi-tenant SDN plugin (default is 'redhat/openshift-ovs-subnet')
    os_sdn_network_plugin_name='redhat/openshift-ovs-multitenant'

    # Disable the OpenShift SDN plugin
    # openshift_use_openshift_sdn=False

    # Configure SDN cluster network and kubernetes service CIDR blocks. These
    # network blocks should be private and should not conflict with network blocks
    # in your infrastructure that pods may require access to. Can not be changed
    # after deployment.
    # 为了实现1000+ pods per node, 才有了这里的更改
    osm_cluster_network_cidr=12.1.0.0/12
    openshift_portal_net=170.16.0.0/16


    # ExternalIPNetworkCIDRs controls what values are acceptable for the
    # service external IP field. If empty, no externalIP may be set. It
    # may contain a list of CIDRs which are checked for access. If a CIDR
    # is prefixed with !, IPs in that CIDR will be rejected. Rejections
    # will be applied first, then the IP checked against one of the
    # allowed CIDRs. You should ensure this range does not overlap with
    # your nodes, pods, or service CIDRs for security reasons.
    #openshift_master_external_ip_network_cidrs=['0.0.0.0/0']

    # IngressIPNetworkCIDR controls the range to assign ingress IPs from for
    # services of type LoadBalancer on bare metal. If empty, ingress IPs will not
    # be assigned. It may contain a single CIDR that will be allocated from. For
    # security reasons, you should ensure that this range does not overlap with
    # the CIDRs reserved for external IPs, nodes, pods, or services.
    #openshift_master_ingress_ip_network_cidr=172.46.0.0/16

    # Configure number of bits to allocate to each host’s subnet e.g. 8
    # would mean a /24 network on the host.
    osm_host_subnet_length=10

    # Configure master API and console ports.
    #openshift_master_api_port=8443
    #openshift_master_console_port=8443

    # set RPM version for debugging purposes
    #openshift_pkg_version=-1.1

    # Configure custom ca certificate
    #openshift_master_ca_certificate={'certfile': '/path/to/ca.crt', 'keyfile': '/path/to/ca.key'}
    #
    # NOTE: CA certificate will not be replaced with existing clusters.
    # This option may only be specified when creating a new cluster or
    # when redeploying cluster certificates with the redeploy-certificates
    # playbook. If replacing the CA certificate in an existing cluster
    # with a custom ca certificate, the following variable must also be
    # set.
    #openshift_certificates_redeploy_ca=true

    # Configure custom named certificates (SNI certificates)
    #
    # https://docs.openshift.org/latest/install_config/certificate_customization.html
    #
    # NOTE: openshift_master_named_certificates is cached on masters and is an
    # additive fact, meaning that each run with a different set of certificates
    # will add the newly provided certificates to the cached set of certificates.
    #
    # An optional CA may be specified for each named certificate. CAs will
    # be added to the OpenShift CA bundle which allows for the named
    # certificate to be served for internal cluster communication.
    #
    # If you would like openshift_master_named_certificates to be overwritten with
    # the provided value, specify openshift_master_overwrite_named_certificates.
    #openshift_master_overwrite_named_certificates=true
    #
    # Provide local certificate paths which will be deployed to masters
    #openshift_master_named_certificates=[{"certfile": "/path/to/custom1.crt", "keyfile": "/path/to/custom1.key", "cafile": "/path/to/custom-ca1.crt"}]
    #
    # Detected names may be overridden by specifying the "names" key
    #openshift_master_named_certificates=[{"certfile": "/path/to/custom1.crt", "keyfile": "/path/to/custom1.key", "names": ["public-master-host.com"], "cafile": "/path/to/custom-ca1.crt"}]

    # Session options
    #openshift_master_session_name=ssn
    #openshift_master_session_max_seconds=3600

    # An authentication and encryption secret will be generated if secrets
    # are not provided. If provided, openshift_master_session_auth_secrets
    # and openshift_master_encryption_secrets must be equal length.
    #
    # Signing secrets, used to authenticate sessions using
    # HMAC. Recommended to use secrets with 32 or 64 bytes.
    #openshift_master_session_auth_secrets=['DONT+USE+THIS+SECRET+b4NV+pmZNSO']
    #
    # Encrypting secrets, used to encrypt sessions. Must be 16, 24, or 32
    # characters long, to select AES-128, AES-192, or AES-256.
    #openshift_master_session_encryption_secrets=['DONT+USE+THIS+SECRET+b4NV+pmZNSO']

    # configure how often node iptables rules are refreshed
    #openshift_node_iptables_sync_period=5s

    # Configure nodeIP in the node config
    # This is needed in cases where node traffic is desired to go over an
    # interface other than the default network interface.
    # openshift_node_set_node_ip=True

    # Force setting of system hostname when configuring OpenShift
    # This works around issues related to installations that do not have valid dns
    # entries for the interfaces attached to the host.
    #openshift_set_hostname=True

    # Configure dnsIP in the node config
    openshift_dns_ip=170.16.0.2

    # Configure node kubelet arguments. pods-per-core is valid in OpenShift Origin 1.3 or OpenShift Container Platform 3.3 and later.
    openshift_node_kubelet_args={'pods-per-core': ['0'], 'max-pods': ['1024'], 'image-gc-high-threshold': ['90'], 'image-gc-low-threshold': ['80']}

    # Configure logrotate scripts
    # See: https://github.com/nickhammond/ansible-logrotate
    #logrotate_scripts=[{"name": "syslog", "path": "/var/log/cron\n/var/log/maillog\n/var/log/messages\n/var/log/secure\n/var/log/spooler\n", "options": ["daily", "rotate 7", "compress", "sharedscripts", "missingok"], "scripts": {"postrotate": "/bin/kill -HUP `cat /var/run/syslogd.pid 2> /dev/null` 2> /dev/null || true"}}]

    # openshift-ansible will wait indefinitely for your input when it detects that the
    # value of openshift_hostname resolves to an IP address not bound to any local
    # interfaces. This mis-configuration is problematic for any pod leveraging host
    # networking and liveness or readiness probes.
    # Setting this variable to true will override that check.
    #openshift_override_hostname_check=true

    # Configure dnsmasq for cluster dns, switch the host's local resolver to use dnsmasq
    # and configure node's dnsIP to point at the node's local dnsmasq instance. Defaults
    # to True for Origin 1.2 and OSE 3.2. False for 1.1 / 3.1 installs, this cannot
    # be used with 1.0 and 3.0.
    #openshift_use_dnsmasq=False
    # Define an additional dnsmasq.conf file to deploy to /etc/dnsmasq.d/openshift-ansible.conf
    # This is useful for POC environments where DNS may not actually be available yet.
    #openshift_node_dnsmasq_additional_config_file=/home/bob/ose-dnsmasq.conf

    # Global Proxy Configuration
    # These options configure HTTP_PROXY, HTTPS_PROXY, and NOPROXY environment
    # variables for docker and master services.
    #openshift_http_proxy=http://USER:PASSWORD@IPADDR:PORT
    #openshift_https_proxy=https://USER:PASSWORD@IPADDR:PORT
    #openshift_no_proxy='.hosts.example.com,some-host.com'
    #
    # Most environments don't require a proxy between openshift masters, nodes, and
    # etcd hosts. So automatically add those hostnames to the openshift_no_proxy list.
    # If all of your hosts share a common domain you may wish to disable this and
    # specify that domain above.
    #openshift_generate_no_proxy_hosts=True
    #
    # These options configure the BuildDefaults admission controller which injects
    # environment variables into Builds. These values will default to the global proxy
    # config values. You only need to set these if they differ from the global settings
    # above. See BuildDefaults
    # documentation at https://docs.openshift.org/latest/admin_guide/build_defaults_overrides.html
    #openshift_builddefaults_http_proxy=http://USER:PASSWORD@HOST:PORT
    #openshift_builddefaults_https_proxy=https://USER:PASSWORD@HOST:PORT
    #openshift_builddefaults_no_proxy=build_defaults
    #openshift_builddefaults_git_http_proxy=http://USER:PASSWORD@HOST:PORT
    #openshift_builddefaults_git_https_proxy=https://USER:PASSWORD@HOST:PORT
    # Or you may optionally define your own serialized as json
    #openshift_builddefaults_json='{"BuildDefaults":{"configuration":{"kind":"BuildDefaultsConfig","apiVersion":"v1","gitHTTPSProxy":"http://proxy.example.com.redhat.com:3128","gitHTTPProxy":"http://proxy.example.com.redhat.com:3128","env":[{"name":"HTTP_PROXY","value":"http://proxy.example.com.redhat.com:3128"},{"name":"HTTPS_PROXY","value":"http://proxy.example.com.redhat.com:3128"},{"name":"NO_PROXY","value":"ose3-master.example.com"}]}}'
    # masterConfig.volumeConfig.dynamicProvisioningEnabled, configurable as of 1.2/3.2, enabled by default
    #openshift_master_dynamic_provisioning_enabled=False

    # Configure usage of openshift_clock role.
    #openshift_clock_enabled=true

    # OpenShift Per-Service Environment Variables
    # Environment variables are added to /etc/sysconfig files for
    # each OpenShift service: node, master (api and controllers).
    # API and controllers environment variables are merged in single
    # master environments.
    #openshift_master_api_env_vars={"ENABLE_HTTP2": "true"}
    #openshift_master_controllers_env_vars={"ENABLE_HTTP2": "true"}
    #openshift_node_env_vars={"ENABLE_HTTP2": "true"}

    # Enable API service auditing, available as of 1.3
    #openshift_master_audit_config={"basicAuditEnabled": true}

    # host group for masters
    [masters]
    192.168.56.106 openshift_ip=192.168.56.106 openshift_hostname=192.168.56.106
    192.168.56.107 openshift_ip=192.168.56.107 openshift_hostname=192.168.56.107
    192.168.56.108 openshift_ip=192.168.56.108 openshift_hostname=192.168.56.108


    # openshift_ip和openshift_hostname可以指定内网IP，防止监听到外网ip上
    192.168.56.106  openshift_ip=192.168.56.106 openshift_hostname=192.168.56.106
    192.168.56.107  openshift_ip=192.168.56.107 openshift_hostname=192.168.56.107
    192.168.56.108  openshift_ip=192.168.56.108 openshift_hostname=192.168.56.108

    # NOTE: Containerized load balancer hosts are not yet supported, if using a global
    # containerized=true host variable we must set to false.
    [lb]
    192.168.56.109 containerized=false

    # NOTE: Currently we require that masters be part of the SDN which requires that they also be nodes
    # However, in order to ensure that your masters are not burdened with running pods you should
    # make them unschedulable by adding openshift_schedulable=False any node that's also a master.
    [nodes]
    192.168.56.106 openshift_ip=192.168.56.106 openshift_hostname=192.168.56.106
    192.168.56.107 openshift_ip=192.168.56.107 openshift_hostname=192.168.56.107 
    192.168.56.108 openshift_ip=192.168.56.108 openshift_hostname=192.168.56.108
    192.168.56.110 openshift_ip=192.168.56.110 openshift_hostname=192.168.56.110 openshift_node_labels="{'region': 'primary', 'zone': 'default'}"



