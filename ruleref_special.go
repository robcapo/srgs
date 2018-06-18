package srgs

import (
	"fmt"
	"strings"
)

type Garbage struct {
	match     string
	scanMatch bool

	currentInd int
}

func (g *Garbage) Match(str string, mode MatchMode) {
	g.currentInd = -1
	g.match = str
}

func (g *Garbage) Next() (string, error) {
	if g.currentInd == len(g.match) {
		return "", NoMatch
	}

	if g.currentInd != -1 {
		ind := strings.Index(g.match[g.currentInd:], " ")

		if ind == -1 {
			ind = len(g.match) - 1 - g.currentInd
		}

		g.currentInd += ind
	}

	g.currentInd++

	return g.match[g.currentInd:], nil
}

func (g *Garbage) Copy(r RuleRefs) Expansion {
	return &Garbage{match: g.match, currentInd: g.currentInd, scanMatch: g.scanMatch}
}
func (g *Garbage) Scan(processor Processor) {
	processor.AppendTag(fmt.Sprintf(`
scopes[scopes.length-1]['GARBAGE'] = "%s";
`, g.match[:g.currentInd]))
}
