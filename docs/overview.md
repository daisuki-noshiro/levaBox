# levaBox 项目总览

## 项目定位

levaBox 是面向 Windows 掌机和桌面环境的 Galgame 管理器。项目目标是形成可实际使用的软件，同时以真实需求推进后端数据建模、模块划分、持久化、第三方 API、文件与进程管理。

## 技术栈

- 后端：Go
- 桌面框架：Wails
- 前端：Vue + TypeScript
- 数据库：SQLite
- 元数据源：VNDB v1、Bangumi v1

## 主要模块

| 模块 | 主要职责 |
| --- | --- |
| `model` | `Game`、`Tag`、`GameMetadataSource` 等核心数据 |
| `database` | SQLite schema、基础持久化与事务 helper |
| `metadata` | 查询 VNDB、Bangumi 并输出统一候选数据 |
| `service` | 导入等业务流程、字段合并和事务编排 |
| `app.go` | Wails 暴露层，连接前端与 service |
| `frontend` | 页面、交互、展示和用户确认 |

模块边界见 [architecture.md](./architecture.md)，导入设计见 [modules/import.md](./modules/import.md)。

## 当前状态

后端“游戏导入”主链 v1 已基本完成，包括：

- `Game`、`Tag`、队列和元数据来源模型；
- Game、Tag、`game_tags`、`game_queue`、`game_metadata_sources` 的 SQLite 持久化；
- `StartImport` 的 EXE 查重、路径处理、工作目录和默认搜索词；
- VNDB、Bangumi v1 查询及统一 `metadata.Result`；
- 多来源查询、失败来源 `Issues` 隔离和 `ImportDraft` 合并；
- `SaveImport` 的最终校验、再次查重和 `model.Game` 构造；
- Cover、Background 下载到本地媒体目录；
- Tag 复用/新增、`game_tags` 和 MetadataSource 保存；
- SaveImport 数据库事务、失败回滚和本次媒体清理；
- 不依赖公网的 database、metadata、service 单元测试。

Wails bridge 已完成。前端 `ImportView` 已接入第一段真实流程：

```text
SelectExecutable → StartImport → 展示基础 ImportDraft → 编辑 SearchKeyword
```

## 下一阶段

将基础 `ImportDraft` 传入 `PrepareImportMetadata`，并展示 VNDB / Bangumi 元数据候选。

## 暂未完成

- 前端元数据确认与最终保存流程；
- 编辑已有游戏的完整业务；
- 游戏启动 service；
- 删除游戏时的本地媒体清理；
- metadata 在线刷新；
- Settings、BGM 等后续能力。

## 当前约束

- `metadata` 只查询和规范化外部数据，不决定最终采用值；
- `service` 负责多来源业务规则、用户确认结果和事务编排；
- `model.Game` 不依赖第三方 API 原始结构；
- 最终采用的远程图片必须下载为本地文件；
- 前端负责交互和展示，不承担核心业务规则。

## 文档约定

- `overview.md`：项目状态和总体范围；
- `architecture.md`：系统分层和模块边界；
- `modules/*.md`：具体模块的设计；
- `decisions/*.md`：重要设计决定及原因。
