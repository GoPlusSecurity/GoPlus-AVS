# AVS 服务的主体

## 入口

cli 启动 `cmd/main.go`，目前支持的命令行参数 `start` `register-with-avs` 

配置使用`-c`参数传递env文件的绝对路径

## 主要的类

- Server: HTTP 服务，用于响应来自 Fanout service 发来的 Task。
- SecwareManager: AVS 服务的后台组织者，会维护 DockerRunner, SecwareMonitor，向 Gateway 上报状态等。
- SecwareMonitor: 监控 Secware 的健康状况。
- DockerRunner: 负责以 docker compose 方式启停 Secware，管理端口映射等。
- SecwareAccessor: 是 Secware 的统一 HTTP 接口。
