package srgs

type Alternative struct {
	items []Expansion

	str string
	currentInd int
}

func NewAlternative(items ...Expansion) *Alternative {
	return &Alternative{items:items}
}

func (a *Alternative) Match(str string, mode MatchMode) {
	a.str = str
	a.currentInd = 0

	for _, i := range a.items {
		i.Match(str, mode)
	}
}
func (a *Alternative) Next() (string, error) {
	outErr := NoMatch
	for i := a.currentInd; i < len(a.items); i++ {
		var str string
		var err error

		str, err = a.items[i].Next()

		if err == PrefixOnly {
			outErr = PrefixOnly
		} else if err == nil {
			a.currentInd = i
			return str, nil
		}
	}

	return "", outErr
}

//func (a Alternative) ConsumeStack(str string, stack *stack.Stack) (string, int, error) {
//	outErr := NoMatch
//	for _, alt := range a.items {
//		out, p, err := alt.ConsumeStack(str, stack)
//
//		if err == nil {
//			return out, p, err
//		}
//
//		if err == PrefixOnly {
//			outErr = PrefixOnly
//		}
//	}
//
//	return "", 0, outErr
//}
func (a *Alternative) AppendToProcessor(p Processor) {}
