###  关闭 selinux  firewalld

setenforce 0

systemctl  stop firewalld


###  安装软件


yum -y install git   perl-Git  httpd  mod_ssl openssl


###  生成证书  并签名

openssl genrsa -out ca.key 2048


openssl  req -new -key ca.key -out ca.csr

Common Name (eg, your name or your server's hostname) []:www.git.com

openssl x509 -req -days 365 -in ca.csr -signkey ca.key -out ca.crt

cp ca.crt /etc/pki/tls/certs/

cp ca.key  /etc/pki/tls/private/

cp ca.csr  /etc/pki/tls/private/


###  修改apache ssl 配置文件  将证书 更改为 上一步生成的证书

vim /etc/httpd/conf.d/ssl.conf


SSLCertificateFile /etc/pki/tls/certs/ca.crt
SSLCertificateKeyFile /etc/pki/tls/private/ca.key

systemctl  restart httpd



### 先测试一下 http 
curl  http://localhost

会输出apache 首页



### 添加一个 https  网站

vim /etc/httpd/conf/httpd.conf

在文件最后处添加

<VirtualHost *:443>
        ServerName www.git.com
        <Location />
        </Location>
</VirtualHost>

systemctl  restart httpd



###  在 linux 上测试 https

修改 在服务器上修改hosts 文件

echo  '172.16.2.33  www.git.com'    >> /etc/hosts

在本机访问测试

curl  https://www.git.com  --insecure


在 linux  上添加信任 自签名证书

cat ca.crt  >>   /etc/pki/tls/certs/ca-bundle.crt

再次访问测试

curl  https://www.git.com

输出网站首页






### 在windowns 测试 https 

修改 hosts  文件 



浏览器上访问

https://www.git.com/

git  bash 命令行 访问

curl.exe   https://www.git.com  --insecure


windowns10  上导入证书
导入SSL证书：

WIN键 +  R  输入 MMC
管理控制台 -> 选择菜单“文件 -〉添加/删除管理单元”->列表中选择“证书”->点击“添加”-> 选择“计算机帐户” -> "本地计算机"->点击完成


浏览器测试 

https://www.git.com/


git bash 命令行测试

curl https://www.git.com

输出如下

$ curl.exe   https://www.git.com
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0
curl: (60) SSL certificate problem: self signed certificate
More details here: https://curl.haxx.se/docs/sslcerts.html

curl failed to verify the legitimacy of the server and therefore could not
establish a secure connection to it. To learn more about this situation and
how to fix it, please visit the web page mentioned above.



### 搭建 https git 并使用账号密码的 服务器

apache 配置文件如下


修改上面的  apache 虚拟主机

<VirtualHost *:443>
ServerName   www.git.com
        <Location />
        </Location>
</VirtualHost>


为如下内容

```
<VirtualHost *:443>
ServerName   www.git.com
        SetEnv GIT_HTTP_EXPORT_ALL
        SetEnv GIT_PROJECT_ROOT /home/git
        ScriptAlias /git/ /usr/libexec/git-core/git-http-backend/
        <Location />
                AuthType Basic
                AuthName "Git"
                AuthUserFile /etc/httpd/conf.d/git-team.htpasswd
                Require valid-user
        </Location>
</VirtualHost>
```



配置 git  客户端用户名密码


htpasswd -m -c /etc/httpd/conf.d/git-team.htpasswd   user1

123

123


chown apache:apache /etc/httpd/conf.d/git-team.htpasswd

chmod 640 /etc/httpd/conf.d/git-team.htpasswd

systemctl  restart httpd




建立测试仓库目录结构

cd /home && mkdir git && cd git

mkdir git-test && cd git-test

git init

设置权限

chown -R apache:apache .


设置允许push 

修改如下 文件
/home/git/git-test/.git/config

添加一下内容

[receive]
denyCurrentBranch = ignore










### 测试 clone 

linux  上测试

git clone https://www.git.com/git/git-test

Username for 'https://www.git.com': yandun
Password for 'https://yandun@www.git.com': 123

测试成功


windows 上测试

客户端

git config --global http.sslVerify false

git clone https://www.git.com/git/git-test

测试成功

### 测试 push 


touch a.md

git commit -m "add a file"

git  push

在 git  server 上查看git 客户端 推送上来的文件


git reset --hard





登录到远程的那个文件夹，使用
git config --bool core.bare true  





#### 报错1 :

remote: error: insufficient permission for adding an object to repository database ./objects


权限问题

cd  git-test

chown -R apache:apache .


####  报错2 :
remote: error: refusing to update checked out branch: refs/heads/master
remote: error: By default, updating the current branch in a non-bare repository


git默认拒绝了push操作，需要进行设置

vim /home/git/git-test/.git/config

[receive]
denyCurrentBranch = ignore



#### 报错3 :

登陆 git  server  看不见 已经推送上来的文件

cd /git-test

git reset --hard


#### 报错4 ：  无法 git  pull  参考链接 http://blog.csdn.net/lindexi_gd/article/details/52554159



git  pull

报错 ：

Unpacking objects: 100% (44/44), done.
From https://www.git.com/git/git-test
 * [new branch]      master     -> origin/master
fatal: refusing to merge unrelated histories



git pull --allow-unrelated-histories


报错：

Auto-merging a.md
CONFLICT (add/add): Merge conflict in a.md
Automatic merge failed; fix conflicts and then commit the result.


git  add .

git commit -m "d"



git pull --allow-unrelated-histories

执行成功

参考链接 


http://blog.csdn.net/jacolin/article/details/44014775


