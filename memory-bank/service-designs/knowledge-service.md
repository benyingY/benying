# knowledge-service 详细设计

## 1. 职责范围

- 管理知识入库任务：抓取、解析、切分、Embedding、索引写入。
- 提供在线检索：ACL 预过滤、Hybrid Search、Rerank、引用拼装。
- 维护文档与分片元数据，支持增量更新与重建索引。

## 2. 非职责范围

- 不承担聊天会话管理。
- 不承担 Agent 主流程编排。
- 不承担工具调用执行。

## 3. 技术设计

## 3.1 技术栈

- Python 3.11+
- FastAPI（HTTP API）
- Celery（异步入库任务）
- PostgreSQL + pgvector
- OpenSearch
- Redis（队列与任务状态）

## 3.2 模块划分

- `api/ingest.py`：入库任务 API
- `api/query.py`：在线查询 API
- `workers/ingest_worker.py`：解析、切分、向量化
- `pipeline/chunker.py`：文本切分策略
- `pipeline/acl_mapper.py`：权限元数据映射
- `search/hybrid.py`：向量+关键词融合检索
- `search/rerank.py`：重排与引用选择
- `store/repository.py`：数据访问层

## 3.3 数据模型

- `documents(doc_id, tenant_id, source, title, acl_json, version, status, updated_at)`
- `chunks(chunk_id, doc_id, tenant_id, content, acl_json, token_count, created_at)`
- `embeddings(chunk_id, tenant_id, vector, model, created_at)`
- `ingest_jobs(job_id, tenant_id, source, status, error, created_at, updated_at)`

## 3.4 API

- `POST /v1/knowledge/ingest-jobs`
- `GET /v1/knowledge/ingest-jobs/{job_id}`
- `POST /v1/knowledge/query`
- `POST /v1/knowledge/reindex`

## 3.5 检索流程（在线）

1. 接收 query + tenant + user ACL。
2. 将 ACL 条件作为 pre-filter 下推到检索层。
3. 执行向量检索 + 关键词检索。
4. 融合召回并 rerank。
5. 返回答案上下文片段与 citation。

## 4. 可控做法

1. 入库与查询资源隔离：独立 worker 池、独立队列。
2. pgvector 按租户/业务线分区，避免单表膨胀。
3. 检索默认先过滤 ACL，禁止后过滤兜底。
4. 解析任务 CPU 限制与熔断，避免拖垮在线查询。
5. 预留拆分点：流量增长后将 `ingest-worker` 独立成单服务。

## 5. 细颗粒度开发计划（每步可独立测试）

| Step | 目标 | 交付物 | 独立测试 | 通过标准 |
| --- | --- | --- | --- | --- |
| 1 | 服务骨架 | FastAPI 启动 + 健康检查 | `pytest tests/api/test_health.py` | 健康检查可用 |
| 2 | 数据表与迁移 | documents/chunks/embeddings/jobs | `pytest tests/store/test_migrations.py` | 迁移可重复执行且无破坏 |
| 3 | 入库任务 API | 创建与查询 ingest job | `pytest tests/api/test_ingest_job_api.py` | job 状态机正确流转 |
| 4 | 文档解析与切分 | parser + chunker 管道 | `pytest tests/pipeline/test_chunker.py` | 切分规则满足 token 上限 |
| 5 | ACL 映射入库 | acl mapper + metadata 写入 | `pytest tests/pipeline/test_acl_mapper.py` | chunk 级 ACL 字段完整 |
| 6 | 向量写入 | embedding 生成与 pgvector 入库 | `pytest tests/store/test_vector_write.py` | 向量写入成功且可查询 |
| 7 | Hybrid 检索 | 向量+关键词融合 | `pytest tests/search/test_hybrid_search.py` | 召回结果包含多源并可排序 |
| 8 | ACL 预过滤 | pre-filter 查询路径 | `pytest tests/search/test_acl_prefilter.py` | 非授权文档零泄漏 |
| 9 | Rerank 与引用 | rerank + citation 输出 | `pytest tests/search/test_rerank_citation.py` | 返回引用可追溯到 chunk_id |
| 10 | 压测与隔离验证 | 入库压测不影响在线查询 | `pytest tests/perf/test_ingest_query_isolation.py` | 在线查询 p95 不超阈值 |

## 6. 拆分阈值

- 入库任务队列延迟持续 > 15 分钟，拆分独立 ingest 服务。
- 在线查询 p95 > 2s 且 QPS > 100，拆分独立 query 服务。
