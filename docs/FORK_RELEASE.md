# 本仓库的发布、部署与升级

本分支以 `lqy007700/sub2api` 为发布来源：

| 用途 | 默认地址 |
| --- | --- |
| 代码与部署脚本 | https://github.com/lqy007700/sub2api |
| 二进制 Release | https://github.com/lqy007700/sub2api/releases |
| 容器镜像 | `ghcr.io/lqy007700/sub2api:latest` |

安装器和本分支编译的后端更新服务均读取本仓库的 GitHub Releases。Compose 的三个发布镜像配置、Apple container 默认镜像及部署文档均使用上述 GHCR 镜像。上游项目为 [Wei-Shaw/sub2api](https://github.com/Wei-Shaw/sub2api)，原有许可和署名应保留；Go 模块路径继续使用 `github.com/Wei-Shaw/sub2api`。

## 首次发布前的设置

1. 在本仓库启用 GitHub Actions，并允许 Release 工作流运行。工作流的发布任务申请 `contents: write` 和 `packages: write`，分别用于 GitHub Release 和 GHCR；组织或企业策略仍可能限制这些权限。
2. GHCR 使用 Actions 自动提供的 `GITHUB_TOKEN`，镜像命名空间根据当前仓库所有者转为小写生成。无需为 GHCR 设置 DockerHub 凭据。如果同名包已存在，确认本仓库有该包的 Actions 访问权限。
3. 首次发布镜像后，在 GitHub Packages 中确认 `sub2api` 包的可见性。匿名部署要求包为 Public；私有镜像需要先在部署机器登录 GHCR，并具备 `read:packages` 权限。仓库公开不代表首次创建的包一定公开。
4. DockerHub 是可选的：只有同时配置仓库 Secrets `DOCKERHUB_USERNAME` 和 `DOCKERHUB_TOKEN`，完整发布才会同步推送 `<DOCKERHUB_USERNAME>/sub2api`。令牌需要该 DockerHub 仓库的推送权限；只配置其中一个时跳过 DockerHub。默认部署仍使用 GHCR。
5. `TELEGRAM_BOT_TOKEN` 和 `TELEGRAM_CHAT_ID` 是可选发布通知配置，不影响默认 GHCR 来源；未配置时跳过通知。
6. 默认分支上的 VERSION 同步任务会尝试推送版本文件。若主分支保护阻止 Actions 直接推送，应由维护者手动同步 `backend/cmd/server/VERSION`，并检查该任务的结果。

GitHub Release 二进制和校验和需要能匿名下载。`UPDATE_GITHUB_TOKEN` 仅用于 GitHub REST API 的更新检查和版本查询，不会附带到 Release 资源下载；`GITHUB_TOKEN`、`GH_TOKEN` 不作为安装器或运行服务的更新令牌。仅提供该令牌不能保证私有 Release 的安装或升级可用。

## 先构建验证，再正式发布

在 Actions 的 **Release → Run workflow** 中，选择包含本次代码的分支，填写 `tag=main` 并启用 `dry_run`，可在不创建正式版本标签的情况下验证选定构建。完整模式会构建前端、五个平台二进制、两种 Linux 架构的 OCI 镜像，并检查归档和校验和；dry run 不登录或推送镜像，不发布 GitHub Release，也不运行 VERSION 同步或 Telegram 通知。

正式发布使用新的 `v` 前缀版本标签，例如 `vX.Y.Z`（请替换为真实且未使用的版本）：

```bash
# 先提交并推送准备发布的代码，再由维护者执行：
git tag -a vX.Y.Z -m "vX.Y.Z" -m "本次发布说明"
git push origin vX.Y.Z
```

标签推送触发 Release 工作流。也可以手动运行工作流，填写已经存在的版本标签并关闭 `dry_run`。发布阶段会确认选定源码提交与该标签一致；发布说明读取附注标签正文。准备任务按标签写入构建用 VERSION，发布成功后再同步默认分支中的版本文件。

完整发布的产物为：

- Linux amd64/arm64、macOS amd64/arm64 的 `sub2api_X.Y.Z_<os>_<arch>.tar.gz`，以及 Windows amd64 的 ZIP。
- GitHub Release 中的 `checksums.txt`，供安装器和在线升级校验二进制。
- GHCR 的 amd64/arm64 镜像及多架构标签 `X.Y.Z`、`latest`、`X.Y`、`X`。镜像标签不带开头的 `v`。

`simple_release=true` 或仓库变量 `SIMPLE_RELEASE=true` 仅发布 Linux amd64 GHCR 镜像和 Release 说明，不上传二进制归档或校验和，也不发布 arm64 或 DockerHub。它还会更新 `latest` 为 amd64 镜像，因此需要 Apple silicon 原生容器或二进制在线升级时应使用完整发布。预发布版本也会更新镜像的移动标签；生产部署建议固定正式版本标签或镜像摘要。

## 部署本仓库的版本

首次使用前，必须先成功发布本仓库的镜像或二进制；仅修改默认地址不会自动生成远端产物。

```bash
# Docker 部署准备：在空部署目录执行，随后检查生成的 .env
curl -fsSL https://raw.githubusercontent.com/lqy007700/sub2api/main/deploy/docker-deploy.sh | bash
docker compose up -d

# Linux 二进制安装（需要 PostgreSQL、Redis 及 root 权限）
curl -fsSL https://raw.githubusercontent.com/lqy007700/sub2api/main/deploy/install.sh | sudo bash
```

要固定 Docker 版本，将所用 Compose 配置中 `sub2api` 服务的 image 改为 `ghcr.io/lqy007700/sub2api:X.Y.Z`。Apple container 则设置 `.env` 中的 `APPLE_CONTAINER_SUB2API_IMAGE=ghcr.io/lqy007700/sub2api:X.Y.Z`。Apple container 需要包含 arm64 的完整发布。

## 升级及从上游部署迁移

先备份数据库、应用数据和 `.env`，保留 JWT 与 TOTP 密钥。已有 Compose 部署需将实际使用的配置中的 `sub2api` image 切换到自己的 GHCR 地址，再运行：

```bash
# 一键脚本生成的文件名为 docker-compose.yml
docker compose pull sub2api
docker compose up -d sub2api

# 若使用仓库中的本地目录版配置，则加上相应参数
docker compose -f docker-compose.local.yml pull sub2api
docker compose -f docker-compose.local.yml up -d sub2api
```

Apple container 修改上述镜像设置后，运行 `./apple-container.sh pull` 和 `./apple-container.sh up`。

Linux 二进制可通过本仓库安装器升级：

```bash
curl -fsSL https://raw.githubusercontent.com/lqy007700/sub2api/main/deploy/install.sh | sudo bash -s -- upgrade
```

已经运行的上游二进制仍使用它原来编译进去的更新源。先安装或部署一次本分支构建的版本，此后管理后台的更新检查、下载和回退版本列表才会来自本仓库。自动更新要求新版本号高于当前版本；需要切换到指定版本时可使用安装器的 `install --version vX.Y.Z`，已安装相同版本时安装器会跳过。在线更新还依赖完整发布中与运行平台匹配的归档和校验和，Simple Release 无法提供这些资源。

本次仓库重定向仅能通过本地检查验证脚本配置和测试。GitHub Actions 权限、GHCR 包访问权限与可见性、真实发布、镜像拉取及实际服务器升级，需要维护者在发布后检查，不能从本地测试结果推定已经成功。
