package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"flag"
	"log"
	"os"
	"regexp"
)

type vaultItem map[string]string

func main() {
	txtFile := flag.String("txtFile", "", "Path to an exported plain text file from 1Password 6")
	flag.Parse()

	if *txtFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	inputFile, err := os.Open(*txtFile)
	if err != nil {
		log.Fatal(err)
	}

	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("Error: No output file specified.")
	}

	outputFile, err := os.Create(args[0])
	if err != nil {
		log.Fatal(err)
	}

	csvWriter := csv.NewWriter(outputFile)
	csvWriter.Write([]string{"title", "url", "username", "password", "notes", "tags"})

	item := vaultItem{}

	// `key=value`, potentially with new lines.
	re := regexp.MustCompile("(?s)^(.+)=(.+)")

	scanner := bufio.NewScanner(inputFile)
	scanner.Split(split)
	for scanner.Scan() {
		line := scanner.Text()

		// Blank line means we reached the end of an item, so we write it to output.
		if line == "" {
			// The order in which we write elements from the map matters. For non-existing
			// values, an empty string is written as a column value.
			if err := csvWriter.Write([]string{
				item["title"],
				item["website"],
				item["username"],
				item["password"],
				item["notesPlain"],
				item["tags"],
			}); err != nil {
				log.Fatal(err)
			}
			item = map[string]string{}
			// Nothing left to process for this line, so go to next one.
			continue
		}

		matches := re.FindStringSubmatch(scanner.Text())
		if len(matches) != 3 {
			log.Fatalf("Error: Invalid content in file: %v", line)
		}
		item[matches[1]] = matches[2]
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error: Could not read from input: %v", err)
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		log.Fatal(err)
	}
}

// Split function for scanner that essentially delimits by `\r\n`, because
// 1Password for Windows uses `\r\n` for line endings between key/param
// combinations, and `\n` for line breaks in values, e.g. notes with line breaks.
func split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\r'); i >= 0 && len(data) >= i+2 {
		// We have a full newline-terminated line.
		return i + 2, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}
