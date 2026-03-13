# Kbx

Portuguese (Brazil) version: [docs/README.pt-BR.md](./docs/README.pt-BR.md)

## Table of Contents

- [Overview](#overview)
- [Current Product Scope](#current-product-scope)
- [Current Operational State](#current-operational-state)
- [Core Capabilities](#core-capabilities)
- [Architecture Overview](#architecture-overview)
- [Repository Layout](#repository-layout)
- [Library-First Design](#library-first-design)
- [Configuration and Runtime Conventions](#configuration-and-runtime-conventions)
- [LLM Provider Registry](#llm-provider-registry)
- [Mail and IMAP Utilities](#mail-and-imap-utilities)
- [Security and Secret Handling](#security-and-secret-handling)
- [CLI Surface](#cli-surface)
- [Role in the Ecosystem](#role-in-the-ecosystem)
- [Current Limitations](#current-limitations)
- [Screenshots](#screenshots)

## Overview

`Kbx` is the shared infrastructure toolkit of the Kubex ecosystem.

It is used as a Go library first and as a CLI/tooling surface second.

Its role is to host reusable capabilities that should not be trapped inside a single product repository, especially:

- configuration and defaults handling
- security helpers
- utility packages
- mail and IMAP helpers
- shared types
- provider registry and provider adapters for multi-LLM runtime paths

`Kbx` is now part of the real runtime path of active products, especially `GNyx`.

## Current Product Scope

At the current state, `Kbx` is strongest in:

- shared config and defaults loading
- reusable helper packages
- security and secret-storage helpers
- mail send and IMAP read utilities
- LLM provider registry and adapters
- shared type definitions used by other repositories

It also includes:

- a Cobra-based CLI surface
- metadata and manifest helpers
- selected service/daemon scaffolding

## Current Operational State

Operationally, `Kbx` is not dormant.

It already participates in live flows:

- `GNyx` uses its provider registry/runtime support
- `GNyx` and `Domus` rely on shared config/default behaviors
- provider selection and runtime behavior in the gateway now depend on `Kbx` infrastructure
- security helpers and utility packages are used in practical code paths

## Core Capabilities

Current practical capabilities include:

- shared configuration loading
- defaults and runtime helpers
- crypto and secret storage support
- provider metadata and instantiation logic
- provider adapters for `openai`, `gemini`, `groq`, and `anthropic`
- mail/SMTP and IMAP support utilities
- common types and supporting helper packages

## Architecture Overview

`Kbx` is intentionally broad but still library-first.

Important areas include:

- `load/` and related helpers for config and defaults
- `tools/providers/` for provider registry and adapters
- `types/` for shared runtime structures
- mail, security, and utility packages
- `cmd/` for the lightweight CLI surface

## Repository Layout

```text
cmd/                    CLI entrypoints
defaults/               runtime defaults and configuration helpers
load/                   loading helpers
tools/providers/        provider registry and concrete provider adapters
types/                  shared ecosystem types
```

## Library-First Design

The main value of `Kbx` is not its CLI. It is the fact that multiple repositories consume the same shared implementation surface.

That has two implications:

- changes in `Kbx` can have immediate blast radius across the ecosystem
- correctness, compatibility, and initialization behavior matter more than appearance

## Configuration and Runtime Conventions

`Kbx` helps standardize runtime conventions such as:

- config loading
- default expansion
- runtime-home aware behavior in dependent applications
- utility behavior that should be consistent across repositories

## LLM Provider Registry

This is one of the most strategically important areas of `Kbx` right now.

Recent consolidation made the registry materially more reliable for real product usage.

Current practical state:

- runtime config is separated from instantiated provider runtime
- registry loading is less fragile than before
- supported providers are surfaced coherently
- `GNyx` now uses this layer in its active provider execution path

Supported practical providers currently include:

- `OpenAI`
- `Gemini`
- `Groq`
- `Anthropic`

## Mail and IMAP Utilities

`Kbx` includes reusable utilities for:

- SMTP send flows
- IMAP read/fetch flows
- mail-related helpers reusable across Kubex projects

These are not the current center of ecosystem evolution, but they remain valuable shared infrastructure.

## Security and Secret Handling

Security-related helpers currently include areas such as:

- encryption helpers
- secret storage support
- encoding/generation helpers

These utilities are not just decorative toolkit code. They are reused where secret and runtime material handling must be consistent.

## CLI Surface

A CLI exists, but the CLI is secondary to the library value of the repository.

The most important engineering question for `Kbx` is usually not “what command exists?” but rather “what shared runtime contract do downstream consumers depend on?”

## Role in the Ecosystem

Today `Kbx` is a real dependency of:

- `GNyx`
- `Domus`
- other Kubex-side runtime or automation code

Its current strategic importance increased significantly once provider runtime became a critical path for product features.

## Current Limitations

Current limitations include:

- maturity is uneven across modules
- some toolkit areas are far more battle-tested than others
- the provider layer is recent enough that continued hardening is still expected
- the CLI surface is less central than the library surface

## Screenshots

Placeholder suggestions:

- `[Screenshot Placeholder: provider registry debug output]`
- `[Screenshot Placeholder: CLI help]`
