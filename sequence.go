package srgs

type Sequence struct {
	exps []Expansion

	str string
	mode MatchMode

	nextInd int
}

func (s *Sequence) Match(str string, mode MatchMode) {
	s.str = str
	s.mode = mode
	s.nextInd = 0
}

func (s *Sequence) Next() (string, error) {
	str, err :=  s.next(s.str, s.nextInd)

	if err != nil {
		s.nextInd--

		if s.nextInd < 0 {
			return str, err
		}

		return s.Next()
	}

	return str, err
}

func (s *Sequence) next(str string, i int) (string, error) {
	if i == len(s.exps) {
		return str, nil
	}

	s.exps[i].Match(str, s.mode)
	str, err := s.exps[i].Next()
	if err == nil {
		s.nextInd = i
		return s.next(str, i + 1)
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

