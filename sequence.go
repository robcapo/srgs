package srgs

type Sequence struct {
	exps []Expansion

	str  string
	mode MatchMode

	nextInd int
}

func (s *Sequence) Copy(g *Grammar) Expansion {
	out := new(Sequence)
	out.exps = make([]Expansion, len(s.exps))

	for ind, e := range s.exps {
		out.exps[ind] = e.Copy(g)
	}

	return out
}

func (s *Sequence) Match(str string, mode MatchMode) {
	s.str = str
	s.mode = mode

	s.nextInd = 0

	s.exps[0].Match(str, mode)
}

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

func (s *Sequence) Scan(p Processor) {
	for _, exp := range s.exps {
		exp.Scan(p)
	}
}
