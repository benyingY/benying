# 企业级 Agent 平台技术栈（定版 v1.1）

本文档给出平台建设的最终技术栈清单。所有模块只列当前确定采用的技术，不再列候选或替代方案。

## 1. 技术栈总览

- 运行平台：Kubernetes 1.29+、Helm、Argo CD
- 主要语言：TypeScript 5.x、Python 3.11+、Go 1.22+
- 核心存储：PostgreSQL 16+、Redis 7+、S3 兼容对象存储
- 跨服务事件系统：Kafka 3.6+
- 可观测：OpenTelemetry、Prometheus、Grafana、Tempo、Loki
- 安全体系：Keycloak、OPA(Rego)、Vault、KMS

## 2. 分层技术栈

### 2.1 入口与产品层

| 能力 | 采用技术栈 |
| --- | --- |
| Web 门户与 Agent Store | Next.js + TypeScript + Ant Design |
| API 网关 | Kong |
| IM/Bot 接入 | 自研 Adapter + 飞书/Slack 官方 SDK |
| 身份认证与单点登录 | Keycloak（OIDC/SAML） |

### 2.2 Agent 构建层（Control Plane）

| 能力 | 采用技术栈 |
| --- | --- |
| Agent 配置与版本管理 | PostgreSQL + GitOps（配置即代码） |
| Agent 图编排 | LangGraph |
| 工作流持久化与容错 | Temporal SDK（TypeScript/Python）+ Temporal Server |
| Prompt/Policy 管理 | Git 仓库 + PR 审批流 |
| Studio 调试能力 | Trace 回放 + 沙箱执行环境 |

### 2.3 Agent 运行时（Data Plane）

| 能力 | 采用技术栈 |
| --- | --- |
| 编排执行引擎 | Temporal Server |
| Agent 执行粒度 | Temporal Workflow（粗粒度）+ LangGraph Activity（细粒度） |
| 异步任务总线 | Kafka |
| 会话状态管理 | Redis + PostgreSQL |
| 工具执行器 | Go 服务 + gVisor/Kata 容器隔离 + Kubernetes NetworkPolicy |
| 人工审批节点 | Temporal Signals/Queries + 审批 API + PostgreSQL |
| Tool 调用协议 | gRPC + OpenAPI 3.1 JSON Schema |
| 执行隔离策略 | Trusted Tools 进程内执行；Untrusted Tools 强制进入 gVisor/Kata |

### 2.4 能力层（模型/知识/工具）

| 能力 | 采用技术栈 |
| --- | --- |
| Model Gateway | LiteLLM + 自研路由策略服务 |
| 提示缓存 | Redis |
| 向量检索 | PostgreSQL + pgvector |
| 关键词检索 | OpenSearch |
| 知识入库任务调度 | Celery + Redis |
| RAG 权限过滤 | ACL Metadata 预过滤（Pre-Filtering）+ 召回参数调优（ef_search） |
| 文档处理与切分 | Unstructured + 自研清洗管道 |
| 重排模型（Rerank） | bge-reranker |
| 工具注册中心 | PostgreSQL + OpenAPI/JSON Schema |

### 2.5 治理与安全层

| 能力 | 采用技术栈 |
| --- | --- |
| 权限模型 | RBAC + ABAC（OPA/Rego） |
| 密钥管理 | Vault + 云 KMS |
| 数据脱敏与 DLP | 自研策略服务 + DLP 规则引擎 |
| 内容安全防护 | Prompt Injection 检测服务 + 输出审查策略 |
| 审计留痕 | PostgreSQL 审计库 + 对象存储归档（WORM 策略） |

### 2.6 可观测与 SRE

| 能力 | 采用技术栈 |
| --- | --- |
| 分布式追踪 | OpenTelemetry + Tempo |
| 指标监控 | Prometheus + Grafana |
| 日志系统 | Loki |
| 告警与值班联动 | Alertmanager + On-call 平台 |
| 成本归因 | OpenCost + 应用层计费服务（Token、检索、工具调用、向量检索 CPU 分摊） |

## 3. 数据分层

- 事务数据：PostgreSQL（配置、权限、审批、审计索引）
- 高速状态：Redis（会话上下文、限流计数、幂等记录）
- 检索数据：PostgreSQL + pgvector（向量）+ OpenSearch（关键词）
- 向量分区策略：pgvector 按租户/业务线分区，避免单表过大导致检索退化
- 权限索引策略：Embedding 入库同步 ACL Metadata（用户组/部门/租户），检索时先做权限过滤再向量召回
- 大文件与产物：S3（文档原件、日志归档、评测样本）

## 4. CI/CD 与工程规范

- 代码托管：GitHub
- 持续集成：GitHub Actions（单测、静态扫描、镜像构建）
- 持续交付：Argo CD（分环境发布、回滚）
- 镜像仓库：Harbor
- 安全扫描：Trivy

## 5. 版本基线

- Kubernetes 1.29+
- PostgreSQL 16+
- Redis 7+
- Kafka 3.6+
- OpenSearch 2.x
- Python 3.11+
- Node.js 20 LTS
- Go 1.22+
- TypeScript 5.x

## 6. 关键架构约束

1. Temporal 与 LangGraph 分工：Temporal 仅编排业务级长流程（如接单、审批、写回系统），不管理 LangGraph 节点跳转细节；LangGraph 仅负责单次 Agent 推理图执行，作为 Temporal Activity 运行。
2. 状态边界：LangGraph 中间状态仅在单次 Activity 生命周期内维护，不与 Temporal 双写；Activity 结束后仅输出结构化结果与证据，由 Temporal 继续后续流程。
3. 审批机制：统一通过 Temporal Signals/Queries 实现人工介入，避免引入第二套 BPM 引擎。
4. Python 与 Go 协议：Python Agent 通过 gRPC 调用 Go Tool Executor，参数与返回统一遵循 OpenAPI 3.1 Schema。

## 7. 评测与上线门禁

- 评测基线：每个业务 Agent 上线前必须建设 20-50 条 Golden Dataset（标准问答与期望动作）。
- 离线评测：每次发布前执行成功率、引用正确性、工具调用正确性评测。
- 上线门禁：离线评测未达阈值禁止发布到生产环境。
