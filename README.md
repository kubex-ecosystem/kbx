# Kbx

Kbx is the shared utility and infrastructure toolkit of the Kubex ecosystem.

It is used as a Go library first, and as a small executable/tooling surface second. Its job is to provide reusable building blocks that should not live inside a single product repository: configuration loading, runtime defaults, mail utilities, security helpers, lightweight concurrency primitives, manifest handling, and the provider registry used for multi-LLM execution.

Today, Kbx is already part of the real runtime path for projects such as `GNyx` and `Domus`. It is not documentation-only scaffolding, and it is not just a bag of helpers. It is shared infrastructure code with active blast radius across multiple repositories.

> Current posture: real and useful, but heterogeneous. Some modules are mature and operationally relevant right now, while others are still more like toolkit surface, compatibility code, or planned CLI expansion.

## Table of Contents

- [What Kbx Is](#what-kbx-is)
- [Current Product Scope](#current-product-scope)
- [Current Operational Status](#current-operational-status)
- [Core Capabilities](#core-capabilities)
- [Architecture Overview](#architecture-overview)
- [Repository Layout](#repository-layout)
- [Library-First Design](#library-first-design)
- [Configuration and Runtime Conventions](#configuration-and-runtime-conventions)
- [LLM Provider Registry](#llm-provider-registry)
- [Mail and IMAP Utilities](#mail-and-imap-utilities)
- [Security and Secret Handling](#security-and-secret-handling)
- [CLI Surface](#cli-surface)
- [Using Kbx as a Go Dependency](#using-kbx-as-a-go-dependency)
- [Ecosystem Role](#ecosystem-role)
- [Current Limitations](#current-limitations)
- [Roadmap Direction](#roadmap-direction)
- [Screenshots](#screenshots)
- [License](#license)

## What Kbx Is

Kbx is the common toolbox layer for Kubex applications.

In practical terms, it currently serves as:

1. a shared configuration and defaults layer
2. a reusable utility library for Go services
3. a mail and IMAP support package
4. a security helper layer for crypto and secret storage
5. a provider registry for LLM integrations
6. a lightweight executable/tooling surface for selected operations

Kbx should not be understood as a standalone product in the same way as `GNyx`. It is infrastructure code for products, services, automation flows, and runtime coordination across the Kubex ecosystem.

## Current Product Scope

As of the current codebase state, Kbx is strongest in these areas:

- shared config/default loading
- generic helper packages such as `get`, `is`, `load`, and `tools`
- security utilities for encryption and secret storage
- SMTP mail sending and IMAP fetch utilities
- LLM provider registry and provider adapters
- shared type definitions reused by other projects

Kbx also includes:

- a Cobra-based CLI surface
- manifest and metadata helpers
- service/daemon-oriented command scaffolding
- documentation/build support scripts

However, not every part of Kbx has the same level of maturity. The library surface is currently more important and more operationally relevant than the CLI surface.

## Current Operational Status

Operationally, Kbx is already being used in real flows.

Examples of current practical value:

- `GNyx` uses Kbx defaults, config helpers, and provider registry support.
- `Domus` uses Kbx config and utility layers.
- the LLM provider registry is now part of the active provider runtime path used by `GNyx`
- security helpers are used to generate, encode, encrypt, and persist secrets
- mail helpers provide a reusable path for SMTP send and IMAP read capabilities

So Kbx is not merely “available for future use”. It already participates in the live runtime of the ecosystem.

## Core Capabilities

### 1. Shared configuration and defaults

Kbx provides reusable defaults and parsing logic for:

- server/runtime settings
- logging configuration
- mail configuration
- OAuth/auth provider configuration
- LLM provider configuration
- path conventions across the ecosystem

Packages central to this layer include:

- `get/`
- `is/`
- `load/`
- `types/`
- `internal/module/kbx/`

### 2. LLM provider registry

Kbx includes a provider registry used to load, normalize, instantiate, and resolve LLM providers from configuration.

The active and implemented provider adapters currently include:

- OpenAI
- Gemini
- Groq
- Anthropic

This registry is no longer hypothetical. It is part of the runtime path already consumed by `GNyx`.

### 3. Mail utilities

Kbx includes:

- SMTP send helpers
- provider-specific mail send integrations
- IMAP mailbox reading
- email rendering and template-loading helpers
- retry-aware mailer orchestration

### 4. Security helpers

Kbx provides reusable security primitives such as:

- symmetric encryption helpers
- file-based encrypted secret storage
- Vault-backed secret storage
- Redis-backed secret storage
- interface-based abstractions for secret handling

### 5. Generic utility primitives

The `tools/` layer includes reusable helpers for:

- retry handling
- finite state machines
- queues
- wait groups
- manifest parsing and validation
- mapper and style helpers

## Architecture Overview

At a high level, Kbx is organized like this:

```text
Public library surface
  -> kbx.go
  -> get/
  -> is/
  -> load/
  -> types/
  -> tools/

Feature toolkits
  -> mailing/
  -> tools/mail/
  -> tools/providers/
  -> tools/security/

Executable / CLI surface
  -> cmd/
  -> cmd/cli/
  -> internal/module/

Build / packaging support
  -> support/
  -> Makefile
```

The key architectural point is this:

- Kbx is primarily a reusable shared library
- the executable exists, but the library surface carries more strategic weight

## Repository Layout

```text
cmd/                    Binary entrypoints and CLI glue
get/                    Generic value/env/default helpers
is/                     Validation and safety helpers
load/                   Config parsing and runtime hydration helpers
mailing/                Higher-level mailer and IMAP flows
tools/                  Generic utility packages and subsystems
tools/mail/             Provider-based mail send helpers
tools/providers/        LLM provider registry and adapters
tools/security/         Crypto and secret storage utilities
types/                  Shared config and domain types
internal/module/        CLI/module metadata and manifest plumbing
support/                Build, install, validation, and docs scripts
```

## Library-First Design

The main way to understand Kbx is as a dependency.

The root package re-exports and simplifies access to several internal facilities, including:

- default path helpers
- config loaders
- LLM config builders
- mail config builders
- server/logging config loaders
- shared types such as `ChatRequest`, `ChatChunk`, `Email`, `MailConnection`, and `LLMConfig`

This means consumers can use Kbx without needing to navigate every internal package directly.

## Configuration and Runtime Conventions

Kbx encodes shared path and runtime conventions used across Kubex projects.

Examples from current defaults include:

- `~/.kubex`
- `~/.kubex/gnyx/...`
- `~/.kubex/domus/...`
- default server host/port values
- default provider API key env names
- default mail/config path helpers

This is useful because it gives multiple repositories a common vocabulary for runtime files, certificates, secrets, and config bootstrapping.

It also means changes in Kbx defaults can affect several consumers, so these defaults are part of the real architecture, not just convenience sugar.

## LLM Provider Registry

The provider subsystem lives under `tools/providers/`.

Its job is to:

- load provider configuration
- normalize provider names and runtime config
- resolve API keys from environment-backed references
- instantiate available providers
- expose a consistent provider interface for chat and capability inspection

Core concepts include:

- `LLMConfig`
- `LLMProviderConfig`
- `Provider`
- `ProviderExt`
- `ChatRequest`
- `ChatChunk`
- `Usage`

### What is actually implemented today

Implemented adapters:

- `openai`
- `gemini`
- `groq`
- `anthropic`

There are also broader config defaults for additional future providers, but the supported runtime registry surface is narrower than the full list of names present in config defaults.

That distinction matters.

### Why this matters

This is now a real shared subsystem with immediate impact on upstream applications. In practice, it is one of the highest-value parts of Kbx because it enables multi-provider AI execution without forcing each product repository to reinvent the same contracts.

## Mail and IMAP Utilities

Kbx contains two related but distinct mail layers.

### `tools/mail/`

This is the lower-level provider-based SMTP sending layer.

Current provider map includes:

- Gmail
- Outlook
- Microsoft
- `sendmail` is present as a conceptual fallback path, but not part of the active provider map in the same way as the others

### `mailing/`

This is the higher-level orchestration layer that adds:

- request-to-email conversion
- retries
- timeouts
- template sending
- IMAP read support

This split is useful:

- `tools/mail/` is the delivery primitive
- `mailing/` is the workflow-oriented layer

## Security and Secret Handling

The security subsystem under `tools/security/` is one of the more reusable parts of Kbx.

It currently includes:

- ChaCha20-Poly1305 based crypto utilities
- file-based encrypted keyring replacement
- Vault-backed secret storage
- Redis-backed secret storage
- interfaces for pluggable secret backends

The file-based secret storage is especially relevant in local and self-hosted Kubex flows because it avoids depending on desktop keyring behavior.

This is infrastructure code, not just helper code.

## CLI Surface

Kbx includes a Cobra-based executable, but the CLI surface is lighter and less strategically central than the library surface.

The active root module/CLI metadata exists, along with command packages for areas such as:

- `mail`
- `daemon`
- `service`
- `version`

Important caveat:

- parts of the CLI are still scaffold-like or only partially realized
- some command descriptions still reflect legacy or copied text and should not be treated as authoritative product documentation

So the executable should be read as a secondary surface, not as the primary identity of the repository.

## Using Kbx as a Go Dependency

Minimal example using the root package:

```go
package main

import (
    "context"

    kbx "github.com/kubex-ecosystem/kbx"
    registry "github.com/kubex-ecosystem/kbx/tools/providers"
)

func main() {
    cfg := kbx.NewLLMConfig()
    _ = cfg

    rg, err := registry.Load("./providers.yaml")
    if err != nil {
        panic(err)
    }

    _, err = rg.Chat(context.Background(), kbx.ChatRequest{
        Provider: "gemini",
        Messages: []kbx.Message{{Role: "user", Content: "hello"}},
    })
    if err != nil {
        panic(err)
    }
}
```

Minimal example using mail configuration helpers:

```go
package main

import kbx "github.com/kubex-ecosystem/kbx"

func main() {
    params := kbx.NewMailSrvParams("./smtp.json")
    _ = params
}
```

## Ecosystem Role

Kbx is a cross-project dependency layer.

### In relation to GNyx

Kbx supplies:

- provider registry support
- runtime defaults
- config loading helpers
- shared auth/logging/server config helpers
- mail and secret-management building blocks

### In relation to Domus

Kbx supplies:

- configuration defaults
- manifest/build helpers
- runtime utility code
- supporting abstractions reused by the data-service layer

### In relation to Logz

Kbx builds on `Logz` for its logging behavior, rather than trying to become a logging framework itself.

## Current Limitations

Kbx is valuable, but it is not yet perfectly uniform.

Current limitations include:

- the repository mixes highly reusable infrastructure with lighter helper code and partially realized CLI ideas
- the CLI surface is less mature and less central than the library surface
- some module metadata and descriptions still carry legacy or copied naming/text
- configuration defaults include more provider names than the actively implemented registry adapters
- some packages still reflect ecosystem evolution rather than a fully consolidated architectural pass

These limitations do not negate the usefulness of the repository. They simply define how it should be read accurately.

## Roadmap Direction

The practical direction for Kbx is clear.

### Near-term consolidation

- keep strengthening Kbx as a shared library first
- continue hardening the LLM provider registry
- keep security and mail utilities stable for reuse
- reduce legacy/ambiguous CLI metadata and copied descriptions

### Mid-term consolidation

- formalize runtime contracts between Kbx and consuming projects
- tighten config/default conventions across the ecosystem
- keep broad helper growth under control so the package remains coherent

### Future expansion

- expand provider support only where there is a real consumer path
- treat new subsystems as infrastructure products, not as arbitrary helper accumulation
- preserve the distinction between “available in config defaults” and “implemented in runtime adapters”

## Screenshots

Placeholders for future documentation assets:

- `[Placeholder] Provider registry flow diagram`
- `[Placeholder] Mail + IMAP usage screenshot`
- `[Placeholder] Secret storage/runtime path diagram`
- `[Placeholder] Library usage example screenshot`

## License

This repository is licensed under the [MIT License](./LICENSE).
