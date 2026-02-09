# Sub-Agent Management

- 文档类型：治理规范（非服务）
- 文档版本：v1.0
- 更新时间：2026-02-09

## 1. 目的

定义用户自定义 Sub-Agent 的准入、发布、授权、运行与审计规范。

## 2. 必填信息

- `sub_agent_id`
- `name`
- `owner`
- `category`
- `goal`
- `trigger_rule`
- `prompt_template`
- `input_schema`
- `output_schema`
- `tools`
- `skills`
- `memory_policy`
- `knowledge_spaces`
- `permissions`
- `timeout_ms`
- `retry_policy`
- `fallback_policy`
- `model_policy`
- `version`

## 3. 用户可维护范围

租户可维护：

- `skills`
- `prompt`
- `memory`
- `knowledge`

平台统一维护：

- JWT 签发与密钥
- 全局 RBAC 核心规则
- 沙箱和网络基础策略

## 4. 生命周期

- `DRAFT`
- `REVIEW`
- `PUBLISHED`
- `ACTIVE`
- `DEPRECATED`
- `DISABLED`

## 5. 准入流程

1. 提交 `sub-agent.yaml` 与执行包
2. Schema/测试/漏洞/签名校验
3. 审核通过后发布 OCI
4. 配置租户授权与灰度
5. Main Agent 调度可用版本

## 6. 运行约束

- 仅允许调度 `ACTIVE` 版本
- 仅可访问声明且授权资源
- 大载荷使用 Claim Check（`context_ref_id`/`result_ref_id`）
- 用户代码默认在 gVisor 沙箱执行

## 7. 审计字段

- `trace_id`
- `tenant_id`
- `task_id`
- `main_agent_id`
- `sub_agent_id`
- `sub_agent_version`
- `prompt_version`
- `memory_profile_id`
- `knowledge_space_ids`
- `eval_score`
- `fallback_path`
- `result_status`

## 8. API

- `POST /v1/sub-agents`
- `POST /v1/sub-agents/{id}/versions`
- `POST /v1/sub-agents/{id}/review`
- `POST /v1/sub-agents/{id}/activate`
- `POST /v1/sub-agents/{id}/disable`
- `PUT /v1/sub-agents/{id}/prompt`
- `PUT /v1/sub-agents/{id}/memory`
- `PUT /v1/sub-agents/{id}/knowledge`
- `POST /v1/sub-agents/{id}/dry-run`
- `POST /v1/sub-agents/{id}/feedback`
