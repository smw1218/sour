# Sour

Sour is an opinionated microservices "framework" for go. I put "framework" in scare quotes because most of the
features can just be done optionally. The goal is to make it trivially easy to build a Go-based microservices
to server REST APIs at a small scale.

## Features

* Project Structure Helpers (TODO)
* Zero-conf local development (Partial)
* Build and Deployment Helpers (TODO)
* JWT-based auth integration and helpers (Partial)
* Error response helpers (Done)
* Deployment Environment/Local support (Done)
* Reasonable logging (WIP)
* Graceful shutdown (TODO)
* OpenTelemetry integration? (TODO)
* OpenAPI integration? (TODO)
* Client Helpers (TODO)

## Opinions
Strong opinions weakly held :D

* Convention over configuration
* Local environment is the most important environment and it should work by default and be easy to use
    * With no configs, `go run main.go` should start a running service
        * Hosted envs should be the ones that are explicitly configured
        * Some initial common setup may be necessary (like running a shared database)
        * The service should initialize anything it needs to internally, including things like DB migrations etc
    * Auth should be optional or disabled
* Routes should be all registered in a single place per service and include the method, full path AND authentication type (preferably in a single line of code)
* slog is fine; we just use it
* JWT is fine (come at me bro)
    * Scopes are too complicated and brittle; just create and check issuers (iss) to segment user groups. More granular checks should be done in business logic in the services. 
* Error responses should be `{"code":500 "message":"Internal error"}` where code mirrors the HTTP response.
* REST endpoints should return a single type. That type is guaranteed to be valid with a 200 response. Errors must respond with non-2xx.
* Version bumps should be restricted to only backwards incompatible changes. This should happen very infrequently or you're doing something wrong.
* The whole API is versioned for any change (ala Facebook).
* OpenAPI code generation sucks, don't ever use it. OpenAPI documentation is pretty nice.
* Each service runs on it's own port for easy local development
* Gin works a bit better than the stdlib
* Route hanlders should be named methods on a struct that contains the dependencies
    * Closures capturing the dependencies into scope is a mess and hard to read/debug
    * Dependencies are initialized at boot and injected into individual structs that have handler methods
    * Each struct can have multiple handlers for domain-based organization

## Motivation
I've had to build out something similar to this at three different startups, and each time I had to start over
from scratch. If I open source this I hopefully won't need to do it again.
