# Model Gateway Service Design

- 服务名：`model-gateway`
- 技术选型：`LiteLLM`（开源）
- 角色：统一模型接入、路由、限流、成本与降级控制
- SLO：可用性 99.95%，P95 < 500ms（不含上游模型推理时长）

## 1. 职责边界

职责：

- 统一模型协议：屏蔽不同厂商 API 差异。
- 模型路由：按租户、场景、预算、延迟策略选择模型。
- 限流与配额：租户级 QPS、并发、token 限额。
- 成本计量：记录 token、费用估算、调用耗时。
- 降级策略：主模型失败时自动切换备模型。
- 审计记录：模型调用请求/响应元数据留痕。

非职责：

- 不做业务编排。
- 不做工具调用。
- 不持有业务主数据。

## 2. 对外接口

gRPC：

- `ChatCompletion`
- `Embedding`
- `Rerank`（可选）
- `ListModels`

HTTP（可选兼容）：

- OpenAI-compatible `/v1/chat/completions`
- OpenAI-compatible `/v1/embeddings`

## 3. 路由策略

- 主路由：`tenant_policy -> workload_type -> model_tier`
- 触发降级条件：
- 超时
- 429/5xx
- 成本超阈值

降级顺序示例：

1. `primary_model`
2. `backup_model`
3. `budget_model`

## 4. LiteLLM 落地要点

- 使用 LiteLLM Proxy 统一管理 provider key 与模型映射。
- 模型映射由配置中心维护（按租户可覆盖）。
- 启用 LiteLLM fallback 与 retry 配置。
- 启用 request logging 并汇总到 OTel/Loki。

## 5. 安全与合规

- API Key 由 Vault/Secret 注入，不在业务日志中打印。
- 请求内容按策略进行脱敏存储。
- 禁止跨租户模型凭据共享。

## 6. 可观测

核心指标：

- `model_request_latency_ms`
- `model_request_success_rate`
- `model_fallback_trigger_rate`
- `token_usage_total`
- `cost_estimate_total`

日志字段：

- `trace_id`、`tenant_id`、`task_id`、`main_agent_id`、`sub_agent_id`
- `model_name`、`provider`、`latency_ms`、`token_usage`、`fallback_path`

## 7. 与其他服务关系

- `agent-runtime` 只调用 `model-gateway`，不直连模型厂商。
- `control-plane` 负责下发租户模型策略与预算配置。
- `knowledge-service` 需要 embedding 时调用 `model-gateway`。

## 8. 故障与降级

- 单 provider 故障：自动切换备用 provider。
- 多 provider 故障：返回结构化错误并触发规则/人工兜底。
- 达到预算阈值：切换低成本模型或降采样上下文。
