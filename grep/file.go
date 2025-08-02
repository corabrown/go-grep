package grep

import (
	"bufio"
	"fmt"
	"os"
)

func MatchFile(filename string, pattern Pattern) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file: %v", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Bytes()

		if match, _ := PatternMatch(line, pattern); match {
			fmt.Fprint(os.Stdout, string(line), "\n")
		}
	}

}
