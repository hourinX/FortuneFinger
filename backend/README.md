# FortuneFinger Backend

FortuneFinger 后端重构版本。当前采用模块化单体结构，先保持一个进程运行，同时按照业务边界组织代码，为后续按模块拆分服务保留空间。

## 目录结构

```text
backend-copy/
├─ cmd/
│  └─ fortune.go                 # 程序入口
├─ etc/
│  └─ fortune.yaml              # 运行时配置
├─ migrations/                  # 数据库迁移及使用说明
├─ internal/
│  ├─ biz/                      # 无 HTTP 路由的可复用业务流程
│  ├─ config/                   # 配置字段映射与加载
│  ├─ mw/                       # 全局 HTTP 中间件
│  ├─ model/                    # 数据库表字段映射
│  ├─ pkg/                      # 不含业务语义的内部技术组件
│  │  ├─ auth/                  # JWT 生成与校验
│  │  ├─ captcha/               # 验证码生成
│  │  ├─ constant/              # 通用常量（待实现）
│  │  ├─ crypto/                # 编解码与加密能力（待实现）
│  │  ├─ errorx/                # 统一错误
│  │  ├─ gormc/                 # GORM 多数据源、连接池与日志适配
│  │  ├─ httpx/                 # HTTP 统一响应
│  │  ├─ logger/                # Zap 日志初始化
│  │  ├─ requestctx/            # 请求头统一注入与 Context 读取入口
│  │  │  └─ userid/             # X-User-ID 请求上下文
│  │  ├─ redisc/                # Redis Get/Set 客户端
│  │  └─ snowflake/             # 雪花 ID 生成器
│  ├─ repository/               # 可被 server 与 biz 复用的数据访问层
│  ├─ server/                   # 对外服务的业务模块
│  │  ├─ core/                  # 菜单、角色、权限、字典等核心能力
│  │  │  ├─ register.go         #   core 统一注册入口
│  │  │  ├─ router/             #   路由定义
│  │  │  ├─ handler/            #   HTTP 参数绑定与响应
│  │  │  ├─ logic/              #   业务逻辑
│  │  │  └─ types/              #   模块全部出入参
│  │  ├─ prize/                 # 奖项、奖品和库存业务
│  │  │  ├─ router/
│  │  │  ├─ handler/
│  │  │  ├─ logic/
│  │  │  └─ types/
│  │  └─ user/                  # 登录、认证和系统用户业务
│  │     ├─ router/
│  │     ├─ handler/
│  │     ├─ logic/
│  │     └─ types/
│  ├─ startup/                  # 应用装配与生命周期
│  │  ├─ app.go                 # 应用启动、运行和关闭入口
│  │  ├─ datasource/            # 数据源装配
│  │  ├─ dependency/            # 基础设施依赖装配与释放
│  │  ├─ httpserver/            # HTTP Server 配置与关闭
│  │  └─ router/                # 全局路由和中间件装配
│  ├─ svc/                      # 进程级基础设施上下文
│  └─ trun/                     # 测试使用的真实运行环境
└─ README.md
```

空目录表示结构已经确定但尚未开始实现。Git 不保存空目录，首次添加代码时再纳入版本控制。

## 各层职责

| 目录 | 职责 |
|---|---|
| `server` | 接收 HTTP 请求并完成面向接口的业务编排，按照 `core`、`user`、`prize` 等业务域划分，不按照管理端、抽奖端或移动端划分 |
| `biz` | 不直接暴露 HTTP 路由的完整业务能力，例如开票、充值和余额查询，可被多个 `server` 模块调用 |
| `model` | 定义数据库表字段映射，一张表对应一个模型文件，不包含业务逻辑 |
| `repository` | 封装模型与数据库操作，集中维护 SQL，供 `server` 和 `biz` 共同复用 |
| `svc` | 持有配置、日志、数据库、Redis、验证码等进程级基础设施，不存放 repository 和请求上下文 |
| `startup` | 负责加载、装配、启动、优雅关闭和资源释放 |
| `pkg` | 与具体业务无关的内部技术组件 |
| `trun` | 仅服务于测试，使测试可以复用与正式程序一致的 GORM、Redis 和日志环境；不属于正式服务启动链路 |

## 依赖方向

```text
server ──> biz ──> repository ──> model
   └─────────────> repository ──> model

server / biz / repository ──> svc 或所需的基础设施接口
startup ──> config + pkg + svc + server
```

约束如下：

1. `server` 可以直接调用 `biz`，也可以在简单 CRUD 场景直接调用 `repository`。
2. `biz` 不得依赖 `server`，也不得调用 `server` 内部的数据访问实现。
3. 相同数据库操作只在 `repository` 实现一次，禁止在 `server` 与 `biz` 中各写一份 SQL。
4. `repository` 只负责持久化操作，不包含 HTTP 处理和业务规则。
5. `pkg` 不得依赖 `server`、`biz` 或 `repository`。
6. `startup` 是装配入口；业务包不主动从全局容器查找依赖。

## 业务模块规划

- `server/core`：菜单、角色、权限、Casbin、字典、系统设置等后台核心能力。
- `server/user`：验证码、登录、刷新凭证和系统用户管理；公钥与登出等接口按后续认证方案补充。
- `server/prize`：奖项、奖品、图片和库存管理。
- 后续抽奖流程形成稳定边界后，再增加 `server/lottery`。

三个前端应用共用同一套后端接口。接口根据业务域归档，不使用 `admin`、`lottery`、`h5` 作为客户端类型目录。

## `biz` 的组织方式

`biz` 按可复用业务能力划分，而不是作为公共工具目录。例如：

```text
internal/biz/invoice/
├─ invoice.go                   # 开票业务流程
├─ types.go                     # 业务出入参
├─ repository.go               # 当前业务需要的数据接口
├─ provider.go                 # 第三方提供方抽象
├─ nucarf_client.go            # 指定提供方的请求实现
├─ nucarf_types.go             # 指定提供方的协议结构
└─ nucarf_sign.go              # 指定提供方的签名实现
```

`biz` 可以声明自己需要的窄接口，由 `repository` 中已有实现满足。Go 接口采用隐式实现，因此无需复制已有 SQL。
