# wc Tool Implementation

This is my Go implementation of the classic Unix `wc` (word count) tool, built as part of the [codingchallenges.fyi](https://codingchallenges.fyi/challenges/challenge-wc) challenge.

## Features

- **`-c`**: Count bytes in a file
- **`-l`**: Count lines in a file
- **`-w`**: Count words in a file
- **`-m`**: Count characters in a file
- **Default**: Display lines, words, and bytes (equivalent to `-l -w -c`)

## Usage

```bash
# Build the tool
go build -o wc main.go

# Count bytes
❯ ./wc -c lorum.txt
445 lorum.txt

# Count lines
❯ ./wc -l lorum.txt
4 lorum.txt

# Count words
❯ ./wc -w lorum.txt
69 lorum.txt

# Count characters
❯ ./wc -m lorum.txt
445 lorum.txt

# Default output (lines, words, bytes)
❯ ./wc lorum.txt
4 69 445 lorum.txt
```

## Implementation Details

The tool uses Go's `flag` package for command-line argument parsing and `io.ReadAll()` to read the entire file into memory. Each counting function processes the content differently:

- **Bytes**: Uses `len(content)` to get raw byte count
- **Lines**: Counts newline characters using `strings.Count()`
- **Words**: Uses `strings.Fields()` to split on whitespace and count tokens
- **Characters**: Uses `len(string(content))` to count runes

## Testing

Run the test suite to verify the implementation:

```bash
go test -v
```

The test suite includes:

- Unit tests for each counting function
- Edge cases (empty content, whitespace, Unicode)
- Integration tests with temporary files
- Specific tests for the `lorum.txt` file
