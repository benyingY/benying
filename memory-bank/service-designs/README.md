# 微服务详细设计总览（V1）

本目录对应 V1 交付拓扑：`1 个前端应用 + 4 个后端微服务`。

## 1. 组件清单

- `web-portal`：前端应用（会话、轨迹、审批中心）
- `platform-api`：统一入口（鉴权、配置管理、会话 API、知识 API、SSE 转发）
- `agent-core`：Temporal + LangGraph 运行时（业务流程编排 + Agent 推理）
- `knowledge-service`：知识入库与检索（ACL 预过滤 + Hybrid Search）
- `tool-executor`：Go 工具执行器（gRPC、Schema 校验、分级沙箱）

## 2. 关键依赖关系

- `web-portal -> platform-api`
- `platform-api -> agent-core`（会话执行与审批 Signal 转发）
- `platform-api -> knowledge-service`（知识管理）
- `agent-core -> knowledge-service`（在线检索）
- `agent-core -> tool-executor`（工具调用）

## 3. 可控做法（统一约束）

1. 服务内部模块化，服务间接口收敛，避免“先合并后混乱”。
2. `knowledge-service` 入库与查询通过队列和独立 worker 池隔离资源。
3. 所有外部调用必须有超时、重试、熔断和幂等键。
4. 每个服务都定义“拆分阈值”（QPS、延迟、团队边界）以便后续平滑拆分。
5. 每个开发步骤都附带独立测试点，不依赖全链路联调才能验证。

## 4. 文档索引

- `web-portal.md`
- `platform-api.md`
- `agent-core.md`
- `knowledge-service.md`
- `tool-executor.md`
