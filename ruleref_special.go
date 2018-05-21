package srgs

import (
	"strings"
)

type Garbage struct {
	match string

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

func (g *Garbage) Copy(gr *Grammar) Expansion { return new(Garbage) }
func (g *Garbage) Scan(processor Processor)   {}
