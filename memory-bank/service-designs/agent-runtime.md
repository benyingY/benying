# Agent Runtime Service Design

- 服务名：`agent-runtime`
- 语言：Python 3.12+
- 角色：推理执行面（Temporal Worker + LangGraph Runtime）
- SLO：任务执行成功率 >= 98%（剔除外部依赖不可用）

## 1. 职责边界

职责：

- Main Agent 推理与编排。
- Sub-Agent 路由与并发调度。
- 用户自定义 Sub-Agent 调度（不可信代码沙箱执行）。
- Context Optimizer（预算、裁剪、摘要）。
- 通过 `model-gateway` 统一发起模型调用。
- Evals（Judge 评分）与反馈数据上报。
- Dry Run 仿真执行。

非职责：

- 不直接写业务元数据主库（仅通过 control-plane 或只读副本）。
- 不直连模型厂商，不维护模型 fallback 策略。
- 不持有长期凭据。

## 2. Temporal 与 LangGraph 边界

- Temporal：生命周期、超时、补偿、挂起/唤醒、回调。
- LangGraph：一次推理循环（Thinking -> Sub-Agent/Tool -> Aggregation）。
- 约束：每个 LangGraph 循环映射为 Temporal Activity，状态持久化以 Temporal 为准。

## 3. 执行流程

1. Worker 收到 `context_ref_id`。
2. 从 Claim Store 拉取上下文。
3. 运行 Context Optimizer（Budget + Pruning）。
4. Main Agent 选择领域 -> Sub-Agent Router 分发。
5. 调用 `knowledge-service` / `tool-gateway` / Skills。
6. 通过 `model-gateway` 完成 chat/embedding/rerank 模型调用。
7. 结果聚合、结构化校验。
8. 写回 `result_ref_id`，上报指标与评估分。

## 4. Sub-Agent 执行模型

- 内置 Sub-Agent：同进程 Worker 执行。
- 用户自定义 Sub-Agent：通过 K8s Job/Pod（gVisor）隔离执行。
- 调度策略：并发上限、优先级、租户配额、超时隔离。

## 5. Context Optimizer

策略：

- `FIFO`：移除最旧上下文。
- `Summary`：对历史对话进行摘要压缩。
- `Importance`：按相关性保留片段。

关键配置：

- `max_context_window`
- `budget_limit`
- `max_history_turns`

## 6. Fallback 机制

- 模型级：由 `model-gateway` 执行主模型 -> 备用模型切换，`agent-runtime` 仅消费结果。
- 规则级：连续 N 次结构化失败 -> Rule Engine。
- 人工级：高风险/连续失败 -> Human Queue。

## 7. Evals 与反馈

在线评估输出：

- `answer_quality_score`
- `ground_truth_match_rate`
- `hallucination_risk`

反馈回流：

- 接收点赞/点踩和纠正答案。
- 上报至评估数据集（离线训练/提示词优化）。

## 8. 安全与隔离

- 用户代码运行在 gVisor 沙箱。
- 默认无网络访问，按策略开放目标域名。
- 资源限制：CPU/Memory/执行时长。

## 9. 可观测性

核心指标：

- `agent_run_latency_ms`
- `sub_agent_success_rate`
- `fallback_trigger_rate`
- `token_cost_per_task`
- `dry_run_count`

日志字段：

- `trace_id`、`task_id`、`main_agent_id`、`sub_agent_id`、`prompt_version`

## 10. 扩缩容建议

- Worker 水平扩容按队列积压和 CPU 使用率触发。
- GPU 依赖逻辑（如 rerank/judge）与普通 worker 分池部署。
