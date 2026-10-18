<p align="center">
  <img src="frontend/admin/src/assets/branding/fortune-finger-wordmark.png" alt="FortuneFinger" width="320">
</p>

<p align="center">
  面向企业年会、庆典和现场活动的开源抽奖系统
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white" alt="Go 1.26">
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-Apache%202.0-blue.svg" alt="Apache License 2.0"></a>
  <img src="https://img.shields.io/badge/status-active%20development-ef4444" alt="Active development">
</p>

FortuneFinger 是一套可自行部署的企业活动抽奖系统，覆盖活动准备、现场签到、抽奖执行和结果留存，适用于企业年会、庆典和其他现场活动。

> [!IMPORTANT]
> FortuneFinger 仍在持续开发，当前版本适合本地体验和参与开发，尚未形成可直接用于正式活动的稳定发行版。

## 功能特性

- **多企业活动管理**：适配集团、子公司和部门结构，分别管理参与人员、奖品和活动数据。
- **扫码签到**：通过现场动态二维码进入签到页面，支持自由签到、名单签到、签到时段和重复签到规则。
- **签到大屏**：实时展示签到人数、最新入场人员和活动照片墙。
- **奖品与轮次**：按企业维护奖项和奖品，为不同时间安排多轮抽奖。
- **现场开奖**：提供适合投屏的抽奖大屏、开奖动画和中奖结果展示。
- **抽奖规则**：控制重复中奖、同奖项中奖、单人中奖上限和阳光普照奖发放方式。
- **结果留存**：保存签到记录、抽奖轮次和中奖名单，便于活动结束后查询。
- **权限管理**：按角色分配菜单和操作权限，区分系统管理员、活动管理员及只读人员。
- **存储与通知**：支持 MinIO、阿里云 OSS 和标准 SMTP 邮件服务。

## 使用流程

1. 创建企业组织并配置活动管理人员。
2. 维护参与名单、奖项、奖品和抽奖轮次。
3. 在现场签到大屏展示二维码，参与者使用手机完成签到。
4. 到达预定时间后通过抽奖大屏执行开奖。
5. 查看并留存签到结果和中奖名单。

## 运行 FortuneFinger

### 环境准备

运行当前开发版本需要：

- Go 1.26+
- Node.js 20.19+ 或 22.12+
- MySQL 8+
- Redis
- npm
- MinIO，或可用的阿里云 OSS

### 获取代码

```bash
git clone https://github.com/hourinX/FortuneFinger.git
cd FortuneFinger
```

### 初始化数据库

先创建一个使用 `utf8mb4` 的 MySQL 数据库，再执行项目迁移：

```bash
migrate -path backend/migrations \
  -database "mysql://USER:PASSWORD@tcp(localhost:3306)/fortunefinger" up
```

数据库迁移工具和命令说明见 [迁移说明](backend/migrations/README.md)。

### 配置服务

打开 `backend/etc/fortune.yaml`，根据本机环境修改：

- MySQL 连接信息
- Redis 连接信息
- JWT 密钥
- MinIO 或阿里云 OSS 连接信息
- 服务监听地址和端口

配置文件中的示例凭据只适用于本地开发。部署前必须替换所有密码、访问密钥和 JWT 密钥，并避免将生产凭据提交到公开仓库。

### 启动后端

```bash
cd backend
go mod download
go run ./cmd/fortune.go -f ./etc/fortune.yaml
```

后端默认监听 `http://localhost:6327`。

### 启动管理端

另开一个终端：

```bash
cd frontend/admin
npm ci
npm run dev
```

根据终端显示的地址打开管理端。开发服务器会将 `/api` 请求转发至本地后端。

> [!NOTE]
> 当前版本尚未提供首次启动向导。数据库迁移只创建表结构，初始超级管理员、角色、菜单和权限数据仍需自行准备。该流程会在正式发行版前完善。

## 部署提示

正式部署时，请更换配置文件中的示例凭据，并为数据库、Redis、对象存储和 SMTP 使用独立账号。公网演示环境不应向演示账号开放菜单管理、系统设置和其他敏感权限。

## 获取帮助

请通过 [GitHub Issues](https://github.com/hourinX/FortuneFinger/issues) 报告问题或提出建议。提交问题时请附上运行环境、复现步骤和必要的错误信息，并删除其中的密码、令牌及个人信息。

## 开源许可

FortuneFinger 使用 [Apache License 2.0](LICENSE) 发布，署名信息见 [NOTICE](NOTICE)。你可以依照许可证使用、修改和分发本项目。
