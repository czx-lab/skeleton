```shell
go test -v -run=TestConfig ./test
```

#### 设置环境变量并下载项目依赖
```shell
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go mod download
```

#### 运行项目
```shell
go run ./cmd/main.go
```

#### 项目编译打包运行
```shell
go build ./cmd/main.go

// 静态编译
CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s'

// 运行项目
./main
```

#### 项目目录结构说明
```text
├─app
│  ├─command ---> 命令行
│  ├─controller
│  │    └─base.go ---> BaseController，主要定义了request参数验证器validator
│  ├─event
│  │  ├─entity ---> 事件实体目录
│  │  ├─listen ---> 事件监听执行脚本目录
│  │  └─event.go ---> 事件注册代码
│  │       
│  ├─middleware ---> 中间件代码目录
│  ├─request ---> 请求参数校验代码目录
│  │   └─request.go ---> 参数验证器
│  └─task ---> 定时任务代码目录
│     └─task.go ---> 注册定时任务脚本
├─cmd ---> 项目入口目录
│  └─cli ---> 项目命令行模式入口目录
├─config
│  └─config.yaml ---> 配置文件
├─internal ---> 包含第三方包的封装
├─router ---> 路由目录
│  └─router.go
├─storage ---> 日志、资源存储目录
│  └─logs
└─test ---> 单元测试目录
```

### 基础功能

---

#### 路由

该骨架的web框架是gin，所以路由定义可直接阅读Gin框架的文档。

在该骨架中定义注册路由需要在`router`文件夹下面的`router.go`文件中的`func (*AppRouter) Add(server *gin.Engine)`方法定义注册：

```go
server.GET("/foo", func(ctx *gin.Context) {
    ctx.String(http.StatusOK, "hello word!")
})
```

也可以通过自己定义路由的定义注册，只需要实现`github.com/czx-lab/skeleton/internal/server/router`下面的`Interface`接口。如下示例：
在router目录下定义了一个`CustomRouter`结构体，该结构体实现了`Interface`接口

```go
package router

import (
	"net/http"
	
	"github.com/czx-lab/skeleton/internal/server"
	"github.com/gin-gonic/gin"
)

type CustomRouter struct {
	server server.HttpServer
}

func NewCustom(srv server.HttpServer) *CustomRouter {
	return &CustomRouter{
		srv,
	}
}

func (*CustomRouter) Add(srv *gin.Engine) {
	srv.GET("/custom", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "custom router")
	})
}
```

> 需要注意的是，如果是自定义路由注册，需要修改项目`cmd`文件夹下面的`main.go`入口文件，通过`http.SetRouters(router.NewCustom(http))`注册给`gin`

#### 中间件

定义中间件与`gin`框架一样，该估计默认实现了panic异常的中间件，可以查看`internal/server/middleware`文件夹中的`exception.go`文件。

如果需要定义其他的中间件并加载注册，可以将定义好的中间件通过`server.HttpServer`接口的`SetMiddleware(middlewares ...middleware.Interface)`方法注册加载，
比如我们实现如下自定义全局中间件`middleware/custom.go`：
```go
type Custom struct{}

func (c *Custom) Handle() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        fmt.Println("Custom middleware exec...")
    }
}
```
然后在定义路由的地方使用`server.SetMiddleware(&middleware.Custom{})`注册中间件。
定义全局路由中间件可以参考`router/router.go`中的`New`方法。

> 如果是局部中间件，可以直接在具体的路由上注册，参考gin路由中间件的用法





