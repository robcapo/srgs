package srgs

type Sequence struct {
	exps []Expansion

	str  string
	mode MatchMode

	nextInd int
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

		s.nextInd = i + 1

		if i+1 < len(s.exps) {
			s.exps[i+1].Match(str, s.mode)
		}
	}

	return str, err
}

//
//func (s Sequence) ConsumeStack(str string, stack *stack.Stack) (string, int, error) {
//	var err error
//	var pushes, p int
//	for _, e := range s.exps {
//		str, p, err = e.ConsumeStack(str, stack)
//
//		if err != nil {
//			for i := pushes; i > 0; i-- {
//				stack.Pop()
//			}
//			return str, 0, err
//		}
//
//		pushes += p
//	}
//
//	return str, pushes, nil
//}

func (s *Sequence) AppendToProcessor(p Processor) {
	for _, exp := range s.exps {
		exp.AppendToProcessor(p)
	}
}
