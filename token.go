package srgs

type Token struct {
	token string

	str  string
	mode MatchMode

	called bool
}

func (t *Token) Copy(g *Grammar) Expansion {
	return &Token{token: t.token}
}

func NewToken(str string) *Token {
	return &Token{token: str}
}

func (t *Token) Match(str string, mode MatchMode) {
	t.str = str
	t.mode = mode
	t.called = false
}

func (t *Token) Next() (string, error) {
	if t.called {
		return "", NoMatch
	}

	t.called = true

	lent := len(t.token)
	lens := len(t.str)

	lim := lent
	if lens < lent {
		lim = lens
	}

	// Ensure all characters match up to their lengths
	for i := 0; i < lim; i++ {
		if t.str[i] != t.token[i] {
			return "", NoMatch
		}
	}

	// If they have the same length, perfect match
	if lent == lens {
		return "", nil
	}

	// If this token is longer than the query string, query string is a prefix
	if lent > lens {
		if t.mode == ModePrefix {
			return "", nil
		}

		return "", PrefixOnly
	}

	// If this token is shorter than the query string, but ends on a word boundary consume it and return the rest
	if t.str[lent] == ' ' {
		return t.str[lent+1:], nil
	}

	// Otherwise it's no match
	return "", NoMatch
}

func (t *Token) AppendToProcessor(p Processor) { p.AppendString(t.token) }

//func (t Token) ConsumeStack(str string, stack *stack.Stack) (string, int, error) {
//	if strings.HasPrefix(string(t), str) && len(string(t)) > len(str) {
//		return "", 0, PrefixOnly
//	}
//
//	if strings.HasPrefix(str, string(t)) {
//		str = str[len(t):]
//
//		if len(str) == 0 {
//			stack.Push(t)
//			return "", 1, nil
//		}
//
//		if str[0] == ' ' {
//			stack.Push(t)
//			return str[1:], 1, nil
//		}
//	}
//
//	return "", 0, NoMatch
//}
//func (t Token) Consume(str string) (string, Sequence, error) {
//	if strings.HasPrefix(string(t), str) && len(string(t)) > len(str) {
//		return "", nil, PrefixOnly
//	}
//
//	if strings.HasPrefix(str, string(t)) {
//		str = str[len(t):]
//
//		if len(str) == 0 {
//			return "", []Expansion{t}, nil
//		}
//
//		if str[0] == ' ' {
//			return str[1:], []Expansion{t}, nil
//		}
//	}
//
//	return "", nil, NoMatch
//}
