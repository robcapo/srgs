package srgs

import (
	"fmt"
	"strings"
)

type Garbage struct {
	str       string
	scanMatch bool

	currentInd int
}

func (g *Garbage) Match(str string, mode MatchMode) {
	g.currentInd = -1
	g.str = str
}

func (g *Garbage) Next() (string, error) {
	if g.currentInd == len(g.str) {
		return "", NoMatch
	}

	if g.currentInd != -1 {
		ind := strings.Index(g.str[g.currentInd:], " ")

		if ind == -1 {
			ind = len(g.str) - 1 - g.currentInd
		}

		g.currentInd += ind
	}

	g.currentInd++

	return g.str[g.currentInd:], nil
}

func (g *Garbage) Copy(r RuleRefs) Expansion {
	return &Garbage{str: g.str, currentInd: g.currentInd, scanMatch: g.scanMatch}
}
func (g *Garbage) Scan(processor Processor) {
	processor.AppendTag(fmt.Sprintf(`
scopes[scopes.length-1]['GARBAGE'] = "%s";
`, g.str[:g.currentInd]))
	processor.AppendString(g.str[:g.currentInd])
}

func (g *Garbage) ScanIDAndMatch(scorer Scorer) {}
