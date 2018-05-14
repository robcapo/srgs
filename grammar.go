package main

import (
	"errors"
	"strings"
)

var NoMatch = errors.New("cannot consume string with token")
var PrefixOnly = errors.New("string matched is a prefix")

// <rule id="foo">
// 	hello
//  <item repeat="3-">foo</item>
//  <one-of>
// 		<item>bar</item>
// 		<item weight="2">baz</item>
// 	</one-of>
// </rule>

type Path interface {
	AddString(str string)
	AddTag(tag Tag)
}

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

	AppendToProcessor(processor Processor)
}

type RuleRef struct {
	rule *Expansion
}

func (r RuleRef) Consume(str string) (string, Sequence, error) {
	return (*r.rule).Consume(str)
}
func (r RuleRef) AppendToProcessor(p Processor) {}

type Alternative []Item

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

func (s Sequence) Consume(str string) (string, Sequence, error) {
	out := make([]Expansion, 0)

	var seq []Expansion
	var err error
	for _, e := range s {
		str, seq, err = e.Consume(str)

		if err != nil {
			return str, out, err
		}

		out = append(out, seq...)
	}

	return str, out, nil
}
func (s Sequence) AppendToProcessor(p Processor) {
	for _, exp := range s {
		exp.AppendToProcessor(p)
	}
}

type Token string

func (t Token) Consume(str string) (string, Sequence, error) {
	if strings.HasPrefix(string(t), str) && len(string(t)) > len(str) {
		return "", nil, PrefixOnly
	}

	if strings.HasPrefix(str, string(t)) {
		str = str[len(t):]

		if len(str) == 0 {
			return "", []Expansion{t}, nil
		}

		if str[0] != ' ' {
			return "", nil, NoMatch
		}

		return str[1:], []Expansion{t}, nil
	}

	return "", nil, NoMatch
}
func (t Token) AppendToProcessor(p Processor) { p.AppendString(string(t)) }

type Tag string

func (t Tag) Consume(str string) (string, Sequence, error) { return str, []Expansion{t}, nil }
func (t Tag) AppendToProcessor(p Processor)                { p.AppendTag(string(t)) }
