# GO项目目录架构
如果你尝试学习 Go，或者你正在为自己建立一个 PoC 或一个玩具项目，这个项目布局是没啥必要的。从一些非常简单的事情开始(一个 main.go 文件绰绰有余)。

当有更多的人参与这个项目时，你将需要更多的结构，包括需要一个 toolkit来方便生成项目的模板，尽可能大家统一的工程目录布局。

首先大家可以看下这一篇：
https://github.com/golang-standards/project-layout/blob/master/README_zh.md

- `/cmd`目录：本项目的主干。每个应用程序的目录名应该与你想要的可执行文件的名称相匹配(例如，/cmd/myapp)。
 `go build` 默认会将 bin 项目编译成 `main.main()` 函数所在文件所在的文件夹名。
如下所示：
```
├── cmd/
│   ├── demo/
│   │   ├── demo    # <- go build 输出
│   │   └── main.go
│   └── demo1/
│       ├── demo1   # <- go build 输出
│       └── main.go
```

但是注意不要在这个目录中放置太多代码。如果你认为代码可以导入并在其他项目中使用，那么它应该位于 /pkg 目录中。如果代码不是可重用的，或者你不希望其他人重用它，请将该代码放到 /internal 目录中。

- `/internal`私有应用程序和库代码。这是你不希望其他人在其应用程序或库中导入代码。

Go 1.4 之后强制保证。引用其他包的 `internal` 子包无法通过编译。
你可以选择向 internal 包中添加一些额外的结构，以分隔共享和非共享的内部代码。这不是必需的(特别是对于较小的项目)，但是最好有有可视化的线索来显示预期的包的用途。

你的实际应用程序代码可以放在 /internal/app 目录下(例如 /internal/app/myapp)，这些应用程序共享的代码可以放在 /internal/pkg 目录下(例如/internal/pkg/myprivlib)。

```
├── internal/
│   ├── app/           # <- 存放各 bin 应用专用的程序代码
│   │   └── myapp/     # <- 存放 myapp 专用的程序代码
│   ├── demo/          # <- 也可忽略 app 层。存放 demo 的专用程序代码。如果只有一个 bin 应用，这个层也可以去除。
│   │   ├── biz/
│   │   ├── data/
│   │   └── service/
│   └── pkg/           # <- 存放各 bin 共享程序代码，但因为有 internal 下，其他项目无法引用。
│       └── myprivlib/ # <- 按功能分 lib 包
```

- `/pkg`外部应用程序可以使用的库代码(例如 `/pkg/mypubliclib`)。  
其他项目可以导入这些库，所以**在这里放东西之前要三思**

要显示地表示目录中的代码对于其他人来说可安全使用的，使用 `/pkg` 目录是一种很好的方式。跟`/internal`恰恰相反，一个公用，一个私用，岂不美哉。

`/pkg` 目录内，可以参考 go 标准库的组织方式，按照功能分类。
 `/internal/pkg` 一般用于项目内的跨多应用的公共共享代码，但其作用域仅在单个项目内。
```
├── pkg/
│   ├── cache/
│   │   ├── memcache/
│   │   └── redis/
│   └── conf/
│       ├── dsn/
│       ├── env/
│       ├── flagvar/
│       └── paladin/
```

当根目录包含大量非 Go 组件和目录时，这也是一种将 Go 代码分组到一个位置的方法，这使得运行各种 Go 工具变得更加容易组织。

```
.
├── README.md
├── docs/
├── example/
├── go.mod
├── go.sum
├── misc/
├── pkg/
├── third_party/
└── tool/
```

## kit 工具包项目

- kit 库：工具包/基础库/框架库

每个公司都应当为不同的微服务建立一个统一的 kit 工具包项目.

基础库 kit 为独立项目，公司级建议只有一个（通过行政手段保证），按照功能目录来拆分会带来不少的管理工作，因此建议合并整合。

> To this end, the Kit project is not allowed to have a vendor folder.
> If any of packages are dependent on 3rd party packages, 
> they must always build against the latest version of those dependences.

kit 项目必须具备的特点:
* 统一
* 标准库方式布局
* 高度抽象
* 支持插件

文件布局可如下：

```
├── cache/      # <- 缓存
│   ├── memcache/
|   |   └── test
│   └── redis/
|       └── test
├── conf/
│       ├── dsn/
│       ├── env/
│       ├── flagvar/
│       └── paladin/
│           └── apollo/
│               └── internal/
│                   └── mockserver/
├── container/ #<-容器池
│   ├── group/
│   ├── pool/
│   └── queue/
│       └── aqm/
├── darabase/ #<-数据库
│   ├── hbase/
│   ├── sql/
│   └── tidb/
├── ecode/
│   └──  types1/   
│   
│      
├── log/  #<- 日志
│   └── internal/
│       ├── core/
│       └── filewriter/
```

## 微服务项目