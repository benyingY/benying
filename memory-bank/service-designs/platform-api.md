# platform-api 详细设计

## 1. 职责范围

- 作为统一北向入口，向 `web-portal` 提供稳定 API。
- 负责鉴权、租户隔离、会话管理、Agent 配置管理、审批动作接入。
- 将聊天请求转发至 `agent-core`，并将流式结果转为前端可消费的 SSE。
- 代理知识管理与检索接口到 `knowledge-service`。
- 写入审计日志与请求追踪上下文。

## 2. 非职责范围

- 不执行 Agent 推理。
- 不执行工具调用。
- 不在本服务内执行知识检索与入库（仅做代理转发）。

## 3. 技术设计

## 3.1 技术栈

- Go 1.22+
- `chi`（HTTP 路由）
- `pgx` + `sqlc`（PostgreSQL）
- `go-redis`（会话态缓存、限流）
- OpenTelemetry（trace/metrics/log correlation）

## 3.2 模块划分

- `cmd/server`：启动与配置
- `internal/http`：路由、middleware、handler
- `internal/auth`：JWT 校验、租户上下文注入
- `internal/session`：会话与消息管理
- `internal/agentconfig`：Agent 配置 CRUD 与版本
- `internal/stream`：SSE 网关与背压处理
- `internal/approval`：审批动作到 Temporal Signal 的适配
- `internal/knowledge`：knowledge-service API 代理与错误映射
- `internal/audit`：审计日志记录
- `internal/store`：PostgreSQL/Redis 数据访问

## 3.3 数据模型

- `sessions(id, tenant_id, user_id, agent_id, status, created_at)`
- `messages(id, session_id, role, content, created_at)`
- `agent_configs(id, tenant_id, name, spec_json, version, status, updated_at)`
- `approval_actions(id, run_id, approver_id, action, reason, created_at)`
- `audit_logs(id, actor_id, tenant_id, resource, action, result, trace_id, created_at)`

## 3.4 北向 API

- `POST /v1/sessions`
- `POST /v1/sessions/{session_id}/messages`
- `GET /v1/sessions/{session_id}/stream`
- `GET /v1/runs/{run_id}`
- `GET /v1/approvals`
- `POST /v1/approvals/{approval_id}:approve`
- `POST /v1/approvals/{approval_id}:reject`
- `GET /v1/agents`
- `POST /v1/agents`
- `PUT /v1/agents/{agent_id}`
- `POST /v1/knowledge/ingest-jobs`
- `GET /v1/knowledge/ingest-jobs/{job_id}`
- `POST /v1/knowledge/query`
- `POST /v1/knowledge/reindex`

## 3.5 南向接口

- `agent-core`：gRPC（create run / stream events / query run state）
- `knowledge-service`：HTTP（ingest jobs / query / reindex）
- `Temporal`：通过 `agent-core` 暴露的审批回调接口，不直接耦合 SDK 到 handler

## 4. 可控做法

1. 所有 handler 只做协议转换，不写业务规则。
2. 请求级幂等：`Idempotency-Key` + Redis 去重窗口。
3. SSE 网关单独 worker pool，防止阻塞普通 API。
4. 对 `agent-core` 调用设置超时、重试、熔断。
5. 预留拆分点：`agentconfig` 与 `chat-gateway` 可按流量拆分。

## 5. 细颗粒度开发计划（每步可独立测试）

| Step | 目标 | 交付物 | 独立测试 | 通过标准 |
| --- | --- | --- | --- | --- |
| 1 | 服务骨架 | 启动、配置加载、`/healthz` | `go test ./... -run TestHealth` | 健康检查返回 200 |
| 2 | 基础中间件 | trace-id、recover、request log | `go test ./internal/http -run TestMiddleware` | 请求链路含 trace-id |
| 3 | 鉴权模块 | Keycloak JWT 校验 + tenant 注入 | `go test ./internal/auth -run TestJWT` | 非法 token 被拒绝，合法 token 注入上下文 |
| 4 | 会话 API | 创建会话与写消息 API | `go test ./internal/session -run TestSessionCRUD` | 会话和消息可持久化 |
| 5 | SSE 网关 | 对接 agent-core 流式事件转发 | `go test ./internal/stream -run TestSSEBridge` | 事件顺序正确，无连接泄漏 |
| 6 | Agent 配置 API | 配置 CRUD + 版本字段 | `go test ./internal/agentconfig -run TestVersioning` | 更新产生版本递增 |
| 7 | 审批 API | approve/reject 转发到工作流 | `go test ./internal/approval -run TestSignal` | 审批动作可落库并触发信号 |
| 8 | 知识 API 代理 | ingest/query/reindex 转发与错误映射 | `go test ./internal/knowledge -run TestKnowledgeProxy` | 透传状态码与错误码一致 |
| 9 | 审计日志 | 关键操作全量审计 | `go test ./internal/audit -run TestAuditLog` | 每个写操作都有审计记录 |
| 10 | 稳定性门禁 | 限流、超时、熔断策略 | `go test ./internal/http -run TestRateLimit` | 高并发下错误率与超时率在阈值内 |

## 6. 拆分阈值

- `stream` 路径 QPS > 200 或连接数 > 5k 时，拆出独立 `chat-gateway`。
- Agent 配置写流量占比 > 20% 且变更审计复杂时，拆出 `config-service`。
