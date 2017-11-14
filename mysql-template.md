cd /opt/mysql-container/5.7/root-common/usr/libexec

scp container-setup  fix-permissions  172.16.2.32:/usr/libexec/



CentOS-Base.repo 


yum install -y centos-release-scl

yum -y install rsync tar gettext hostname bind-utils groff-base shadow-utils rh-mysql57








mkdir  -p /usr/share/container-scripts/mysql

cd  /opt/mysql-container/5.7/root-common/usr/share/container-scripts/mysql


CONTAINER_SCRIPTS_PATH=/usr/share/container-scripts/mysql

MYSQL_PREFIX=/opt/rh/rh-mysql57/root/usr

ENABLED_COLLECTIONS=rh-mysql57

BASH_ENV=${CONTAINER_SCRIPTS_PATH}/scl_enable


ENV=${CONTAINER_SCRIPTS_PATH}/scl_enable


PROMPT_COMMAND=". ${CONTAINER_SCRIPTS_PATH}/scl_enable"






scp 


echo "CONTAINER_SCRIPTS_PATH=/usr/share/container-scripts/mysql"   >> /etc/profile


echo "MYSQL_PREFIX=/opt/rh/rh-mysql57/root/usr"   >> /etc/profile

echo  "ENABLED_COLLECTIONS=rh-mysql57"   >> /etc/profile


