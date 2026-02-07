# agent-core 详细设计

## 1. 职责范围

- 承载 `Temporal Workflow`（粗粒度业务流程）。
- 在 Activity 中执行 `LangGraph`（细粒度 Agent 推理图）。
- 处理人工审批等待（Signals）和恢复执行。
- 编排知识检索调用与工具调用，并输出结构化运行结果。

## 2. 非职责范围

- 不提供北向鉴权与租户入口。
- 不直接处理文档入库。
- 不承担工具执行沙箱能力（由 `tool-executor` 提供）。

## 3. 技术设计

## 3.1 技术栈

- Python 3.11+
- Temporal Python SDK
- LangGraph
- Pydantic（状态与协议 schema）
- gRPC client（调用 tool-executor）

## 3.2 模块划分

- `workflows/chat_workflow.py`：业务主流程
- `activities/run_graph_activity.py`：LangGraph Activity 包装
- `graph/nodes/*`：Planner/Retriever/Tool/Answerer 节点
- `integrations/knowledge_client.py`：检索客户端
- `integrations/tool_client.py`：工具执行客户端
- `integrations/model_gateway.py`：模型调用抽象
- `approvals/signal_handler.py`：审批信号处理
- `runtime/state_store.py`：运行态存取

## 3.3 状态边界

- `WorkflowState`：业务状态（run_id、approval_state、final_result）
- `GraphState`：推理中间状态，仅在单次 Activity 生命周期内有效
- 禁止双写：GraphState 不写入 Workflow History，只输出最终结构化结果

## 3.4 对外接口（供 platform-api）

- gRPC `CreateRun(session_id, message, agent_id, context)`
- gRPC `StreamRunEvents(run_id)`
- gRPC `SignalApproval(run_id, action, reason, approver)`
- gRPC `GetRun(run_id)`

## 3.5 核心流程

1. `CreateRun` 启动 Temporal Workflow。
2. Workflow 执行 `RunGraphActivity`。
3. Activity 内按 LangGraph 节点运行（检索、工具、生成）。
4. 若命中审批策略，Workflow 进入等待 Signal。
5. 收到 `approve/reject` 后继续/终止，输出最终结果。

## 4. 可控做法

1. Workflow 代码保持确定性，外部 I/O 仅在 Activity 内执行。
2. Activity 设置超时、重试上限与幂等输入。
3. Graph 节点统一 Pydantic schema，避免 JSON 漂移。
4. `tool_client` 和 `knowledge_client` 都有降级路径与熔断器。
5. 预留拆分点：高流量时将 Activity Worker 与 Workflow Worker 分离部署。

## 5. 细颗粒度开发计划（每步可独立测试）

| Step | 目标 | 交付物 | 独立测试 | 通过标准 |
| --- | --- | --- | --- | --- |
| 1 | Worker 框架 | Temporal worker 启动与注册 | `pytest tests/workflows/test_worker_bootstrap.py` | Worker 可拉起并注册 workflow/activity |
| 2 | Workflow 骨架 | `ChatWorkflow` 状态机 | `pytest tests/workflows/test_chat_workflow.py::test_happy_path` | 无外部依赖即可跑通基础流程 |
| 3 | Graph State Schema | Pydantic 状态定义与校验 | `pytest tests/graph/test_state_schema.py` | 非法状态输入被拒绝 |
| 4 | LangGraph Activity | Activity 封装图执行 | `pytest tests/activities/test_run_graph_activity.py` | Activity 输出结构化结果 |
| 5 | 知识检索节点 | knowledge-client 集成节点 | `pytest tests/graph/test_retriever_node.py` | 可返回带引用的检索结果 |
| 6 | 工具调用节点 | tool-client 集成节点 | `pytest tests/graph/test_tool_node.py` | 工具调用失败可重试并记录错误 |
| 7 | 审批信号 | wait/signal/recover 逻辑 | `pytest tests/workflows/test_approval_signal.py` | approve/reject 分支行为正确 |
| 8 | 事件流输出 | run events 映射为流式事件 | `pytest tests/runtime/test_event_stream.py` | 事件顺序稳定、字段完整 |
| 9 | 回放与恢复 | workflow replay 与断点恢复 | `pytest tests/workflows/test_replay.py` | 历史回放无 nondeterministic error |

## 6. 拆分阈值

- Workflow backlog 持续 > 10k 且执行延迟 > 5 分钟，拆分 worker 集群。
- Tool 调用密集型场景占比 > 60% 时，拆出独立 `tool-orchestration-activity` worker。
