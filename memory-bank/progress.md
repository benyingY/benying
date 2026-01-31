这个文档记录

1. 今天做了什么
2. 今天学到了什么 / 改变了什么想法
## 进展记录
2026-01-31: 初始化项目文档（README、架构、技术选型、计划）。
2026-01-31: 完成 Step 1 骨架（Python FastAPI + Go Gin），并通过测试（pytest / go test ./...）。
2026-01-31: 完成 Step 2 本地编排与健康检查（docker-compose + /health）。
2026-01-31: 运行测试，pytest 通过；go test ./... 通过（因 Go build cache 权限，使用提升权限运行）。
2026-01-31: 尝试 docker compose up -d：Docker Hub TLS handshake timeout；切换阿里源后 pull denied，未完成 Docker 启动测试。
2026-01-31: 完成 Step 3 接入服务路由与上下文（/invoke 固定响应；request-id 贯穿），新增 Go 测试并通过 go test ./...（提升权限运行）；python3 -m pytest 通过。
2026-01-31: 更新实现计划为真实功能落地版本；技术选型补充 JWT、Redis 限流、OpenTelemetry/Prometheus/Grafana、Postgres、Elasticsearch 与多模型 Provider 支持。
2026-01-31: Step 3 接入层 /invoke 结构化请求/响应落地（校验 input、错误结构化、request-id 贯穿），补充 Go 测试并通过 go test ./...；python3 -m pytest 通过。
2026-01-31: 优化 access 性能：Gin 默认 Release 模式、request-id 生成改为轻量原子计数、/invoke JSON 解析改为直接 Decoder；go test ./... 与 python3 -m pytest 通过。
2026-01-31: access /invoke 支持 SSE（Accept: text/event-stream 或 stream=true），新增 SSE 测试并通过 go test ./... 与 python3 -m pytest。
2026-01-31: access /invoke SSE 改为 OpenAI 风格（data: chat.completion.chunk + [DONE]），更新测试并通过 go test ./... 与 python3 -m pytest。
2026-01-31: Step 4 接入层 JWT 鉴权（静态公钥、RS256 校验、可选 issuer/audience 校验），tenant 解析预留空函数；新增鉴权测试并通过 go test ./... 与 python3 -m pytest。
