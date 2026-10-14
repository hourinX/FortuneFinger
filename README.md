<p align="center">
  <img src="frontend/admin/src/assets/branding/fortune-finger-wordmark.png" alt="FortuneFinger" width="300">
</p>

<h1 align="center">FortuneFinger</h1>

<p align="center">
  面向企业年会、庆典和现场活动的开源抽奖系统
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white" alt="Go 1.26">
  <img src="https://img.shields.io/badge/Vue-3-42B883?logo=vuedotjs&logoColor=white" alt="Vue 3">
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-Apache%202.0-blue.svg" alt="Apache License 2.0"></a>
</p>

FortuneFinger 计划覆盖活动配置、参与人员管理、奖项与奖品管理、抽奖策略、现场抽奖大屏和移动端互动。项目采用 Go 与 Vue 3 构建，以模块化单体作为当前形态，并按照业务边界组织代码，为后续扩展和按模块拆分服务保留空间。

> [!IMPORTANT]
> 项目目前处于积极开发阶段。后台认证、RBAC、菜单、角色、用户、组织和字典管理等基础能力已经建立；奖品、活动、抽奖流程、大屏和移动端仍在开发中，暂不建议直接用于正式活动。

## 功能

当前已经实现：

- JWT 登录、刷新凭证和退出登录
- 图片验证码与可切换的增强验证
- 登录失败次数控制和阶梯式安全限制
- 基于角色、菜单和权限编码的 RBAC
- 动态菜单与前端动态路由
- 角色层级、角色权限配置和权限预览
- 系统用户、角色分配和用户状态管理
- 集团、公司、部门三级组织结构
- 字典类型与字典数据管理，支持统一配置选项名称、显示颜色、排序和启用状态
- 统一错误码、HTTP 响应、日志和请求超时处理
- MySQL、Redis、雪花 ID 和数据库迁移基础设施
- 管理端多页签和页面缓存控制

正在规划或开发：

- 系统设置与文件存储
- 活动、参与人员和签到管理
- 奖项、奖品与库存管理
- 随机抽奖、权重策略和抽奖记录
- 抽奖大屏与移动端参与页面
- 操作日志、定时任务和首次启动引导

## 技术栈

| Layer | Technologies |
| --- | --- |
| Backend | Go, Gin, GORM, Viper, Zap |
| Authentication | JWT, bcrypt, captcha, enhanced verification |
| Persistence | MySQL, Redis |
| Frontend | Vue 3, Vite, Pinia, Vue Router |
| UI | Ant Design Vue, Less |
| Migration | golang-migrate compatible SQL migrations |

数据库基础层支持 MySQL、PostgreSQL、SQL Server 和 Oracle 数据源；FortuneFinger 当前默认并主要使用 MySQL。

## 架构

后端采用模块化单体。HTTP 接口按照业务域组织，而不是按照管理端、大屏端或移动端分别复制业务逻辑。

```text
HTTP request
    │
    ▼
router ──> handler ──> logic ──> repository ──> model
                           │
                           └────> reusable biz workflow

startup ──> config + infrastructure + service context + router
```

主要职责：

| Directory | Responsibility |
| --- | --- |
| `backend/internal/server` | 按业务域组织 HTTP 路由、参数处理和业务逻辑 |
| `backend/internal/biz` | 不直接暴露路由、可被多个模块复用的完整业务流程 |
| `backend/internal/repository` | 数据库访问和可复用持久化操作 |
| `backend/internal/model` | 数据库表映射 |
| `backend/internal/pkg` | 与具体业务无关的内部基础组件 |
| `backend/internal/startup` | 配置加载、依赖装配、服务启动和优雅关闭 |
| `backend/internal/svc` | 进程级配置与基础设施上下文 |
| `frontend/admin` | Vue 3 管理端 |

更详细的后端分层说明见 [backend/README.md](backend/README.md)。

## 快速开始

### 环境要求

- Go 1.26+
- Node.js 20.19+ 或 22.12+
- MySQL 8+
- Redis
- npm

### 1. 获取代码

```bash
git clone https://github.com/hourinX/FortuneFinger.git
cd FortuneFinger
```

### 2. 准备数据库

创建使用 `utf8mb4` 的 MySQL 数据库，然后通过 [golang-migrate](https://github.com/golang-migrate/migrate) 应用迁移：

```bash
migrate -path backend/migrations \
  -database "mysql://USER:PASSWORD@tcp(localhost:3306)/fortunefinger" up
```

迁移规则和常用命令见 [backend/migrations/README.md](backend/migrations/README.md)。

> 当前尚未实现首次管理员创建向导。首次运行前仍需准备管理员、角色和权限策略数据，这部分将在首次启动流程中统一解决。

### 3. 配置后端

编辑 `backend/etc/fortune.yaml`，至少检查以下配置：

- HTTP 监听地址和端口
- MySQL 地址、账号、密码和数据库名
- Redis 地址、密码和数据库编号
- JWT Secret 与有效期
- Snowflake Machine ID

不要在公开仓库、日志或问题反馈中提交生产密码、JWT Secret 和其他凭据。

### 4. 启动后端

```bash
cd backend
go mod download
go run ./cmd/fortune.go -f ./etc/fortune.yaml
```

默认监听地址为 `http://localhost:6327`。

### 5. 启动管理端

打开另一个终端：

```bash
cd frontend/admin
npm ci
npm run dev
```

开发服务器默认由 Vite 启动，并将 `/api` 请求代理到 `http://localhost:6327`。

## 项目结构

```text
FortuneFinger/
├── backend/
│   ├── cmd/                    # 服务入口
│   ├── etc/                    # 运行配置
│   ├── migrations/             # 数据库迁移
│   └── internal/
│       ├── biz/                # 可复用业务流程
│       ├── config/             # 配置加载与映射
│       ├── model/              # 数据库模型
│       ├── mw/                 # HTTP 中间件
│       ├── pkg/                # 内部基础组件
│       ├── repository/         # 持久化层
│       ├── server/             # 业务接口模块
│       ├── startup/            # 应用装配与生命周期
│       ├── svc/                # 基础设施上下文
│       └── trun/               # 测试运行环境
├── frontend/
│   ├── admin/                  # 管理端
│   └── shared/                 # 前端共享能力
├── docs/                       # 项目文档和数据库说明
├── LICENSE
└── NOTICE
```

## 开发与检查

后端检查：

```bash
cd backend
go test ./...
go vet ./...
go build ./...
```

前端检查：

```bash
cd frontend/admin
npm run lint
npm run build
```

提交代码前请阅读 [AGENTS.md](AGENTS.md) 中的项目结构和编码约定。数据库结构变更必须通过新迁移完成，不应修改已经发布的迁移文件。

## 参与贡献

欢迎通过 Issue 报告问题、提出功能建议，也欢迎提交 Pull Request。

提交问题时请提供必要的复现步骤、运行环境和错误日志，并删除密码、令牌、个人信息及其他敏感数据。提交代码时请保持变更范围清晰，并确保相关构建和检查通过。

## 许可证

FortuneFinger 使用 [Apache License 2.0](LICENSE) 开源。署名信息见 [NOTICE](NOTICE)。
