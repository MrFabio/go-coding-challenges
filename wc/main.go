package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		fmt.Println("Error closing file: ", file.Name(), err)
		os.Exit(1)
	}
}

func countBytes(content []byte) int {
	return len(content)
}

func countLines(content []byte) int {
	return strings.Count(string(content), "\n")
}

func countWords(content []byte) int {
	return len(strings.Fields(string(content)))
}

func countCharacters(content []byte) int {
	return len(string(content))
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	// process arguments
	// -c count bytes
	// -l count lines
	// -w count words
	// -m count characters

	bytesFlag := flag.Bool("c", false, "count bytes")
	linesFlag := flag.Bool("l", false, "count lines")
	wordsFlag := flag.Bool("w", false, "count words")
	charsFlag := flag.Bool("m", false, "count characters")

	flag.Parse()

	// parse file from the last argument
	filename := flag.Args()[len(flag.Args())-1]
	filename = filepath.Clean(filename)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return err
	}
	defer closeFile(file)

	// read file
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file: ", err)

		return err
	}

	// process flags
	switch {
	case *bytesFlag:
		fmt.Println(countBytes(content), filename)
	case *linesFlag:
		fmt.Println(countLines(content), filename)
	case *wordsFlag:
		fmt.Println(countWords(content), filename)
	case *charsFlag:
		fmt.Println(countCharacters(content), filename)
	default:
		fmt.Println(countBytes(content), countLines(content), countWords(content), filename)
	}

	return nil
}
