package srgs

import (
	"errors"
)

type MatchMode int

const (
	ModePrefix MatchMode = iota
	ModeExact
)

var (
	NoMatch = errors.New("cannot consume string with token")
	PrefixOnly = errors.New("string matched is a prefix")
	Exhausted = errors.New("exhausted expansion")
)

// An expansion is any part of a grammar that can match a string
type Expansion interface {
	// Check if this expansion covers part or all of string `str`.
	// Returns a consumed version of str, with the part that this expansion covers removed from the beginning.
	// Uses a pre-allocated stack to push the Expansions involved in consuming the str
	//
	// May return two errors:
	//
	// - PrefixOnly if the entire text is only a prefix in the expansion.
	// - NoMatch if the string does not match the expansion at all
	//ConsumeStack(str string, stack *stack.Stack) (string, int, error)

	// Set the string to match on the expansion. Match either in ModePrefix (return nil error as soon as a prefix is
	// found) or ModeExact (return nil only if there is an exact match -- otherwise return PrefixOnly error)
	Match(str string, mode MatchMode)

	// If there are other ways of matching the prefix (e.g. if multiple alternatives or repeats match), return the
	// next version of the consumed string. Otherwise return "" and Exhausted error
	Next() (string, error)

	// Append this expansion to a processor. This will be enable the Processor to provide the output for a given path
	AppendToProcessor(processor Processor)
}

type RuleRef struct {
	rule   *Expansion
	ruleId string
}

func (r RuleRef) Match(str string, mode MatchMode) {
	(*r.rule).Match(str, mode)
}
func (r RuleRef) Next() (string, error) {
	return (*r.rule).Next()
}
//func (r RuleRef) ConsumeStack(str string, stack *stack.Stack) (string, int, error) {
//	PreProcessTag(r.ruleId).ConsumeStack(str, stack)
//	out, p, err := (*r.rule).ConsumeStack(str, stack)
//	PostProcessTag(r.ruleId).ConsumeStack(str, stack)
//	return out, p + 2, err
//}
func (r RuleRef) AppendToProcessor(p Processor) {}
