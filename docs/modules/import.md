# 游戏导入模块

## 目标

将用户选择的 EXE 与外部元数据整理为可确认草稿，并在确认后持久化为本地 Game。

## 流程

```text
选择 EXE
↓
StartImport
↓
基础 ImportDraft
↓
用户确认 SearchKeyword
↓
PrepareImportMetadata
↓
VNDB / Bangumi
↓
BuildImportDraft
↓
用户编辑最终字段
↓
SaveImportRequest
↓
SaveImport
↓
Game + Tags + MetadataSources + 本地媒体
```

## StartImport

负责 EXE 查重、绝对路径与清理、WorkingDirectory 和默认 SearchKeyword，不联网。

## PrepareImportMetadata

按配置顺序查询多个来源，隔离单来源错误，并通过 `BuildImportDraft` 形成可编辑草稿。

## SaveImport

负责：

- 最终输入校验和再次 EXE 查重；
- 下载用户确认的 Cover、Background；
- 构造 `model.Game`；
- 在同一 transaction 中保存 Game、Tags、`game_tags` 和 `game_metadata_sources`；
- 数据库失败时 Rollback 并删除本次媒体目录。

媒体目录为：

```text
<UserConfigDir>/levaBox/media/<gameID>/
```

## Transaction 边界

HTTP 下载和文件准备在 transaction 前完成。Transaction 内只执行数据库操作，以避免网络等待期间长期持有事务。

## Wails 暴露层

`app.go` 暴露 `SelectExecutable`、`StartImport`、`PrepareImportMetadata` 和 `SaveImport`。该层只负责文件选择和参数转发，导入业务规则仍由 `ImportService` 负责。

## 前端接入进度

已接入：

```text
选择 EXE
↓
StartImport
↓
展示基础 ImportDraft
↓
编辑 SearchKeyword
↓
PrepareImportMetadata
↓
元数据编辑
↓
候选图片选择
```

尚未接入：

```text
最终 Tag 选择/编辑
↓
SaveImportRequest
↓
SaveImport
↓
导入完成
```

`ImportMetadataResult.Issues` 只表示单个来源查询失败。页面会保留并展示其他来源成功形成的草稿和候选结果，同时以轻量警告列出失败来源；只有 Wails Promise reject 才视为整次查询错误。

## 当前不负责

- 视频背景和 BGM；
- 编辑已有游戏；
- 删除游戏媒体；
- metadata 在线刷新。
