# 各层详细技术选型 (Technical Stack)

我们以 Python 和 Go 为主（AI 生态主要在 Python，基础设施与性能敏感服务以 Go 承载），构建一套企业级、可私有化部署的方案。
全局约束：所有组件必须可私有化部署（可自托管），不依赖外部托管的关键能力。

## 1. 接入与多租户层 (Access & Multitenancy)
目标：统一入口、鉴权、租户识别与流量治理。
- API Gateway：APISIX 或 Kong（外部入口）
  - 负责流量分发、API Key 校验、基础限流与租户识别。
- PoC 过渡方案：先用 Go Gin 实现轻量接入服务/网关能力，验证路由、上下文与链路追踪；后续再切换到 APISIX。
- 认证与租户：JWT Token（tenant-id/role 从 token claims 解析；必要时映射表辅助）
- 限流与配额：Redis + 令牌桶/滑动窗口（按 tenant/user 维度限流；策略可配置）

## 2. 模型接入与抽象层 (Model Gateway)
目标：统一管理 LLM，解耦上层应用与底层模型。
- 开源/标准组件：LiteLLM（强烈推荐）
  - Python 库（也有 Proxy 模式），统一 100+ 模型为 OpenAI API 兼容接口。
- 私有化推理服务：vLLM 或 Xinference
  - 自建 Llama / Qwen 等模型时使用，vLLM 在推理速度与显存利用率上优势明显。
- 真实模型 Provider：支持多厂商可插拔（由 LiteLLM 统一适配与路由）

## 3. 编排与推理引擎 (Orchestration)
目标：处理对话逻辑、状态流转与任务编排。
- 核心框架：LangGraph（LangChain 生态）
  - 图 + 状态机，适合循环、分支与多步推理场景。
- Web 框架：FastAPI（核心编排与 API 服务）
  - 异步能力强，适合 LLM 调用与流式响应。
- 多 Agent 框架（可选）：Microsoft AutoGen 或 CrewAI
  - 若强调“多角色协作”，可集成。

## 4. 记忆与知识库 (RAG & Memory)
目标：存向量、存历史、解析文档。
- 向量数据库：Qdrant（首选）或 Milvus
  - Qdrant 部署轻量、性能高；Milvus 适合超大规模向量场景。
- ETL/文档解析：Unstructured.io 或 LlamaParse
  - 解析 PDF / PPT / Excel 等非结构化文档。
- 缓存/会话存储：Redis
  - 短期记忆与流式响应状态存储。
- 多租户数据与配置存储：Postgres

## 5. 工具与行动层 (Tools & Sandbox)
目标：安全地执行代码和调用 API。
- 代码沙箱：Docker 容器 + gVisor（或 Firecracker）
  - 必须隔离运行 Agent 生成的代码，支持用完即焚，满足私有化部署。
- API 定义：OpenAPI (Swagger)
  - 工具接口标准化，无需额外选型。

## 6. 开发与调试环境 (IDE / Playground)
目标：提供可视化编排与调试界面。
- 前端框架：Next.js (React) + TypeScript
- 流程图编排 UI：React Flow（拖拽式工作流画布）
- UI 组件库：ShadcnUI + TailwindCSS

## 7. 运行时与基础设施 (Runtime)
目标：支持异步执行、可观测、可治理。
- 任务队列：Celery + RabbitMQ
  - 处理长耗时 Agent 任务。
- 基础设施与高性能服务：Go（网关扩展、计费、审计、调度等）
- 可观测与调试：Langfuse（强烈推荐）
  - 记录调用耗时、Token 消耗与输入输出（可脱敏）。
- 链路追踪与指标：OpenTelemetry + Prometheus + Grafana
- 审计日志存储与检索：Elasticsearch
- 安全护栏：NVIDIA NeMo Guardrails 或 Guardrails AI
  - 输入输出合规拦截。
- 容器与编排：Docker + Kubernetes
- 配置与密钥：Vault / 云 KMS
