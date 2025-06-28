# Go Coding Challenges

This repository contains my solutions to real-world coding challenges from [codingchallenges.fyi](https://codingchallenges.fyi/), implemented in Go.

## 🚀 Completed Challenges

### 1. [wc Tool](/wc/)

> Recreate the classic Unix `wc` tool to count lines, words, characters, and bytes in files.

### 2. [url shortener](/url-shortener/)

> Build my own URL shortening service with URL deduplication and multiple backends.

## 🔧 CI/CD

This repository uses GitHub Actions for continuous integration. The workflow automatically:

- Runs tests and builds for all Go projects
- Runs linting with golangci-lint

### Adding New Projects

To add a new Go project to the CI pipeline, add `your-new-project` name to the matrix in `.github/workflows/go.yml`:

```yaml
strategy:
  matrix:
    project: [url-shortener, wc, your-new-project]
```

1. Ensure your project has:
   - A valid `go.mod` file
   - Tests (run with `go test ./...`)
   - A buildable main package (run with `go build ./...`)

All projects run concurrently for faster CI execution.
