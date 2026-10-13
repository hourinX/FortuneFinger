# FortuneFinger

Go + Gin 分层单体。配置直接读取 backend/etc/fortune.yaml。

第一次启动前必须阅读 [FIRST_START.md](docs/FIRST_START.md)。

## 启动

填写 YAML 中的 MySQL / Redis 连接信息，准备 fortunefinger 数据库。在 backend 目录执行：

    go run fortune.go -f etc/fortune.yaml

默认监听 0.0.0.0:6327，Redis 默认开启。当前用户路由包含验证码、登录和刷新令牌接口。

响应：{"code":0,"msg":"ok","data":{"message":"FortuneFinger"}}

## 基础包

从本机 drama-go-refactor 迁入以下能力，未复制测试和微服务依赖：

- internal/pkg/gormc：源项目实际拼写为 gormc。连接池、zap SQL 日志，实例注入；启动 Ping 使用 context，失败关闭连接。SQL 日志不展开绑定参数。
- internal/pkg/logger：指定的 zap 编码和文件轮转方案，沿用当前配置及关闭接口，不引入全局单例、环境推导或 tracing。
- internal/pkg/response：提供 Ok / Error，返回 code、msg、可选 data。
- internal/pkg/errs：response 必要的错误码依赖，未知错误不暴露内部详情。
- internal/pkg/utils/idutil：雪花 ID；Snowflake.Enabled=true 时创建 svc.IDGen。
- internal/pkg/utils/pageutil：分页规范化、偏移量。
- internal/pkg/utils/timeutil：时间格式化、可选时间解析。

源项目 logger 定义了 Sanitize，但 admin 没有调用。它只是手动 map 掩码工具，不是 zap 的运行依赖；当前没有调用，已删除。不要主动记录密码或完整令牌。

## 路由

router.New 使用 gin.New()，只挂自定义 middleware.Recovery。没有 SetTrustedProxies、Gin 默认日志或额外 Observe。每个模块一个路由文件，在 RegisterXxx 内创建 group，在 New 中集中注册。JWT 按需挂载。

startup 使用 Application 管理生命周期：NewApplication 加载配置、创建 ServiceContext、组装 router、server、listener；Run 运行 HTTP 服务；Close 关闭 HTTP、listener 和进程依赖。config 只保存、加载、校验配置；MySQL DSN 在 ServiceContext 初始化时组装。

异常响应通过 response.Error 输出。请求超时返回 HTTP 503：

    {"code":10009,"msg":"请求超时"}

后端请求执行期限默认 15 秒，HTTP 写超时 30 秒；超时取消 context 并拒绝迟到响应。后续业务继续传递 c.Request.Context()。超时不能强杀 goroutine 或撤销已提交写入，当前响应缓冲方式不用于流式接口。

frontend/shared/request.mjs 保留默认 20 秒超时。Captcha 目前只有配置，没有登录验证码接口。

## 检查

在 backend 目录执行 go build ./... 和 go vet ./...。未进行真实 MySQL / Redis 联调。
