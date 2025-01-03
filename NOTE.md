# go-backend
## 1. Gin Session 存储的实现
单机单实例部署，`memstore`；多实例部署，`redis`，确保不同实例都能访问到 `session`
## 2. 刷新登陆状态
- 固定间隔时间刷新
- 长短token
## 3. 参数如何设定
性能测试
## 4. JWT JSON WEB TOKEN
Header Payload Signature 三部分组成
## 5. JWT 优缺点
优点：不依赖第三方存储、适合在分布式环境下使用、性能好（没有 Redis 访问之类）
缺点：对加密依赖较大，比 Session 任意泄露、不要再 JWT 中存放敏感信息
## 6. 保护系统
1. 正常用户会不会搞崩你的系统？
2. 如果有人攻击你的系统，如何解决？
### 限流
1. 如何标识对象？限流对象可以用 IP，APP 端可以考虑使用设备序列号
2. 限流阈值？限制的阈值不是很小，就可以解决用一个 IP 多个用户的问题

- 为限流添加对应的监控和警告
- 对不需要登陆就可以访问的接口限流
- 对核心业务接口限流

为 Gin 插件库添加限流插件
- 单机限流
    - 令牌桶
    - 漏桶
    - 滑动窗口
    - 固定窗口
- 基于 Redis 限流
- 基于 Redis IP 限流
## 7. Kubernetes 容器编排平台
- Pod 实例
- Service 服务
- Deployment 管理
### 安装
[参考文档](https://www.qikqiak.com/post/deploy-k8s-on-win-use-wsl2/)
[参考视频](https://www.bilibili.com/video/BV1Ru41137s2/?spm_id_from=333.1007.top_right_bar_window_history.content.click&vd_source=2cb41caee9551fbf13c606149026e31c)
```bash
kubectl apply -f k8s-webook-deployment.yaml
kubectl apply -f k8s-webook-service.yaml
kubectl get deployments
kubectl get pods
kubectl get services
```
## 8. k8s 部署 MySQL
## 9. k8s 部署 Redis
```bash
cache-cli -h localhost -p 16379
```
## 10.
```bash
kubectl delete deployment webook
kubectl delete deployment webook-mysql
kubectl delete deployment webook-cache

kubectl delete service webook
kubectl delete service webook-mysql
kubectl delete service webook-cache

kubectl delete pv mysql-pv
kubectl delete pvc mysql-claim

kubectl delete ingress webook-ingress

kubectl get pod
kubectl get deployment
kubectl get service
kubectl get pv
kubectl get pvc

kubectl apply -f k8s-webook-deployment.yaml
kubectl apply -f k8s-webook-service.yaml
kubectl apply -f k8s-cache-service.yaml
kubectl apply -f k8s-cache-deployment.yaml
kubectl apply -f k8s-mysql-service.yaml
kubectl apply -f k8s-mysql-pv.yaml
kubectl apply -f k8s-mysql-pvc.yaml
kubectl apply -f k8s-mysql-deployment.yaml
kubectl apply -f k8s-ingress-nginx.yaml
```
## 11. 启动配置
```bash
go build -tags=k8s -o webook .
```
## 12. 压测
```bash
wrk -t4 -d5s -c50 -s ./script/wrk/register.lua http://localhost:8080/user/register
# t 线程 d 持续时间 c 并发数 s 后接测试脚本
wrk -t4 -d5s -c50 -s ./script/wrk/login.lua http://localhost:8080/user/login
```
## 13. 性能优化
### 1. 缓存
崩了？
1. 加载数据库，做好兜底，数据库限流
2. 不加载数据库，用户体验差

- 主从集群

## 14. 多种登陆方式
### 需求分析
1. 参考竞品
2. 从不同角度分析：
    1. 功能：
    2. 非功能： 安全 - 拓展 - 性能
3. 从正常、异常流程两个角度分析
### 腾讯云 SDK
```bash
go get -v -u github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common
```
### 验证码服务
1. 安全问题：发送频率、验证码有效期（超时 or 已使用）、不能被暴力破解（超过一定错误次数则失效）

#### 发送逻辑 (redis)： code:biz:number
1. 没有key，发送
2. 有key
    - 没有过期时间，说明系统异常，拒绝发送
    - 有过期时间（需要判定发生频率 即1分钟发生一次吧 15-14=1）
        - 多于14分钟（设定key生命15分钟），发送频繁，拒绝发送
        - 少于14分钟，发送
#### 验证逻辑
1. 验证码不存在
    - 提示发送验证码
2. 验证码存在
    - 验证次数小于3
        - 匹配，确认
        - 不匹配，返回验证码错误
    - 验证次数大于3，返回验证码错误

## 15. wire
```bash
go install github.com/google/wire/cmd/wire@latest
```
## 16. 快慢路径
快路径 触发降级操作，只走快路径 即系统资源不不足，只服务已经注册过的用户
## 17. 数据库唯一索引问题
设置唯一字段为 null, sql.NullString, 为什么不使用指针，需要解引用，需要判空
## 18. 依赖注入
A 依赖于 B，也就是 A 要调用 B 上面的方法，那么 A 在初始化的时候要求传进去一个创建好的 B

不使用依赖注入缺点：
1. 深度耦合依赖的初始化过程
2. 需要定义额外的配置信息
3. 缺乏拓展性
4. 测试不友好
5. 难以复用，如 DB、Redis 客户端

## 19. 控制反转 IOC
依赖注入 是 控制反转 的一种实现方式，还有一种叫 依赖查找

## 20. mock
```bash
go install go.uber.org/mock/mockgen@latest

mockgen -source=./utils/service/user.go -package=svcmock -destination=./utils/service/mocks/user_gen.go
```

## 21. 技术选型
1. 功能性：覆盖需求
2. 社区和支持度：活跃社区、文档齐全、搜索引擎..
3. 非功能性：易用、用户友好、扩展性（定制）、性能

## 22. 配置模块
```bash
go get github.com/spf13/viper
```
## 23. etcd
```bash
git clone git@github.com:etcd-io/etcd.git
cd etcdctl/
go install .

etcdctl --endpoints=127.0.0.1:23790 put /webook "$(<dev.yaml)"
```

## 24. TDD
测试驱动开发, 先写测试再写实现（专注于某个功能的实现）
1. 理清楚接口如何定义，体会**用户**使用起来是否舒适
2. 考虑主流程、异常流程
### 核心流程
1. 根据对需求的理解，**初步定义接口**。
2. **根据接口定义测试**
3. 执行核心循环
    - 增加测试用例
    - 提供、修改实现
    - 执行测试用例

## 25. Cron

## 26. 压测 k6

## 27. 第三方服务治理
针对一切跟第三方打交道的地方，都要做好容错
### 核心思路
- 尽量不要搞垮第三方
- 万一第三方奔溃，你自己的服务还能够稳定运行
具体到短信服务这里：
- 短信服务商都有保护自己系统的机制，要小心不要触发。比如短信服务商的限流机制
- 短信服务商可能崩溃 比如网络崩溃 你要做好容错机制

### 自动切换服务商
问题在于什么时候服务商出现问题：
- 频繁收到超时响应
- 收到 EOF 响应或者 UnexpectedEOF 响应
- 响应时间很长

#### 第一次策略：fail over
如果出现错误直接切换服务商重试
#### 第二种策略：动态判定服务商状态
错误率、响应时间增长率、CPU、内存、网络 IO 等等

## 数据迁移
### 不停机数据迁移
- 一边迁移数据一边产生新数据或者老数据被更新
- 迁移数据的时候 不能对数据库造成太大的压力 否则会影响应用的正常运行

难点：数据始终处于变动之中

四个阶段：
1. **业务读写源表**，在此阶段要完成目标表的数据初始化过程
2. **双写阶段，以源表为准**，在此阶段数据会被双写到源表和目标表中，并且读是读源表，如果出现不一致，以源表的数据为准
3. **双写阶段，以目标表为准**，在此阶段数据保持双写，但是读以目标表为准，并且修复数据的时候以目标表为准
4. **业务读写目标表** 

具体步骤比较复杂：
- 创建目标表
- 用源表的数据初始化目标表
- 执行一次校验并且修复数据，此时用源表数据修复目标表数据
- 业务代码开启双写，此时读源表，并且先写源表，数据以源表为准
- 开启增量校验和数据恢复，保持一段时间
- 切换双写顺序，此时读目标表，并且先写目标表，数据以目标表为准
- 继续保持增量校验和数据修复
- 切换为目标单写，读写都只操作目标表 

#### 全量校验与修复
难点：一条一条数据比对，数据量特别大，怎么尽快完成全量校验与修复

答案：并发

整个流程可以划分为2个步骤：校验，如果发现不一致，则修复

比较经典的做法：
- 如果发现不一致则立刻修复。这些操作都是同一个 goroutine 来执行
- 如果发现数据不一致，立刻交给另外一个 goroutine，也可以引入 channel。
- 如果发现数据不一致，则发送消息到消息队列，消费者消费了再去修复数据

最好采用消息队列的方案，因为我们要保护好目标表，也就是通过 Kafka 解耦和削峰。可以很容易控制住消费者的消费速率，也就间接控制住了目标表的写入速率

校验的基本思路：从源表中取出数据，再根据主键去目标表中找出对应的数据，比较所有的字段是否相等。

坑点：
1. 数据库类型众多，能转成 Go 语言类型吗
2. 转成的 Go 语言类型，都是可比较的吗
3. 浮点数之类的，从数据库都出来之后，转化为 Go 的数据类型，精度有没有损失？还能不能比较？

方案选型：
- 针对每一张表都写一个 DAO 查询方法，而后写一个比较方法。
- 借助泛型来写一个通用的查询方法，而后要求实体实现具体的比较方法
- 直接用最底层的 []byte 来接受任何数据，而后直接比较 []byte

