# levaBox 项目总览

## 1. 项目定位

levaBox 是一个面向 Windows 掌机和桌面环境的 Galgame 管理器，技术栈为 Go + Wails + Vue。

项目有两个目标：

- 做成实际可用的软件，而不是只完成演示功能；
- 作为后端工程学习项目，逐步练习需求分析、数据建模、模块划分、持久化、第三方 API、文件与进程管理等能力。

当前开发顺序：

> 需求 → 数据 → 业务规则 → 模块职责 → 接口 → 实现

后端设计不以“某个 API 能返回什么”或“前端已经有什么页面”为起点。

## 2. 技术栈

- 后端：Go
- 桌面框架：Wails
- 前端：Vue + TypeScript
- 数据库：SQLite
- 外部元数据：VNDB；Bangumi 计划接入，尚未实现

## 3. 主要模块

| 模块 | 主要职责 |
| --- | --- |
| `model` | levaBox 核心领域数据，例如 `Game`、`Tag`、`Background` |
| `database` | SQLite 表结构及持久化操作 |
| `metadata` | 查询 VNDB、Bangumi 等外部数据源并整理候选数据 |
| `service` | 导入、编辑、启动等业务流程与业务规则 |
| `app.go` | Wails 暴露层，连接前端与后端服务 |
| `frontend` | 页面、交互、展示和用户选择 |

模块边界见 [architecture.md](./architecture.md)。

## 4. 当前阶段

### 已完成或已具备

- `Game`、`Tag`、`Background`、`LaunchConfig` 等核心模型已建立；
- SQLite 的游戏、标签和队列表结构已建立；
- `StartImport` 已具备 EXE 查重、路径规范化、工作目录和默认搜索词生成能力；
- VNDB 按游戏名搜索已可用；
- VNDB 基础信息查询已可用，可按 VNDB ID 获取标题、开发商、年份和简介；
- VNDB 基础信息支持按字段用途单独查询；
- 已使用 Summer Pockets（`v20424`）进行 VNDB 联网测试。

### 正在进行

当前重点是完成 `metadata` 层的 VNDB 能力，并稳定数据源层的职责边界。

计划顺序：

1. VNDB 基础信息；
2. VNDB Tag；
3. VNDB 图片；
4. Bangumi 元数据；
5. 两个真实数据源都明确后，再设计 `service` 如何组合、选择和编辑候选数据。

### 暂未完成

- Bangumi 接入；
- VNDB Tag 正式实现；
- VNDB `dig` 图片正式实现；
- 外部数据源条目 ID 的持久化方案；
- 导入流程与 metadata 的完整业务编排；
- 用户确认后的最终导入、图片下载和数据库写入；
- 游戏启动、游玩时间记录等完整业务流程。

## 5. 当前约束

- `metadata` 只负责查询和整理外部数据，不决定用户最终采用什么；
- `service` 负责导入、编辑、字段选择和多个数据源之间的业务组合；
- `model.Game` 不直接依赖 VNDB、Bangumi 的 API 数据结构；
- 外部图片 URL 只作为候选数据，最终使用的图片应下载到本地并保存本地路径；
- 前端负责交互和展示，不承担核心业务规则。

## 6. 文档约定

- `overview.md`：项目状态和总体范围；
- `architecture.md`：系统分层和模块边界；
- `modules/*.md`：具体模块的业务设计与接口意图；
- `decisions/*.md`：重要设计决定及原因。

文档记录设计意图和约束，不重复解释代码中已经清楚表达的实现细节。
