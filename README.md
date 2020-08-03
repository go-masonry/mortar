# Mortar 

Mortar is a GO framework/library for building gRPC (and REST) web services.
Mortar has out-of-the-box support for configuration, application metrics, logging, tracing, profiling and much more.
While it comes with predefined defaults Mortar gives you total control to fully customize it. 

## Overview of Building ~~Blocks~~ Bricks

TODO a drawing

## Motivation

- Focus on business logic
- All services speak the same "design" language
- Reduce boilerplate code
- Swap/Update dependencies/libraries easily
- Test friendly

## Documentation

Mortar is not a drop-in replacement. It will probably change the way you code and build services.
We think it's important to read its documentation first, starting with the [Tutorial](https://github.com/go-masonry/tutorial). 

## Core Concepts

- Bricks
- [Builders](wiki/builder.md)
- [Middleware](wiki/middleware.md)
- [Structural typing](https://en.wikipedia.org/wiki/Structural_type_system)
- [Dependency Injection](https://github.com/uber-go/fx)

## Features

- REST via [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
- Monitoring (Metrics) 
- JWT (Consumer)
- Logging
- Configuration
- Tracing
- Profiling, debug, gc stats

## Moulds
Are basically template examples