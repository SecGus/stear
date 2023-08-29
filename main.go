package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var silent bool

func main() {
	deleteFlag := flag.Bool("d", false, "Delete lines from output.txt if they are not in the current command output")
	flag.BoolVar(&silent, "s", false, "Run in silent mode (suppress output)")

	flag.Parse()

	if len(flag.Args()) != 1 {
		printError("Usage: go-trace [-d] [-s/--silent] <output_file>")
		os.Exit(1)
	}

	outputFile := flag.Args()[0]

	previousOutput, err := readLines(outputFile)
	if err != nil {
		printError(fmt.Sprintf("Error reading previous output: %v", err))
		os.Exit(1)
	}

	currentOutput := readFromStdin()

	currentOutputSet := make(map[string]bool)
	for _, domain := range currentOutput {
		currentOutputSet[domain] = true
	}

	var newOutput []string
	for _, domain := range previousOutput {
		if _, exists := currentOutputSet[domain]; exists {
			newOutput = append(newOutput, domain)
		} else {
			printError(fmt.Sprintf("%s", domain))
			if !*deleteFlag {
				newOutput = append(newOutput, domain)
			}
		}
	}

	err = writeLines(outputFile, newOutput)
	if err != nil {
		printError(fmt.Sprintf("Error writing to output file: %v", err))
		os.Exit(1)
	}
}

func printError(message string) {
	if !silent {
		fmt.Println(message)
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
