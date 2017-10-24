6  Git  工具

git 工具下载和安装

设置git bash

光标颜色 

默认 白色 
改成 绿色

光标形状

默认 下划线

改成 块状

光标闪烁

默认闪烁

改成 不闪

字体大小

默认 9号字

改成 14 号字




点击命令行窗口 左上角的 图标

options 选项

looks   Cursor  光标颜色

Cursor   Block  光标形态  块状

Blinking  闪烁勾 去掉




Text 

Select  12   保存

字符集 zh_CN  utf8


在指定目录下启动 git bash

右键选中   滑轮 粘贴


目录下  鼠标右键 git bash  here 

git --version

设置git 参数

显示当前git 配置

git config --list

设置提交仓库时的用户名信息

git config --global user.name "study777"

git config --global user.email "1569092380@qq.com"


用户家目录下 隐藏文件 时配置文件

cd  ~

cat ~/.gitconfig



7 补充 Git Bash 操作





git  log

j  下一行

k  上一行

G  最下面

gg  最上面


/搜索字符

查找时 上下翻  u   n

退出 q



查看git 状态

git status



变动一个文件
vim README.md

提交带暂存区

git add .

提交并输入提交信息

gti commit

git status
On branch master






7 Git  命令1

工作区    workspace
 
暂存区  index /stage

仓库区   repository  


git 本地命令

在当前目录新建一个Git 代码库

git  init

下载一个项目和它的整个代码历史
git clone   



添加指定文件到暂存区

git  add  file1  file2


删除工作区文件，并且将这次删除放入暂存区

git  rm  file1  file2 


改文件名字，并且将这个改名 放入暂存区

git mv  file-origin  file-renamed

提交暂存区 到仓库

git  commit  -m  message



直接从工作区 提交到仓库

前提是  该文件已经有仓库中的历史版本

git  commit  -a  -m  message 

显示变更信息

git  status

显示当前分支的历史版本


git  log

git  log  --oneline







本地仓库仓库：

cd  /c/my-file/git-test

mkdir demo

cd demo/


git init
Initialized empty Git repository in C:/my-file/git-test/demo/.git/




git status
On branch master   master  分支

No commits yet


vim README.md


git status

git add README.md
warning: LF will be replaced by CRLF in README.md.


git status


git commit -m "add a file"

git status
On branch master
nothing to commit, working tree clean



echo "second add someting to this file" >> README.md




git status


对于在仓库从已经存在的文件 可以直接 从工作区 提交到 仓库
git commit -a -m "add something"


git status
On branch master
nothing to commit, working tree clean



查看提交信息


git log
commit c81c88981f045e45a9d36f537b02467a58aca3b3 (HEAD -> master)
Author: study777 <1569092380@qq.com>
Date:   Tue Oct 24 18:11:48 2017 +0800

    add something

commit 7b92c619fd47f6ccaff402272ccd139d14a1450e
Author: study777 <1569092380@qq.com>
Date:   Tue Oct 24 18:08:41 2017 +0800

    add a file



git  log


git  show  hashcode



8 Git  命令2


Remote 远程仓库
workspace  工作区
index/stage  暂存区
repository  仓库区/本地仓库


对于remote

remote    pull >   workspace

remote  fetch/clone  >   repository

repository  push  >   remote

对于 repository
repository  checkout > workspace
index  commit  >  repository
remote  fetch/clone  >   repository
repository  push  >   remote



对于 index/stage

index  commit >   repository
workspace add  > index

对于 workspace
workspace add  > index
repository  checkout > workspace
remote    pull >   workspace





增加远程仓库 并命名
git remote add shortname  url

将本地的提交推送到远程仓库
 
 
git push  remote  branch

将远程仓库的提交拉下到本地

git pull remote branch
 
 
 
github  点击加号  新建一个仓库  demo-test

将本地文件夹  demo  重命名为 demo-test

按照 github 提示 在本地添加远程仓库

git remote add origin https://github.com/study777/demo-test.git

查看远程仓库

git remote -v
 
 
将本地的master 分支推送到远程仓库  origin  

git push origin master


 
 
 






















