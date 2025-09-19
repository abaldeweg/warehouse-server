# Documentation

This directory contains documentation for the warehouse-server project.

## Contents

- **[Architecture Decision Records (ADRs)](adr/)** - Important architectural decisions and their rationale

## Architecture Overview

The warehouse-server is a modular microservices system built with Go, consisting of:

- **Gateway Service** - Main API for warehouse management (books, reservations, inventory)
- **Static Service** - File server for uploads and assets
- **Logs Import Service** - Background log processing
- **Logs Web Service** - Web interface for log viewing  
- **Framework Module** - Shared utilities and patterns

See the [ADRs](adr/) for detailed architectural decisions and their context.

## Getting Started

Refer to the main [README.md](../README.md) for setup and development instructions.