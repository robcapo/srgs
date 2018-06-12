package srgs

// Sequence is any sequence of legal expansions (see https://www.w3.org/TR/speech-grammar/#S2.3)
type Sequence struct {
	exps []Expansion

	str  string
	mode MatchMode

	state SequenceState

	nextInd int
}

// Implements Expansion Copy method
func (s *Sequence) Copy(g *Grammar) Expansion {
	out := new(Sequence)
	out.exps = make([]Expansion, len(s.exps))

	for ind, e := range s.exps {
		out.exps[ind] = e.Copy(g)
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

func (s *Sequence) GetState() State {
	state := make(SequenceState, len(s.exps))

	for i, exp := range s.exps {
		state[i] = exp.GetState()
	}

	return state
}

func (s *Sequence) SetState(state State) {
	seqState, ok := state.(SequenceState)

	if !ok {
		panic("Expecting sequence state")
	}

	if len(seqState) != len(s.exps) {
		panic("Sequence state did not have the correct length")
	}

	s.state = seqState
}

func (s *Sequence) TrackState(t bool) {
	for _, exp := range s.exps {
		exp.TrackState(t)
	}
}

// Implements Expansion Scan method
func (s *Sequence) Scan(p Processor) {
	for i, exp := range s.exps {
		if s.state != nil {
			exp.SetState(s.state[i])
		}

		exp.Scan(p)
	}
}
