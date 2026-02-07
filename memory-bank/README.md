# Memory Bank

`memory-bank` 用于沉淀企业级 Agent 平台的长期知识，作为架构设计、实施推进和进度管理的统一文档入口。

## 目录说明

- `architecture.md`：平台总体架构蓝图与能力分层（含 Temporal/LangGraph 边界、RAG ACL 约束）
- `tech-stack.md`：技术栈定版（v1.1）
- `implementation-plan.md`：细颗粒度实施计划（每步可独立测试）
- `service-designs/`：V1 各服务详细设计文档（职责、接口、数据模型、可控做法、开发步骤）
- `progress.md`：项目执行进展、风险、阻塞项、下一步动作

## 当前状态

- 已完成：`architecture.md`、`tech-stack.md`、`implementation-plan.md`、`service-designs/*`、`progress.md`（模板已初始化）
- 持续维护：`progress.md`（建议按周更新执行状态、风险和阻塞）

## 推荐阅读顺序

1. `architecture.md`：先统一架构边界和核心能力定义
2. `tech-stack.md`：确认技术定版与关键架构约束
3. `service-designs/README.md`：查看服务拆分与依赖关系
4. `service-designs/*.md`：阅读单服务详细设计和可单测开发步骤
5. `implementation-plan.md`：按步骤推进和验收
6. `progress.md`：持续记录执行状态并驱动复盘闭环

## 文档维护约定

- 每次重大调整（架构、技术路线、里程碑）后，同步更新对应文档
- 在 `progress.md` 中记录更新时间、变更摘要、负责人
- 新增文档时，请在本 `README.md` 更新目录说明与用途

## 下一步建议

1. 以 `implementation-plan.md` 的 Step 1 开始执行，逐步过 Gate A-E
2. 为每个服务建立对应代码目录与 CI 任务，映射文档中的测试命令
3. 持续维护 `progress.md`，按周追踪步骤完成率、风险和阻塞
