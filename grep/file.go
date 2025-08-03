package grep

import (
	"bufio"
	"fmt"
	"os"
)

// MatchFile takes in a file name and matches a copy of the already parsed pattern against each 
// line of that file. No linebreaks here, we assume any pattern will be fully contained in a single line
func MatchFile(filename string, pattern Pattern) {
	p := pattern // create a copy of the pattern because certain info is saved during matching

	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file: %v", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Bytes()

		if match, _ := MatchParsedPattern(line, p); match {
			fmt.Fprint(os.Stdout, string(line), "\n")
		}
	}

}
