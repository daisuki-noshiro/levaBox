# Metadata 模块设计

## 职责

`metadata` 查询外部信息源并规范化结果，不决定最终采用值。当前 v1 信息源为：

- VNDB
- Bangumi

两者由 `service` 的轻量 handler registry 注册，不引入复杂 Provider 接口。

## 公共结构

`metadata.Result` 是单个来源的统一输出，包含：

- `Source`、`ExternalID`
- `CompanyCandidates`
- 可空的 `Year`、`Description`
- `Tags`
- `Covers`、`Backgrounds`

`ImageCandidate` 保存候选图片的来源、URL、缩略图和可选尺寸。

`Result` 不包含 Title。最终 Title 使用用户确认的 SearchKeyword，由 `service` 形成。

## 来源能力

- VNDB：公司、年份、简介、Cover 和 Background 候选；当前不采用 VNDB Tag。
- Bangumi：公司、年份、简介、Tag 和 Cover 候选；v1 不提供 Background。

VNDB `dig` 图片按尺寸分类：竖图为 Cover，横图为 Background，正方形忽略。

## service 合并规则

以下是业务规则，不属于 `metadata` 自身职责：

- VNDB Year 优先，无值时按结果顺序回退；
- Bangumi Description 优先，无值时按结果顺序回退；
- Company 只在提供公司的来源存在共同候选时自动采用；
- Tags 忽略大小写去重并保持首次出现顺序；
- Cover、Background 按来源和来源内部顺序合并并按 URL 去重；
- 单一来源失败记录为 `MetadataSourceIssue`，不阻止其他来源形成 `ImportDraft`。

## 外部条目关联

成功结果中的 `Result.Source + ExternalID` 经 `service.ResolvedSource` 进入 `SaveImport`，最终保存到 `game_metadata_sources`。具体导入事务见 [import.md](./import.md)。

## 边界

`metadata` 不负责用户最终选择、`Game` 构造、数据库写入或来源间优先级。入库后的在线刷新尚未实现。
