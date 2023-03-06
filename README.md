# Wox

`Wox` 是一个插件化系统，它作为插件中心使用。

[toc]

## 用法

如果你需要使用`Wox`，首先，你应该引入该程序包

```bash
go get github.com/sumwai/wox
```

然后才可以在程序中使用

```go
package main

import (
    "github.com/sumwai/wox"
)

func main () {
    NewPM().Load("plugins").Run()
}
```

以上代码将会新建一个`PlugManager`，并调用`Load`方法，遍历`./plugins/`下的所有`.so`插件，最后使用`Run`运行所有插件

## 插件定义

插件的定义需要实现`Plug`接口

```go
type Plug interface {
    Run() // 入口函数
}
```

并导出变量`Plugin`

```go
var Plugin
```

如果你需要定义你的插件版本，插件名称，可以使用导出变量

```go
var (
    Name = "Plugin Name"
    Version = "1.0.0"
)
```

以下是一个简单的插件定义(仅实现Plug接口)，该插件每秒在日志中输出一条消息`Timer ticker`

```go
package main

import (
    "log"
    "time"
)

var (
    Plugin      Timer
    Name        = "Timer"
    Description = "a simple plugin, print 'time ticker' every second"
    Version     = "1.0.0"
)

type (
    Timer struct{}
)

func (t Timer) Run() {
    ticker := time.NewTicker(time.Second)
    for {
        <-ticker.C
        log.Println("time ticker")
    }
}
```

该插件运行后将会输出

```text
2023/03/06 20:54:30 loaded [Timer] (v 1.0.0) - a simple plugin, print 'time ticker' every second
2023/03/06 20:54:30 run [Timer] (v 1.0.0)
2023/03/06 20:54:31 time ticker
...
...
^Csignal: interrupt
```
