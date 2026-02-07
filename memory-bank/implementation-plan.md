# V1 实施计划（细颗粒度、可单测）

## 1. 目标与范围

V1 交付范围：`web-portal + platform-api + agent-core + knowledge-service + tool-executor`。

实施原则：

1. 每个步骤都可以独立测试并给出通过/失败结论。
2. 先打通最短业务闭环，再增加治理与优化能力。
3. 先模块内稳定，再做跨服务联调。

## 2. 里程碑

- M1：单轮会话闭环（无审批、无工具）
- M2：审批闭环（Temporal Signals）
- M3：RAG 闭环（ACL 预过滤 + 引用）
- M4：工具调用闭环（Trusted/Untrusted）
- M5：稳定性门禁通过（压测、回放、审计）

## 3. 跨服务开发步骤

| Step | 服务 | 目标 | 前置依赖 | 独立测试 | 通过标准 |
| --- | --- | --- | --- | --- | --- |
| 1 | platform-api | 服务骨架 + 健康检查 | 无 | `go test ./... -run TestHealth` | API 健康检查通过 |
| 2 | agent-core | Temporal worker 基础框架 | 无 | `pytest tests/workflows/test_worker_bootstrap.py` | worker 可启动并注册 |
| 3 | web-portal | 路由骨架与登录页壳 | 无 | `pnpm --filter web-portal test route-smoke` | 核心页面可访问 |
| 4 | platform-api | 鉴权中间件 | Step 1 | `go test ./internal/auth -run TestJWT` | token 校验与租户注入正确 |
| 5 | platform-api | 会话创建/消息写入 API | Step 4 | `go test ./internal/session -run TestSessionCRUD` | 会话与消息持久化正确 |
| 6 | agent-core | ChatWorkflow 骨架 | Step 2 | `pytest tests/workflows/test_chat_workflow.py::test_happy_path` | 基础流程可跑通 |
| 7 | agent-core | gRPC 接口（CreateRun/StreamRunEvents/GetRun） | Step 6 | `pytest tests/runtime/test_grpc_contract.py tests/runtime/test_event_stream.py` | gRPC 协议字段一致且事件流顺序稳定 |
| 8 | platform-api | 对接 agent-core（非流式） | Step 5,7 | `go test ./tests/integration -run TestCreateRun` | API 可成功触发 run |
| 9 | web-portal | 聊天页发送与结果展示 | Step 3,8 | `pnpm --filter web-portal test chat-send` | UI 单轮会话可完成 |
| 10 | agent-core | LangGraph Activity 接入 | Step 6 | `pytest tests/activities/test_run_graph_activity.py` | Activity 产出结构化结果 |
| 11 | platform-api | SSE 流式桥接 | Step 8,10 | `go test ./internal/stream -run TestSSEBridge` | 流事件顺序正确 |
| 12 | web-portal | SSE 增量渲染 | Step 9,11 | `pnpm --filter web-portal test sse-stream` | 流式输出无乱序 |
| 13 | agent-core | 审批等待与 Signal 恢复 | Step 10 | `pytest tests/workflows/test_approval_signal.py` | approve/reject 分支正确 |
| 14 | platform-api | 审批 API（approve/reject） | Step 13 | `go test ./internal/approval -run TestSignal` | 审批动作成功触发 workflow |
| 15 | web-portal | 审批中心页面 | Step 14 | `pnpm --filter web-portal test approval` | 审批列表与动作可用 |
| 16 | knowledge-service | 入库任务 API + 数据模型 | 无 | `pytest tests/api/test_ingest_job_api.py` | 任务状态机可用 |
| 17 | knowledge-service | 检索 API（Hybrid） | Step 16 | `pytest tests/search/test_hybrid_search.py` | 混合检索结果可返回 |
| 18 | knowledge-service | ACL 预过滤 | Step 17 | `pytest tests/search/test_acl_prefilter.py` | 非授权数据零泄漏 |
| 19 | platform-api | 知识 API 代理（ingest/query/reindex） | Step 16,17,18 | `go test ./internal/knowledge -run TestKnowledgeProxy` | 状态码与错误码透传一致 |
| 20 | agent-core | 检索节点接入 knowledge-service | Step 10,18 | `pytest tests/graph/test_retriever_node.py` | 检索节点返回引用 |
| 21 | tool-executor | gRPC 骨架 + Schema 校验 | 无 | `go test ./internal/validator -run TestSchemaValidation` | 入/出参校验可用 |
| 22 | tool-executor | Trusted 执行路径 | Step 21 | `go test ./internal/executor/trusted -run TestExecute` | trusted 工具执行正确 |
| 23 | tool-executor | Untrusted 沙箱路径 | Step 21 | `go test ./internal/executor/sandbox -run TestSandboxTimeout` | 超时中断和回收正确 |
| 24 | agent-core | 工具节点接入 tool-executor | Step 22,23 | `pytest tests/graph/test_tool_node.py` | 工具节点成功与失败路径都可控 |
| 25 | platform-api | 审计日志全链路 | Step 14,20,24 | `go test ./internal/audit -run TestAuditLog` | 写操作均可审计 |
| 26 | 全链路 | E2E 主流程 | Step 12,15,20,24,25 | `make test-e2e-mainflow` | 聊天->审批->检索->工具 全链路通过 |

## 4. 集成测试门禁

- Gate A（M1）：Step 1-12 全部通过。
- Gate B（M2）：Step 13-15 全部通过。
- Gate C（M3）：Step 16-20 全部通过。
- Gate D（M4）：Step 21-24 全部通过。
- Gate E（M5）：Step 25-26 通过，且关键链路 p95、错误率、审计覆盖率达标。

## 5. 稳定性验收指标（V1）

- API 可用性 >= 99.9%
- 会话主链路 p95 <= 2.5s（不含模型首 token 等待）
- 工具调用失败自动重试成功率 >= 95%
- ACL 误放行率 = 0
- 审批动作可追溯率 = 100%

## 6. 风险与缓解

- 风险：跨服务联调晚，导致集中爆雷。
- 缓解：严格按 Step 独立测试先过，再进入 Gate。

- 风险：knowledge-service 入库任务挤压在线查询。
- 缓解：独立队列与 worker 池，设置资源配额和隔离。

- 风险：Temporal Workflow 非确定性错误。
- 缓解：所有外部 I/O 下沉 Activity，固定 replay 回归测试。
