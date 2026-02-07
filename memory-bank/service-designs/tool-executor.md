# tool-executor 详细设计

## 1. 职责范围

- 提供统一工具执行入口（gRPC）。
- 按工具风险等级执行分级隔离（Trusted / Untrusted）。
- 校验工具输入输出 schema（OpenAPI 3.1）。
- 记录工具执行审计、耗时和错误信息。

## 2. 非职责范围

- 不做 Agent 规划与推理。
- 不做业务审批流程。
- 不做知识检索。

## 3. 技术设计

## 3.1 技术栈

- Go 1.22+
- gRPC + Protobuf
- OpenAPI 3.1 schema validator
- gVisor/Kata（非可信执行）
- Redis（限流/熔断状态）

## 3.2 模块划分

- `cmd/server`：gRPC server 启动
- `internal/registry`：工具元数据与版本
- `internal/validator`：输入输出 schema 校验
- `internal/executor/trusted`：进程内执行器
- `internal/executor/sandbox`：沙箱执行器
- `internal/policy`：风险分级与网络策略
- `internal/audit`：审计与指标

## 3.3 gRPC 接口

- `ExecuteTool(ExecuteToolRequest) returns (ExecuteToolResponse)`
- `GetTool(GetToolRequest) returns (ToolSpec)`
- `ListTools(ListToolsRequest) returns (ListToolsResponse)`

请求核心字段：

- `tool_name`
- `tool_version`
- `risk_level`
- `input_json`
- `tenant_id`
- `trace_id`

## 3.4 执行策略

1. 读取工具注册信息。
2. 校验输入 schema。
3. 按 `risk_level` 路由：
- `trusted` -> 进程内执行
- `untrusted` -> gVisor/Kata 沙箱执行
4. 超时控制、重试、熔断。
5. 校验输出 schema 并写审计日志。

## 4. 可控做法

1. Tool 注册与执行解耦，先校验再执行。
2. 执行路径强制超时，防止卡死。
3. Untrusted 工具默认零外网，按白名单开通。
4. 工具级 QPS 限流，避免单工具拖垮系统。
5. 预留拆分点：高风险工具可独立部署专用执行集群。

## 5. 细颗粒度开发计划（每步可独立测试）

| Step | 目标 | 交付物 | 独立测试 | 通过标准 |
| --- | --- | --- | --- | --- |
| 1 | gRPC 骨架 | proto + server + health rpc | `go test ./... -run TestGRPCHealth` | gRPC health 可用 |
| 2 | 工具注册中心 | registry 读写与版本管理 | `go test ./internal/registry -run TestRegistryCRUD` | 工具可按版本读取 |
| 3 | Schema 校验 | OpenAPI 3.1 入/出参校验 | `go test ./internal/validator -run TestSchemaValidation` | 非法输入输出均被拒绝 |
| 4 | Trusted 执行路径 | 进程内执行器 | `go test ./internal/executor/trusted -run TestExecute` | 成功/失败路径正确返回 |
| 5 | Untrusted 执行路径 | 沙箱执行器 + 超时控制 | `go test ./internal/executor/sandbox -run TestSandboxTimeout` | 超时可中断且资源回收 |
| 6 | 风险路由策略 | trusted/untrusted 路由 | `go test ./internal/policy -run TestRiskRouting` | 风险等级映射正确 |
| 7 | 限流与熔断 | 工具级流控与熔断器 | `go test ./internal/policy -run TestRateLimit` | 高并发下系统稳定 |
| 8 | 审计与指标 | 审计日志 + Prom 指标 | `go test ./internal/audit -run TestAuditRecord` | 每次调用都有可追踪记录 |
| 9 | 兼容性回归 | 与 agent-core 协议联测 | `go test ./tests/integration -run TestAgentCoreContract` | 协议字段全量兼容 |

## 6. 拆分阈值

- Untrusted 调用占比 > 40% 且平均执行时长 > 3s，拆分独立沙箱集群。
- 单工具 QPS > 100 时，拆分专用执行池和独立限流策略。
