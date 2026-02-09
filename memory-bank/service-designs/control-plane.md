# Control Plane Service Design

- 服务名：`control-plane`
- 语言：Go 1.24+
- 角色：系统守门与控制中枢（OLTP）
- SLO：可用性 99.9%，P95 < 300ms（不含文件直传链路）

## 1. 职责边界

职责：

- API Gateway：统一入口、鉴权、限流、租户识别。
- Auth & Policy：JWT 校验、RBAC、策略拦截。
- Task Management：任务创建、状态查询、幂等控制。
- Temporal Client：启动/查询 Workflow、回调状态推进。
- Sub-Agent Registry：子 Agent 元数据、版本、审批、授权、灰度。
- PMK Config：用户维护 skills/prompt/memory/knowledge。
- Context Claim：大上下文/结果的引用管理（`context_ref_id`/`result_ref_id`）。
- Dry Run API：仿真执行入口（不写生产状态，不计费）。

非职责：

- 不执行 LLM 推理与复杂计算。
- 不直接执行用户代码。

## 2. 关键约束

- 单请求同步路径禁止超过 500ms 的计算逻辑。
- 大文件和大 Context 不经服务内存中转，使用 MinIO Presigned URL。
- Claim Check 阈值默认 128KB（可配置）。

## 3. 外部接口（HTTP）

- `POST /v1/tasks`：创建任务，返回 `task_id`。
- `GET /v1/tasks/{task_id}`：查询任务状态与结果引用。
- `POST /v1/tasks/{task_id}/cancel`：取消任务。
- `POST /v1/sub-agents`：创建子 Agent 草稿。
- `POST /v1/sub-agents/{id}/versions`：注册版本。
- `POST /v1/sub-agents/{id}/review`：提交审核。
- `POST /v1/sub-agents/{id}/activate`：激活版本。
- `PUT /v1/sub-agents/{id}/prompt`：更新 prompt。
- `PUT /v1/sub-agents/{id}/memory`：更新 memory。
- `PUT /v1/sub-agents/{id}/knowledge`：更新 knowledge。
- `POST /v1/sub-agents/{id}/dry-run`：仿真执行。

## 4. 内部接口（gRPC）

- 调用 `agent-runtime`：`StartAgentActivity`、`DryRunAgent`。
- 调用 `tool-gateway`：仅在策略预检阶段做能力探测（可选）。
- 调用 `knowledge-service`：不直接调用，统一由 `agent-runtime` 调用。

## 5. 数据模型（PostgreSQL）

核心表：

- `tasks`：任务主表。
- `task_events`：任务状态流转事件。
- `sub_agents`：子 Agent 定义。
- `sub_agent_versions`：子 Agent 版本。
- `sub_agent_bindings`：租户授权绑定。
- `pmk_profiles`：prompt/memory/knowledge 配置。
- `audit_logs`：审计日志索引。

关键字段：

- 所有表必须包含 `tenant_id`。
- 任务链路字段：`trace_id`、`task_id`、`main_agent_id`、`sub_agent_id`。

## 6. Claim Check 流程

1. 请求到达后判断 payload 大小。
2. 超过阈值：写入 Redis/MinIO，生成 `context_ref_id`。
3. Temporal payload 只传引用 ID。
4. 结果同理回写 `result_ref_id`。

过期策略：

- Redis：短期缓存，默认 TTL 2h。
- MinIO：归档对象，默认 TTL 7d（可按租户策略调整）。

## 7. 安全与合规

- JWT RS256 验签，支持 token 轮换。
- 策略引擎前置，默认拒绝未授权资源。
- 审计记录覆盖：配置变更、子 Agent 激活、任务取消、策略命中。
- 输入防护：请求体大小限制、字段白名单校验。

## 8. 可观测性

指标：

- `http_request_latency_ms`
- `task_create_success_rate`
- `claim_check_hit_rate`
- `policy_deny_rate`

日志字段：

- `trace_id`、`tenant_id`、`task_id`、`operator_id`、`action`

## 9. 故障与降级

- Temporal 不可用：进入重试队列并返回“稍后重试”。
- 数据库慢查询：熔断部分非关键写路径（如异步审计扩展字段）。
- MinIO 不可用：降级为 Redis 短期缓存（容量受限）。

## 10. 扩缩容建议

- 水平扩容优先，按 QPS 线性扩展。
- 强制无状态，Session 放 Redis。
