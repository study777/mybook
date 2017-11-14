相关目录

/opt/mysql-container/5.7/root-common/usr/share/container-scripts/mysql/cnf

/opt/mysql-container/root-common/usr/share/container-scripts/mysql/common.sh

export MYSQL\_LOWER\_CASE\_TABLE\_NAMES=${MYSQL\_LOWER\_CASE\_TABLE\_NAMES:-0}

制造错误 查看 提示信息

```
docker run -it my-mysql
=> sourcing 20-validate-variables.sh ...
You must either specify the following environment variables:
  MYSQL_USER (regex: '^[a-zA-Z0-9_]+$')
  MYSQL_PASSWORD (regex: '^[a-zA-Z0-9_~!@#$%^&*()-=<>,.?;:|]+$')
  MYSQL_DATABASE (regex: '^[a-zA-Z0-9_]+$')
Or the following environment variable:
  MYSQL_ROOT_PASSWORD (regex: '^[a-zA-Z0-9_~!@#$%^&*()-=<>,.?;:|]+$')
Or both.
Optional Settings:
  MYSQL_LOWER_CASE_TABLE_NAMES (default: 0)
  MYSQL_LOG_QUERIES_ENABLED (default: 0)
  MYSQL_MAX_CONNECTIONS (default: 151)
  MYSQL_FT_MIN_WORD_LEN (default: 4)
  MYSQL_FT_MAX_WORD_LEN (default: 20)
  MYSQL_AIO (default: 1)
  MYSQL_KEY_BUFFER_SIZE (default: 32M or 10% of available memory)
  MYSQL_MAX_ALLOWED_PACKET (default: 200M)
  MYSQL_TABLE_OPEN_CACHE (default: 400)
  MYSQL_SORT_BUFFER_SIZE (default: 256K)
  MYSQL_READ_BUFFER_SIZE (default: 8M or 5% of available memory)
  MYSQL_INNODB_BUFFER_POOL_SIZE (default: 32M or 50% of available memory)
  MYSQL_INNODB_LOG_FILE_SIZE (default: 8M or 15% of available memory)
  MYSQL_INNODB_LOG_BUFFER_SIZE (default: 8M or 15% of available memory)

For more information, see https://github.com/sclorg/mysql-container
```

vim /opt/mysql-container/5.7/root-common/etc/my.cnf

\[mysqld\]  
lower\_case\_table\_names=1

docker pull centos/mysql-57-centos7

git clone --recursive [https://github.com/sclorg/mysql-container.git](https://github.com/sclorg/mysql-container.git)

cd mysql-container

git submodule update --init

make build TARGET=centos7 VERSIONS=5.7

报错内容

    make[1]: Entering directory `/opt/mysql-container'
    mkdir -p 5.7/root
    go-md2man -in "5.7/README.md" -out "5.7/root/help.1"
    /bin/sh: go-md2man: command not found
    make[1]: *** [5.7/root/help.1] Error 127
    make[1]: Leaving directory `/opt/mysql-container'
    make: *** [build-serial] Error 2

解决

yum install golang-github-cpuguy83-go-md2man

some output

    Removing intermediate container e49d7238204b
    Step 7 : ENV CONTAINER_SCRIPTS_PATH /usr/share/container-scripts/mysql MYSQL_PREFIX /opt/rh/rh-mysql57/root/usr ENABLED_COLLECTIONS rh-mysql57
     ---> Running in dd258b2c6c7b
     ---> 6234f57cccc6
    Removing intermediate container dd258b2c6c7b
    Step 8 : ENV BASH_ENV ${CONTAINER_SCRIPTS_PATH}/scl_enable ENV ${CONTAINER_SCRIPTS_PATH}/scl_enable PROMPT_COMMAND ". ${CONTAINER_SCRIPTS_PATH}/scl_enable"
     ---> Running in 2ba03b8bf454
     ---> 713a6ca2a44f
    Removing intermediate container 2ba03b8bf454
    Step 9 : COPY 5.7/root-common /
     ---> 4ca62cb3210b
    Removing intermediate container 2af9130dfd8f
    Step 10 : COPY 5.7/s2i-common/bin/ $STI_SCRIPTS_PATH
     ---> 68fe725f616d
    Removing intermediate container b8342cd32f4b
    Step 11 : COPY 5.7/root /
     ---> 3fab60f52958
    Removing intermediate container 3994ea901045
    Step 12 : RUN rm -rf /etc/my.cnf.d/*
     ---> Running in db26eacfc5f2
     ---> 743efc834960
    Removing intermediate container db26eacfc5f2
    Step 13 : RUN /usr/libexec/container-setup
     ---> Running in 67a7b7c8280b
     ---> 123c1da1ddd7
    Removing intermediate container 67a7b7c8280b
    Step 14 : VOLUME /var/lib/mysql/data
     ---> Running in 000d87fa0078
     ---> 3deb5e212e17
    Removing intermediate container 000d87fa0078
    Step 15 : USER 27
     ---> Running in 5de09ec9c801
     ---> f3e8ec5d37f1
    Removing intermediate container 5de09ec9c801
    Step 16 : ENTRYPOINT container-entrypoint
     ---> Running in 4d82c119ceb2
     ---> 5cc55cbbe478
    Removing intermediate container 4d82c119ceb2
    Step 17 : CMD run-mysqld
     ---> Running in 76c722df5272
     ---> 8d8f102fc13d
    Removing intermediate container 76c722df5272
    Step 18 : LABEL "io.openshift.builder-version" "\"3da6f72\""
     ---> Running in 1cb3ab5a0ff4
     ---> d045b7be9c1c
    Removing intermediate container 1cb3ab5a0ff4
    Successfully built d045b7be9c1c
    mysql 5.7 => d045b7be9c1c
    make[1]: Leaving directory `/opt/mysql-container'
    You have new mail in /var/spool/mail/root



# config pod  success

```
docker pull docker.io/centos/mysql-57-centos7

oc new-project  test-mysql-0

oc policy add-role-to-user admin dev -n  test-mysql-0

oc project test-mysql-0

oc new-app -e MYSQL\_USER=yandun -e MYSQL\_PASSWORD=yandun -e MYSQL\_DATABASE=yandun -e MYSQL\_LOWER\_CASE\_TABLE\_NAMES=1 docker.io/centos/mysql-57-centos7


oc rsh  mysql-57-centos7-1-92rw3

mysql -u $MYSQL\_USER -p$MYSQL\_PASSWORD -h $HOSTNAME $MYSQL\_DATABASE

show variables like "%case%";
mysql> show variables like "%case%";
+------------------------+-------+
| Variable_name          | Value |
+------------------------+-------+
| lower_case_file_system | OFF   |
| lower_case_table_names | 1     |
+------------------------+-------+
2 rows in set (0.01 sec)

```





















