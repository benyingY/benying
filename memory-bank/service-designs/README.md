# Service Designs

本目录包含服务设计文档与治理规范文档。

## 服务清单（5 个）

- `control-plane.md`：统一入口与控制面（Go）。
- `agent-runtime.md`：Main Agent 与 Sub-Agent 执行面（Python）。
- `model-gateway.md`：模型统一网关（LiteLLM）。
- `tool-gateway.md`：工具调用网关与副作用防火墙（Go）。
- `knowledge-service.md`：知识接入与混合检索服务（Python）。

## 治理规范（非服务）

- `sub-agent-management.md`：用户自定义 Sub-Agent 准入与治理规范。

## 跨服务约定

- 同步调用协议：`gRPC + Protobuf`
- 异步事件协议：`Kafka + CloudEvents`
- 统一追踪字段：`trace_id`、`tenant_id`、`task_id`、`main_agent_id`、`sub_agent_id`
- 大载荷传输：`Claim Check`（引用 ID，不在 RPC 里直传大 payload）
- 编排边界：`Temporal` 负责宏观生命周期，`LangGraph` 负责微观推理
