package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go-trace <output_file>")
		os.Exit(1)
	}

	outputFile := os.Args[1]

	previousOutput, err := readLines(outputFile)
	if err != nil {
		fmt.Printf("Error reading previous output: %v\n", err)
		os.Exit(1)
	}

	currentOutput := readFromStdin()

	// Use map for quick lookup
	currentOutputSet := make(map[string]bool)
	for _, domain := range currentOutput {
		currentOutputSet[domain] = true
	}

	var newOutput []string
	for _, domain := range previousOutput {
		if _, exists := currentOutputSet[domain]; exists {
			newOutput = append(newOutput, domain)
		} else {
			fmt.Printf("Domain no longer returned: %s\n", domain)
		}
	}

	err = writeLines(outputFile, newOutput)
	if err != nil {
		fmt.Printf("Error writing to output file: %v\n", err)
		os.Exit(1)
	}
}

func readLines(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func readFromStdin() []string {
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func writeLines(file string, lines []string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, line := range lines {
		_, err := f.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
