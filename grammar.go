package srgs

import (
	"errors"
	"fmt"
	"github.com/beevik/etree"
	"strconv"
	"strings"
)

type MatchMode int

const (
	ModePrefix MatchMode = iota
	ModeExact
)

var (
	NoMatch    = errors.New("cannot consume string with token")
	PrefixOnly = errors.New("string matched is a prefix")
)

var (
	InvalidGrammar     = errors.New("invalid grammar document")
	NoRoot             = errors.New("no root present in grammar")
	RootNotFound       = errors.New("unable to find root rule")
	UnidentifiableRule = errors.New("rules must have an id")
	EmptyRuleRefUri    = errors.New("rulerefs must have a non-empty uri")
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
	Scan(processor Processor)

	Copy(g *Grammar) Expansion
}

type Grammar struct {
	Root *RuleRef
	Xml  string

	root     Expansion
	rules    Rules
	ruleRefs map[string][]*RuleRef
}

func NewGrammar() *Grammar {
	return new(Grammar)
}

func (g *Grammar) HasPrefix(str string) bool {
	g.Root.Match(str, ModePrefix)
	str, err := g.Root.Next()

	for {
		if err != nil {
			return false
		}

		if len(str) == 0 {
			return true
		}

		str, err = g.Root.Next()
	}
}

func (g *Grammar) HasMatch(str string) bool {
	g.Root.Match(str, ModeExact)
	str, err := g.Root.Next()

	for {
		if err != nil {
			return false
		}

		if len(str) == 0 {
			return true
		}

		str, err = g.Root.Next()
	}
}

func (g *Grammar) GetMatch(str string, p Processor) error {
	g.Root.Match(str, ModeExact)
	str, err := g.Root.Next()

	for {
		if err != nil {
			return err
		}

		if len(str) == 0 {
			break
		}

		str, err = g.Root.Next()
	}

	p.AppendTag("var scopes = [{'rules':{}}];")
	g.Root.Scan(p)
	p.AppendTag(fmt.Sprintf("var root = scopes[0]['rules']['%s'];", g.Root.ruleId))

	return nil
}

func (g *Grammar) LoadXml(xml string) error {
	g.Xml = xml

	doc := etree.NewDocument()

	if err := doc.ReadFromString(xml); err != nil {
		return err
	}

	grammar := doc.SelectElement("grammar")

	if grammar == nil {
		return InvalidGrammar
	}

	rootId := grammar.SelectAttrValue("root", "")

	if rootId == "" {
		return NoRoot
	}

	var root Expansion

	g.rules = Rules{}

	// holds references to a given rule id so that they can be filled in once all rules have been processed
	g.ruleRefs = make(map[string][]*RuleRef)

	for _, rule := range grammar.SelectElements("rule") {
		id, exp, err := g.decodeRule(rule)

		if err != nil {
			return err
		}

		if id == rootId {
			root = exp
		}

		g.rules[id] = exp

		if refs, ok := g.ruleRefs[id]; ok {
			for _, ref := range refs {
				ref.rule = exp.Copy(g)
			}

			delete(g.ruleRefs, id)
		}
	}

	if root == nil {
		return RootNotFound
	}

	if len(g.ruleRefs) > 0 {
		refs := ""
		for ref := range g.ruleRefs {
			refs += ref + ", "
		}
		return errors.New("unresolved rule refs: " + strings.TrimSuffix(refs, ", "))
	}

	g.Root = &RuleRef{
		ruleId: rootId,
		rule:   root,
	}

	return nil
}

func (g *Grammar) decodeRule(rule *etree.Element) (string, Expansion, error) {
	id := rule.SelectAttrValue("id", "")

	if id == "" {
		return "", nil, UnidentifiableRule
	}

	exp, err := g.decodeElement(rule)

	return id, exp, err
}

func (g *Grammar) decodeElement(element *etree.Element) (Expansion, error) {
	out := new(Sequence)

	for _, tok := range element.Child {
		if data, ok := tok.(*etree.CharData); ok {
			str := strings.ToLower(strings.TrimSpace(data.Data))

			if len(str) == 0 {
				continue
			}

			out.exps = append(out.exps, decodeCharData(str))
		} else if el, ok := tok.(*etree.Element); ok {
			if el.Tag == "ruleref" {
				special := el.SelectAttrValue("special", "")

				if special == "GARBAGE" {
					out.exps = append(out.exps, new(Garbage))
					continue
				}

				ref := el.SelectAttrValue("uri", "")

				if ref == "" {
					return nil, EmptyRuleRefUri
				}

				if ref[0] != '#' {
					return nil, errors.New("cannot understand ruleref uri " + ref + " because it is not local")
				}

				ruleRef := new(RuleRef)
				ruleRef.ruleId = ref[1:]

				out.exps = append(out.exps, ruleRef)

				if rule, ok := g.rules[ruleRef.ruleId]; ok {
					ruleRef.rule = rule.Copy(g)
				} else {
					g.ruleRefs[ruleRef.ruleId] = append(g.ruleRefs[ruleRef.ruleId], ruleRef)
				}
			} else if el.Tag == "item" {
				exp, err := g.decodeElement(el)

				if err != nil {
					return nil, err
				}

				out.exps = append(out.exps, exp)
			} else if el.Tag == "one-of" {
				alt := new(Alternative)
				for _, item := range el.SelectElements("item") {
					exp, err := g.decodeElement(item)

					if err != nil {
						return nil, err
					}

					alt.items = append(alt.items, exp.(*Item))
				}

				out.exps = append(out.exps, alt)
			} else if el.Tag == "tag" {
				out.exps = append(out.exps, NewTag(el.Text()))
			} else if el.Tag == "example" {
				// ignore
			} else {
				return nil, errors.New("unable to parse tag " + el.Tag + " " + el.SelectAttrValue("id", "no id"))
			}
		} else {
			fmt.Println("Unable to process", tok, "Ignoring.")
		}

	}

	if element.Tag == "item" {
		repeat := element.SelectAttrValue("repeat", "1-1")

		minMax := strings.Split(repeat, "-")

		var min, max int
		var err error

		if len(minMax) == 1 {
			if min, err = strconv.Atoi(minMax[0]); err != nil {
				return nil, err
			}

			max = min
		} else if len(minMax) == 2 {
			if minMax[0] == "" {
				min = 0
			} else if min, err = strconv.Atoi(minMax[0]); err != nil {
				return nil, err
			}

			if max, err = strconv.Atoi(minMax[1]); err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("invalid repeat")
		}

		return NewItem(out, min, max), nil
	}

	return out, nil
}

func decodeCharData(data string) Expansion {
	if len(data) == 0 {
		return nil
	}

	return NewToken(data)
}
