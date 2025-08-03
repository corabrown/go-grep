package grep

// MatchStringPattern takes in the raw string pattern, parses it, and then matches it against the
// input bytes.
func MatchStringPattern(line []byte, pattern string) (found bool, matchedString string) {
	pat := Parse(pattern, nil)
	return MatchParsedPattern(line, pat)
}

// MatchParsedPattern takes in a parsed pattern and matches it against the input bytes.
// It iterates through the input line to be matched and the pattern itself, reseting the pattern
// index any time a match is not found.
func MatchParsedPattern(l []byte, p Pattern) (bool, string) {

	// create matching for this line and pattern 
	m := &matching{line: l, pattern: &p}

	// iterate through line and pattern until we've reached the end of either
	for m.lineOrPatternLeft() {

		// if we're at the end of the pattern and haven't found a match, return false
		if m.endOfPatternReached() {
			return false, ""
		}

		// check for a match at the current pattern index 
		if ok, matchedLine := m.checkForMatching(m.pix); ok {
			// if a match was found, increment the matching and the index 
			m.increment(matchedLine, m.pix)
			m.pix += 1
		} else {
			// if the previous pattern index contains a repeatable pattern, we also check for a match against that 
			if (m.pix != 0) && m.pattern.matchers[m.pix-1].isRepeated() {
				if ok, matchedLine := m.checkForMatching(m.pix - 1); ok {
					// if a match is found, we increment the matching and continue. Do not update the pattern index
					// because we want to continue to check for matches against m.pix  
					m.increment(matchedLine, m.pix-1)
					continue
				}
			}
			// if we didn't find a match and we have a beginning anchor, return false 
			if p.beginningAnchor {
				return false, ""
			}
			m.reset()
		}
	}

	// after iteration, check to make sure all pattern matchers have been matched 
	allPatternMatched := true
	for _, p := range m.pattern.matchers {
		if !p.isMatched() {
			allPatternMatched = false
		}
	}

	// if we have an end anchor and the last character matched index isn't the last one, return false 
	if p.endAnchor && m.lastMatchedCharacterIndex != len(m.line)-1 {
		return false, ""
	}

	// return whether each pattern matcher is matched, along with the matched string 
	return allPatternMatched, m.matchedString
}

// matching contains the bulk of the logic for the iteration. It keeps track both of the index of the input line 
// and the index of the pattern we're currently at. It tracks the last matched index from the input, in addition to 
// the full matched string so far.
type matching struct {
	line                      []byte   // the line to be matched
	pattern                   *Pattern // contains the full pattern
	lix                       int      // current index of the line being matched
	pix                       int      // current index of pattern
	lastMatchedCharacterIndex int      // index of the last character which was a match
	matchedString             string   // full string that has been matched so far
}

// increment moves the matching on to the next indices when we find a match for the given matched line and pattern index 
func (m *matching) increment(matchedLine string, patternIndex int) {
	m.pattern.matchers[patternIndex].setMatched(true)
	m.lix += len(matchedLine)
	m.lastMatchedCharacterIndex = m.lix - 1
	m.matchedString = m.matchedString + matchedLine
}

// if we don't find a match at a given set of indices, we want to reset the pattern but continue to iterate through 
// the input line to check for matches later on in the string 
func (m *matching) reset() {
	for _, matcher := range m.pattern.matchers {
		matcher.setMatched(false)
	}
	m.pattern.capturedGroupMatches.c = make([]string, len(m.pattern.capturedGroupMatches.c))
	m.pix = 0
	m.lix += 1
	m.matchedString = ""
}

// lineOrPatternLeft determines whether we've made it all the way through the line and/or pattern
// if we've finished either we should go on to check whether the entire pattern was matched 
func (m *matching) lineOrPatternLeft() bool {
	return (m.lix < len(m.line)) && (m.pix < len(m.pattern.matchers))
}

// checkForMatching runs the matcher at the pattern index against the line from this point forward 
func (m *matching) checkForMatching(patternIndex int) (bool, string) {
	return m.pattern.matchers[patternIndex].match(m.line[m.lix:], m.pattern)
}

// endOfPatternReached is true if the current pattern index is greater than the length of the pattern 
func (m *matching) endOfPatternReached() bool {
	return m.pix > len(m.pattern.matchers)-1
}
