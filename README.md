# Examples
[English](./README.md) | [简体中文](./README.zh_CN.md)

## Overview
`examples` collects runnable sample projects for
[**Golaxy Distributed Service Development Framework**](https://github.com/pangdogs/framework)
and [**Golaxy Core**](https://github.com/pangdogs/core).

This repository is intended to answer a practical question: how the framework
and core packages are meant to be assembled in real programs. The examples
cover both the low-level execution model from `core` and the distributed
service, gateway, routing, RPC, and infrastructure add-ins provided by
`framework`.

The projects here are functional demos rather than full products. If you need a
larger scaffold closer to production structure, see
[SIMHA](https://github.com/pangdogs/simha) and
[scaffold](https://github.com/pangdogs/scaffold).

## What This Repository Provides
- End-to-end examples that show how services, runtimes, entities, components,
  and add-ins fit together.
- Small focused demos for individual official add-ins such as broker,
  distributed entities, discovery, distributed services, distributed sync,
  gateway, and RPC.
- A larger chat sample that combines `gate`, `router`, `dent`, `dsvc`, and
  `rpc` in one workflow.
- Reference startup commands and local dependency setup for trying the examples
  directly.

## Repository Layout
The repository is organized into three top-level areas:

- `app`: more complete application-style examples.
- `core`: focused examples for the Golaxy Core execution model.
- `official_addins`: focused examples for individual framework add-ins.

### Package Guide
| Path | Responsibility |
| --- | --- |
| [`./app/demo_chat`](./app/demo_chat) | Anonymous chat demo with a gate service, a chat service, distributed user entities, channel routing, and an interactive CLI. |
| [`./core/demo_addin`](./core/demo_addin) | Minimal add-in example built on Golaxy Core. |
| [`./core/demo_ec`](./core/demo_ec) | Minimal entity-component example built on Golaxy Core. |
| [`./official_addins/demo_broker`](./official_addins/demo_broker) | Broker add-in example. |
| [`./official_addins/demo_dent`](./official_addins/demo_dent) | Distributed entity registry/query example. |
| [`./official_addins/demo_discovery`](./official_addins/demo_discovery) | Service discovery example. |
| [`./official_addins/demo_dsvc`](./official_addins/demo_dsvc) | Distributed service addressing and messaging example. |
| [`./official_addins/demo_dsync`](./official_addins/demo_dsync) | Distributed synchronization example. |
| [`./official_addins/demo_gate`](./official_addins/demo_gate) | Gateway and client example. |
| [`./official_addins/demo_rpc`](./official_addins/demo_rpc) | RPC processor and forwarding example. |

## Example Notes
### `app/demo_chat`
`demo_chat` is the most complete sample in this repository.

It demonstrates:

- bootstrapping multiple services in one application
- creating global user entities on both `gate` and `chat`
- mapping sessions to entities through `router`
- querying distributed entity locations through `dent`
- forwarding service and client RPC across `gate` and `chat`
- running a terminal client with `rpcli` and Bubble Tea

Directory structure:

| Path | Responsibility |
| --- | --- |
| [`./app/demo_chat/server`](./app/demo_chat/server) | Service bootstrap and service assemblers. |
| [`./app/demo_chat/server/comps`](./app/demo_chat/server/comps) | Gate-side and chat-side user/channel components. |
| [`./app/demo_chat/cli`](./app/demo_chat/cli) | Interactive terminal client. |
| [`./app/demo_chat/consts`](./app/demo_chat/consts) | Shared service names, entity names, and channel constants. |
| [`./app/demo_chat/bin`](./app/demo_chat/bin) | Demo key files used by the sample client/server. |
| [`./app/demo_chat/docker-compose.yaml`](./app/demo_chat/docker-compose.yaml) | Local dependency setup for ETCD and NATS. |

### `core/*`
The `core` demos are intentionally small. They are useful when you want to
understand entity/component lifecycle and add-in wiring without the extra
distributed infrastructure.

### `official_addins/*`
The `official_addins` demos isolate one add-in at a time. They are useful when
you want to verify the startup model or configuration shape of a single
capability before composing it into a larger project.

## Quick Start
### Requirements
- Go `1.25+`
- Docker and Docker Compose if you want to run dependency-backed examples
- ETCD for examples that use discovery, distributed entities, routing, or
  distributed services
- NATS for examples that use the default broker layer

### Run `demo_chat`
1. Start dependencies:

```bash
cd app/demo_chat
docker compose up -d etcd nats
cd ../..
```

2. Start the server:

```bash
go run ./app/demo_chat/server \
  --cli_pub_key ./app/demo_chat/bin/cli.pub \
  --serv_priv_key ./app/demo_chat/bin/serv.pem
```

3. Start the client in another terminal:

```bash
go run ./app/demo_chat/cli \
  --cli_priv_key ./app/demo_chat/bin/cli.pem \
  --serv_pub_key ./app/demo_chat/bin/serv.pub
```

4. Use the client console:

- `create <channel>`
- `remove <channel>`
- `join <channel>`
- `leave <channel>`
- `switch <channel>`
- `rtt`
- any other input sends a chat message to the current channel

### Run Smaller Demos
Most smaller demos can be started directly with `go run`, for example:

```bash
go run ./core/demo_ec
go run ./core/demo_addin
go run ./official_addins/demo_rpc
```

Dependency-backed demos usually need the corresponding service such as ETCD or
NATS running first.

## Installation
The module currently targets the Go version declared in [`go.mod`](./go.mod).

```bash
go get git.golaxy.org/examples@latest
```

## Related Repositories
- [Golaxy Core](https://github.com/pangdogs/core)
- [Golaxy Distributed Service Development Framework](https://github.com/pangdogs/framework)
- [Golaxy Game Server Scaffold](https://github.com/pangdogs/scaffold)
- [SIMHA](https://github.com/pangdogs/simha)
