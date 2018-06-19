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

type SLM struct {
	str       string
	scanMatch bool

	currentInd int
}

func (s *SLM) Match(str string, mode MatchMode) {
	s.currentInd = -1
	s.str = str
}

func (s *SLM) Next() (string, float64, error) {
	if s.currentInd == len(s.str) {
		return "", 0, NoMatch
	}

	if s.currentInd != -1 {
		ind := strings.Index(s.str[s.currentInd:], " ")

		if ind == -1 {
			ind = len(s.str) - 1 - s.currentInd
		}

		s.currentInd += ind

		//matchProb := callKLM(s.str[s.currentInd:])
		var matchProb float64

		if matchProb != 0 {
			s.currentInd++
			return s.str[s.currentInd:], matchProb, nil
		}
	}

	return "", 0, NoMatch
}

func (s *SLM) Copy(r RuleRefs) Expansion {
	return &SLM{str: s.str, currentInd: s.currentInd, scanMatch: s.scanMatch}
}

func (s *SLM) Scan(processor Processor) {
	processor.AppendTag(fmt.Sprintf(`
scopes[scopes.length-1]['SLM'] = "%s";
`, s.str[:s.currentInd]))
	processor.AppendString(s.str[:s.currentInd])
}

func (g *Garbage) Match(str string, mode MatchMode) {
	g.currentInd = -1
	g.str = str
}

func (g *Garbage) Next() (string, float64, error) {
	if g.currentInd == len(g.str) {
		return "", 0, NoMatch
	}

	if g.currentInd != -1 {
		ind := strings.Index(g.str[g.currentInd:], " ")

		if ind == -1 {
			ind = len(g.str) - 1 - g.currentInd
		}

		g.currentInd += ind
	}

	g.currentInd++

	return g.str[g.currentInd:], 1, nil
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
