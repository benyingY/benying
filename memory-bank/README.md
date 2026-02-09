# 企业级通用 Agent 平台 - Memory Bank

本目录用于沉淀平台级设计与执行信息，目标是让产品、架构、研发、运维、安全在同一套文档上下文中协作。

## 1. 文档目标

- 统一平台边界、术语和关键决策，降低跨团队沟通成本。
- 保证“从方案到落地”可追踪：为什么做、做什么、做到哪一步。
- 支撑企业级要求：安全合规、稳定性、可观测、可运营。

## 2. 平台定位与范围

一句话定位：构建一个可扩展、可治理、可观测的通用 Agent 平台，为企业内部多业务场景提供智能化执行能力。

平台内能力：

- Main Agent（Python + LangGraph）：负责任务理解、拆解、调度子 Agent、汇总最终答案。
- Sub-Agent Runtime（Python）：支持多个内置子 Agent 并行/串行执行。
- Model Gateway（LiteLLM）：统一模型协议、路由、限流、成本计量与模型级 fallback。
- User-Defined Sub-Agent：支持用户自定义子 Agent 的注册、审批、发布、授权、回滚。
- User-Managed Config：用户可自维护 `skills`、`prompt`、`memory`、`knowledge`。
- Tool Integration：统一接入内部系统与第三方服务。
- Knowledge/RAG：知识检索、召回、权限过滤、引用追踪。
- Skills System：Skill 注册、分发、授权、运行时隔离执行。
- Governance：JWT 鉴权、权限控制、审计、策略护栏。
- Data Transport：Claim Check 大载荷传输与上下文指针化。
- Quality Evals：Judge 评分 + 用户反馈回流的效果级可观测。
- Developer Experience：Dry Run/Playground 联调与仿真执行。
- Observability：链路追踪、成本统计、性能指标、错误分析。
- Delivery：API/SDK/Web Portal 多入口交付（Go + Python 双栈）。

平台外能力（不在本仓库直接实现）：

- 业务系统本体改造（仅定义接入契约）。
- 特定业务流程编排（由业务侧仓库承载）。

## 3. 目录说明

- `/architecture.md`：平台架构与关键组件交互说明。
- `/tech-stack.md`：技术栈锁定决策、版本基线、实施约束。
- `/implementation-plan.md`：分阶段实施计划、里程碑、验收标准。
- `/progress.md`：执行进展、风险、阻塞项、下一步动作。
- `/service-designs/README.md`：服务清单（5 个）与治理规范（非服务）说明。
- `/service-designs/model-gateway.md`：模型网关（LiteLLM）详细设计。
- `/service-designs/sub-agent-management.md`：主/子 Agent 与用户自定义子 Agent 治理规范（非服务）。
- `/service-designs/`：按服务拆分的详细设计（接口、数据模型、扩缩容、容灾）与治理文档。

## 4. 阅读顺序（建议）

1. `architecture.md`：先确认主 Agent 与多子 Agent 的边界和交互。
2. `tech-stack.md`：确认技术路线和治理约束。
3. `service-designs/README.md`：确认服务边界与治理文档定位。
4. `service-designs/sub-agent-management.md`：确认自定义子 Agent 的准入与生命周期。
5. `implementation-plan.md`：对齐阶段目标与资源投入。
6. `progress.md`：跟踪当前状态与风险。

## 5. 统一写作标准

- 先写“决策结论”，再写“原因与权衡”，最后给“风险与回滚方案”。
- 所有关键结论必须可验证，尽量附指标或验收条件。
- 术语统一，避免同一概念多命名。
- 任何方案变更都要更新对应文档与 `progress.md` 记录。
- 设计文档需明确 owner、更新时间、依赖项。

## 6. 企业级验收基线（建议）

- 可用性：核心链路满足约定 SLA/SLO。
- 安全性：支持 JWT Token、RBAC、审计日志、敏感信息治理。
- 稳定性：具备限流、隔离、降级、重试与告警闭环。
- 可观测性：可按请求追踪 Main Agent、Sub-Agent、Skill、工具调用。
- 可治理性：支持用户自定义子 Agent 的审批、灰度、回滚和审计。
- 可配置性：支持用户在租户边界内自维护 `skills/prompt/memory/knowledge` 并可审计回滚。
- 可运维性：支持 GitLab CI + Argo CD、Kustomize 配置管理、灰度发布与回滚。

## 7. 当前状态与下一步

当前主文档（`README.md`、`architecture.md`、`tech-stack.md`、`implementation-plan.md`、`progress.md`）已对齐完成，下一步建议按以下顺序推进：

1. 在 `service-designs/` 细化 5 个服务的接口契约、错误码与容量规划。
2. 落地跨语言契约（`gRPC + Protobuf`）与错误码规范。
3. 建立 GitLab CI、Kustomize overlays、Argo CD 发布清单。
4. 按 `implementation-plan.md` 启动主 Agent + 多子 Agent MVP 开发与验收。

---

当前已锁定关键决策：`Go + Python`、`Main Agent + Multi Sub-Agent`、`User-Defined Sub-Agent`、`User-Managed skills/prompt/memory/knowledge`、`Temporal Macro + LangGraph Micro`、`LiteLLM Model Gateway`、`Claim Check`、`Evals`、`Dry Run`、`JWT Token`、`自研 API Gateway`、`GitLab CI + Argo CD`、`Kubernetes + Kustomize`。

维护约定：本目录是平台“单一事实来源”（Single Source of Truth）。任何口头结论若未落文档，不视为已达成一致。
