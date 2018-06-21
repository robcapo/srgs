package srgs

type Scorer interface {
	AppendIDAndMatch(id string, match string)
	GetScore() float64
}

type ScorerImplementation struct {
	matches []idAndMatch

	lms map[string]*kenlm
}

type idAndMatch struct {
	id    string
	match string
}

func NewScorer() Scorer {
	return &ScorerImplementation{
		matches: make([]idAndMatch, 0, 1),
	}
}

func (s *ScorerImplementation) AppendIDAndMatch(id string, match string) {
	s.matches = append(s.matches, idAndMatch{id, match})
}

func (s *ScorerImplementation) GetScore() float64 { return 0.0 }
