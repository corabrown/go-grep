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
			"bac is",
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
		{
			"single-backreference3",
			"once a dreaaamer, alwayszzz a dreaaamer",
			"once a (drea+mer), alwaysz? a \\1",
			false,
		},
		{
			"single-backreference4",
			"that starts and ends with this",
			"^(this) starts and ends with \\1$",
			false,
		},
		{
			"single-backreference5",
			"bugz here and bugs there",
			"(b..s|c..e) here and \\1 there",
			false,
		},
		{
			"multiple-backreferences",
			"3 red squares and 3 red circles",
			"(\\d+) (\\w+) squares and \\1 \\2 circles",
			true,
		},
		{
			"multiple-backreferences1",
			"apple pie, apple and pie",
			"^(\\w+) (\\w+), \\1 and \\2$",
			true,
		},
		{
			"nested-backreferences",
			"'cat and cat' is the same as 'cat and cat'",
			"('(cat) and \\2') is the same as \\1",
			true,
		},
		{
			"nested-backreferences1",
			"abc-def is abc-def, not efg, abc, or def",
			"(([abc]+)-([def]+)) is \\1, not ([^xyz]+), \\2, or \\3",
			true,
		},
		{
			"nested-backreferences2",
			"cat and fish, cat with fish, cat and fish",
			"((c.t|d.g) and (f..h|b..d)), \\2 with \\3, \\1",
			true,
		},
		{
			"simplify-failing-nested-backreference-test",
			"abc-def is",
			"(([abc]+)-([def]+)) is",
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
