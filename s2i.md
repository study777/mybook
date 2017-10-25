#  S2I  环境搭建
cd  /opt

wget https://github.com/openshift/source-to-image/releases/download/v1.1.0/source-to-image-v1.1.0-9350cd1-linux-amd64.tar.gz


tar zxvf source-to-image-v1.1.0-9350cd1-linux-amd64.tar.gz  -C /usr/bin/


s2i version

s2i v1.1.0

# 准备 创建 builder image 所需的base 镜像
docker pull maven:3.3-jdk-7

cd  /opt

s2i   create    tomcat-s2i      tomcat-s2i



编辑 docker  file 
> /opt/tomcat-s2i/Dockerfile 

vim  /opt/tomcat-s2i/Dockerfile


```
#tomcat-s2i
FROM maven:3.3-jdk-7
MAINTAINER xxx
LABEL io.openshift.s2i.scripts-url=image:///usr/libexec/s2i \
      io.k8s.description="Tomcat S2I Builder" \
      io.k8s.display-name="tomcat s2i builder 1.0" \
      io.openshift.expose-services="8080:http" \
      io.openshift.tags="builder,tomcat"
WORKDIR /opt

ADD ./apache-tomcat-8.5.5.tar.gz /opt

RUN useradd -m tomcat -u 1001 && \

chmod -R a+rw /opt && \
chmod a+rwx /opt/apache-tomcat-8.5.5/* && \
chmod +x /opt/apache-tomcat-8.5.5/bin/*.sh && \
rm -rf /opt/apache-tomcat-8.5.5/webapps/*
COPY ./.s2i/bin/ /usr/libexec/s2i
USER 1001
EXPOSE 8080
ENTRYPOINT []
CMD ["usage"]
```




编辑S2I 脚本文件
大多为默认配置  斜体为 最后追加内容

```
>  /opt/tomcat-s2i/.s2i/bin/assemble

vim /opt/tomcat-s2i/.s2i/bin/assemble

#!/bin/bash -e
#
#S2I assemble script for the 'tomcat-s2i' image.
#The 'assemble' script builds your application source so that it is ready to run.
#
#For more information refer to the documentation:
#https://github.com/openshift/source-to-image/blob/master/docs/builder_image.md
#

if [[ "$1" == "-h" ]]; then
        # If the 'tomcat-s2i' assemble script is executed with '-h' flag,
        # print the usage.
        exec /usr/libexec/s2i/usage
fi

#Restore artifacts from the previous build (if they exist).
#
if [ "$(ls /tmp/artifacts/ 2>/dev/null)" ]; then
  echo "---> Restoring build artifacts..."
  mv /tmp/artifacts/. ./
fi

echo "---> Installing application source..."
cp -Rf /tmp/src/. ./

echo "---> Building application from source..."
# TODO: Add build steps for your application, eg npm install, bundle install

** cp -Rf /tmp/src/.  ./ **
mvn -Dmaven.test.skip=true package
find . -type f -name '*.war'|xargs -i cp {} /opt/apache-tomcat-8.5.5/webapps/
mvn clean





> /opt/tomcat-s2i/.s2i/bin/run
vim /opt/tomcat-s2i/.s2i/bin/run 

#!/bin/bash -e
bash -c "/opt/apache-tomcat-8.5.5/bin/catalina.sh run"
```



将所需的tomcat 实现至 S2I 工作目录


wget -P /opt/tomcat-s2i/ http://archive.apache.org/dist/tomcat/tomcat-8/v8.5.5/bin/apache-tomcat-8.5.5.tar.gz



在 s2i 工作目录下 执行 make   生成builder image

cd /opt/tomcat-s2i/

make

出现 Successfully built 4ef7d4e84678  类似输出 表示

builder image 成功创建


查看 生成的 buildere 

docker images | grep tomcat-s2i
tomcat-s2i                                                     latest              4ef7d4e84678        About a minute ago   617.8 MB

#  测试builderImage

指定 一个 remote  git 仓库 将代码注入 builder image   并生成一个 带有代码运行时 环境的 新 image




创建一个附有代码的 image  test-app     他的builder image 是  tomcat-s2i  
代码位置 是   https://github.com/nichochen/mybank-demo-maven 


命令执行位置 可不在 上节的工作目录  任何目录下都可



s2i build  https://github.com/nichochen/mybank-demo-maven      tomcat-s2i     test-app

如果成功创建 test-app  会有如下类似输出

[INFO] ------------------------------------------------------------------------
[INFO] BUILD SUCCESS
[INFO] ------------------------------------------------------------------------
[INFO] Total time: 14.630 s
[INFO] Finished at: 2017-10-18T14:21:43+00:00
[INFO] Final Memory: 11M/262M
[INFO] ------------------------------------------------------------------------






也可以将 代码 事先 git  clone 到本地 

删除之前的 test-app  

docker rmi test-app

cd /opt/tomcat-s2i/


git clone https://github.com/nichochen/mybank-demo-maven

s2i build  /opt/tomcat-s2i/mybank-demo-maven        tomcat-s2i     test-app



#  测试由  builder image 生成的test-app 

docker run -it -p 8080:8080 test-app



当出现如下信息时 需要等待一会

18-Oct-2017 14:37:54.844 INFO [localhost-startStop-1] org.apache.jasper.servlet.TldScanner.scanJars At least one JAR was scanned for TLDs yet contained no TLDs. Enable debug logging for this logger for a complete list of JARs that were scanned but no TLDs were found in them. Skipping unneeded JARs during scanning can improve startup time and JSP compilation time.



之后会有输出 即可正常测试 访问

curl http://172.16.2.31:8080


# 在openshift 上使用自己创建的builder image

如下是在和openshift 集群 网络互通的一个服务器上 的操作

具体可灵活配置



导入之前创建的builder image

docker load -i tomcat-s2i.tar 


创建一个用于测试 builder image 的  registry



docker run -d -p 5000:5000 --restart=always --name registry registry:2


在 openshift 集群个节点 和 本机都需要如下操作

vim /etc/sysconfig/docker

添加如下

INSECURE_REGISTRY='--insecure-registry 172.16.2.31:5000'

systemctl  restart docker



为builder image  tag  并推搡至 新创建的仓库


docker  tag tomcat-s2i  172.16.2.31:5000/tomcat-s2i


docker push 172.16.2.31:5000/tomcat-s2i


在 openshift  master 上 如下操作


oc  import-image 172.16.2.31:5000/tomcat-s2i  -n openshift --confirm --insecure 


编辑 新导入的 S2I 使 openshift 识别为 builder  image

红色位置 为更改处

oc edit is tomcat-s2i  -n openshift

将


spec:
  lookupPolicy:
    local: false
  tags:
  - annotations: null
    from:
      kind: DockerImage
      name: 172.16.2.31:5000/tomcat-s2i
      
      
改为

tags:
  - annotations:
      tags: builder
    from:
      kind: DockerImage
      name: 172.16.2.31:5000/tomcat-s2i
      

# 在openshift 上测试builder image
<font face="黑体">我是黑体字</font>

<font face="微软雅黑">我是微软雅黑</font>

<font face="STCAIYUN">我是华文彩云</font>

<font color=#0099ff size=12 face="黑体">黑体</font>

<font color=#00ffff size=3>null</font>

<font color=gray size=5>gray</font>
