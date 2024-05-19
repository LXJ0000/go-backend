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
1. 安全问题：发送频率、验证码有效期、不能被暴力破解

#### 发送逻辑：
1. 没有key，发送
2. 有key
    - 没有过期时间，系统异常，拒绝发送
    - 有过期时间
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
