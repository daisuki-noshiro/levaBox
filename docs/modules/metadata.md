# Metadata 模块设计

## 职责

`metadata` 负责查询外部信息源并将结果规范化，不决定最终采用的数据。`service` 负责合并候选、形成导入草稿和处理最终选择。

第一阶段信息源：

- VNDB
- Bangumi

当前不引入复杂的 Provider 抽象。

## 公共结构

`metadata/types.go` 定义公共结构：

- `Source`：信息源标识。
- `ImageCandidate`：图片候选，包含 `Source`、`URL`、`Thumbnail`、`Width`、`Height`。
- `Result`：单个信息源能够提供的原材料。

`Result` 包含：

- `Source`
- `ExternalID`
- `CompanyCandidates`
- `Year`
- `Description`
- `Tags`
- `Covers`
- `Backgrounds`

`Result` 不包含 Title。用户最终确认的搜索词默认作为 `Game.Title`，metadata 不提供或决定最终 Title。

Company 以候选数组返回，最终值由 service 和用户确认。Year、Description 在来源无数据时允许为 `nil`。

## Tag

第一阶段只使用：

- Bangumi Tag
- 用户手动输入或本地已有 Tag

不采用 VNDB Tag。Title、Company、Year、Tag 入库后直接在本地编辑。

## 图片候选

候选界面后续展示信息源和原始分辨率，`ImageCandidate` 使用 `Source`、`URL`、`Thumbnail`、`Width`、`Height` 表达所需信息。

VNDB 图片规则：

- 从 `/release` 查询 `images.type == dig` 的图片。
- 竖图作为 Cover 候选。
- 横图作为 Background 候选。
- 正方形图片忽略。

Bangumi 第一阶段只提供 Cover，不承担 Background。

## 导入流程

第一版各信息源按相关度搜索，并默认采用排序第一的条目作为匹配结果；用户可通过修改搜索词重新搜索进行纠正。

```text
选择 EXE
→ 推测 SearchKeyword
→ 用户可修改搜索词
→ 查询启用的信息源
→ metadata 分别得到 Result
→ service 整理为 ImportDraft
→ 用户确认
→ 保存
```

入库后的 metadata 在线重新获取只计划用于：

- Description
- Cover
- Background

Title、Company、Year、Tag 直接在本地编辑。
