包信号实现对传入信号的访问。

信号主要用于类unix系统。为了使用这个在Windows和计划9上的包，见下面。

## 类型的信号

SIGKILL和SIGSTOP信号可能不会被程序捕获因此不受此包的影响。

同步信号是由程序中的错误触发的信号执行:SIGBUS, SIGFPE, SIGSEGV。这些只被考虑同步时引起的程序执行，而不是发送时使用os.Process。杀死或杀死程序或一些类似的机制。在一般来说，除了下面讨论的，Go程序将转换a同步信号变成运行时恐慌。

其余的信号是异步信号。他们不是由程序错误触发，而不是从内核或从其他项目。

在异步信号中，SIGHUP信号是在程序运行时发送的失去控制终端。时发送SIGINT信号用户在控制终端按下中断字符，默认情况下是^C (Control-C)SIGQUIT信号发送时在控制终端上的用户按下quit字符默认为^\ (control -反斜杠)。一般来说，你可以导致a程序简单地退出按^C，你可以使它退出通过按下^\来转储堆栈。

## Go程序中信号的默认行为

默认情况下，同步信号被转换为运行时恐慌。一个SIGHUP、SIGINT或SIGTERM信号会导致程序退出。一个SIGQUIT SIGILL SIGTRAP SIGABRT SIGSTKFLT SIGEMT或SIGSYS信号使程序退出并产生堆栈转储。一个SIGTSTP，或SIGTTOU信号获取系统认行为(这些信号是由外壳用于作业控制)。SIGPROF信号被处理直接通过Go运行时实现runtime. cpuprofile。其他信号会被捕捉，但不会采取行动。

如果启动Go程序时忽略SIGHUP或SIGINT(信号处理程序设置为SIG_IGN)，它们仍将被忽略。

如果Go程序是以一个非空信号掩码开始的，那就会通常是荣幸。然而，一些信号是明确解除阻塞的:同步信号SIGILL, SIGTRAP, SIGSTKFLT, SIGCHLD, SIGPROF，在GNU/Linux上，信号32 (SIGCANCEL)和33 (SIGSETXID)(SIGCANCEL和SIGSETXID在glibc内部使用)。子流程由操作系统。或由os/ Exec包继承修改后的信号掩码。

## 在Go程序中改变信号的行为

这个包中的函数允许程序改变运行方式程序处理信号。

Notify禁用给定异步集合的默认行为信号，而不是传递他们在一个或多个注册频道。具体来说，它适用于信号SIGHUP, SIGINT，SIGQUIT, SIGABRT和SIGTERM。它也适用于作业控制信号SIGTSTP、SIGTTIN和SIGTTOU，在这种情况下系统默认行为不会发生。它也适用于一些信号否则无法执行:SIGUSR1, SIGUSR2, SIGPIPE, SIGALRM，SIGCHLD SIGCONT SIGURG SIGXCPU SIGXFSZ SIGVTALRM SIGWINCHSIGIO SIGPWR SIGSYS SIGINFO SIGTHR SIGWAITING SIGLWP SIGFREEZESIGTHAW, SIGLOST, SIGXRES, SIGJVM1, SIGJVM2，以及任何实时信号在系统上使用。注意，并非所有这些信号都可用在所有系统。

如果程序启动时忽略了SIGHUP或SIGINT，并发送通知，则将为?安装一个信号处理程序这一信号将不再被忽视。如果，稍后，重置或为该信号调用Ignore，或在所有通道上调用Stop将该信号传递给Notify，该信号将再次为忽略了。属性将恢复系统默认行为信号，而忽略会导致系统忽略信号完全。

如果程序以非空信号掩码开始，则会出现一些信号将如上所述显式解除阻塞。如果调用了Notify对于一个被阻塞的信号，它将被解除阻塞。如果，稍后，重置是调用该信号，或停止被调用的所有通道传递通知该信号，该信号将再次被阻止。

## SIGPIPE

当一个Go程序写入一个损坏的管道时，内核将引发aSIGPIPE信号。

如果程序没有调用Notify来接收SIGPIPE信号，则行为取决于文件描述符号。写信给
破管文件描述符1或2(标准输出或标准输出)错误)将导致程序退出SIGPIPE信号。一个写其他文件描述符不会对损坏的管道采取任何操作SIGPIPE信号，写操作将失败，并出现EPIPE错误。

如果程序已经调用通知接收SIGPIPE信号，文件描述符的数量并不重要。SIGPIPE信号为传递到通知通道，写操作将会失败错误。

这意味着，在默认情况下，命令行程序的行为类似于典型的Unix命令行程序，而其他程序不会当写入一个封闭的网络连接时SIGPIPE崩溃。

## 使用cgo或SWIG的Go程序

在一个包含非Go代码的Go程序中，通常是C/ c++代码使用cgo或SWIG访问，Go的启动代码通常首先运行。它之前按照Go运行时的期望配置信号处理程序运行非go启动代码。如果非运行启动代码希望这样做安装自己的信号处理程序后，它必须采取一定的步骤来保持运行工作得很好。本节记录这些步骤和总体内容效果改变信号处理程序设置的非走代码可以在走程序。在极少数情况下，非Go代码可能会在Go之前运行代码，在这种情况下，下一节也适用。

如果Go程序调用的非Go代码没有改变任何信号处理程序或掩码，那么行为是相同的，为一个纯粹的走程序。

如果非执行代码安装了任何信号处理程序，则必须使用带有sigaction的SA_ONSTACK标志。如果不这样做，很可能会导致如果收到信号，程序就会崩溃。项目经常去使用有限的堆栈运行，因此设置一个备用信号堆栈。此外，Go标准库期望任何信号处理程序将使用SA_RESTART标志。如果不这样做，可能会导致一些库返回“中断的系统调用”错误的调用。

如果非执行代码为任何同步信号(SIGBUS, SIGFPE, SIGSEGV)，则应记录现有的Go信号处理程序。当这些信号出现时执行Go代码时，它应该调用Go信号处理程序(无论是否执行Go代码时发生的信号可以通过查看来确定PC传递到信号处理器)。否则一些运行时恐慌不会像预期的那样发生。

如果非执行代码为任何异步信号，它可以调用Go信号处理程序或不作为它选择。当然，如果它不调用Go信号处理程序，则上面描述的Go行为将不会发生。这可能是一个问题特别是SIGPROF信号。

非执行代码不应该改变任何线程上的信号掩码由Go运行时创建。如果非运行代码启动新的线程它自己可以随意设置信号掩码。

如果非执行代码启动一个新线程，则更改信号掩码，并且然后在该线程中调用Go函数，Go运行时将自动解除某些信号的阻塞:同步信号，SIGILL, SIGTRAP, SIGSTKFLT, SIGCHLD, SIGPROF, SIGCANCEL，和SIGSETXID。当Go函数返回时，非Go信号掩码将返回被恢复。

如果在未运行Go的非运行线程上调用Go信号处理程序代码时，处理程序通常将信号转发给非执行代码，如遵循。如果信号是SIGPROF，那么Go处理器就是SIGPROF什么都没有。否则，Go处理程序将删除自身，解除阻塞信号，并再次引发它，以调用任何non-Go处理程序或默认值系统处理程序。如果程序不退出，则Go处理程序重新安装自身并继续程序的执行。

## 调用Go代码的非Go程序

当Go代码使用-buildmode=c-shared这样的选项构建时，它会这样做作为现有非go程序的一部分运行。非go代码可能在Go代码启动时已经安装了信号处理程序(即在使用cgo或SWIG时，可能在不寻常的情况下发生;在这种情况下,这里的讨论适用于。For -buildmode=c-存档Go运行时将在全局构造函数时初始化信号。为-buildmode=c-shared运行时将初始化信号加载共享库。

如果Go运行时看到已有的SIGCANCEL或SIGSETXID信号(只在GNU/Linux上使用)，它将打开使用SA_ONSTACK标志，否则保留信号处理程序。

对于同步信号和SIGPIPE, Go运行时将安装a信号处理程序。它将保存任何现有的信号处理程序。如果一个同步信号在执行非Go代码时到达，即Go运行时会调用现有的信号处理程序而不是Go信号吗处理程序。

使用-buildmode=c-archive或-buildmode=c-shared来构建代码默认情况下不安装任何其他信号处理程序。如果有现有的信号处理程序，Go运行时将打开SA_ONSTACK标记，否则保留信号处理程序。如果Notify被调用异步信号，一个Go信号处理程序将为此安装信号。如果，稍后，重置该信号，原始的对该信号的处理将重新安装，恢复不运行如果有信号处理程序。

运行不使用-buildmode=c-archive或-buildmode=c-shared的代码为上面列出的异步信号安装一个信号处理器，并保存任何现有的信号处理程序。如果一个信号被发送到非go线程，它将像上面描述的那样工作，除非有一个现有的非go信号处理程序，该处理程序将被安装在发出信号之前。

## Window

在Windows上，a ^C (Control-C)或^BREAK (Control-Break)通常引起退出的程序。如果os调用Notify。中断，或中断将导致操作系统。中断被发送到通道上，程序将
不退出。如果在所有通过的通道上调用Reset或Stop若要通知，则将恢复默认行为。

另外，如果调用了Notify，并且Windows发送CTRL_CLOSE_EVENT，CTRL_LOGOFF_EVENT或CTRL_SHUTDOWN_EVENT通知进程syscall.SIGTERM返回。与Control-C和Control-Break不同，Notify可以当CTRL_CLOSE_EVENT，接收到CTRL_LOGOFF_EVENT或CTRL_SHUTDOWN_EVENT——进程会接收到仍然会被终止，除非它退出。但是接收系统调用。SIGTERM将在终止之前给流程一个清理的机会。

## 计划9

在计划9中，信号有类型系统调用。注意，它是一个字符串。调用使用系统调用进行通知。将导致该值被发送到通道时，该字符串被张贴作为一个注意。