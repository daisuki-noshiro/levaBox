# levaBox 架构说明

## 分层

```text
frontend
   ↓
app.go
   ↓
service
   ↓
├── metadata
├── database
└── model
```

这是适配当前项目规模的简单分层，不追求完整的 Clean Architecture。未被真实需求验证的抽象不提前引入。

## model

`model` 定义 levaBox 自己的核心数据，包括 `Game`、`Tag`、`Background`、`LaunchConfig`、`GameQueue` 和 `GameMetadataSource`。

第三方 API 的原始结构不能进入 `model`。外部数据必须先由 `metadata` 规范化、由 `service` 选择和转换。

## database

`database` 负责 SQLite：

- Game；
- Tag；
- `game_tags`；
- `game_queue`；
- `game_metadata_sources`；
- 普通 CRUD 和小型 transaction helper。

`DBTX` 只包含 `*sql.DB` 与 `*sql.Tx` 共同需要的最小查询/写入能力，用于 SaveImport 事务内复用必要 helper。普通读取函数不因此全部抽象化。

`database` 不决定来源优先级或最终字段值，这些属于 `service`。

## metadata

`metadata` 当前接入 VNDB v1 和 Bangumi v1，负责：

- 按关键词匹配外部条目；
- 按 ExternalID 查询基础信息、Tag 和图片候选；
- 处理第三方请求、响应和字段转换；
- 输出统一的 `metadata.Result`。

`metadata` 不构造最终 `Game`，不写数据库，也不决定多来源业务优先级。

## service

导入相关的 `service` 当前负责：

- `StartImport` 本地导入准备；
- 多来源查询编排和单来源失败隔离；
- `metadata.Result` → `ImportDraft` 合并；
- `SaveImportRequest` 最终校验和再次查重；
- 用户确认数据 → `model.Game`；
- Cover、Background 媒体文件准备；
- Game、Tags、`game_tags`、`game_metadata_sources` 的数据库事务编排；
- 数据库失败后的本次媒体清理。

`service` 可以调用 `metadata` 和 `database`，但不理解第三方 API 原始 JSON。

## app.go 与 frontend

`app.go` 是 Wails 暴露层，只转发前端需要的 service 能力，不承载复杂业务。前端负责展示、输入、候选选择和页面状态，核心规则仍在 Go 后端。

当前完整导入能力尚未全部通过 `app.go` 暴露给前端。

## 元数据来源持久化

外部条目关联使用 `game_metadata_sources` 保存：

- `GameID`
- `Source`
- `ExternalID`

数据路径为：

```text
service.ResolvedSource
↓
SaveImport
↓
model.GameMetadataSource
↓
database
```

同一 Game 的同一 Source 只保留一个 ExternalID；不同本地 Game 可以引用同一外部条目。

## 当前不做的抽象

目前不引入通用 `MetadataProvider`、Repository 大接口或未知数据源插件体系。VNDB 与 Bangumi 通过轻量 handler registry 接入，多来源业务规则保留在 `service`。
