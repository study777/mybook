svn co [http://172.16.4.1/shtsvn/doc/06\_运维共享目录/supp](http://172.16.4.1/shtsvn/doc/06_运维共享目录/supp)  --no-auth-cache

Authentication realm: &lt;[http://172.16.4.1:80&gt;](http://172.16.4.1:80&gt); SHT-svn

Password for 'root':  
  直接回车

Authentication realm: &lt;[http://172.16.4.1:80&gt;](http://172.16.4.1:80&gt); SHT-svn

Username: wangzhichao  
  输入用户名

Password for 'wangzhichao':   输入密码

cd supp/

cat /supp/src/main/resources/jdbc.properties

jdbc.driverClassName=com.mysql.cj.jdbc.Driver

jdbc.url=jdbc:mysql://139.129.226.254:3306/shengxian?useSSL=false&serverTimezone=UTC

jdbc.username=shengxian

jdbc.password=shengxian2

mysqldump  -ushengxian -h139.129.226.254 shengxian -pshengxian2  &gt; shengxian.sql





oc project supp-data



chmod +x  /root/shengxian.sql

oc cp /root/shengxian.sql mysql-1-jfls6:/tmp/



oc rsh mysql-1-jfls6 

mysql -u $MYSQL_USER -p$MYSQL_PASSWORD -h $HOSTNAME $MYSQL_DATABASE

use shengxian;

source /tmp/shengxian.sql



















