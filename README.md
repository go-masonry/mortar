# Mortar

![Go](https://github.com/go-masonry/mortar/workflows/Go/badge.svg)

<p align="center">
    <img src=wiki/logo.svg align="center" height=256>
</p>

Mortar is a lightweight GO framework/library for building gRPC (and REST) web services.
Mortar has out-of-the-box support for configuration, application metrics, logging, tracing, profiling and much more.
While it comes with predefined defaults Mortar gives you total control to fully customize it.

## Overview of Building ~~Blocks~~ Bricks

![Diagram](wiki/diagram.png)

## Motivation

- Focus on business logic
- All services speak the same "design" language
- Reduce boilerplate code
- Swap/Update dependencies/libraries easily
- Test friendly

## [Documentation](https://github.com/go-masonry/tutorial)

Mortar is not a drop-in replacement.

It's important to read its documentation first, starting with the step by step [Tutorial](https://github.com/go-masonry/tutorial) which is also a bit-of-everything example.

### Core Concepts

- [Builders](wiki/builder.md)
- [Middleware](wiki/middleware.md)
- [Dependency Injection](wiki/di.md)
- [Multiple Web Servers](wiki/multiweb.md)

### Everything else (Work in progress)

To understand better some of the internals (without browsing the code) have a look [here](wiki/features.md)

## Scaffolds

To help you bootstrap your services with Mortar [here](https://github.com/go-masonry/scaffolds) you can find different templates.

## Bricks

Mortar defines different interfaces, without implementing all of them.
There are lot of great libraries that can be used to implement them.
They just need to be {[(wrapped)]} first.

We call them [Bricks](wiki/bricks.md).

### Logger

- [zerolog](https://github.com/go-masonry/bzerolog)

### Configuration

- [viper](https://github.com/go-masonry/bviper)

### Monitoring/Metrics

- [prometheus](https://github.com/go-masonry/bprometheus)

### Tracing

- [jaeger](https://github.com/go-masonry/bjaeger)  
