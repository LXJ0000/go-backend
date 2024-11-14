## 目录结构
```
.env
api/
    controller/
    middleware/
    route/
assets/
bootstrap/
    app.go
    cache.go
    cron_job.go
    db.go
    env.go
    ...
docker-compose.yaml
internal/
main.go
pkg/
script/
utils/
var/
```
基本介绍：
- .env 配置：具体通过 bootstrap/env.go 加载到内存中，此外还将其全部加载到环境变量
- api 对外暴露的服务即 http 接口
 - router 路由 调用 controller 
 - controller 控制器 调用 internal/usecase
 - middleware 中间件包括鉴权、限流等 
- internal 内部服务
  - domain 领域 包括 model、usecase接口、repository接口
  - usecase service层 具体的实现 domain 中的 usecase 接口
  - repository 存储层 具体的实现 domain 中的 repository 接口
  - event 消息队列使用到的消息体
- assets 静态资源文件
- bootstrap 加载环境配置以及各种服务的启动包括Mysql、Redis、Kafka、Cron
- main.go 项目入口
- pkg 封装了对第三方服务的使用
- script 项目使用到的脚本
- utils 工具集合
- var 日志记录