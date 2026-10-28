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

FortuneFinger 是一套可自行部署的企业活动签到与抽奖系统，适用于企业年会、庆典和其他现场活动。

> [!IMPORTANT]
> FortuneFinger 仍在持续开发，当前仓库已经包含管理端、系统配置、奖品管理和签到大屏。移动端签到 H5、完整签到记录及正式抽奖流程尚未形成稳定发行版，正式活动使用前请完成充分测试。

## 功能特性

- **企业与权限管理**：维护集团、企业和部门结构，通过角色控制菜单、页面和操作权限。
- **字典与系统设置**：统一维护奖项等级、奖品类型和签到模式，并配置签到时段、大屏模板及邮件服务。
- **奖品管理**：按所属企业维护奖品、奖项等级、数量、状态和图片。
- **签到大屏**：提供七套差异化现场模板，根据后台配置加载，实时同步签到状态和人数。
- **动态签到码**：后端签发30秒短期凭证并生成带品牌标识的二维码，未开始或已截止时停止展示和刷新。
- **对象存储**：支持 MinIO 和阿里云 OSS，用于统一管理业务图片。
- **邮件服务**：支持标准 SMTP、STARTTLS 和 SSL/TLS，可在系统设置中保存配置并发送测试邮件。
- **访问控制**：敏感系统配置按权限开放，适合为公网演示账号提供受限功能。

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

### 启动签到大屏

再开启一个终端：

```bash
cd frontend/web
npm ci
npm run dev
```

根据终端显示的地址打开签到大屏。大屏会读取后台保存的签到模板、签到时段和实时状态，并在允许签到的时间内获取动态二维码。

需要使用手机在局域网中预览时，让电脑和手机连接同一网络，并将开发服务监听到局域网：

```bash
npm run dev -- --host 0.0.0.0
```

随后使用终端显示的 Network 地址访问。正式公网部署建议为管理端、签到大屏和后续移动签到页配置 HTTPS。

> [!NOTE]
> 当前版本尚未提供首次启动向导。数据库迁移只创建表结构，初始超级管理员、角色、菜单和权限数据仍需自行准备。该流程会在正式发行版前完善。

## 部署提示

正式部署时，请更换配置文件中的示例凭据，并为数据库、Redis、对象存储和 SMTP 使用独立账号。公网演示环境不应向演示账号开放菜单管理、系统设置和其他敏感权限。

## 获取帮助

请通过 [GitHub Issues](https://github.com/hourinX/FortuneFinger/issues) 报告问题或提出建议。提交问题时请附上运行环境、复现步骤和必要的错误信息，并删除其中的密码、令牌及个人信息。

## 开源许可

FortuneFinger 使用 [Apache License 2.0](LICENSE) 发布，署名信息见 [NOTICE](NOTICE)。你可以依照许可证使用、修改和分发本项目。
