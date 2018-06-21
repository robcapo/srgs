package srgs

// Sequence is any sequence of legal expansions (see https://www.w3.org/TR/speech-grammar/#S2.3)
type Sequence struct {
	exps []Expansion

	str  string
	mode MatchMode

	nextInd int
}

// Implements Expansion Copy method
func (s *Sequence) Copy(r RuleRefs) Expansion {
	out := &Sequence{
		exps:    make([]Expansion, len(s.exps)),
		str:     s.str,
		mode:    s.mode,
		nextInd: s.nextInd,
	}

	for ind, e := range s.exps {
		out.exps[ind] = e.Copy(r)
	}

	return out
}

// Implements Expansion Match method
func (s *Sequence) Match(str string, mode MatchMode) {
	s.str = str
	s.mode = mode

	s.nextInd = 0

	s.exps[0].Match(str, mode)
}

// Implements Expansion Next method
func (s *Sequence) Next() (string, error) {
	if s.nextInd < 0 {
		return "", NoMatch
	}

	var str string
	var err error

	for i := s.nextInd; i < len(s.exps); i++ {
		str, err = s.exps[i].Next()

		if err != nil {
			s.nextInd--
			return s.Next()
		}

		if i+1 < len(s.exps) {
			s.nextInd = i + 1
			s.exps[s.nextInd].Match(str, s.mode)
		}
	}

	return str, err
}

// Implements Expansion Scan method
func (s *Sequence) Scan(p Processor) {
	for _, exp := range s.exps {
		exp.Scan(p)
	}
}

func (s *Sequence) ScanIDAndMatch(scorer Scorer) {
	for _, exp := range s.exps {
		exp.ScanIDAndMatch(scorer)
	}
}
