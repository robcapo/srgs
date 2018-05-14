package main

import (
	"errors"
	"github.com/beevik/etree"
)

var (
	InvalidGrammar = errors.New("invalid grammar document")
)

func ParseXml(xml string) (Expansion, error) {
	doc := etree.NewDocument()

	if err := doc.ReadFromString(xml); err != nil {
		return nil, err
	}

	grammar := doc.SelectElement("grammar")

	if grammar == nil {
		return nil, InvalidGrammar
	}

	rules := Rules{}

	for _, rule := range grammar.SelectElements("rule") {
		id, exp, err := ProcessRule(rule)

		if err != nil {
			return nil, err
		}

		rules[id] = exp
	}

}

func ProcessRule(rule *etree.Element) (id string, root Expansion, err error) {
	id = rule.SelectAttrValue("id", "")

}
