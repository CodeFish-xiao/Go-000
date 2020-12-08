# Week03 作业题目：
基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。

问题来了：
## signal 是什么？
信号(Signal)是Linux, 类Unix和其它POSIX兼容的操作系统中用来进程间通讯的一种方式。一个信号就是一个异步的通知，发送给某个进程，或者同进程的某个线程，告诉它们某个事件发生了。

当信号发送到某个进程中时，操作系统会中断该进程的正常流程，并进入相应的信号处理函数执行操作，完成后再回到中断的地方继续执行。
如果目标进程先前注册了某个信号的处理程序(signal handler),则此处理程序会被调用，否则缺省的处理程序被调用。

## Golang中怎么调用它呢？
golang中提供了signal对于信号的操作
~~~go
func Ignore(sig ...os.Signal)//忽略
func Ignored(sig os.Signal) bool//判断是否被忽略
func Notify(c chan<- os.Signal, sig ...os.Signal)//唤醒
func Reset(sig ...os.Signal)//重置
func Stop(c chan<- os.Signal)//停止
~~~
**Ignore**
~~~go
func Ignore(sig ...os.Signal)
~~~
Ignore将忽略提供的信号。如果程序接收到它们，则不会发生任何事情。Ignore将撤消先前对提供的信号进行通知的任何调用的影响。如果未提供信号，则所有输入信号都将被忽略。

**Ignored**
~~~go
func Ignored(sig os.Signal) bool
~~~
Ignored报告当前信号是否忽略。

**Notify**
~~~go
func Notify(c chan<- os.Signal, sig ...os.Signal)
~~~
Notify使包信号将传入信号转发给c，如果没有信号，则将所有传入信号转发给c，否则仅发送提供的信号。

包信号不会阻塞发送给c:调用者必须确保c有足够的缓冲区空间来保持预期的信号速率。对于只用于通知一个信号值的通道，大小为1的缓冲区就足够了。

允许使用同一通道多次调用Notify:每次调用扩展发送到该通道的信号集。从集合中移除信号的唯一方法是调用Stop。

允许使用不同的通道和相同的信号多次调用通知:每个通道分别接收传入信号的副本。

**Reset**
~~~go
func Reset(sig ...os.Signal)
~~~
Reset将解除任何先前调用的效果，以通知所提供的信号。如果没有提供信号，所有信号处理程序将被重置。

**Stop**
~~~go
func Stop(c chan<- os.Signal)
~~~
Stop使包信号停止向c中继传入信号，解除之前所有调用的效果，使用c通知。当Stop返回时，保证c不会再收到信号。


[这里有包的具体描述](./singnal包描述.md)也可看官方文档（链接是我机翻的）