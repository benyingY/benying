# Knowledge Service Design

- 服务名：`knowledge-service`
- 语言：Python 3.12+（锁定）
- 角色：知识接入、索引、混合检索、重排
- SLO：检索 P95 < 600ms（缓存命中场景）

## 1. 为什么锁定 Python

- 文档解析生态成熟（PDF/Office/HTML）。
- Embedding/Rerank 生态与模型支持更完整。
- 与 RAG/LLM 工具链整合成本更低。

## 2. 职责边界

职责：

- 文档接入（上传、解析、清洗、切片）。
- 通过 `model-gateway` 发起 embedding/rerank 模型调用。
- 向量化与索引构建（Milvus）。
- 关键词索引（OpenSearch）。
- 混合检索（向量 + 关键词）与重排。
- 租户权限过滤与引用证据返回。

非职责：

- 不做任务编排。
- 不做工具副作用调用。

## 3. 数据流水线

1. 文档上传（MinIO）。
2. 解析与切片（chunking）。
3. 通过 `model-gateway` 执行向量化（embedding）。
4. 写入 Milvus + OpenSearch。
5. 更新索引版本与可见性元数据。

## 4. 查询流程

1. 接收 query + tenant_scope + knowledge_spaces。
2. Query Rewrite（可选）。
3. OpenSearch BM25 初筛 + Milvus ANN 召回。
4. 通过 `model-gateway` 执行重排（reranker）。
5. 返回 Top-K + citation + confidence。

## 5. 接口

gRPC：

- `IngestDocument`
- `BuildIndex`
- `Search`
- `ListKnowledgeSpaces`

依赖调用：

- `model-gateway`：`Embedding`、`Rerank`（可选）。

## 6. 关键配置

- `top_k`
- `recall_k`
- `rerank_model`
- `max_chunk_size`
- `embedding_model`

## 7. 安全与隔离

- 所有查询强制 `tenant_id` 过滤。
- 知识空间绑定必须命中 PMK 配置。
- 文档与索引操作全量审计。

## 8. 可靠性与性能

- 索引构建异步化（Kafka 任务队列）。
- 热点查询缓存（Redis）。
- 索引版本化，支持回滚。
- GPU 节点池独立部署（rerank/embedding）。

## 9. 可观测

指标：

- `ingest_success_rate`
- `search_latency_ms`
- `retrieval_hit_rate`
- `rerank_latency_ms`
- `citation_coverage_rate`

日志字段：

- `trace_id`、`tenant_id`、`knowledge_space_id`、`query_hash`

## 10. 扩缩容建议

- 检索 API 与 Ingestion Pipeline 分离部署。
- GPU 工作负载（embedding/rerank）与 CPU 检索节点分池。
