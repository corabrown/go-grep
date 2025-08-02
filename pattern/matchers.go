package pattern

import (
	"strings"
)

type matcher interface {
	match(lines []byte, grep *Grep) (bool, string)
	isRepeated() bool
	isMatched() bool
	setMatched(bool)
}

type exactMatch struct {
	b       byte
	matched bool
}

func (v *exactMatch) match(line []byte, grep *Grep) (matchFound bool, matchedString string) {
	if len(line) == 0 {
		return
	}
	if v.b == line[0] {
		return true, string(v.b)
	}
	return
}

func (v *exactMatch) isRepeated() bool  { return false }
func (v *exactMatch) isMatched() bool   { return v.matched }
func (v *exactMatch) setMatched(m bool) { v.matched = m }

type matchDigits struct{ matched bool }

func (v matchDigits) match(line []byte, grep *Grep) (bool, string) {
	if len(line) == 0 {
		return false, ""
	}
	return (line[0] >= '0') && (line[0] <= '9'), string(line[0])
}
func (v *matchDigits) isRepeated() bool  { return false }
func (v *matchDigits) isMatched() bool   { return v.matched }
func (v *matchDigits) setMatched(m bool) { v.matched = m }

type matchAlphanumeric struct{ matched bool }

func (v matchAlphanumeric) match(line []byte, grep *Grep) (bool, string) {
	if len(line) == 0 {
		return false, ""
	}
	b := line[0]
	return (b >= '0' && b <= '9') || (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b == '_'), string(b)
}
func (v *matchAlphanumeric) isRepeated() bool  { return false }
func (v *matchAlphanumeric) isMatched() bool   { return v.matched }
func (v *matchAlphanumeric) setMatched(m bool) { v.matched = m }

type matchCharacterGroup struct {
	s        string
	negative bool
	matched  bool
}

func (v matchCharacterGroup) match(line []byte, grep *Grep) (bool, string) {
	if len(line) == 0 {
		return false, ""
	}

	return strings.Contains(v.s, string(line[0])) == !v.negative, string(line[0])
}
func (v *matchCharacterGroup) isRepeated() bool  { return false }
func (v *matchCharacterGroup) isMatched() bool   { return v.matched }
func (v *matchCharacterGroup) setMatched(m bool) { v.matched = m }

type matchOneOrMore struct {
	matcher        matcher
	nMatches       int
	matchedPattern string
}

func (v *matchOneOrMore) match(line []byte, grep *Grep) (bool, string) {
	if len(line) == 0 {
		return false, ""
	}
	if ok, res := v.matcher.match(line, grep); ok {
		v.nMatches += 1
		v.matchedPattern = v.matchedPattern + res
		return true, res
	}
	return false, ""
}
func (v *matchOneOrMore) isRepeated() bool { return true }
func (v *matchOneOrMore) isMatched() bool  { return v.nMatches > 0 }
func (v *matchOneOrMore) setMatched(m bool) {
	if !m {
		v.nMatches = 0
	}

}

type matchZeroOrOne struct {
	matcher        matcher
	nMatches       int
	matchedPattern string
}

func (v *matchZeroOrOne) match(line []byte, grep *Grep) (bool, string) {
	if (len(line) == 0) || (v.nMatches > 0) {
		return false, ""
	}
	if ok, res := v.matcher.match(line, grep); ok {
		v.nMatches += 1
		v.matchedPattern = v.matchedPattern + res
		return true, v.matchedPattern
	} else if v.nMatches == 0 {
		return true, ""
	}
	return false, ""
}
func (v *matchZeroOrOne) isRepeated() bool  { return true }
func (v *matchZeroOrOne) isMatched() bool   { return true }
func (v *matchZeroOrOne) setMatched(m bool) {}

type wildcard struct{}

func (v *wildcard) match(line []byte, grep *Grep) (bool, string) {
	if len(line) == 0 {
		return false, ""
	}
	return true, string(line[0])
}

func (v *wildcard) isRepeated() bool  { return false }
func (v *wildcard) isMatched() bool   { return true }
func (v *wildcard) setMatched(m bool) {}

type matchAlternatingGroup struct {
	groupNumber        int
	matched            bool
	subPatterns        []Pattern
	matchedString      string
	currentStringStack string
	repeatingCheck     bool
}

func (v *matchAlternatingGroup) match(line []byte, grep *Grep) (bool, string) {
	defer func() {
		v.repeatingCheck = false
	}()

	for _, subPat := range v.subPatterns {
		subPat.beginningAnchor = true
		if v.repeatingCheck {
			subPat.matchers = subPat.matchers[len(subPat.matchers)-1:]
		}
		if ok, m := grep.patternMatch(line, subPat); ok {
			v.matched = true
			v.matchedString = v.matchedString + m
			grep.capturedGroupMatches[v.groupNumber] = grep.capturedGroupMatches[v.groupNumber] + m
			return true, m
		}
	}
	return false, ""
}

func (v *matchAlternatingGroup) isRepeated() bool {
	for _, p := range v.subPatterns {
		if len(p.matchers) > 0 {
			if p.matchers[len(p.matchers)-1].isRepeated() {
				v.repeatingCheck = true
				return true
			}
		}
	}
	return false
}

func (v *matchAlternatingGroup) isMatched() bool { return v.matched }
func (v *matchAlternatingGroup) setMatched(m bool) {
	v.matched = m
	if !m {
		v.matchedString = ""
	}
}

type backreference struct {
	groupNumber int
	matched     bool
}

func (v *backreference) match(line []byte, grep *Grep) (bool, string) {
	g := NewGrep("^" + grep.capturedGroupMatches[v.groupNumber])
	if ok, m := g.Match(line); ok {
		v.matched = true
		return true, m
	}
	return false, ""
}

func (v *backreference) isRepeated() bool  { return false }
func (v *backreference) isMatched() bool   { return v.matched }
func (v *backreference) setMatched(m bool) { v.matched = m }
