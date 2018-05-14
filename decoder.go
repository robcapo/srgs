package main

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

type Parser struct {
	Root  Expansion
	rules Rules
	Xml   string
}

func NewParser() *Parser {
	return new(Parser)
}

func (p *Parser) LoadXml(xml string) error {
	p.Xml = xml

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

	p.rules = Rules{}

	for _, rule := range grammar.SelectElements("rule") {
		id, exp, err := p.decodeRule(rule)

		if err != nil {
			return err
		}

		if id == rootId {
			p.Root = exp
		}

		p.rules[id] = exp
	}

	if p.Root == nil {
		return RootNotFound
	}

	return nil
}

func (p *Parser) decodeRule(rule *etree.Element) (string, Expansion, error) {
	id := rule.SelectAttrValue("id", "")

	if id == "" {
		return "", nil, UnidentifiableRule
	}

	exp, err := p.decodeElement(rule)

	return id, exp, err
}

func (p *Parser) decodeElement(element *etree.Element) (Expansion, error) {
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

				if rule, ok := p.rules[ref[1:]]; ok {
					out = append(out, rule)
				}
			} else if el.Tag == "item" {
				exp, err := p.decodeElement(el)

				if err != nil {
					return nil, err
				}

				out = append(out, exp)
			} else if el.Tag == "one-of" {
				alt := Alternative{}
				for _, item := range el.SelectElements("item") {
					exp, err := p.decodeElement(item)

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
				return nil, errors.New("unable to parse tag " + el.Tag)
			}
		} else {
			fmt.Println(tok)
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
