package main

import "fmt"

func match(line []byte, pattern string) (found bool, matchedString string) {

	defer func() {
		capturedGroups := capturedGroupMatches
		_ = capturedGroups

		fmt.Println(matchedString)
	}()

	var beginningAnchor bool
	if pattern[0] == '^' {
		beginningAnchor = true
		pattern = pattern[1:]
	}

	var endAnchor bool
	if pattern[len(pattern)-1] == '$' {
		endAnchor = true
		pattern = pattern[:len(pattern)-1]
	}

	pat := parse(pattern)
	pix := 0
	lix := 0
	lastMatchedCharacterIndex := 0

	capturedGroups := capturedGroupMatches
	_ = capturedGroups

	for (lix < len(line)) && (pix < len(pat)) {
		l := string(line[lix])
		_ = l

		if pix > len(pat)-1 {
			return false, ""
		}

		if ok, matchedLine := pat[pix].match(line[lix:]); ok {
			pat[pix].setMatched(true)
			lix += len(matchedLine)
			lastMatchedCharacterIndex = lix - 1
			matchedString = matchedString + matchedLine
			pix += 1
		} else {
			if (pix != 0) && pat[pix-1].isRepeated() {
				if ok, matchedLine := pat[pix-1].match(line[lix:]); ok {
					pat[pix-1].setMatched(true)
					lix += len(matchedLine)
					lastMatchedCharacterIndex = lix - 1
					matchedString = matchedString + matchedLine
					continue
				}
			}
			if beginningAnchor {
				return false, ""
			}
			resetPattern(pat)
			pix = 0
			lix += 1
			matchedString = ""
		}
		if pix == len(pat) && (!endAnchor || lix == len(line)) {
			return true, matchedString
		}
	}

	allPatternMatched := true
	for _, p := range pat {
		if !p.isMatched() {
			allPatternMatched = false
		}
	}

	if endAnchor && lastMatchedCharacterIndex != len(line)-1 {
		return false, ""
	}

	return allPatternMatched, matchedString
}

func resetPattern(pat []matcher) {
	for _, matcher := range pat {
		matcher.setMatched(false)
	}
}
