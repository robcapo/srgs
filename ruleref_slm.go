package srgs

import (
	"fmt"
	"strings"
)

type SLM struct {
	str       string
	scanMatch bool

	uri string
	lm  *kenlm

	currentInd int
}

func (s *SLM) Match(str string, mode MatchMode) {
	s.currentInd = -1
	s.str = str
}

func (s *SLM) CallLM(str string, uri string) float64 {
	if str != "" {
		return -5.5
	}
	return -1000
}

func (s *SLM) Next() (string, error) {
	if s.currentInd == len(s.str) {
		return "", NoMatch
	}

	if s.currentInd != -1 {
		ind := strings.Index(s.str[s.currentInd:], " ")

		if ind == -1 {
			ind = len(s.str) - 1 - s.currentInd
		}

		s.currentInd += ind

		matchProb := s.CallLM(s.str[s.currentInd:], s.uri)

		if matchProb <= -1000 {
			return "", NoMatch
		}
	}

	s.currentInd++

	return s.str[s.currentInd:], nil
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

func (s *SLM) ScanIDAndMatch(scorer Scorer) {
	scorer.AppendIDAndMatch("SLM-"+s.uri, s.str[:s.currentInd])
}
