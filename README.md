# rapid-go

## motivation

rapid-go is a boilerplate that accelerates API development based on layered architecture and clarifying responsibilities.

## what is this

```
the boilerplate for monorepo application (support only http protocol)
```

- Base project is https://github.com/golang-standards/project-layout

## Apps

| Package                                     | Localhost             | Production |
| :------------------------------------------ | :-------------------- | :-------- |
| **[[REST] api server](./cmd/app/http_server_cmd.go)** | http://localhost:8080 | api.\*    |

## Documentation

### Getting Started

- **[Development Setup](./docs/development-setup/README.md)** - Environment setup, running the application, and common development tasks
- **[Project Overview](./.claude/CLAUDE.md)** - Architecture, tech stack, directory structure, and coding guidelines

### CLI Tools

- **[create-root-admin CLI](./docs/tools/create-root-admin-cli/README.md)** - Create initial root administrator accounts
- **[init-new-repository](./docs/tools/init-new-repository/README.md)** - Initialize a new repository from rapid-go template

### Specifications

- **[Specifications](./docs/specifications/)** - System specification documents

### Infrastructure

- **[Infra Architectures](./docs/infra-architectures/)** - AWS/GCP infrastructure architecture documents

### Development Guidelines

Detailed coding rules and patterns are organized in `.claude/rules/`:

- **[Domain Model Guidelines](./.claude/rules/domain-model.md)** - Entity patterns, constructors, state transitions
- **[Repository Guidelines](./.claude/rules/repository.md)** - Data access layer patterns
- **[Usecase Interactor Guidelines](./.claude/rules/usecase-interactor.md)** - Business logic layer patterns
- **[gRPC Handler Guidelines](./.claude/rules/grpc-handler.md)** - API handler patterns
- **[Testing Guidelines](./.claude/rules/testing.md)** - Unit testing conventions
- **[Proto Definition Guidelines](./.claude/rules/proto-definition.md)** - Protocol Buffers style guide
- **[Migration Guidelines](./.claude/rules/migration.md)** - Database migration patterns
- **[CLI Command Pattern](./.claude/rules/cli-command-pattern.md)** - CLI implementation guidelines
