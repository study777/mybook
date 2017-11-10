#  环境准备  搭建一个svn 并载入测试代码

搭建  http  svn


yum install httpd

httpd  -v

Server version: Apache/2.4.6 (CentOS)

Server built:   Nov 14 2016 18:04:44



yum install subversion

svnserve --version

svnserve, version 1.7.14 


yum install mod_dav_svn





find / -name mod_dav_svn.so


/usr/lib64/httpd/modules/mod_dav_svn.so


find / -name mod_authz_svn.so

/usr/lib64/httpd/modules/mod_authz_svn.so


mkdir /var/www/svn

svnadmin create /var/www/svn/spring-hello-world



chown -R apache:apache /var/www/svn/spring-hello-world/





touch /var/www/svn/passwd  

htpasswd /var/www/svn/passwd admin  


cp /var/www/svn/spring-hello-world/conf/authz /var/www/svn/authz


cat /var/www/svn/authz 
[/]
admin = rw




touch /etc/httpd/conf.d/subversion.conf


cat /etc/httpd/conf.d/subversion.conf 
```
<Location /svn>
    DAV svn
    SVNParentPath /var/www/svn
    AuthType Basic
    AuthName "Authorization SVN"
    AuthzSVNAccessFile /var/www/svn/authz
    AuthUserFile /var/www/svn/passwd
    Require valid-user
</Location>
```

systemctl start httpd.service

测试 svn

浏览器 访问

http://172.16.2.30/svn/spring-hello-world/

账号密码都是 

admin


导入测试代码

git clone https://github.com/nichochen/mybank-demo-maven


svn import -m "First SVN Repo"  mybank-demo-maven   http://172.16.2.30/svn/spring-hello-world


命令行测试导入的代码

命令行测试

svn co http://172.16.2.30/svn/spring-hello-world --no-auth-cache









cd /opt


创建的builder image 名字是 my-tomcat8-svn   工作空间是 my-tomcat8-svn-workspace

s2i create my-tomcat8-svn   my-tomcat8-svn-workspace

cd my-tomcat8-svn-workspace/


编辑Dockerfile 
vim /opt/my-tomcat8-svn-workspace/Dockerfile

```
# openshift-tomcat8-svn
FROM docker.io/centos
# TODO: Put the maintainer name in the image metadata
MAINTAINER huliaoliao
# TODO: Rename the builder environment variable to inform users about application you provide them
ENV BUILDER_VERSION 1.0
ENV LC_CTYPE en_US.UTF-8
# TODO: Set labels used in OpenShift to describe the builder image
LABEL io.openshift.s2i.scripts-url=image:///usr/libexec/s2i \
      io.k8s.description="Platform for building tomcat" \
      io.k8s.display-name="builder tomcat" \
      io.openshift.expose-services="8080:http" \
      io.openshift.tags="builder,tomcat,java,etc."
# TODO: Install required packages here:
#COPY ./CentOS-Base.repo /etc/yum.repos.d/CentOS-Base.repo
RUN yum makecache &&  yum install -y java-1.8.0-openjdk subversion maven && yum clean all -y
COPY ./.s2i/bin/ /usr/libexec/s2i
# TODO (optional): Copy the builder files into /opt/app-root
COPY ./tomcat8/ /opt/app-root/tomcat8
# TODO: Copy the S2I scripts to /usr/local/s2i, since openshift/base-centos7 image sets io.openshift.s2i.scripts-url label that way, or update that label
#COPY ./s2i/bin/ /usr/libexec/s2i
# TODO: Drop the root user and make the content of /opt/app-root owned by user 1001
RUN useradd -m tomcat -u 1002 && \
    chmod -R a+rw /opt && \
    chmod -R a+rw /opt/app-root && \
    chmod a+rwx /opt/app-root/tomcat8/* && \
    chmod +x /opt/app-root/tomcat8/bin/*.sh && \
    rm -rf /opt/app-root/tomcat8/webapps/* && \
    rm -rf /usr/share/maven/conf/settings.xml
#ADD ./settings.xml /usr/share/maven/conf/
# This default user is created in the openshift/base-centos7 image
USER 1002
# TODO: Set the default port for applications built using this image
EXPOSE 8080
ENTRYPOINT []
# TODO: Set the default CMD for the image
CMD ["/usr/libexec/s2i/usage"]
```




创建tomcat 目录

mkdir  /opt/my-tomcat8-svn-workspace/tomcat8


wget -P /opt/my-tomcat8-svn-workspace/   http://archive.apache.org/dist/tomcat/tomcat-8/v8.5.5/bin/apache-tomcat-8.5.5.tar.gz



cd /opt/my-tomcat8-svn-workspace/

tar xzvf apache-tomcat-8.5.5.tar.gz

mv /opt/my-tomcat8-svn-workspace/apache-tomcat-8.5.5/*   /opt/my-tomcat8-svn-workspace/tomcat8/

编写S2I脚本

vim /opt/my-tomcat8-svn-workspace/.s2i/bin/assemble



```
#!/bin/bash -e
#
# S2I assemble script for the 'nico-tomcat' image.
# The 'assemble' script builds your application source ready to run.
#
# For more information refer to the documentation:
#       https://github.com/openshift/source-to-image/blob/master/docs/builder_image.md
#
# Restore artifacts from the previous build (if they exist).
#
if [ "$1" = "-h" ]; then
        # If the 'nico-tomcat' assemble script is executed with '-h' flag,
        # print the usage.
        exec /usr/libexec/s2i/usage
fi
# Restore artifacts from the previous build (if they exist).
#
if [ "$(ls /tmp/artifacts/ 2>/dev/null)" ]; then
  echo "---> Restoring build artifacts"
  mv /tmp/artifacts/. ./
fi
echo "---> Installing application source"
WORK_DIR=/tmp/src;
cd $WORK_DIR;
if [ ! -z ${SVN_URI} ] ; then
  echo "Fetching source from Subversion repository ${SVN_URI}"
  svn co ${SVN_URI} --username=${SVN_USERNAME} --password=${SVN_PASSWORD} --no-auth-cache
  export SRC_DIR=`basename $SVN_URI`
  echo "Finished fetching source from Subversion repository ${SVN_URI}"
  cd $WORK_DIR/$SRC_DIR/
  mvn package -Dmaven.test.skip=true;
else
  echo "SVN_URI not set, skip Subverion source download";
fi
find /tmp/src/ -name '*.war'|xargs -i mv -v {} /opt/app-root/tomcat8/webapps/ROOT.war
echo "---> Building application from source"
```


vim /opt/my-tomcat8-svn-workspace/.s2i/bin/run

```
#!/bin/bash -e
exec /opt/app-root/tomcat8/bin/catalina.sh run
```


构建builder 镜像
cd /opt/my-tomcat8-svn-workspace/

开始构建

make


docker images | grep my-tomcat8-svn
my-tomcat8-svn                                                 latest              d3f41ec740d3        26 seconds ago      519 MB

构建成功


为builder 进行 打 标记 并推送至仓库

docker tag my-tomcat8-svn  172.16.2.31:5000/my-tomcat8-svn

docker push 172.16.2.31:5000/my-tomcat8-svn

在 openshift master 上导入镜像 至 平台的 IS

oc  import-image  172.16.2.31:5000/my-tomcat8-svn   -n openshift --confirm --insecure

输出内容如下

```
The import completed successfully.

Name:			my-tomcat8-svn
Namespace:		openshift
Created:		Less than a second ago
Labels:			<none>
Annotations:		openshift.io/image.dockerRepositoryCheck=2017-10-26T08:16:26Z
Docker Pull Spec:	docker-registry.default.svc:5000/openshift/my-tomcat8-svn
Image Lookup:		local=false
Unique Images:		1
Tags:			1

latest
  tagged from 172.16.2.31:5000/my-tomcat8-svn
    will use insecure HTTPS or HTTP connections

  * 172.16.2.31:5000/my-tomcat8-svn@sha256:b8588664d04287d8c3e7259c44e5bd108807d3ceff44a9457b1adf26142eef36
      Less than a second ago

Image Name:	my-tomcat8-svn:latest
Docker Image:	172.16.2.31:5000/my-tomcat8-svn@sha256:b8588664d04287d8c3e7259c44e5bd108807d3ceff44a9457b1adf26142eef36
Name:		sha256:b8588664d04287d8c3e7259c44e5bd108807d3ceff44a9457b1adf26142eef36
Created:	Less than a second ago
Image Size:	204.2 MB (first layer 7.853 MB, last binary layer 73.39 MB)
Image Created:	5 minutes ago
Author:		huliaoliao
Arch:		amd64
Command:	/usr/libexec/s2i/usage
Working Dir:	<none>
User:		1002
Exposes Ports:	8080/tcp
Docker Labels:	build-date=20170911
		io.k8s.description=Platform for building tomcat
		io.k8s.display-name=builder tomcat
		io.openshift.expose-services=8080:http
		io.openshift.s2i.scripts-url=image:///usr/libexec/s2i
		io.openshift.tags=builder,tomcat,java,etc.
		license=GPLv2
		name=CentOS Base Image
		vendor=CentOS
Environment:	PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
		BUILDER_VERSION=1.0
```






编辑 template   并关联 IS  builder image

vim my-tomcat8-svn-template.json

```
cat my-tomcat8-svn-template.json 
{
    "kind": "Template",
    "apiVersion": "v1",
    "metadata": {
        "annotations": {
            "iconClass" : "icon-tomcat",
            "description": "Application template for JavaEE WAR deployment with Tomcat 8 from svn repo."
        },
        "name": "my-tomcat8-svn-template"
    },
    "labels": {
        "template": "my-tomcat8-svn-template"
    },
    "parameters": [
        {
            "description": "Tomcat 8.5.5",
            "name": "IMG_VERSION",
            "displayName":"Image Version",
            "value": "latest",
            "required": true
        },
        {
            "description": "The name for the application.",
            "name": "APPLICATION_NAME",
            "displayName":"Application Name",
            "value": "",
            "required": true
        },
        {
            "description": "Custom hostname for service routes.  Leave blank for default hostname, e.g.: <application-name>.<project>.<default-domain-suffix>",
            "name": "APPLICATION_HOSTNAME",
            "displayName":"Application Hostname",
            "value": ""
        },
        {
            "description": "Subversion source URI for application",
            "name": "SVN_URI",
            "displayName":"Subversion source URI",
            "value": "",
            "required": true
        },
        {
            "description": "Subversion Username",
            "name": "SVN_USERNAME",
            "displayName":"Subversion Username",
            "value": "",
            "required": true
        },
        {
            "description": "Subversion Password",
            "name": "SVN_PASSWORD",
            "displayName":"Subversion Password",
            "value": "",
            "required": true
        }
    ],
    "objects": [
        {
            "kind": "Service",
            "apiVersion": "v1",
            "spec": {
                "ports": [
                    {
                        "port": 8080,
                        "targetPort": 8080
                    }
                ],
                "selector": {
                    "deploymentConfig": "${APPLICATION_NAME}"
                }
            },
            "metadata": {
                "name": "${APPLICATION_NAME}",
                "labels": {
                    "application": "${APPLICATION_NAME}"
                },
                "annotations": {
                    "description": "The web server's http port."
                }
            }
        },
        {
            "kind": "Route",
            "apiVersion": "v1",
            "id": "${APPLICATION_NAME}-http-route",
            "metadata": {
                "name": "${APPLICATION_NAME}-http-route",
                "labels": {
                    "application": "${APPLICATION_NAME}"
                },
                "annotations": {
                    "description": "Route for application's http service."
                }
            },
            "spec": {
                "host": "${APPLICATION_HOSTNAME}",
                "to": {
                    "name": "${APPLICATION_NAME}"
                }
            }
        },
        {
            "kind": "ImageStream",
            "apiVersion": "v1",
            "metadata": {
                "name": "${APPLICATION_NAME}",
                "labels": {
                    "application": "${APPLICATION_NAME}"
                }
            }
        },
        {
            "kind": "BuildConfig",
            "apiVersion": "v1",
            "metadata": {
                "name": "${APPLICATION_NAME}",
                "labels": {
                    "application": "${APPLICATION_NAME}"
                }
            },
            "spec": {
                "strategy": {
                    "type": "Source",
                    "sourceStrategy": {
                        "from": {
                            "kind": "ImageStreamTag",
                            "namespace": "openshift",
                            "name": "my-tomcat8-svn:latest"
                        },
                        "env": [
                                    {
                                        "name": "SVN_URI",
                                        "value": "${SVN_URI}"
                                    },
                                    {
                                        "name": "SVN_USERNAME",
                                        "value": "${SVN_USERNAME}"
                                    },
                                    {
                                        "name": "SVN_PASSWORD",
                                        "value": "${SVN_PASSWORD}"
                                    }

                       ]
                    }
                },
                "output": {
                    "to": {
                        "kind": "ImageStreamTag",
                        "name": "${APPLICATION_NAME}:latest"
                    }
                },
                "triggers": [
                    {
                        "type": "GitHub",
                        "github": {
                            "secret": "${GITHUB_TRIGGER_SECRET}"
                        }
                    },
                    {
                        "type": "Generic",
                        "generic": {
                            "secret": "${GENERIC_TRIGGER_SECRET}"
                        }
                    },
                    {
                        "type": "ImageChange",
                        "imageChange": {}
                    }
                ]
            }
        },
        {
            "kind": "DeploymentConfig",
            "apiVersion": "v1",
            "metadata": {
                "name": "${APPLICATION_NAME}",
                "labels": {
                    "application": "${APPLICATION_NAME}"
                }
            },
            "spec": {
                "strategy": {
                    "type": "Recreate"
                },
                "triggers": [
                    {
                        "type": "ImageChange",
                        "imageChangeParams": {
                            "automatic": true,
                            "containerNames": [
                                "${APPLICATION_NAME}"
                            ],
                            "from": {
                                "kind": "ImageStream",
                                "name": "${APPLICATION_NAME}"
                            }
                        }
                    }
                ],
                "replicas": 1,
                "selector": {
                    "deploymentConfig": "${APPLICATION_NAME}"
                },
                "template": {
                    "metadata": {
                        "name": "${APPLICATION_NAME}",
                        "labels": {
                            "deploymentConfig": "${APPLICATION_NAME}",
                            "application": "${APPLICATION_NAME}"
                        }
                    },
                    "spec": {
                        "containers": [
                            {
                                "name": "${APPLICATION_NAME}",
                                "image": "${APPLICATION_NAME}",
                                "imagePullPolicy": "Always",
                                "readinessProbe": {
                                    "exec": {
                                        "command": [
                                            "/bin/bash",
                                            "-c",
                                            "curl http://localhost:8080"
                                        ]
                                    }
                                },
                                "ports": [
                                    {
                                        "name": "http",
                                        "containerPort": 8080,
                                        "protocol": "TCP"
                                    }
                                ],
                                "env": [
                                    {
                                        "name": "SVN_URI",
                                        "value": "${SVN_URI}"
                                    }

                                ]
                            }
                        ]
                    }
                }
            }
        }
    ]
}
```

从 json 文件创建template


oc create -n openshift -f /root/my-tomcat8-svn-template.json


web console 选择


my-tomcat8-svn-template 

svn url  

http://172.16.2.30/svn/spring-hello-world

账号

admin

密码

admin


