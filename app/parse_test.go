package main

import (
	"testing"
)

func TestMatching(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		pattern  string
		expected bool
	}{
		{
			"characters",
			"cat",
			"cat",
			true,
		},
		{
			"digit",
			"this is 9",
			"\\d",
			true,
		},
		{
			"alphanumeric",
			"fool101",
			"\\w",
			true,
		},
		{
			"non-match-alphanumeric",
			"$!?",
			"\\w",
			false,
		},
		{
			"positive-character-group",
			"cat",
			"[abc]",
			true,
		},
		{
			"non-match-positive-character-group",
			"fsdf",
			"[abc]",
			false,
		},
		{
			"negative-character-group",
			"e",
			"[^abc]",
			true,
		},
		{
			"non-match-negative-character-group",
			"b",
			"[^abc]",
			false,
		},
		{
			"combining-character-classes",
			"1 apple",
			"\\d apple",
			true,
		},
		{
			"combining-character-classes1",
			"101 apple",
			"\\d\\d\\d apple",
			true,
		},
		{
			"combining-character-classes2",
			"(bac is)",
			"([abcd]+) is",
			true,
		},
		{
			"non-matching-combining-character-classes",
			"1 apple",
			"\\d\\d\\d apple",
			false,
		},
		{
			"beginning-anchor",
			"log",
			"^log",
			true,
		},
		{
			"non-matching-beginning-anchor",
			"this is a log",
			"^log",
			false,
		},
		{
			"end-of-string-anchor",
			"dog",
			"dog$",
			true,
		},
		{
			"non-matching-end-of-string-anchor",
			"dogs",
			"dog$",
			false,
		},
		{
			"one-or-more",
			"caats",
			"ca+ts",
			true,
		},
		{
			"one-or-more-with-repeat",
			"caaats",
			"ca+at",
			true,
		},
		{
			"zero-or-one",
			"dogs",
			"dogs?",
			true,
		},
		{
			"zero-or-one1",
			"dog",
			"dogs?",
			true,
		},
		{
			"zero-or-one2",
			"dg",
			"do?g",
			true,
		},
		{
			"wildcard",
			"dog",
			"d.g",
			true,
		},
		{
			"alternating-group",
			"cat",
			"(cat|dog)",
			true,
		},
		{
			"alternating-group1",
			"a dog and cats",
			"a (cat|dog) and (cat|dog)s",
			true,
		},
		{
			"single-backreference",
			"cat and cat",
			"(cat) and \\1",
			true,
		},
		{
			"single-backreference1",
			"grep 101 is doing grep 101 times",
			"(\\w\\w\\w\\w \\d\\d\\d) is doing \\1 times",
			true,
		},
		{
			"single-backreference2",
			"abcd is abcd, not efg",
			"([abcd]+) is \\1, not [^xyz]+",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			capturedGroupMatches = make([]string, 0)
			result, _ := match([]byte(tt.input), tt.pattern)
			if result != tt.expected {
				t.Errorf("incorrect result for %v, %v", string(tt.input), tt.pattern)
			}
		})
	}
}
