# 企业级通用 Agent 平台进展跟踪

- 文档版本：v1.3
- 更新时间：2026-02-09
- 当前阶段：Phase 1（MVP）
- 总体状态：In Progress

## 1. 已锁定决策（基线）

- 后端：Go + Python 双栈
- Agent 架构：Main Agent + Multi Sub-Agent
- 自定义能力：支持用户自定义 Sub-Agent
- 用户配置：支持用户自维护 `skills/prompt/memory/knowledge`
- 编排边界：Temporal（Macro）+ LangGraph（Micro）
- 大载荷处理：Claim Check Pattern
- 上下文管理：Context Optimizer + Token Budget
- 质量可观测：Evals（Judge + 用户反馈）
- 调试能力：Dry Run / Playground
- 路由策略：Hierarchical Sub-Agent Routing
- 模型网关：LiteLLM Model Gateway
- 兜底策略：模型级 + 规则级 + 人工级 Fallback
- 认证：JWT Token（RS256）+ Refresh Token
- 网关：自研 API Gateway（Go）
- CI/CD：GitLab CI + Argo CD
- 部署：Kubernetes + Kustomize

## 2. 已完成项

- 完成 `README.md`，补齐主/子 Agent 与自定义子 Agent 平台定位。
- 完成 `architecture.md`，加入 8 项增强设计（Claim Check、沙箱、边界、Context、Evals、Dry Run、分层路由、Fallback）。
- 完成 `tech-stack.md`，锁定对应技术选型与运行策略。
- 完成 `implementation-plan.md`，将 8 项增强映射到分阶段交付与验收指标。
- 完成文档口径统一：`skills/prompt/memory/knowledge` 由用户在租户边界内自维护。
- 完成 5 个核心服务设计文档：`control-plane.md`、`agent-runtime.md`、`model-gateway.md`、`tool-gateway.md`、`knowledge-service.md`。
- 完成 `service-designs/sub-agent-management.md` 作为子 Agent 治理规范入口（非服务）。
- 完成 `service-designs/README.md`，明确“服务文档”与“治理规范文档”边界。

## 3. 进行中

- 细化 5 个服务设计文档中的接口细节、错误码与容量规划。
- 细化用户自维护 PMK（prompt/memory/knowledge）配置 API 与审计模型。
- 拆解 MVP 开发任务并分配 owner。

## 4. 下一步（按优先级）

1. 完善现有服务文档的接口契约与示例请求：
- `service-designs/control-plane.md`
- `service-designs/agent-runtime.md`
- `service-designs/model-gateway.md`
- `service-designs/tool-gateway.md`
- `service-designs/knowledge-service.md`

2. 建立跨语言接口契约：
- 定义 `gRPC + Protobuf` 接口规范。
- 固化错误码、追踪字段、主子 Agent 调用语义。

3. 建立 DevOps 骨架：
- GitLab CI 模板（Go、Python、Sub-Agent 包）。
- Kustomize overlays（dev/staging/prod）。
- Argo CD 应用清单与发布策略。
- gVisor RuntimeClass 与 NetworkPolicy 模板。

## 5. 风险与阻塞

- 风险 R1：Temporal 与 LangGraph 边界不清导致重复状态管理。
- 缓解：以 Activity 颗粒度固定边界，禁止平行持久化。

- 风险 R2：大 Context 导致吞吐与存储膨胀。
- 缓解：Claim Check + payload 阈值 + 生命周期清理策略。

- 风险 R3：用户自定义 Sub-Agent 增加运行时安全风险。
- 缓解：gVisor 沙箱、NetworkPolicy 默认拒绝、租户白名单与审计全量留痕。

- 风险 R4：用户自维护 PMK 配置变更过快引发线上波动。
- 缓解：配置版本化、灰度发布、快速回滚与变更告警。

- 风险 R5：缺乏质量评估导致业务价值不可证明。
- 缓解：Evals 评分 + 用户反馈回流 + Ground Truth 对比。

## 6. 变更记录

- 2026-02-09：新增 `model-gateway`（LiteLLM）服务设计，并同步到架构/技术栈/计划/总览文档。
- 2026-02-09：明确 `sub-agent-management.md` 为治理规范文档（非独立服务）。
- 2026-02-09：同步主 Agent + 多子 Agent + 用户自定义子 Agent 决策至所有主文档。
- 2026-02-09：补充用户可自维护 `skills/prompt/memory/knowledge` 的治理边界与实施计划。
- 2026-02-09：采纳并落地 8 项架构增强建议（Claim Check、沙箱、边界、Context、Evals、Dry Run、分层路由、Fallback）。

---

更新规则：每周至少更新一次；若发生架构或技术栈决策变更，需在 24 小时内更新本文件。
