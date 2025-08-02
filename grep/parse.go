package grep

// var CapturedGroupMatches []string

type Pattern struct {
	matchers             []matcher
	beginningAnchor      bool
	endAnchor            bool
	capturedGroupMatches *capturedGroup
}

func Parse(pattern string, capturedGroupMatches *capturedGroup) Pattern {
	if capturedGroupMatches == nil {
		v := newCapturedGroup()
		capturedGroupMatches = v
	}

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

	matchers := make([]matcher, 0, len(pattern))

	var escaped bool
	var characterGroupMatcher *matchCharacterGroup
	var alternatingGroupMatcher *matchAlternatingGroup
	var skipNext bool
	alternatingGroupCount := 0

	for i := range pattern {

		if skipNext {
			skipNext = false
			continue
		}

		if pattern[i] == '(' {
			alternatingGroupCount += 1
		}

		if pattern[i] == ')' {
			alternatingGroupCount -= 1
			if alternatingGroupCount == 0 {
				alternatingGroupMatcher.subPatterns = append(alternatingGroupMatcher.subPatterns, Parse(alternatingGroupMatcher.currentStringStack, capturedGroupMatches))
				matchers = append(matchers, alternatingGroupMatcher)
				alternatingGroupMatcher = nil
				continue
			}
		}

		if alternatingGroupMatcher != nil {
			if (pattern[i] == '|') && (alternatingGroupCount == 1) {
				alternatingGroupMatcher.subPatterns = append(alternatingGroupMatcher.subPatterns, Parse(alternatingGroupMatcher.currentStringStack, capturedGroupMatches))
				alternatingGroupMatcher.currentStringStack = ""
				continue
			}
			alternatingGroupMatcher.currentStringStack = alternatingGroupMatcher.currentStringStack + string(pattern[i])
			continue
		}

		if pattern[i] == '[' {
			characterGroupMatcher = &matchCharacterGroup{}
			continue
		}
		if pattern[i] == ']' {
			matchers = append(matchers, characterGroupMatcher)
			characterGroupMatcher = nil
			continue
		}

		if pattern[i] == '(' {
			if alternatingGroupMatcher == nil {
				alternatingGroupMatcher = &matchAlternatingGroup{groupNumber: capturedGroupMatches.len(), subPatterns: make([]Pattern, 0)}
				capturedGroupMatches.add("")
			}
			continue
		}

		if characterGroupMatcher != nil {
			if (len(characterGroupMatcher.s) == 0) && (pattern[i] == '^') {
				characterGroupMatcher.negative = true
			}

			characterGroupMatcher.s = characterGroupMatcher.s + string(pattern[i])
			continue
		}

		if pattern[i] == '\\' {
			escaped = true
			continue
		}

		if escaped {
			if pattern[i] == 'd' {
				matchers = append(matchers, &matchDigits{})
				escaped = false
				continue
			}
			if pattern[i] == 'w' {
				matchers = append(matchers, &matchAlphanumeric{})
				escaped = false
				continue
			}
			if pattern[i] == '1' {
				matchers = append(matchers, &backreference{groupNumber: 0})
				escaped = false
				continue
			}
			if pattern[i] == '2' {
				matchers = append(matchers, &backreference{groupNumber: 1})
				escaped = false
				continue
			}
			if pattern[i] == '3' {
				matchers = append(matchers, &backreference{groupNumber: 2})
				escaped = false
				continue
			}
		}

		if (pattern[i] == '+') && (len(matchers) > 0) {
			prevMatcher := matchers[len(matchers)-1]
			matchers[len(matchers)-1] = &matchOneOrMore{matcher: prevMatcher}
			if (len(pattern) > i+1) && (i > 1) {
				if pattern[i+1] == pattern[i-1] {
					skipNext = true
				}
			}
			continue
		}

		if (pattern[i] == '?') && (len(matchers) > 0) {
			prevMatcher := matchers[len(matchers)-1]
			matchers[len(matchers)-1] = &matchZeroOrOne{matcher: prevMatcher}
			continue
		}

		if pattern[i] == '.' {
			matchers = append(matchers, &wildcard{})
			continue
		}

		matchers = append(matchers, &exactMatch{b: pattern[i]})
	}

	return Pattern{matchers, beginningAnchor, endAnchor, capturedGroupMatches}
}

type capturedGroup struct {
	c []string
}

func newCapturedGroup() *capturedGroup {
	return &capturedGroup{c: make([]string, 0)}
}

func (v *capturedGroup) add(s string) {
	if v.c == nil {
		return
	}
	v.c = append(v.c, s)
}

func (v *capturedGroup) len() int {
	return len(v.c)
}

func (v *capturedGroup) modifyVal(i int, s string) {
	v.c[i] = s
}

func (v *capturedGroup) getVal(i int) string {
	return v.c[i]
}
