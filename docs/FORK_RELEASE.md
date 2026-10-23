# 本仓库的发布、部署与升级

本分支以 `lqy007700/sub2api` 为发布来源，只发布 Ubuntu 等 Linux 系统使用的 **amd64 / x86_64** 安装包和镜像。

| 用途 | 默认地址 |
| --- | --- |
| 代码与部署脚本 | https://github.com/lqy007700/sub2api |
| 二进制 Release | https://github.com/lqy007700/sub2api/releases |
| 容器镜像 | `ghcr.io/lqy007700/sub2api:latest` |

安装器和本分支编译的后端更新服务均读取本仓库的 GitHub Releases。三个发布镜像 Compose 配置及部署文档均使用上述 GHCR 镜像。上游项目为 [Wei-Shaw/sub2api](https://github.com/Wei-Shaw/sub2api)，原有许可和署名应保留；Go 模块路径继续使用 `github.com/Wei-Shaw/sub2api`。

## 发布目标与版本号

当前主分支后续发布只生成 Linux amd64 产物，不再构建 Linux arm64、macOS 或 Windows，也不再提供 Apple container 部署入口。Ubuntu 服务器运行 `uname -m` 应返回 `x86_64`；`aarch64` 服务器无法使用这些原生产物。

版本号沿用 `v<官方版本>-lqy007700.<修订号>`，例如 `v0.2.18-lqy007700.3`。同一官方版本再次发布时递增最后的修订号；合并下一官方版本后从 `.1` 开始。`lqy007700.N` 表示本分支的正式修订号；GoReleaser 显式配置 `prerelease: false` 和 `make_latest: "true"`，避免根据连字符后缀自动判定为预发布。

**默认安装器与后台自动更新查询 GitHub 的最新正式 Release。** 正式发布必须能通过 `https://api.github.com/repos/lqy007700/sub2api/releases/latest` 查询，并包含 Linux amd64 安装包和校验文件。后台先比较官方版本，再按数字比较本分支修订号，例如 `.10` 高于 `.2`。

`0.2.18-lqy007700.2` 及更早构建会忽略 fork 修订号差异。同一官方版本内首次升级到包含修复的 `.3` 时，需要手动更新镜像，或显式安装指定版本；升级后后台可正确发现后续修订。页面手动刷新会跳过默认 20 分钟的更新缓存。

## 首次发布前的设置

1. 在本仓库启用 GitHub Actions，并允许 Release 工作流运行。发布任务申请 `contents: write` 和 `packages: write`，分别用于 GitHub Release 和 GHCR；组织或企业策略仍可能限制这些权限。
2. GHCR 使用 Actions 自动提供的 `GITHUB_TOKEN`，命名空间根据当前仓库所有者转为小写生成。无需 DockerHub 凭据。如果同名包已存在，确认本仓库有该包的 Actions 访问权限。
3. 在 GitHub Packages 中确认 `sub2api` 包的可见性。匿名部署要求 Public；私有镜像需要先在部署机器登录 GHCR，并具备 `read:packages` 权限。仓库公开不代表包一定公开。
4. DockerHub 可选：只有同时配置 Secrets `DOCKERHUB_USERNAME` 和 `DOCKERHUB_TOKEN`，完整发布才会同步推送 `<DOCKERHUB_USERNAME>/sub2api`。只配置其中一个时跳过，默认部署仍使用 GHCR。
5. `TELEGRAM_BOT_TOKEN` 和 `TELEGRAM_CHAT_ID` 为可选通知配置，未配置时跳过。
6. 发布成功后会尝试将版本写入默认分支的 `backend/cmd/server/VERSION`。若分支保护阻止 Actions 推送，应由维护者手动同步并检查该任务结果。

GitHub Release 安装包和校验和需要能匿名下载。`UPDATE_GITHUB_TOKEN` 仅用于 GitHub REST API 更新检查与版本查询，不会附带到 Release 资源下载；`GITHUB_TOKEN`、`GH_TOKEN` 不作为安装器或运行服务的更新令牌。仅提供该令牌不能保证私有 Release 的安装或升级可用。

## 先构建验证，再发布

普通 `main` 推送触发 CI，不会自动发布。Release 工作流不以 CI 成功为依赖，应先确认测试通过再创建发布标签。

在 Actions 的 **Release → Run workflow** 中选择待发布分支，填写 `tag=main` 并启用 `dry_run`，可先验证构建。完整模式构建前端、一个 Linux amd64 二进制安装包和一个同架构 OCI 镜像，并检查归档和校验和。dry run 不登录或推送镜像、不发布 GitHub Release，也不运行 VERSION 同步或 Telegram 通知。

发布使用新的、尚未使用的标签，例如：

```bash
# 仅示例：先提交并推送准备发布的代码，再创建标签。
git tag -a v0.2.18-lqy007700.3 -m "v0.2.18-lqy007700.3" -m "合并官方 v0.2.18；保留隐藏倍率；仅发布 Linux amd64。"
git push origin v0.2.18-lqy007700.3
```

标签推送触发 Release。也可手动运行工作流，填写已存在的版本标签并关闭 `dry_run`。准备任务锁定标签对应的源码提交，写入构建用 VERSION；发布说明读取附注标签正文，成功后再同步默认分支的版本文件。

完整发布包括：

- GitHub Release：`sub2api_X.Y.Z-lqy007700.N_linux_amd64.tar.gz` 和 SHA-256 `checksums.txt`。
- GHCR：单一 `linux/amd64` 镜像，标签为 `X.Y.Z-lqy007700.N-amd64`、`X.Y.Z-lqy007700.N`、`latest`、`X.Y`、`X`。镜像标签不带开头的 `v`。
- 配置了完整 DockerHub 凭据时，同步同架构镜像及上述标签。

每次只构建一份 Linux amd64 镜像并直接标记各个标签，不需要 QEMU 或跨架构清单合并。

`simple_release=true` 或仓库变量 `SIMPLE_RELEASE=true` 仍表示**仅镜像发布**：只发布 Linux amd64 GHCR 镜像和 Release 说明，不上传二进制或校验和，也不推送 DockerHub。默认完整模式同样只有 amd64，但保留安装器和指定版本升级需要的安装包。发布会更新镜像的移动标签，生产部署建议固定具体版本或摘要。

## 部署本仓库的版本

必须先成功发布对应镜像或二进制。以下示例版本应替换为实际发布版本。

```bash
# Docker 部署准备：在空部署目录执行，随后检查生成的 .env。
curl -fsSL https://raw.githubusercontent.com/lqy007700/sub2api/main/deploy/docker-deploy.sh | bash
# 将 Compose 中 sub2api 的 image 固定为以下地址，然后启动：
# ghcr.io/lqy007700/sub2api:0.2.18-lqy007700.3
docker compose up -d

# Ubuntu / Linux amd64 二进制安装，需要 PostgreSQL、Redis 与 root 权限。
curl -fsSL https://raw.githubusercontent.com/lqy007700/sub2api/main/deploy/install.sh | sudo bash -s -- install --version v0.2.18-lqy007700.3
```

## 升级及从上游部署迁移

先备份数据库、应用数据和 `.env`，保留 JWT 与 TOTP 密钥。已有 Compose 部署需将实际使用配置中的 `sub2api` image 切换为本仓库的具体 GHCR 版本，再运行：

```bash
# 一键脚本生成的配置。
docker compose pull sub2api
docker compose up -d sub2api

# 使用仓库本地目录版配置时。
docker compose -f docker-compose.local.yml pull sub2api
docker compose -f docker-compose.local.yml up -d sub2api
```

Linux amd64 二进制切换到指定版本：

```bash
curl -fsSL https://raw.githubusercontent.com/lqy007700/sub2api/main/deploy/install.sh | sudo bash -s -- install --version v0.2.18-lqy007700.3
```

已经运行的上游二进制仍使用原来编译进去的更新源。先安装或部署一次本分支构建的版本，之后后台的版本查询与回退列表才来自本仓库。默认自动更新仍只查询最新正式 Release；Simple Release 不提供二进制在线升级需要的资源。

发布范围精简只影响采用新配置的后续发布。先前多平台版本的远端标签、安装包和镜像予以保留，不追溯删除。
