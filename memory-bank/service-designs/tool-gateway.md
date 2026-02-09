# Tool Gateway Service Design

- 服务名：`tool-gateway`
- 语言：Go 1.24+
- 角色：工具调用防火墙与副作用治理层
- SLO：可用性 99.95%，P95 < 400ms（不含下游系统耗时）

## 1. 职责边界

职责：

- 工具协议适配（HTTP/gRPC/DB/SaaS）。
- 凭据托管与注入（从 Vault/Secrets 读取）。
- 权限白名单与参数校验。
- 出网控制与审计留痕。
- 幂等保护与副作用操作确认。

非职责：

- 不执行推理逻辑。
- 不存储业务主数据。

## 2. 接口

gRPC：

- `InvokeTool`：执行工具调用。
- `GetToolSchema`：查询工具输入输出 schema。
- `ValidateToolPermission`：调用前权限校验。

## 3. 工具注册模型

关键字段：

- `tool_id`
- `name`
- `protocol`
- `input_schema`
- `timeout_ms`
- `retry_policy`
- `side_effect_level`（read/write/high_risk）
- `allowed_tenants`

## 4. 安全设计

- 默认 deny，按租户+子 Agent 白名单放行。
- SSRF 防护：禁止内网保留地址直接访问。
- Egress 限制：只能访问白名单目标。
- 高风险工具（write/high_risk）需审批标志。

## 5. 可靠性设计

- 超时/重试：按工具配置执行。
- 熔断：按 tool_id 和目标 endpoint 粒度。
- 幂等：写操作需携带 `idempotency_key`。
- 降级：工具不可用时返回可解释错误并触发 Main Agent fallback。

## 6. 审计与可观测

审计字段：

- `trace_id`、`tenant_id`、`task_id`、`sub_agent_id`
- `tool_id`、`target`、`input_hash`、`result_code`

指标：

- `tool_call_success_rate`
- `tool_call_latency_ms`
- `tool_circuit_breaker_open`
- `tool_denied_count`

## 7. 部署建议

- 与 control-plane 独立部署，避免故障域耦合。
- 可按工具类型拆分 worker 池（高风险/低风险）。
