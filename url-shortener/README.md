# URL Shortener

A flexible URL shortening service with a generic database interface, built as part of the [codingchallenges.fyi](https://codingchallenges.fyi/challenges/challenge-url-shortener) challenge.

## Features

- **Generic Database Interface**: Plug-and-play database implementations
- **In-Memory Storage**: Fast, temporary storage using Go maps
- **Redis Storage**: Persistent storage with Redis backend
- **URL Deduplication**: Prevents duplicate URLs using SHA256 hashing
- **Random ID Generation**: 6-character alphanumeric short IDs
- **CRUD Operations**: Add, get, delete, and count entries

## Database Implementations

The service uses a generic `Database` interface that supports multiple storage backends:

- **In-Memory**: `db/in_mem`
- **Redis**: `db/redis`

## Usage

```bash
# Build the service
go build -o url-shortener main.go

# Run with in-memory database
./url-shortener

# Example output of `database.String()`
db:
EMWfO8 https://www.github.com
y0eyIh https://www.sapo.pt
```

## Implementation Details

The service uses a clean architecture with:

- **`Database` Interface**: Defines CRUD operations for URL entries
- **`Entry` Struct**: Represents a URL with its hash and short ID
- **Multiple Backends**: In-memory maps and Redis storage

## Testing

Run the comprehensive test suite:

```bash
go test ./db/...
```

The test suite includes:

- Generic tests that work with any Database implementation
- Implementation-specific tests for each backend

## Adding New Database Backends

1. Implement the `Database` interface
2. Add to `TestAllImplementations()` in `database_test.go`
3. Add implementation-specific tests
