# levaBox 架构说明

## 1. 目标

本文档描述 levaBox 当前后端的分层方式和模块边界。重点是说明代码应该放在哪里、各层可以依赖什么，以及哪些职责不能混在一起。

当前架构仍处于早期阶段。未验证的抽象不提前引入。

## 2. 分层

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

这不是严格的 Clean Architecture 实现，而是当前项目规模下用于控制职责边界的简单分层。

## 3. model

`model` 定义 levaBox 自己的核心数据。

当前包括：

- `Game`
- `Tag`
- `Background`
- `LaunchConfig`
- `GameQueue`

### 约束

`model` 不应该出现第三方数据源专有概念，例如：

- `VNDBResult`
- `BangumiSubject`
- VNDB API 的 `developers`、`screenshots` 等结构

第三方数据进入 `model` 前，必须先经过业务层选择和转换。

## 4. database

`database` 负责 SQLite 持久化。

职责包括：

- schema；
- Game CRUD；
- Tag CRUD；
- Game 与 Tag 关系；
- 游戏队列持久化。

### 约束

`database` 不决定：

- 用户选择哪个元数据；
- VNDB 和 Bangumi 谁优先；
- 某个候选字段是否应该采用。

这些属于 `service`。

## 5. metadata

`metadata` 负责访问外部元数据源，并把第三方 API 数据整理成 levaBox 能消费的候选数据。

当前数据源：

- VNDB

计划数据源：

- Bangumi

### metadata 应负责

- 根据关键词搜索外部游戏条目；
- 根据外部条目 ID 查询基础信息；
- 查询 Tag 候选；
- 查询图片候选；
- 处理第三方 API 的请求、响应和字段转换。

### metadata 不负责

- 决定用户最终采用哪个字段；
- 把候选直接写入 `Game`；
- 决定导入流程；
- 修改数据库；
- 决定多个数据源之间的优先级。

## 6. service

`service` 负责业务流程和业务规则。

例如：

- 游戏导入；
- 游戏资料编辑；
- 游戏启动；
- 多数据源候选的组合与选择；
- 把用户最终确认的数据转换为 `Game`；
- 调用 database 完成持久化。

`service` 可以调用 `metadata` 和 `database`，但不应该理解第三方 API 的原始 JSON 结构。

## 7. app.go

`app.go` 是 Wails 暴露层。

职责是把需要提供给前端的后端能力暴露出去，并把调用转给对应 service。

原则上不在这里编写复杂业务逻辑。

## 8. frontend

前端负责：

- 展示；
- 用户输入；
- 候选选择；
- 页面状态；
- 动效和交互。

核心业务规则仍以 Go 后端为准。

## 9. 外部元数据的基本流程

```text
搜索关键词
↓
metadata 搜索外部条目
↓
用户确认正确条目
↓
记录外部条目 ID
↓
metadata 根据 ID 查询字段候选
↓
service 组织候选和业务流程
↓
用户选择 / 修改
↓
service 形成最终 Game 数据
↓
database 持久化
```

其中“外部条目 ID 如何持久化”尚未最终设计。

## 10. 当前不做的抽象

目前不提前实现：

- 通用 `MetadataProvider` interface；
- VNDB/Bangumi 的统一超大数据结构；
- 为未来未知数据源准备的复杂插件体系。

原因是当前只有 VNDB 已实际接入，Bangumi 尚未实现。等两个真实数据源都完成后，再根据实际共性抽象。
