# Architecture Decision Records (ADRs)

This directory contains Architecture Decision Records (ADRs) for the warehouse-server project.

## What are ADRs?

Architecture Decision Records are documents that capture important architectural decisions made along with their context and consequences. They help teams:

- Understand why decisions were made
- Track the evolution of the architecture
- Onboard new team members
- Avoid revisiting settled decisions

## Format

ADRs follow this template:
- **Status**: Proposed, Accepted, Deprecated, Superseded
- **Context**: The situation that led to the decision
- **Decision**: What was decided
- **Consequences**: The positive and negative outcomes

## Index

| ADR | Title | Status |
|-----|-------|--------|
| [0001](0001-modular-microservices-architecture.md) | Modular Microservices Architecture with Go Workspace | Accepted |

## Future ADRs

These ADRs would provide additional architectural documentation:

- **ADR-0002**: Repository Pattern for Data Access
- **ADR-0003**: API Key Authentication System  
- **ADR-0004**: Storage Abstraction Layer
- **ADR-0005**: Multi-tenant Branch-based Data Isolation
- **ADR-0006**: CORS and Security Middleware Strategy
- **ADR-0007**: Container Deployment Strategy
- **ADR-0008**: Database Technology Choices (SQLite/MySQL vs MongoDB)