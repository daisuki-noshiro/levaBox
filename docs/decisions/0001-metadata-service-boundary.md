# ADR-0001：Metadata 与 Service 的职责边界

- 状态：Accepted
- 日期：2026-08-13

## 背景

levaBox 计划从 VNDB、Bangumi 等外部数据源获取游戏资料。

早期设计中曾考虑在 `ImportService` 中直接接入 VNDB 搜索和详情查询。但当前只有 VNDB 已实际接入，Bangumi 尚未实现，如果过早围绕 VNDB 设计完整导入流程，后续可能使 service 与某个具体数据源耦合。

同时，导入和后续资料编辑都会使用元数据查询能力，因此元数据查询本身不属于导入流程专用能力。

## 决定

将外部元数据查询能力放在 `metadata` 包中，业务流程和字段选择放在 `service` 包中。

### metadata

负责：

- 调用 VNDB、Bangumi 等第三方 API；
- 搜索外部游戏条目；
- 根据 External ID 查询基础信息、Tag、图片等候选数据；
- 把第三方 API 数据转换为 levaBox 可消费的候选结构。

不负责：

- 决定采用哪个数据源；
- 决定采用哪个字段候选；
- 构造最终 `Game`；
- 写数据库。

### service

负责：

- 游戏导入和资料编辑流程；
- 组合多个数据源提供的候选；
- 处理用户最终选择；
- 将确认后的数据转换为 `Game`；
- 调用 database 完成持久化。

## 原因

这样可以保证：

1. 导入和编辑能够复用相同的元数据查询能力；
2. service 不需要理解 VNDB/Bangumi 的原始 API 结构；
3. 新增数据源时不需要把数据源专用逻辑塞进 ImportService；
4. 在两个真实数据源都实现前，不必猜测它们应该如何统一。

## 影响

当前阶段先完成 `metadata` 中的 VNDB 和 Bangumi 能力，再回到 `service` 设计多数据源的组合、选择和编辑流程。

暂不引入通用 `MetadataProvider` interface。是否需要统一接口，待 Bangumi 实现后根据实际重复和差异决定。
