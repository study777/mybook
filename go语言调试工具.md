<<<<<<< HEAD
#  安装 delve 调试工具

先设置$GOPATH
=======
# 先设置$GOPATH

vim /etc/profile

```
export GOROOT=/usr/local/go
export GOBIN=$GOROOT/bin
export PATH=$PATH:$GOBIN
export GOPATH=/go
```


>>>>>>> ccddf11a254065670cbd656d4a36729480087994


在GOPATH下创建三个目录  bin   pkg   src

# 在线安装

执行go get github.com/derekparker/delve/cmd/dlv

# 离线安装



mkdir -p /go/src/github.com/derekparker

cd  /go/src/github.com/derekparker

git clone https://github.com/derekparker/delve.git


/go/src/github.com/derekparker/delve

go install github.com/derekparker/delve/cmd/dlv


mkae install



dlv会被安装在$GOPATH/bin目录下（推荐）




go get github.com/derekparker/delve/cmd/dlv


delve 命令格式

dlv [command]

dlv [command]  --help



子命令讲解
debug   编译并调试源代码

命令格式  dlv debug [package] [flags]


./dlv debug --helo

dlv debug        [package]                [flags]
               包名或源代码文件名
               
     
               
               
exec —— 对编译过的可执行程序进行调试
  命令格式：dlv exec [./path/to/binary] [flags]
                              
attach —— 对后台的进程进行调试
  命令格式： dlv attach pid [flags]

connect —— 远程调试
  命令格式： dlv connect addr [flags]


其他命令

trace    通做正则表达式匹配包下的函数名 如果匹配成功此函数将会被跟踪调试  是dlv 的一个子命令 

version   查看当前dlv 的版本


run     过时 被debug 命令替代

test    和debug 差不多 专用于 调试 后缀为test 的测试文件




单点调试实战

./dlv --help

Delve是一个源代码级别的Go程序调试器。
			
			Delve通过控制流程的执行、变量评估和线程或协程态的信息提供以及CPU寄存器状态等等，使你能与程序交互。
			
			此工具的亮点在于，它为调试Go程序，提供了一个简单但强大的交互界面（非图形界面）

               
        
在 src 先创建my 文件夹

mkdir my

创建文件  main.go


package main

var b [2]byte

//长度为2的数组
func main() {
	for i := 0; i < 3; i++ {
		b[i] = byte(i)
	}
}

//针对这个数组进行3次循环赋值，而这个数组的长度只有2
//所以当执行到b[2]的时候就超出了数组的界限
//会报一个越界的错误





go run /gopath1/src/my/main.go 
panic: runtime error: index out of range    越界

goroutine 1 [running]:
panic(0x45b020, 0xc42000a0d0)               报恐慌
	/usr/local/go/src/runtime/panic.go:500 +0x1a1
main.main()
	/gopath1/src/my/main.go:8 +0x36          弟8行
exit status 2




调试看到更详细的信息


/gopath1/bin/dlv debug /gopath1/src/my/main.go

(dlv) help

print   打印某一个变量或者表达式

exit   退出

continue  让程序继续执行

break   设置断点

breakpoints  查看断点

clear  清除某个断点


list  列出源代码


step  单步调试














/gopath1/bin/dlv debug /gopath1/src/my/main.go  使用dlv打开一个程序
这个程序已经被编译过了， 这个编译的结果不会被保存，只是现在在调试的阶段，编译的结果还在




(dlv) ls

发现有一个箭头指向第八行   里面的内容是汇编语言  因为GO语言是基于汇编语言和C语言的
所以它程序的真正入口并不是 main函数，而从GO语言的逻辑看源代码，它的入口才在main包的main函数

所以第一步将main包的main函数设置为断点处


=>   8:		LEAQ	8(SP), SI // argv


(dlv) b main.main                  将main包的main函数设置为断点
Breakpoint 1 set at 0x40100f for main.main() ./main.go:6



(dlv) ls     再次ls 一下 还是在刚才的位置  因为程序还没有开始执行

=>   8:		LEAQ	8(SP), SI // argv   

(dlv) continue

=>   6:	func main() {


(dlv) ls   再次查看  
=>   6:	func main() {



接下来 继续执行程序 

(dlv) s
在 7 8 行循环执行

(dlv) p i   打印变量

当i  时2的时候

(dlv) p i
2


之后接着 s  会一直到最后 知道抛出异常


查看设置哪些断点 

(dlv) bp

Breakpoint unrecovered-panic  系统设置好的
Breakpoint 1 at 0x40100f for main.main() ./main.go:6 (1)    刚刚设置的

删除断点 1

(dlv) clear 1
Breakpoint 1 cleared at 0x40100f for main.main() ./main.go:6

再次查看断点
(dlv) bp

(dlv) clearall  删除所有断点   但是系统断点是不会被删除的

(dlv) bp
Breakpoint unrecovered-panic at 0x422f40 for runtime.startpanic() /usr/local/go/src/runtime/panic.go:543 (0)


重新设置断点 

(dlv) b main.main
Breakpoint 2 set at 0x40100f for main.main() ./main.go:6

查看断点
(dlv) bp
Breakpoint unrecovered-panic at 0x422f40 for runtime.startpanic() /usr/local/go/src/runtime/panic.go:543 (0)
Breakpoint 2 at 0x40100f for main.main() ./main.go:6 (0)

删除断点
(dlv) clearall
Breakpoint 2 cleared at 0x40100f for main.main() ./main.go:6


继续运行程序
(dlv) c


(dlv) c

Process 1448 has exited with status 2   表明程序退出  退出的装填是2




重新启动程序
(dlv) restart
Process restarted with PID 1467

查看cpu寄存器

(dlv) regs

Rip = 0x000000000044ccb0    Rip 的内容是一个内存地址  就是刚才箭头所指向的代码所保存的在的内存地址   44ccb0

(dlv) ls

=>   8:		LEAQ	8(SP), SI // argv   这位置的指令就保存在   44ccb0  的位置


设置断点    go逻辑入口点

(dlv) b main.main
Breakpoint 1 set at 0x40100f for main.main() ./main.go:6      


(dlv) continue   会在断点处停下来

(dlv) ls

=>   6:	func main() {


(dlv) s  
=>   7:		for i := 0; i < 3; i++ {   

查看  =>   7:		for i := 0; i < 3; i++ {      内存位置

(dlv) regs
Rip = 0x000000000040101d


dlv) p b

