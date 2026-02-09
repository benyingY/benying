# 企业级通用 Agent 平台实施计划

- 文档版本：v1.3
- 状态：Execution Ready
- 更新时间：2026-02-09
- 关联文档：`architecture.md`、`tech-stack.md`

## 1. 锁定范围

本计划仅基于以下唯一技术路线执行：

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
- 模型接入：LiteLLM Model Gateway
- 兜底策略：模型级 + 规则级 + 人工级 Fallback
- 身份认证：JWT Token（RS256）+ Refresh Token
- 网关：自研 API Gateway（Go）
- CI/CD：GitLab CI + Argo CD
- 部署：Kubernetes + Kustomize

## 2. 实施阶段

## 2.1 Phase 1 - MVP（主链路可用）

目标：打通 Main Agent 编排多个内置 Sub-Agent 的端到端主链路，并控制大载荷与上下文成本。

交付项：

- Go 控制面上线：API Gateway、Agent API、Policy Engine、Tool Gateway、Sub-Agent Registry、PMK Config Service、Context Claim Service。
- Python 执行面上线：Main Agent Coordinator、3 个内置 Sub-Agent、Skill Runtime、Context Optimizer（FIFO + Summary）。
- Model Gateway 上线：LiteLLM 路由、限流、成本计量、模型级 fallback。
- Temporal + LangGraph 嵌套链路跑通，明确 Activity 边界。
- Claim Check 上线：请求和结果大载荷通过 `context_ref_id/result_ref_id` 传递。
- JWT 鉴权闭环：Access/Refresh、租户隔离、权限校验。
- 基础模型级 Fallback 上线：主模型失败切备模型。
- 可观测基础闭环：指标、日志、追踪、告警。

验收标准：

- 主链路任务成功率 >= 95%（staging 压测场景）。
- Sub-Agent 并发调度成功率 >= 95%。
- 基础接口 P95 <= 800ms（不含长任务执行时间）。
- Claim Check 覆盖率 100%（超过阈值载荷）。
- 所有请求均带 `trace_id`、`main_agent_id`、`sub_agent_id`。

## 2.2 Phase 2 - Beta（治理增强）

目标：支持用户自定义 Sub-Agent 与效果级可观测，并提升安全隔离与开发效率。

交付项：

- 自定义 Sub-Agent 准入流程：创建、校验、审核、发布、授权、停用。
- Sub-Agent 包管理：OCI 分发、版本管理、灰度策略、快速回滚。
- 用户自维护配置：skills/prompt/memory/knowledge 的变更、灰度、回滚、审计。
- 安全隔离：gVisor RuntimeClass、NetworkPolicy 默认拒绝、Egress 白名单。
- Evals 上线：Judge Agent 在线评分 + 用户反馈回流。
- Dry Run/Playground 上线：Mock Context、联调仿真。
- 弹性增强：故障注入演练、熔断隔离、DLQ 回放流程。

验收标准：

- 自定义 Sub-Agent 接入时长 <= 1 个工作日（流程完备前提下）。
- 子 Agent 执行失败可定位率 >= 99%。
- 线上评估分可用率 >= 99%。
- 故障注入场景恢复时间 <= 15 分钟。

## 2.3 Phase 3 - GA（企业级稳定）

目标：完成大规模场景下的路由效率、兜底能力与高可用达标。

交付项：

- 分层路由上线：Domain Router + Domain Sub-Agent 分组调度。
- 规则级 Fallback：结构化输出连续失败自动回退规则引擎。
- 人工级 Fallback：高风险任务转人工队列闭环。
- 多可用区高可用与跨区灾备演练。
- 全链路安全合规验收（认证、授权、审计、数据保护）。

验收标准：

- 关键链路可用性 >= 99.9%。
- 达成 RTO <= 30 分钟，RPO <= 5 分钟。
- 自定义子 Agent 治理与审计检查零阻断项。

## 3. 工作流与依赖关系

- 先完成 Temporal/LangGraph 边界和 Claim Check，再扩展复杂推理能力。
- 先完成 Registry、审批、审计，再开放生产租户自定义 Sub-Agent。
- 先上线系统级可观测，再上线效果级可观测（Evals）。
- 先控制上下文成本，再推进多子 Agent 复杂路由。

## 4. 工作分解（WBS）

## 4.1 控制面（Go）

- 自研 API Gateway：路由、限流、鉴权、追踪注入。
- Agent API：任务 API、幂等、租户隔离。
- Workflow Orchestrator：Temporal 生命周期、补偿与超时。
- Sub-Agent Registry：元数据、版本、审批、授权、灰度策略。
- PMK Config Service：prompt、memory、knowledge 配置与版本管理。
- Context Claim Service：大载荷存取与引用 ID 管理。
- Dry Run API：仿真执行入口。

## 4.2 执行面（Python）

- Main Agent Coordinator：任务分解与子 Agent 调度。
- Context Optimizer：Budget、裁剪、摘要。
- Sub-Agent Router：分层路由与并发策略。
- 内置 Sub-Agent 集：检索、工具执行、分析。
- 自定义 Sub-Agent Runtime：包拉取、沙箱执行、输出标准化。
- Evals Judge Worker：在线质量评分。

## 4.3 模型服务（LiteLLM）

- `model-gateway` 服务部署与高可用配置。
- 多 provider 模型路由与 fallback 策略配置。
- 租户级模型预算、限流、成本统计接入。
- 模型调用审计字段标准化。

## 4.4 平台基础设施

- PostgreSQL、Redis、Kafka、OpenSearch、Milvus、MinIO 部署。
- OCI Registry 子 Agent 包管理。
- gVisor RuntimeClass、NetworkPolicy、Egress Gateway。
- Kustomize 环境分层（dev/staging/prod）。
- GitLab CI + Argo CD 持续交付链路。

## 4.5 可靠性与安全

- JWT 生命周期管理与密钥轮换。
- Sub-Agent Schema 校验、签名校验、白名单机制。
- PMK 配置变更审计、版本回滚与权限校验机制。
- 模型级/规则级/人工级 Fallback 机制。
- 熔断、限流、超时、重试、DLQ。

## 5. 质量门禁

- Go：`golangci-lint` + `go test` 必过。
- Python：`ruff` + `mypy` + `pytest` 必过。
- Sub-Agent：`sub-agent.yaml` Schema 校验、签名校验、漏洞扫描必过。
- Model Gateway：路由策略、fallback 策略、限流策略回归测试必过。
- PMK：prompt/memory/knowledge 变更必须通过权限、审计与版本校验。
- Evals：线上质量评分链路可用性 >= 99%。
- 安全：高危漏洞未豁免不得发布。
- 性能：P95、成功率、错误率必须满足阶段性 SLO。

## 6. 风险与缓解

- 风险 R1：Temporal 与 LangGraph 职责重叠。
- 缓解：强制 Macro/Micro 边界，流程通过 Activity 封装审查。

- 风险 R2：大 Context 导致吞吐与存储膨胀。
- 缓解：Claim Check + Payload 阈值控制 + 生命周期清理。

- 风险 R3：用户自定义 Sub-Agent 安全风险。
- 缓解：gVisor 沙箱 + NetworkPolicy + 签名校验 + 审批准入。

- 风险 R4：上下文增长导致成本和准确率下降。
- 缓解：Context Optimizer + Budget + Pruning 策略。

- 风险 R5：质量不可观测导致业务不信任。
- 缓解：Judge 评分 + 用户反馈 + Ground Truth 回流。

## 7. 里程碑检查清单

- M1：Main Agent + 3 个内置 Sub-Agent 主链路可用（含 Claim Check）。
- M2：用户自定义 Sub-Agent 与 Dry Run/Evals 全流程可用。
- M3：分层路由与三层 Fallback 达标，完成高可用与合规验收。

---

执行要求：所有功能建设与发布必须遵循本计划，不得偏离已锁定技术路线。
