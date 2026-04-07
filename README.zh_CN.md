# Examples
[English](./README.md) | [简体中文](./README.zh_CN.md)

## 简介
`examples` 仓库收录了
[**Golaxy 分布式服务开发框架**](https://github.com/pangdogs/framework)
和 [**Golaxy Core**](https://github.com/pangdogs/core) 的可运行示例。

这个仓库主要解决一个实践问题：框架和核心模块在真实项目里应该如何组装。这里既包含 `core` 的底层执行模型示例，也包含 `framework` 提供的分布式服务、网关、路由、RPC 和基础设施 add-in 示例。

这些项目属于功能演示，而不是完整产品。如果你需要更接近实际业务结构的工程，可以继续参考
[SIMHA](https://github.com/pangdogs/simha) 和
[scaffold](https://github.com/pangdogs/scaffold)。

## 这个仓库提供什么
- 展示 service、runtime、entity、component 和 add-in 如何组合的端到端示例。
- 针对单个官方 add-in 的小型独立示例，例如 broker、分布式实体、服务发现、分布式服务、分布式同步、网关和 RPC。
- 一个把 `gate`、`router`、`dent`、`dsvc`、`rpc` 串起来的聊天示例。
- 可直接运行的启动命令和本地依赖说明，便于快速验证。

## 仓库结构
仓库分为三个顶层区域：

- `app`：更接近应用形态的完整示例。
- `core`：针对 Golaxy Core 执行模型的精简示例。
- `official_addins`：针对单个 framework add-in 的精简示例。

### 目录说明
| 路径 | 职责 |
| --- | --- |
| [`./app/demo_chat`](./app/demo_chat) | 匿名聊天示例，包含 gate 服务、chat 服务、分布式用户实体、频道路由和交互式 CLI。 |
| [`./core/demo_addin`](./core/demo_addin) | 基于 Golaxy Core 的最小 add-in 示例。 |
| [`./core/demo_ec`](./core/demo_ec) | 基于 Golaxy Core 的最小实体组件示例。 |
| [`./official_addins/demo_broker`](./official_addins/demo_broker) | Broker add-in 示例。 |
| [`./official_addins/demo_dent`](./official_addins/demo_dent) | 分布式实体注册与查询示例。 |
| [`./official_addins/demo_discovery`](./official_addins/demo_discovery) | 服务发现示例。 |
| [`./official_addins/demo_dsvc`](./official_addins/demo_dsvc) | 分布式服务寻址与消息传递示例。 |
| [`./official_addins/demo_dsync`](./official_addins/demo_dsync) | 分布式同步示例。 |
| [`./official_addins/demo_gate`](./official_addins/demo_gate) | 网关与客户端示例。 |
| [`./official_addins/demo_rpc`](./official_addins/demo_rpc) | RPC 处理与转发示例。 |

## 示例说明
### `app/demo_chat`
`demo_chat` 是这个仓库里最完整的样例。

它主要演示：

- 在一个应用里启动多个服务
- 在 `gate` 和 `chat` 上创建同一个全局用户实体
- 通过 `router` 将 session 映射到实体
- 通过 `dent` 查询分布式实体所在节点
- 在 `gate` 和 `chat` 之间转发服务 RPC 与客户端 RPC
- 使用 `rpcli` 和 Bubble Tea 构建终端交互客户端

目录结构如下：

| 路径 | 职责 |
| --- | --- |
| [`./app/demo_chat/server`](./app/demo_chat/server) | 服务启动入口和 service assembler。 |
| [`./app/demo_chat/server/comps`](./app/demo_chat/server/comps) | gate 侧和 chat 侧的用户/频道组件。 |
| [`./app/demo_chat/cli`](./app/demo_chat/cli) | 交互式终端客户端。 |
| [`./app/demo_chat/consts`](./app/demo_chat/consts) | 共享服务名、实体名和频道常量。 |
| [`./app/demo_chat/bin`](./app/demo_chat/bin) | 示例客户端和服务端使用的密钥文件。 |
| [`./app/demo_chat/docker-compose.yaml`](./app/demo_chat/docker-compose.yaml) | ETCD 与 NATS 的本地依赖启动配置。 |

### `core/*`
`core` 下的示例刻意保持精简，适合在没有分布式基础设施干扰的情况下理解实体/组件生命周期和 add-in 装配方式。

### `official_addins/*`
`official_addins` 下的示例会把单个 add-in 独立出来，适合先理解某一项能力的启动方式、依赖和配置形态，再把它组合进更大的项目。

## 快速开始
### 环境要求
- Go `1.25+`
- 如果要运行依赖外部基础设施的示例，建议准备 Docker 和 Docker Compose
- 使用服务发现、分布式实体、路由、分布式服务的示例需要 ETCD
- 使用默认 broker 的示例需要 NATS

### 运行 `demo_chat`
1. 启动依赖：

```bash
cd app/demo_chat
docker compose up -d etcd nats
cd ../..
```

2. 启动服务端：

```bash
go run ./app/demo_chat/server \
  --cli_pub_key ./app/demo_chat/bin/cli.pub \
  --serv_priv_key ./app/demo_chat/bin/serv.pem \
  --etcd.address localhost:2379
```

3. 在另一个终端启动客户端：

```bash
go run ./app/demo_chat/cli \
  --endpoint localhost:9090 \
  --cli_priv_key ./app/demo_chat/bin/cli.pem \
  --serv_pub_key ./app/demo_chat/bin/serv.pub
```

4. 在客户端控制台中使用：

- `create <channel>`
- `remove <channel>`
- `join <channel>`
- `leave <channel>`
- `switch <channel>`
- `rtt`
- 其他任意输入都会作为当前频道消息发送

### 运行更小的示例
大多数精简示例可以直接通过 `go run` 启动，例如：

```bash
go run ./core/demo_ec
go run ./core/demo_addin
go run ./official_addins/demo_rpc
```

如果示例依赖 ETCD、NATS 等外部服务，请先把对应依赖启动起来。

## 安装
当前模块以 [`go.mod`](./go.mod) 中声明的 Go 版本为准。

```bash
go get git.golaxy.org/examples@latest
```

## 相关仓库
- [Golaxy Core](https://github.com/pangdogs/core)
- [Golaxy 分布式服务开发框架](https://github.com/pangdogs/framework)
- [Golaxy 游戏服务器脚手架](https://github.com/pangdogs/scaffold)
- [SIMHA](https://github.com/pangdogs/simha)
