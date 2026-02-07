# web-portal 详细设计

## 1. 职责范围

- 提供业务用户界面：聊天页、执行轨迹页、审批中心、Agent 配置页。
- 处理前端会话状态、UI 鉴权状态、流式输出渲染。
- 所有业务数据通过 `platform-api` 获取，不直连后端内部服务。

## 2. 非职责范围

- 不实现业务编排逻辑。
- 不直连 Temporal、PostgreSQL、Redis、OpenSearch。
- 不做敏感权限决策（仅做展示层控制，权限以后端返回为准）。

## 3. 技术设计

## 3.1 技术栈

- Next.js（App Router）
- TypeScript
- React Query（服务端状态）
- Zustand（轻量 UI 状态）
- SSE（流式消息渲染）
- Playwright（E2E）

## 3.2 页面与模块

- `app/chat`：聊天与流式响应渲染
- `app/runs/[runId]`：执行轨迹与步骤详情
- `app/approvals`：审批列表与审批动作
- `app/agents`：Agent 配置管理
- `lib/api-client`：统一请求封装（重试、超时、错误映射）
- `lib/auth`：OIDC 登录态管理

## 3.3 关键交互

1. 用户在聊天页发送消息。
2. 前端调用 `platform-api` 创建会话并发起流式请求。
3. 接收 SSE 事件并增量渲染。
4. 若返回 `approval_required` 事件，跳转审批中心。

## 4. 接口契约（对 platform-api）

- `POST /v1/sessions`
- `POST /v1/sessions/{session_id}/messages`
- `GET /v1/sessions/{session_id}/stream`（SSE）
- `GET /v1/runs/{run_id}`
- `GET /v1/approvals`
- `POST /v1/approvals/{approval_id}:approve`
- `POST /v1/approvals/{approval_id}:reject`

## 5. 可控做法

- UI 模块按页面边界拆分，不做跨页面共享复杂状态。
- SSE 采用事件类型白名单，未知事件直接降级为日志提示。
- 前端缓存设置短 TTL，防止陈旧审批状态。
- 预留拆分点：`chat` 与 `ops-console` 可拆成独立前端子应用。

## 6. 细颗粒度开发计划（每步可独立测试）

| Step | 目标 | 交付物 | 独立测试 | 通过标准 |
| --- | --- | --- | --- | --- |
| 1 | 搭建应用骨架与路由 | 基础 layout + 4 个页面路由 | `pnpm --filter web-portal test route-smoke` | 4 个路由可访问，构建通过 |
| 2 | 接入登录态管理 | OIDC 登录/登出、token 刷新 | `pnpm --filter web-portal test auth` | 未登录跳转登录，登录后可访问受保护页面 |
| 3 | 封装 API Client | 统一错误处理与超时重试 | `pnpm --filter web-portal test api-client` | 4xx/5xx/超时都能映射到统一错误码 |
| 4 | 聊天页基础版 | 消息列表 + 输入框 + 发送流程 | `pnpm --filter web-portal test chat-send` | 可完成单轮消息发送与展示 |
| 5 | SSE 流式渲染 | token 增量渲染 + 结束事件处理 | `pnpm --filter web-portal test sse-stream` | 流式文本无乱序、可正确结束 |
| 6 | 执行轨迹页 | run timeline 与步骤详情面板 | `pnpm --filter web-portal test run-trace` | 步骤状态可视化正确 |
| 7 | 审批中心 | 列表查询 + approve/reject 操作 | `pnpm --filter web-portal test approval` | 审批动作提交后状态即时更新 |
| 8 | E2E 主流程 | 登录 -> 发起会话 -> 审批 -> 完成 | `pnpm --filter web-portal test:e2e main-flow` | 主流程 E2E 稳定通过 |

## 7. 拆分阈值

- Chat 页面月活 > 1 万且运维后台需求激增时，拆分 `ops-console`。
- 前端构建时间 > 10 分钟或 E2E > 30 分钟时，按业务域拆分工程。
