package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/codecrafters-io/grep-starter-go/pattern"
)

func MatchFile(filename string, grep pattern.Grep) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file: %v", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Bytes()

		if match, _ := grep.Match(line); match {
			fmt.Fprint(os.Stdout, string(line), "\n")
		}
	}

}
