# ADR 001: Database Choice

## Status
Accepted

## Context
We need to select a database for the application. The "Database Setup" feature is tracked in FEATURES.yml and blocks the "Login" feature.

## Decision
We will use **SQLite** for development and **PostgreSQL** for production.

## Rationale
- **SQLite**: Zero-config, file-based, perfect for local development and testing
- **PostgreSQL**: Production-grade, ACID compliant, excellent JSON support, strong ecosystem
- Both support the same SQL dialect (mostly), enabling easy migration
- Prisma ORM supports both seamlessly

## Alternatives Considered
- **MySQL**: Similar to PostgreSQL but less advanced JSON support
- **MongoDB**: NoSQL, but adds complexity for relational data
- **Pure PostgreSQL**: Would require local Postgres setup for all developers

## Consequences
- Development uses SQLite (file: `./dev.db`)
- Production uses PostgreSQL (managed service recommended)
- Schema migrations managed via Prisma
- Need to handle minor SQL dialect differences in migrations