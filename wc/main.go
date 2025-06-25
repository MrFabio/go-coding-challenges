package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

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
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		os.Exit(1)
	}
	defer file.Close()

	// read file
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file: ", err)
		os.Exit(1)
	}

	// process flags
	if *bytesFlag {
		fmt.Println(countBytes(content), filename)
	} else if *linesFlag {
		fmt.Println(countLines(content), filename)
	} else if *wordsFlag {
		fmt.Println(countWords(content), filename)
	} else if *charsFlag {
		fmt.Println(countCharacters(content), filename)
	} else {
		fmt.Println(countBytes(content), countLines(content), countWords(content), filename)
	}
}
