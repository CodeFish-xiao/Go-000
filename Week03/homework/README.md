# Week03 作业题目：
基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。

问题来了：
## signal 是什么？
信号(Signal)是Linux, 类Unix和其它POSIX兼容的操作系统中用来进程间通讯的一种方式。一个信号就是一个异步的通知，发送给某个进程，或者同进程的某个线程，告诉它们某个事件发生了。

当信号发送到某个进程中时，操作系统会中断该进程的正常流程，并进入相应的信号处理函数执行操作，完成后再回到中断的地方继续执行。
如果目标进程先前注册了某个信号的处理程序(signal handler),则此处理程序会被调用，否则缺省的处理程序被调用。

## 在Go程序中signal的默认行为

默认情况下，同步信号会转换为运行时紧急情况。SIGHUP，SIGINT或SIGTERM信号导致程序退出。SIGQUIT，SIGILL，SIGTRAP，SIGABRT，SIGSTKFLT，SIGEMT或SIGSYS信号会导致程序以堆栈转储退出。SIGTSTP，SIGTTIN或SIGTTOU信号获得系统默认行为（shell将这些信号用于作业控制）。SIGPROF信号由Go运行时直接处理以实现runtime.CPUProfile。其他信号将被捕获，但不会采取任何措施。

如果Go程序在忽略SIGHUP或SIGINT的情况下启动（信号处理程序设置为SIG_IGN），则它们将保持忽略状态。

如果Go程序以非空的信号掩码启动，通常会很荣幸。但是，某些信号被明确地解除了阻塞：同步信号SIGILL，SIGTRAP，SIGSTKFLT，SIGCHLD，SIGPROF，以及在GNU / Linux上，信号32（SIGCANCEL）和33（SIGSETXID）（glibc在内部使用了SIGCANCEL和SIGSETXID）。由os.Exec或os / exec程序包启动的子进程将继承修改后的信号掩码。