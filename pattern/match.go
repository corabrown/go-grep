package pattern

func (v *Grep) Match(line []byte) (found bool, matchedString string) {
	return v.patternMatch(line, v.pat)
}

func (v *Grep) patternMatch(line []byte, p Pattern) (found bool, matchedString string) {

	capturedGroups := v.capturedGroupMatches
	_ = capturedGroups

	pix := 0
	lix := 0
	lastMatchedCharacterIndex := 0

	pat := p.matchers

	for (lix < len(line)) && (pix < len(pat)) {
		l := string(line[lix])
		_ = l

		if pix > len(pat)-1 {
			return false, ""
		}

		if ok, matchedLine := pat[pix].match(line[lix:], v); ok {
			pat[pix].setMatched(true)
			lix += len(matchedLine)
			lastMatchedCharacterIndex = lix - 1
			matchedString = matchedString + matchedLine
			pix += 1
		} else {
			if (pix != 0) && pat[pix-1].isRepeated() {
				if ok, matchedLine := pat[pix-1].match(line[lix:], v); ok {
					pat[pix-1].setMatched(true)
					lix += len(matchedLine)
					lastMatchedCharacterIndex = lix - 1
					matchedString = matchedString + matchedLine
					continue
				}
			}
			if p.beginningAnchor {
				return false, ""
			}
			resetPattern(pat)
			pix = 0
			lix += 1
			matchedString = ""
		}
	}

	allPatternMatched := true
	for _, p := range pat {
		if !p.isMatched() {
			allPatternMatched = false
		}
	}

	if p.endAnchor && lastMatchedCharacterIndex != len(line)-1 {
		return false, ""
	}

	return allPatternMatched, matchedString
}

func resetPattern(pat []matcher) {
	for _, matcher := range pat {
		matcher.setMatched(false)
	}
}
