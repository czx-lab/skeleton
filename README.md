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




