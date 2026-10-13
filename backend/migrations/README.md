# 数据库迁移（migrations）

用 [golang-migrate](https://github.com/golang-migrate/migrate) 管理表结构演进。

## 常用命令

```bash
# 应用全部未执行的迁移
migrate -path migrations -database "mysql://root:密码@tcp(localhost:3306)/fortunefinger" up

# 回滚最近一次迁移
migrate -path migrations -database "..." down 1

# 查看当前版本与脏状态
migrate -path migrations -database "..." version

# 新建一对迁移文件（-seq 使用 000002 递增序号）
migrate create -ext sql -dir migrations -seq add_lottery_round
```

> 实际执行由 `setup` 初始化/升级流程调用（见 [docs/FIRST_START.md](../../docs/FIRST_START.md)），
> 正常启动不自动改表结构。

## 文件命名

```
NNNNNN_动词_对象.up.sql      向前：应用变更
NNNNNN_动词_对象.down.sql    向后：精确撤销
```

- `NNNNNN`：6 位递增序号。**新迁移永远取下一个号，绝不修改已发布的号。**
- 动词：`init` / `create` / `add` / `alter` / `drop` / `seed`。

## 拆分原则：按「变更意图」，不按「表数量」

回滚的粒度等于迁移文件的粒度，所以拆分的标准是**变更意图**是否相同：

| 情况 | 做法 |
|---|---|
| 多张表同时诞生（项目初始建库） | **一个文件** |
| 多张表为**同一个目的**修改（如统一加审计字段） | **一个文件** |
| 多张表**互不相关**的改动 | **分开多个文件** |

反面教材：把「抽奖加轮次」「奖品加库存」「菜单加图标」三个无关改动塞进一个文件——
一旦其中一个出问题，`down` 会把三个一起退掉，中途失败还会让库卡在脏状态。

## 铁律

1. **先写迁移，再改库**。永远不手敲 DDL 改生产库。
2. **已发布的迁移不可修改**，只能追加新迁移，否则已部署环境会漂移。
3. **up 必须能重复执行**（`CREATE TABLE IF NOT EXISTS`、加字段前先判断）。
4. **down 精确撤销 up**（up 建表，down 就 `DROP TABLE`）。
5. 表名/字段名/类型与 GORM model 严格对齐。
6. **种子数据**：系统启动必需的数据（如内置角色）走迁移；测试用数据不走。

## 关于 dirty 状态

MySQL 的 DDL 是隐式提交的，多条 DDL **不构成一个事务**。如果一个迁移文件里有 3 条
`ALTER`，第 2 条失败时数据库会停在「做了一半」的状态，migrate 会把 `dirty` 置 1 并
**拒绝继续执行**，直到人工核实并修复。

```bash
migrate ... version          # 查看是否 dirty
# 手工把数据库改到与某个版本一致后：
migrate ... force <版本号>    # 强制标记版本，清除 dirty
```

这也是「一个文件只做一件事」的另一个理由——**爆炸半径小**。

## 生产注意事项

- 大表 `MODIFY COLUMN` 可能重建整张表并锁表数分钟。有量之后改用
  `gh-ost` 或 `pt-online-schema-change` 做在线变更。
- `down` 主要用于开发/测试环境。生产环境出问题通常用**新迁移向前修复**，
  而不是回滚——应用代码已经依赖新结构，回滚 schema 会打断线上服务。
