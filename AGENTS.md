# AGENTS.md - Agentic Coding Guidelines

This file provides guidelines for agents operating on the ride-sharing codebase.

## Project Overview

This is a Go-based microservices ride-sharing application with a Next.js web frontend. The backend includes:
- `services/api-gateway` - HTTP gateway service
- `services/trip-service` - Trip management service

The project uses Docker, Kubernetes (Tilt for local development), MongoDB, RabbitMQ, and Jaeger.

---

## Build, Lint, and Test Commands

### Go Backend

**Run all tests:**
```bash
go test ./...
```

**Run a single test:**
```bash
go test -run TestCreateTrip_Success ./services/trip-service/internal/service
```

**Build all services:**
```bash
go build ./...
```

**Run go vet (linting):**
```bash
go vet ./...
```

**Format code:**
```bash
gofmt -w .
```

**Generate protobuf:**
```bash
make generate-proto
```

### Web Frontend

```bash
cd web && npm install
cd web && npm run dev   # development server
cd web && npm run build # production build
cd web && npm run lint  # run linter
```

---

## Code Style Guidelines

### Go

**Imports:** Standard library first, then third-party, then project packages.
```go
import (
    "context"
    "fmt"
    "net/http"
    
    "ride-sharing/services/trip-service/internal/domain"
    "ride-sharing/shared/types"
    
    "go.mongodb.org/mongo-driver/bson/primitive"
)
```

**Formatting:** Use `gofmt` for automatic formatting, tabs for indentation, max 100 chars/line.

**Types:** Use meaningful type names (e.g., `TripModel`, `RideFareModel`). Use interfaces for repository abstractions. Use `context.Context` as first parameter for service methods. Define domain models in `internal/domain/`.

**Naming:**
- Packages: lowercase (e.g., `service`, `domain`)
- Interfaces: `Noun` + `er` suffix (e.g., `TripRepository`)
- Structs: PascalCase (e.g., `TripModel`)
- Variables: camelCase (e.g., `httpAddr`)
- Acronyms: Keep original casing (e.g., `URL`, `ID`)

**Error Handling:** Return errors directly, use `fmt.Errorf` with %v for wrapping:
```go
if err != nil {
    return nil, fmt.Errorf("failed to create trip: %v", err)
}
```

**Project Structure:**
```
services/api-gateway/main.go
services/trip-service/
  cmd/main.go
  internal/domain/      # Domain models and interfaces
  internal/service/     # Business logic
  internal/infrastructure/  # HTTP handlers, repositories
shared/types/           # Shared types
shared/env/             # Environment utilities
```

**Testing:** Test files: `*_test.go` in same package. Name tests: `Test<Method>_<Scenario>`. Use mock repositories.

---

### TypeScript / Next.js (Web)

**Imports:** React/Next first, then third-party, then relative imports.
```typescript
import { useState } from 'react'
import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { formatDate } from '@/lib/utils'
```

**Formatting:** ESLint (next/core-web-vitals), 2 spaces, single quotes, trailing commas.

**Types:** Explicit types for params/returns, interfaces for objects, avoid `any`.

**Naming:** Components: PascalCase (e.g., `TripCard.tsx`), Hooks: `useTrip`, Utils: camelCase.

**Components:** Functional components with TypeScript, use Tailwind classes with `cn()` utility.

---

## Environment Configuration

- Go: Use `shared/env` package for environment variables
- Web: Next.js environment variables (`.env.local`)

**Required tools:** Go 1.23+, Node.js (see `.nvmrc`), Docker, Tilt, protoc

---

## Common Development Workflows

**Running the full stack locally:**
```bash
tilt up
```

**Running a specific service:**
```bash
go run services/trip-service/cmd/main.go
```

**Checking dependencies:**
```bash
go mod tidy
```

---

## Development Workflow

This project follows a test-driven development approach. When fixing bugs or adding features:

1. **Understand existing code** - Read source files to understand current behavior
2. **Run existing tests first** - Always run `go test ./...` to see current test status
3. **Analyze failures** - If tests fail, understand WHY they fail before making changes
4. **Fix through understanding** - Fix the underlying issue, not just the test
5. **Iterate** - Run tests again after each change until all pass
6. **Make code testable** - If tests fail due to hardcoded values (e.g., URLs), make them configurable via environment variables
7. **Use mocks appropriately** - For unit tests, mock external dependencies rather than removing functionality

### Running Tests

```bash
# Run all unit tests
go test ./...

# Run integration tests (if available)
go test -tags=integration ./...

# Run both unit and integration tests
make test
```

### Commit Process

After all tests pass:
1. Run `go vet ./...` to check for issues
2. Run `gofmt -w .` to format code
3. Use conventional commit messages

---

## Notes for Agents

- Always run tests after changes: `go test ./...`
- Format code before committing: `gofmt -w .`
- Check for vet issues: `go vet ./...`
- Web uses Next.js 15 with React 19 and Tailwind CSS
- Do not commit secrets or credentials
