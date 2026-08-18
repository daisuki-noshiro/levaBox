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

## 当前不负责

- 前端交互与 Wails 完整接入；
- 视频背景和 BGM；
- 编辑已有游戏；
- 删除游戏媒体；
- metadata 在线刷新。
