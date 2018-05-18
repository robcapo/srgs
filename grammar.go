package srgs

import (
	"errors"
	"fmt"
	"strings"
	"github.com/golang-collections/collections/stack"
)

var NoMatch = errors.New("cannot consume string with token")
var PrefixOnly = errors.New("string matched is a prefix")

// An expansion is any part of a grammar that can
// match a string
type Expansion interface {
	// Check if this expansion covers part or all of string `str`.
	// Returns a consumed version of str, with the part that this expansion covers removed from the beginning.
	// Returns a slice of Expansions that were involved in covering the string
	//
	// May return two errors:
	//
	// - PrefixOnly if the entire text is only a prefix in the expansion.
	// - NoMatch if the string does not match the expansion at all
	Consume(str string) (string, Sequence, error)

	// Does the same thing as Consume, but uses a preallocated stack to push the path onto
	// this saves on allocating a Sequence for each Expansion in the path
	ConsumeStack(str string, stack *stack.Stack) (string, int, error)

	// Check if the expansion matches a certain string. Will return the consumed version of the string
	// as well as an error value of either nil (perfect match), PrefixOnly (str was a prefix), or
	// NoMatch (str did not match)
	Match(str string) (string, error)

	// Check if str is a prefix of the Expansion only. This is an optimized version of Match, with the caveat
	// that it does not distinguish between a prefix and a perfect match. Prefix search decoders can leverage
	// this method to very quickly check if a candidate utterance is a prefix of the grammar.
	MatchPrefix(str string) (string, error)

	// Append this expansion to a processor. This will be enable the Processor to provide the output for a
	// given path
	AppendToProcessor(processor Processor)
}

type RuleRef struct {
	rule   *Expansion
	ruleId string
}

func (r RuleRef) Match(str string) (string, error) {
	return (*r.rule).Match(str)
}
func (r RuleRef) MatchPrefix(str string) (string, error) {
	return (*r.rule).MatchPrefix(str)
}
func (r RuleRef) ConsumeStack(str string, stack *stack.Stack) (string, int, error) {
	PreProcessTag(r.ruleId).ConsumeStack(str, stack)
	out, p, err := (*r.rule).ConsumeStack(str, stack)
	PostProcessTag(r.ruleId).ConsumeStack(str, stack)
	return out, p + 2, err
}

func (r RuleRef) Consume(str string) (string, Sequence, error) {
	str, seq, err := (*r.rule).Consume(str)

	if err != nil {
		return str, seq, err
	}

	return str, Sequence{PreProcessTag(r.ruleId), seq, PostProcessTag(r.ruleId)}, nil
}
func (r RuleRef) AppendToProcessor(p Processor) {}

type Alternative []Item

func (a Alternative) Match(str string) (string, error) {
	outErr := NoMatch
	for _, alt := range a {
		str, err := alt.Match(str)

		if err == nil {
			return str, nil
		}

		if err == PrefixOnly {
			outErr = PrefixOnly
		}
	}
	return "", outErr
}
func (a Alternative) MatchPrefix(str string) (string, error) {
	for _, alt := range a {
		str, err := alt.MatchPrefix(str)

		if err == nil {
			return str, nil
		}
	}

	return "", NoMatch
}

func (a Alternative) ConsumeStack(str string, stack *stack.Stack) (string, int, error) {
	outErr := NoMatch
	for _, alt := range a {
		out, p, err := alt.ConsumeStack(str, stack)

		if err == nil {
			return out, p, err
		}

		if err == PrefixOnly {
			outErr = PrefixOnly
		}
	}

	return "", 0, outErr
}
func (a Alternative) Consume(str string) (string, Sequence, error) {
	outErr := NoMatch
	for _, alt := range a {
		out, seq, err := alt.Consume(str)

		if err == nil {
			return out, seq, err
		}

		if err == PrefixOnly {
			outErr = PrefixOnly
		}
	}

	return "", nil, outErr
}

func (a Alternative) AppendToProcessor(p Processor) {}

type Item struct {
	Sequence
	repeatMin int
	repeatMax int
}

func (i Item) Match(str string) (string, error) {
	return i.match(str, i.repeatMin, i.repeatMax)
}
func (i Item) match(str string, min, max int) (string, error) {
	if max == 0 {
		return str, nil
	}

	outStr, err := i.Sequence.Match(str)

	if err != nil {
		if min <= 0 {
			return str, nil
		}

		return "", err
	}

	return i.match(outStr, min - 1, max - 1)
}

func (i Item) MatchPrefix(str string) (string, error) { return i.matchPrefix(str, i.repeatMin, i.repeatMax) }
func (i Item) matchPrefix(str string, min, max int) (string, error) {
	if max == 0 {
		return str, nil
	}

	outStr, err := i.Sequence.MatchPrefix(str)

	if err != nil {
		if min <= 0 {
			return str, nil
		}

		return "", err
	}

	return i.matchPrefix(outStr, min - 1, max - 1)
}

func (i Item) ConsumeStack(str string, stack *stack.Stack) (string, int, error) {
	return i.consumeStack(str, stack, i.repeatMin, i.repeatMax)
}
func (i Item) consumeStack(str string, stack *stack.Stack, min, max int) (string, int, error) {
	if max == 0 {
		return str, 0, nil
	}

	outStr, p, err := i.Sequence.ConsumeStack(str, stack)

	if err != nil {
		if min <= 0 {
			return str, p, nil
		}

		return "", p, err
	}

	out2, p2, err2 := i.consumeStack(outStr, stack, min - 1, max - 1)

	return out2, p + p2, err2
}

func (i Item) Consume(str string) (string, Sequence, error) {
	return i.consume(str, i.repeatMin, i.repeatMax, nil)
}
func (i Item) consume(str string, min int, max int, seq Sequence) (string, Sequence, error) {
	if max == 0 {
		return str, seq, nil
	}

	outStr, seq1, err := i.Sequence.Consume(str)

	if err != nil {
		if min <= 0 {
			return str, seq, nil
		}

		return "", nil, err
	}

	seq = append(seq, seq1...)

	outStr, seq2, err := i.consume(outStr, min-1, max-1, seq)

	if err != nil {
		if min > 0 {
			return str, nil, err
		}

		return "", nil, err
	}

	return outStr, seq2, nil
}

type Sequence []Expansion

func (s Sequence) Match(str string) (string, error) {
	var err error
	for _, e := range s {
		str, err = e.Match(str)

		if err != nil {
			return str, err
		}
	}

	return str, nil
}
func (s Sequence) MatchPrefix(str string) (string, error) {
	var err error
	for _, e := range s {
		str, err = e.MatchPrefix(str)

		if err != nil {
			return str, err
		}
	}

	return str, nil
}

func (s Sequence) ConsumeStack(str string, stack *stack.Stack) (string, int, error) {
	var err error
	var pushes, p int
	for _, e := range s {
		str, p, err = e.ConsumeStack(str, stack)

		if err != nil {
			for i := pushes; i > 0; i-- {
				stack.Pop()
			}
			return str, 0, err
		}

		pushes += p
	}

	return str, pushes, nil
}

func (s Sequence) Consume(str string) (string, Sequence, error) {
	out := make([]Expansion, 0, len(s))

	var seq Expansion
	var err error
	for _, e := range s {
		str, seq, err = e.Consume(str)

		if err != nil {
			return str, out, err
		}

		out = append(out, seq)
	}

	return str, out, nil
}
func (s Sequence) AppendToProcessor(p Processor) {
	for _, exp := range s {
		exp.AppendToProcessor(p)
	}
}

type Token string

func (t Token) Match(str string) (string, error) {
	if strings.HasPrefix(string(t), str) && len(string(t)) > len(str) {
		return "", PrefixOnly
	}

	if strings.HasPrefix(str, string(t)) {
		str = str[len(t):]

		if len(str) == 0 {
			return "", nil
		}

		if str[0] == ' ' {
			return str[1:], nil
		}
	}

	return "", NoMatch
}

func (t Token) MatchPrefix(str string) (string, error) {
	if strings.HasPrefix(string(t), str) {
		return "", nil
	}

	if strings.HasPrefix(str, string(t)) {
		str = str[len(t):]

		if len(str) == 0 {
			return "", nil
		}

		if str[0] == ' ' {
			return str[1:], nil
		}
	}

	return "", NoMatch
}
func (t Token) ConsumeStack(str string, stack *stack.Stack) (string, int, error) {
	if strings.HasPrefix(string(t), str) && len(string(t)) > len(str) {
		return "", 0, PrefixOnly
	}

	if strings.HasPrefix(str, string(t)) {
		str = str[len(t):]

		if len(str) == 0 {
			stack.Push(t)
			return "", 1, nil
		}

		if str[0] == ' ' {
			stack.Push(t)
			return str[1:], 1, nil
		}
	}

	return "", 0, NoMatch
}
func (t Token) Consume(str string) (string, Sequence, error) {
	if strings.HasPrefix(string(t), str) && len(string(t)) > len(str) {
		return "", nil, PrefixOnly
	}

	if strings.HasPrefix(str, string(t)) {
		str = str[len(t):]

		if len(str) == 0 {
			return "", []Expansion{t}, nil
		}

		if str[0] == ' ' {
			return str[1:], []Expansion{t}, nil
		}
	}

	return "", nil, NoMatch
}
func (t Token) AppendToProcessor(p Processor) { p.AppendString(string(t)) }

type Tag string

func (t Tag) Match(str string) (string, error) { return str, nil }
func (t Tag) MatchPrefix(str string) (string, error) { return str, nil }
func (t Tag) ConsumeStack(str string, stack *stack.Stack) (string, int, error) {
	stack.Push(t)
	return str, 1, nil
}

func (t Tag) Consume(str string) (string, Sequence, error) { return str, []Expansion{t}, nil }
func (t Tag) AppendToProcessor(p Processor)                {
	p.AppendTag(string(t))
}


type PreProcessTag string
func (t PreProcessTag) MatchPrefix(str string) (string, error) { return str, nil }
func (t PreProcessTag) Match(str string) (string, error) { return str, nil }
func (t PreProcessTag) ConsumeStack(str string, stack *stack.Stack) (string, int, error) {
	stack.Push(t)
	return str, 1, nil
}

func (t PreProcessTag) Consume(str string) (string, Sequence, error) { return str, []Expansion{t}, nil }
func (t PreProcessTag) AppendToProcessor(p Processor) {
	p.AppendTag(fmt.Sprintf(`
rules.%s = {};

if (root == undefined) {
	root = rules.%s;
}

(function() {
	var out;
`, string(t), string(t)))
}

type PostProcessTag string
func (t PostProcessTag) Match(str string) (string, error) { return str, nil }
func (t PostProcessTag) MatchPrefix(str string) (string, error) { return str, nil }
func (t PostProcessTag) ConsumeStack(str string, stack *stack.Stack) (string, int, error) {
	stack.Push(t)
	return str, 1, nil
}
func (t PostProcessTag) Consume(str string) (string, Sequence, error) { return str, []Expansion{t}, nil }
func (t PostProcessTag) AppendToProcessor(p Processor) {
	p.AppendTag(fmt.Sprintf(`
	rules.%s.out = out;
})();
`, string(t)))
}
