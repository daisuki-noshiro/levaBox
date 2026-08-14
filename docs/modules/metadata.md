# Metadata 模块设计

## 1. 目的

`metadata` 为 levaBox 提供外部游戏资料查询能力。

该模块只负责回答：

> 某个数据源能为这个游戏提供哪些信息？

它不负责决定用户最终采用哪些信息。

## 2. 数据源

当前：

- VNDB

计划：

- Bangumi

Bangumi 尚未实现，因此目前不为了统一两个数据源提前设计通用接口。

## 3. 通用业务流程

一个数据源的典型使用方式：

```text
Search(keyword)
↓
返回多个外部游戏条目
↓
用户确认正确条目
↓
得到 External ID
↓
以后按 External ID 查询各类候选数据
```

第一次通过名称搜索的目的只有一个：定位正确的外部条目。

条目确定后，后续编辑不应再次依赖名称匹配，而应直接使用已经记录的 External ID。

## 4. VNDB

### 4.1 条目搜索

当前已有：

```go
SearchVNDB(keyword)
```

用途：

- 根据游戏名称搜索 VNDB；
- 返回多个候选 VN 条目；
- 让用户确认哪个条目对应本地游戏。

搜索阶段只需要足以辨认游戏的信息：

- VNDB ID；
- 标题；
- alternate title；
- 发售日期；
- 开发商；
- 缩略图。

搜索阶段不获取简介、Tag 和大量图片。

### 4.2 基础信息

条目确认后，根据 VNDB ID 查询：

- Title；
- Company；
- Year；
- Description。

当前已有：

```go
GetVNDBBasicInfo(vndbID)
GetVNDBTitleCandidates(vndbID)
GetVNDBCompanyCandidates(vndbID)
GetVNDBYearCandidate(vndbID)
GetVNDBDescription(vndbID)
```

设计原因：

- 初次导入时需要一次获取全部基础信息；
- 后续只编辑某一栏时，应只查询对应信息；
- 对外按业务用途提供函数，内部可共享 HTTP 请求实现。

### 4.3 Title

候选来源：

- VNDB `title`；
- VNDB `alttitle`；
- 用户手动输入。

第一版不处理完整 `titles[]`。

最终仍保存到：

```text
Game.Title string
```

### 4.4 Company

VNDB 可能返回多个 `developers`。

业务规则：

```text
VNDB developers
↓
开发商候选
↓
用户选一个或手动填写
↓
Game.Company
```

`Game.Company` 保持 `string`，不因 VNDB 支持多个 developer 而改成数组。

### 4.5 Year

VNDB `released` 用于产生年份候选。

例如：

```text
2018-06-29 → 2018
2018-06    → 2018
2018       → 2018
无法确定   → 不提供候选
```

最终保存到：

```text
Game.Year int
```

### 4.6 Description

VNDB description 只作为可编辑候选。

流程：

```text
VNDB description
↓
用户采用
↓
填入编辑框
↓
用户可继续修改
↓
保存最终文本
```

待确认事项：VNDB description 中的 formatting codes 是否需要在 metadata 层清理或转换。

## 5. VNDB Tag

状态：待实现。

已确定方向：

- 来源：`tags.name`、`tags.rating`、`tags.spoiler`；
- Tag 是候选，不自动全部导入；
- `spoiler == 0` 可作为普通候选；
- `spoiler > 0` 第一版忽略；
- `rating` 可用于筛选或排序，但不进入 levaBox Tag 数据库；
- 最终 Tag 数据仍只有 `ID` 和 `Name`；
- 保存时需要与现有同名 Tag 去重并建立 Game ↔ Tag 关系。

具体筛选规则在实现前继续确认。

## 6. VNDB 图片

状态：实验完成，正式实现待进行。

已实测：

### screenshots

大多为游戏内截图，不适合作为 levaBox 大厅资源。

第一版不使用。

### pkgfront

经常带有 Switch 等平台包装元素，不适合作为主要封面来源。

第一版不使用。

### dig

实测图片质量较好，确定作为第一版主要图片来源。

规则：

```text
/release → images.type == dig
↓
根据宽高分类

高 > 宽 → Cover 候选
宽 > 高 → Background 候选
宽 == 高 → 第一版忽略
```

VNDB 只提供图片候选，不直接写入 `Game.CoverPath` 或 `Game.Background`。

用户选择网络图片后，最终应下载到本地，Game 保存本地路径。

本地文件仍然是有效来源：

- Cover：本地图片；
- Background：本地图片或本地视频。

## 7. 文件组织

当前建议：

```text
metadata/
├── vndb.go
├── vndb_basic.go
├── vndb_tag.go       # 待新增
└── vndb_image.go     # 待新增
```

暂不拆出过多 client/request 文件。只有当重复请求逻辑在多个文件中实际出现并影响维护时，再抽公共实现。

## 8. 与 service 的边界

metadata 返回候选数据。

service 负责：

- 导入时调用哪些候选查询；
- 编辑某一栏时调用哪个查询；
- VNDB 与 Bangumi 的候选如何组合；
- 用户最终选什么；
- 如何形成和保存 Game。

因此在 Bangumi 尚未实现前，暂不继续设计 metadata 与 ImportService 的最终统一接口。
