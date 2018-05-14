package srgs

import (
	"errors"
	"fmt"
	"github.com/beevik/etree"
	"strconv"
	"strings"
)

var (
	InvalidGrammar     = errors.New("invalid grammar document")
	NoRoot             = errors.New("no root present in grammar")
	RootNotFound       = errors.New("unable to find root rule")
	UnidentifiableRule = errors.New("rules must have an id")
	EmptyRuleRefUri    = errors.New("rulerefs must have a non-zero uri")
)

type Rules map[string]Expansion

type Grammar struct {
	Root Expansion
	Xml  string

	root     Expansion
	rules    Rules
	ruleRefs map[string][]*RuleRef
}

func NewGrammar() *Grammar {
	return new(Grammar)
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
				ref.rule = &exp
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

	g.Root = RuleRef{
		ruleId: rootId,
		rule:   &root,
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
	var out Sequence

	for _, tok := range element.Child {
		if data, ok := tok.(*etree.CharData); ok {
			str := strings.ToLower(strings.TrimSpace(data.Data))

			if len(str) == 0 {
				continue
			}

			out = append(out, decodeCharData(str))
		} else if el, ok := tok.(*etree.Element); ok {
			if el.Tag == "ruleref" {
				ref := el.SelectAttrValue("uri", "")

				if ref == "" {
					return nil, EmptyRuleRefUri
				}

				if ref[0] != '#' {
					return nil, errors.New("cannot understand ruleref uri " + ref + " because it is not local")
				}

				ruleRef := new(RuleRef)
				ruleRef.ruleId = ref[1:]

				out = append(out, ruleRef)

				if rule, ok := g.rules[ruleRef.ruleId]; ok {
					ruleRef.rule = &rule
				} else {
					g.ruleRefs[ruleRef.ruleId] = append(g.ruleRefs[ruleRef.ruleId], ruleRef)
				}
			} else if el.Tag == "item" {
				exp, err := g.decodeElement(el)

				if err != nil {
					return nil, err
				}

				out = append(out, exp)
			} else if el.Tag == "one-of" {
				alt := Alternative{}
				for _, item := range el.SelectElements("item") {
					exp, err := g.decodeElement(item)

					if err != nil {
						return nil, err
					}

					alt = append(alt, exp.(Item))
				}

				out = append(out, alt)
			} else if el.Tag == "tag" {
				out = append(out, Tag(el.Text()))
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

		return Item{out, min, max}, nil
	}

	return out, nil
}

func decodeCharData(data string) Expansion {
	if len(data) == 0 {
		return nil
	}

	return Token(data)
}
