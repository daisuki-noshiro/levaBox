# ADR-0002：导入事务与媒体准备边界

- 状态：Accepted
- 日期：2026-08-18

## 背景

SaveImport 同时涉及 HTTP 下载、文件写入和 SQLite 持久化。网络请求不应长期占用数据库 transaction，但数据库失败后也不能留下本次导入的孤立媒体。

## 决定

1. 先将用户确认的 Cover、Background 下载到 `media/<gameID>`；
2. 文件准备完成后再开启数据库 transaction；
3. Transaction 内保存 Game、Tags、`game_tags` 和 `game_metadata_sources`；
4. 任一数据库步骤或 Commit 失败时 Rollback，并删除本次 `media/<gameID>`；
5. 不实现文件系统事务或两阶段提交。

## 原因

该边界缩短了 transaction 持有时间，并以较小复杂度为 v1 提供足够可靠的数据库与媒体一致性。

## 影响

数据库提交成功后媒体才视为导入结果的一部分。删除已有游戏时的媒体清理由未来的 `GameService.DeleteGame` 单独定义，database 层只处理数据库记录。
