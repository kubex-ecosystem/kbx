# Kbx

English version: [../README.md](../README.md)

## Sumário

- [Visão Geral](#visão-geral)
- [Escopo Atual do Produto](#escopo-atual-do-produto)
- [Estado Operacional Atual](#estado-operacional-atual)
- [Capacidades Principais](#capacidades-principais)
- [Visão Geral da Arquitetura](#visão-geral-da-arquitetura)
- [Estrutura do Repositório](#estrutura-do-repositório)
- [Design Library-First](#design-library-first)
- [Convenções de Configuração e Runtime](#convenções-de-configuração-e-runtime)
- [Registry de Providers LLM](#registry-de-providers-llm)
- [Utilitários de Mail e IMAP](#utilitários-de-mail-e-imap)
- [Segurança e Tratamento de Secrets](#segurança-e-tratamento-de-secrets)
- [Superfície de CLI](#superfície-de-cli)
- [Papel no Ecossistema](#papel-no-ecossistema)
- [Limitações Atuais](#limitações-atuais)
- [Screenshots](#screenshots)

## Visão Geral

`Kbx` é o toolkit de infraestrutura compartilhada do ecossistema Kubex.

Ele é usado primeiro como biblioteca Go e depois como superfície de CLI/tooling.

Seu papel é hospedar capacidades reutilizáveis que não devem ficar presas dentro de um único repositório de produto, especialmente:

- carregamento de config e defaults
- helpers de segurança
- pacotes utilitários
- utilitários de mail e IMAP
- tipos compartilhados
- registry de providers e adapters para paths multi-LLM

O `Kbx` agora faz parte do caminho real de runtime de produtos ativos, especialmente o `GNyx`.

## Escopo Atual do Produto

No estado atual, o `Kbx` é mais forte em:

- carregamento compartilhado de config e defaults
- helper packages reutilizáveis
- helpers de segurança e secret storage
- utilitários de envio de mail e leitura IMAP
- registry e adapters de providers LLM
- definições de tipos compartilhados usadas por outros repositórios

Ele também inclui:

- uma superfície de CLI baseada em Cobra
- helpers de manifest e metadata
- scaffolding selecionado de serviços/daemons

## Estado Operacional Atual

Operacionalmente, o `Kbx` não está dormente.

Ele já participa de fluxos reais:

- o `GNyx` usa seu support de runtime/registry de providers
- `GNyx` e `Domus` dependem de comportamentos compartilhados de config/defaults
- a seleção e o comportamento de runtime dos providers no gateway agora dependem de infraestrutura do `Kbx`
- helpers de segurança e utilitários são usados em caminhos práticos do código

## Capacidades Principais

Capacidades práticas atuais incluem:

- carregamento compartilhado de configuração
- helpers de defaults e runtime
- suporte a crypto e secret storage
- metadata e lógica de instanciação de providers
- adapters para `openai`, `gemini`, `groq` e `anthropic`
- utilitários de SMTP e IMAP
- tipos comuns e helper packages de apoio

## Visão Geral da Arquitetura

O `Kbx` é propositalmente amplo, mas continua sendo library-first.

Áreas importantes incluem:

- `load/` e helpers relacionados para config e defaults
- `tools/providers/` para registry e adapters concretos
- `types/` para estruturas compartilhadas de runtime
- pacotes de mail, segurança e utilidades
- `cmd/` para a superfície leve de CLI

## Estrutura do Repositório

```text
cmd/                    entrypoints da CLI
defaults/               defaults de runtime e helpers de config
load/                   helpers de loading
tools/providers/        registry de providers e adapters concretos
types/                  tipos compartilhados do ecossistema
```

## Design Library-First

O principal valor do `Kbx` não é sua CLI. É o fato de múltiplos repositórios consumirem a mesma superfície compartilhada de implementação.

Isso implica duas coisas:

- mudanças no `Kbx` podem ter blast radius imediato no ecossistema
- correção, compatibilidade e comportamento de inicialização importam mais do que aparência

## Convenções de Configuração e Runtime

O `Kbx` ajuda a padronizar convenções de runtime como:

- loading de config
- expansão de defaults
- comportamento sensível a runtime-home em aplicações dependentes
- utilidades que precisam ser consistentes entre repositórios

## Registry de Providers LLM

Esta é uma das áreas mais estrategicamente importantes do `Kbx` agora.

A consolidação recente tornou o registry materialmente mais confiável para uso real em produto.

Estado prático atual:

- a config de runtime ficou separada do runtime de provider instanciado
- o loading do registry está menos frágil do que antes
- providers suportados são expostos de forma coerente
- o `GNyx` usa essa camada no caminho ativo de execução de providers

Providers suportados na prática hoje incluem:

- `OpenAI`
- `Gemini`
- `Groq`
- `Anthropic`

## Utilitários de Mail e IMAP

O `Kbx` inclui utilitários reutilizáveis para:

- envio SMTP
- leitura/fetch via IMAP
- helpers relacionados a mail que podem ser reaproveitados em vários projetos Kubex

Essas áreas não são hoje o centro da evolução do ecossistema, mas continuam sendo infraestrutura compartilhada valiosa.

## Segurança e Tratamento de Secrets

Os helpers de segurança hoje incluem áreas como:

- helpers de criptografia
- suporte a secret storage
- helpers de encoding/geração

Esses utilitários não são apenas toolkit decorativo. Eles são reaproveitados quando o tratamento de segredos e material de runtime precisa ser consistente.

## Superfície de CLI

Existe uma CLI, mas ela é secundária em relação ao valor de biblioteca do repositório.

A pergunta de engenharia mais importante para o `Kbx` normalmente não é “qual comando existe?”, mas sim “de qual contrato compartilhado de runtime os consumidores dependem?”.

## Papel no Ecossistema

Hoje o `Kbx` é dependência real de:

- `GNyx`
- `Domus`
- outros códigos de runtime e automação do lado Kubex

Sua importância estratégica cresceu bastante quando o runtime de providers passou a ser caminho crítico para features de produto.

## Limitações Atuais

Limitações atuais incluem:

- a maturidade é desigual entre módulos
- algumas áreas do toolkit são muito mais battle-tested que outras
- a camada de providers é recente o bastante para ainda exigir hardening contínuo
- a superfície de CLI é menos central do que a superfície de biblioteca

## Screenshots

Sugestões de placeholders:

- `[Screenshot Placeholder: debug do provider registry]`
- `[Screenshot Placeholder: ajuda da CLI]`
