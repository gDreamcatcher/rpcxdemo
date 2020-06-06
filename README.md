## RPCX-Demo

本仓库中实现了rpcx的基本功能, server, client, httpclient, 以及启动一个HTTP Server接受HTTP请求转发到tcp  server

### 运行
```
# 启动tcp  server
go run main.go -method server

# 启动 tcp client
go run main.go -method client

# 启动 http client
go run main.go -mehod httpclient

# 启动 http server转发到tcp server
go run main.go -method httpserver
```