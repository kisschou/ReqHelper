# ReqHelper
---

![go](https://img.shields.io/github/go-mod/go-version/kisschou/ReqHelper?color=green&style=flat-square) ![commit](https://img.shields.io/github/last-commit/kisschou/ReqHelper?color=green&style=flat-square) ![license](https://img.shields.io/github/license/kisschou/ReqHelper?color=green&style=flat-square)

## 关于

从上行的请求中获取出请求的参数.

因为上行请求的类型不同, 会导致数据参数的接收方式不同, 所以有了这个脚本做一个汇总.

<br />

## 用法

首先你得依靠自己的力量在你的项目上安装好`Golang环境`, 并且确认你即将导入包的`Golang项目`是可用的, 然后:

1. 在需要使用的脚本中引入包:

```go
import "github.com/kisschou/ReqHelper"
```

<br />

2. 在需要使用的地方使用:

```go
import (
    "fmt"
    "net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
    rh := ReqHelper.New(r)
    fmt.Println("请求的地址: ", rh.Host)
    fmt.Println("请求的子地址: ", rh.Path)
    fmt.Println("请求的IP: ", rh.IpAddr)
    fmt.Println("请求的头部: ", rh.Header)
    fmt.Println("请求的参数: ", rh.Params)
    // ...
}

func main() {
    http.HandleFunc("/", indexHandler)
    http.ListenAndServe(":8000", nil)
}

r := ReqHelper.New(v.Req)
```

<br />

## Licence

Copyright (c) 2022-present Kisschou.
