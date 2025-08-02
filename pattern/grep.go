package pattern

type Grep struct {
	capturedGroupMatches []string
	pat                  Pattern
}

func NewGrep(pattern string) Grep {
	g := Grep{capturedGroupMatches: make([]string, 0)}
	g.pat = g.Parse(pattern)
	return g
}
