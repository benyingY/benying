# 企业级通用 Agent 平台技术栈（最终锁定版）

- 文档版本：v1.5
- 状态：Locked（已确认）
- 更新时间：2026-02-09
- 关联文档：`README.md`、`architecture.md`

## 1. 选型目标与约束

### 1.1 目标

- 以企业级稳定性、安全性、可观测性为第一优先级。
- 在 MVP 阶段采用可长期演进的生产级技术栈，避免二次重构。
- 支持 Main Agent 编排多个 Sub-Agent，并支持用户自定义 Sub-Agent。
- 支持租户用户自维护 `skills/prompt/memory/knowledge`。
- 支持大载荷、质量评估、调试仿真、降级兜底等 AI 原生能力。

### 1.2 约束

- 默认部署在 Kubernetes。
- 技术栈必须支持私有化部署。
- 所有核心组件必须可横向扩展并支持高可用。
- 平台只保留一套技术路线，不引入平行实现。

## 2. 最终技术决策（一览）

- 后端语言：`Go 1.24+` + `Python 3.12+`
- Go 服务框架：`Gin`
- Python 服务框架：`FastAPI`
- Agent 拓扑：`Main Agent + Multi Sub-Agent`
- Agent 图编排：`LangGraph`
- 工作流编排：`Temporal`
- 模型网关：`LiteLLM (Model Gateway)`
- 编排边界：`Temporal Macro + LangGraph Micro`
- 大载荷传输：`Claim Check Pattern (context_ref_id)`
- 上下文存储：`Redis + MinIO`
- 子 Agent 注册与版本：`Sub-Agent Registry Service + PostgreSQL`
- 子 Agent 包分发：`GitLab Container Registry (OCI Artifact)`
- 子 Agent 规范：`sub-agent.yaml + JSON Schema`
- Prompt 配置：`Prompt Registry Service + Versioning`
- Memory 配置：`Memory Profile Service`
- Knowledge 配置：`Knowledge Space Binding Service`
- 运行时沙箱：`gVisor RuntimeClass`（默认）
- 网络隔离：`K8s NetworkPolicy + Egress Gateway`
- 消息系统：`Kafka`
- 主数据库：`PostgreSQL 16+`
- 缓存与会话：`Redis 7.2+`
- 全文检索：`OpenSearch 2.x`
- 向量检索：`Milvus 2.4+`
- 对象存储：`MinIO`（S3 兼容）
- 身份认证：`JWT Token (RS256)` + `Refresh Token`
- 网关与流量治理：`自研 API Gateway (Go)`
- 质量评估：`Evaluation Service + Judge Agent`
- 可观测：`OpenTelemetry + Prometheus + Grafana + Loki + Tempo`
- 密钥与配置：`External Secrets Operator + HashiCorp Vault`
- 前端门户：`React + Next.js 15+`
- CI/CD：`GitLab CI + Argo CD`
- 部署与发布：`Kubernetes + Kustomize`

## 3. Go + Python 职责边界（锁定）

- Go（平台控制面）
- `Agent API`：鉴权、配额、幂等、计费、请求编排入口。
- `Workflow Orchestrator`：Temporal 工作流生命周期控制。
- `Sub-Agent Registry Service`：子 Agent 元数据、审批、授权、版本治理。
- `PMK Config Service`：prompt/memory/knowledge 配置与版本审计。
- `Context Claim Service`：大载荷上下文存取与引用 ID 管理。
- `Dry Run API`：仿真调试入口。

- Python（智能执行面）
- `Main Agent Coordinator`：任务分解、子 Agent 选择、结果聚合。
- `Context Optimizer`：上下文预算与裁剪。
- `Sub-Agent Router`：分层路由与并发调度。
- `Sub-Agent Workers`：内置与用户自定义子 Agent 执行。
- `Skill Runtime`：Skill 加载、隔离执行、结果标准化回传。
- `Model Client`：统一调用 `model-gateway`，不直连模型厂商。
- `Evaluation Service`：在线质量评估与反馈回流。

- 跨语言通信（唯一标准）
- 同步调用：`gRPC + Protobuf`
- 异步事件：`Kafka + CloudEvents`
- 统一追踪：`trace_id` + `main_agent_id` + `sub_agent_id`

## 4. 分层技术栈矩阵（单一方案）

| 层级 | 最优选型 | 版本基线 | 锁定原因 |
|---|---|---|---|
| 控制面服务 | Go + Gin | Go 1.24+ | 高并发、低资源占用、稳定性高 |
| 智能执行面 | Python + FastAPI + LangGraph | Python 3.12+ | AI 生态成熟，主/子 Agent 图编排效率高 |
| 宏观编排 | Temporal | 1.25+ | 生命周期、补偿、挂起/唤醒、重试稳定 |
| 微观推理 | LangGraph | 稳定版 | Main Agent 推理回路与子 Agent 调度清晰 |
| 模型网关 | LiteLLM | 稳定版 | 多模型统一协议、路由、限流与 fallback 能力成熟 |
| 大载荷处理 | Claim Check + Redis/MinIO | Redis 7.2+ / MinIO 稳定版 | 避免 gRPC 与 Temporal 历史膨胀 |
| 上下文优化 | Context Optimizer | Python 3.12+ | 控制 token 成本并提升回答质量 |
| 分层路由 | Hierarchical Router | Python 3.12+ | 子 Agent 数量增长时保持路由稳定 |
| 子 Agent 注册管理 | Registry Service + PostgreSQL | Go 1.24+ / PostgreSQL 16+ | 版本可追溯，授权与审计可控 |
| Prompt 管理 | Prompt Registry Service | Go 1.24+ | 支持用户维护 prompt 与回滚 |
| Memory 管理 | Memory Profile Service | Go 1.24+ | 支持用户维护记忆策略 |
| Knowledge 管理 | Knowledge Space Binding Service | Go 1.24+ | 支持用户维护知识空间绑定 |
| 子 Agent 包分发 | GitLab OCI Artifact | GitLab 17+ | 版本化分发、灰度与回滚 |
| 运行时隔离 | gVisor RuntimeClass | 稳定版 | 用户代码强隔离，降低逃逸风险 |
| 网络隔离 | K8s NetworkPolicy | K8s 1.31+ | 默认拒绝，白名单出网 |
| 质量评估 | Judge Agent Service | 稳定版 | 补齐效果级可观测 |
| 队列与事件 | Kafka | 3.8+ | 高吞吐异步解耦 |
| 主数据库 | PostgreSQL | 16+ | 事务能力强、企业成熟度高 |
| 全文检索 | OpenSearch | 2.x | 检索能力成熟，便于权限过滤 |
| 向量检索 | Milvus | 2.4+ | 大规模向量检索性能更优 |
| 身份认证 | JWT Token（RS256） | JWT RFC 7519 | 无状态鉴权，跨服务校验简单 |
| 网关治理 | 自研 API Gateway（Go） | Go 1.24+ | 与租户策略、审计和路由深度集成 |
| 可观测 | OTel + Prometheus + Grafana + Loki + Tempo | 稳定版 | 指标、日志、追踪全链路闭环 |
| CI/CD | GitLab CI + Argo CD | 稳定版 | GitOps 可追溯、可灰度、可回滚 |
| 部署平台 | Kubernetes + Kustomize | K8s 1.31+ / Kustomize 5+ | 原生配置管理与环境一致性 |

## 5. Temporal 与 LangGraph 分工

- Temporal（Macro-Orchestration）
- 请求生命周期管理。
- 总超时控制与补偿。
- 人机交互挂起/唤醒。
- 异步回调与计费状态。

- LangGraph（Micro-Reasoning）
- Main Agent 推理回路。
- 子 Agent 选择与顺序控制。
- 工具调用决策与结果聚合。

- 锁定规则
- Main Agent 一次推理循环封装为 Temporal Activity。
- 禁止使用 LangGraph Checkpointer 替代 Temporal 持久化。

## 6. 上下文管理与 Token 预算

- 租户级配置：`max_context_window`、`budget_limit`。
- 裁剪策略：`FIFO`、`Summary`、`Importance`。
- 执行策略：每次 Main Agent 执行前先运行 Context Optimizer。
- 目标：降低 token 成本、缓解 “Lost in the Middle”。

## 7. 质量可观测（Evals）

- 在线评估：Judge Agent 评分（相关性、准确性、无害性）。
- 用户反馈：点赞/点踩、标准答案回填。
- 指标沉淀：`answer_quality_score`、`ground_truth_match_rate`、`hallucination_rate`。
- 回流机制：反馈进入评估数据集，驱动 prompt/策略优化。

## 8. 安全与隔离基线

- 认证：JWT Token（Access/Refresh）。
- 授权：RBAC + 资源级策略（Main Agent、Sub-Agent、Tool、Skill、Knowledge）。
- 沙箱：用户自定义 Sub-Agent 与 Skill Runtime 默认运行在 gVisor。
- 网络：NetworkPolicy 默认拒绝，Egress Gateway 白名单出网。
- 配额：租户级 ResourceQuota + LimitRange。
- 审计：策略命中、子 Agent 调度、工具调用、PMK 变更全量留痕。

## 9. 可靠性与兜底

- 模型级 Fallback：主模型超时或限流时切换备用模型。
- 规则级 Fallback：连续 N 次结构化解析失败，回退 Rule Engine。
- 人工级 Fallback：高风险/多次失败任务转人工队列。
- 故障恢复：熔断、重试、DLQ 重放。

## 10. 开发体验（Dry Run / Debug）

- Dry Run API：不入生产队列的仿真执行。
- Playground：支持 Mock Context 单测 Sub-Agent。
- Remote Debug：开发环境联调 Main Agent 与子 Agent。
- 约束：Dry Run 不写生产状态，不计入正式计费。

## 11. 分阶段实施要求（同一技术栈）

### 11.1 MVP

- Main Agent + 3 个内置 Sub-Agent 上线。
- Claim Check 主链路上线（context_ref_id）。
- Context Optimizer 基础策略（FIFO + Summary）上线。
- 基础模型级 Fallback 上线。

### 11.2 Beta

- 用户自定义 Sub-Agent 准入、审批、灰度、回滚上线。
- gVisor + NetworkPolicy + Egress 白名单落地。
- Evals 服务与反馈回流上线。
- Dry Run/Playground 上线。

### 11.3 GA

- 分层路由（Domain Router）全量上线。
- 规则级/人工级 Fallback 与运维流程闭环。
- 多可用区高可用与灾备演练达标。

## 12. 锁定决策清单

以下决策已确认，作为唯一执行路径：

- 后端固定为 `Go + Python` 双栈。
- Agent 拓扑固定为 `Main Agent + Multi Sub-Agent`。
- 用户可自定义 Sub-Agent，并通过 `Registry + OCI + 审批` 管理。
- 用户可自维护 `skills/prompt/memory/knowledge` 并通过审计治理。
- 编排边界固定为 `Temporal Macro + LangGraph Micro`。
- 大载荷传输固定为 `Claim Check Pattern`。
- 运行时隔离固定为 `gVisor + NetworkPolicy`。
- 降级固定为 `模型级 + 规则级 + 人工级 Fallback`。
- 可观测固定为 `系统级 + 效果级（Evals）` 双层闭环。
- 模型接入固定为 `LiteLLM Model Gateway`。

---

执行原则：Go 负责平台治理与稳定性，Python 负责 AI 执行效率；Main Agent 管编排，Sub-Agent 管执行；用户自定义子 Agent 与用户维护配置必须通过统一治理链路。
