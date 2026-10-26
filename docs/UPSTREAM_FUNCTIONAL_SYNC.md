# 官方功能同步记录

## 2026-10-26

从 fork `20ae2955fc972cde78ffb31c2540e5e281635877`（`v0.2.18-lqy007700.4`）出发，按主线父提交依次 cherry-pick：

| 官方提交 | 功能 |
| --- | --- |
| `37067e8e340e9440d173a92b4faa05e8ce58368f` | [#8319](https://github.com/Wei-Shaw/sub2api/pull/8319)：按操作、分组、平台及账户类型授权的管理员 API Key |
| `ba81c7fd758e425d67f856d601147ef0c5263d3f` | [#8323](https://github.com/Wei-Shaw/sub2api/pull/8323)：专用图片端点按实际出图张数计费 |

官方版本号、赞助素材及宣传文案不在同步范围，功能同步阶段保留 fork 的 `.4`。本次集成以 `v0.2.18-lqy007700.5` 发起完整发布，包含 Linux amd64 安装包、校验和及 GHCR 镜像；生产部署分别处理。

### 与 fork 的适配

- 认证快照统一为版本 `26`，同时保存隐藏倍率和按张计费开关；拒绝双方原有的 `25` 快照，避免读取缺失字段。
- 合并后的 Group schema 重新生成 Ent 代码，修正运行时字段描述符和排序索引的列位置；直接自动合并会造成类型断言 panic。
- 隐藏倍率继续覆盖公开基础倍率和用户专属倍率；图片独立倍率仍按原规则优先。按张计费不叠加 token 高峰因子，公开展示倍率仍保存请求时快照。
- 按张计费只对允许生图的 OpenAI/Composite 分组及专用图片端点生效；混合 Responses/Chat/WS 请求沿用原有计费选择。

### 迁移

保留已发布的 `247_group_hidden_rate_multiplier.sql` 内容及校验和。迁移以完整文件名记录，所以它可与官方 `247_scoped_admin_api_keys.sql` 共存；随后执行官方 `248_group_image_count_billing.sql` 和 fork `249_group_auth_cache_image_count_billing.sql`。249 更新认证缓存失效触发器，同时保留隐藏倍率、利润控制等已有监测字段。

### 验证方式

后端使用仓库要求的 Go 版本，运行 `go test -tags unit ./...` 和 `go build -tags embed ./cmd/server`；前端执行管理员密钥、分组隐藏倍率、分组复制、组合平台及用量页面的定向 Vitest，并运行完整构建。

数据库测试只连接临时隔离 PostgreSQL。`SCOPED_ADMIN_TEST_DSN` 启用上游管理员密钥的真实数据库测试；`FORK_MIGRATION_TEST_DSN` 还需允许创建数据库，用于分别验证空库与已应用 fork 247 的升级路径、迁移重复执行、直接 SQL 缓存失效及 Ent/认证投影的组合字段读写。每个迁移用例创建并清理独立数据库。

本次结果：Go 全量单测含子用例 24,047 项通过、18 项按环境条件跳过；本次管理员密钥、计费和迁移用例均实际运行。前端定向测试共 97 项通过，完整前端构建、追加类型检查、相关 ESLint、后端嵌入构建及相关包 `go vet` 均通过。原 247 迁移、版本号与发布配置的差异检查为空。

首轮全量运行中，原有 `TestRecordCyberPolicyEvent_RuntimeSnapshotRefreshFailureKeepsStaleScope` 超时一次；其相关源码未改动，合并前基线上连续 10 次及合并后定向复跑通过，最终全量复跑通过。
