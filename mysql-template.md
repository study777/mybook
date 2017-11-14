cd /opt/mysql-container/5.7/root-common/usr/libexec

scp container-setup  fix-permissions  172.16.2.32:/usr/libexec/



CentOS-Base.repo 


yum install -y centos-release-scl

yum -y install rsync tar gettext hostname bind-utils groff-base shadow-utils rh-mysql57











CONTAINER_SCRIPTS_PATH=/usr/share/container-scripts/mysql

MYSQL_PREFIX=/opt/rh/rh-mysql57/root/usr

ENABLED_COLLECTIONS=rh-mysql57

BASH_ENV=${CONTAINER_SCRIPTS_PATH}/scl_enable


ENV=${CONTAINER_SCRIPTS_PATH}/scl_enable


PROMPT_COMMAND=". ${CONTAINER_SCRIPTS_PATH}/scl_enable"





scp /opt/mysql-container/5.7/root-common/etc/my.cnf    172.16.2.32:/etc/


scp /opt/mysql-container/5.7/root-common/usr/bin/*     172.16.2.32:/usr/bin/


scp /opt/mysql-container/5.7/root-common/usr/libexec/*       172.16.2.32:/usr/libexec/



mkdir  -p /usr/share/container-scripts/mysql

cd  /opt/mysql-container/5.7/root-common/usr/share/container-scripts/mysql


scp -r /opt/mysql-container/5.7/root-common/usr/share/container-scripts/mysql/*     172.16.2.32:/usr/share/container-scripts/mysql






mkdir /usr/libexec/s2i
scp /opt/mysql-container/5.7/s2i-common/bin/*    172.16.2.32:/usr/libexec/s2i
cd /usr/libexec/s2i/
ln -s /bin/run-mysqld ./run






scp 


echo "CONTAINER_SCRIPTS_PATH=/usr/share/container-scripts/mysql"   >> /etc/profile


echo "MYSQL_PREFIX=/opt/rh/rh-mysql57/root/usr"   >> /etc/profile

echo  "ENABLED_COLLECTIONS=rh-mysql57"   >> /etc/profile







```
RUN /usr/libexec/container-setup

scl enable rh-mysql57  -- my_print_defaults --help --verbose  | grep --after=1 '^Default options' | tail -n 1  | grep -o '[^ ]*opt[^ ]*my.cnf'




```



COPY 5.7/root /



scp /opt/mysql-container/5.7/root/usr/share/container-scripts/mysql/README.md   172.16.2.32:/usr/share/container-scripts/mysql/













