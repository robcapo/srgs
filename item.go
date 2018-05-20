package srgs

type Item struct {
	child     Expansion
	repeatMin int
	repeatMax int

	str  string
	mode MatchMode

	currentRepeat int
}

func (it *Item) Copy(g *Grammar) Expansion {
	return &Item{
		child:     it.child.Copy(g),
		repeatMin: it.repeatMin,
		repeatMax: it.repeatMax,
	}
}

func NewItem(child Expansion, repeatMin, repeatMax int) *Item {
	return &Item{
		child:     child,
		repeatMin: repeatMin,
		repeatMax: repeatMax,
	}
}

func (it *Item) Match(str string, mode MatchMode) {
	it.str = str
	it.mode = mode
	it.currentRepeat = it.repeatMin - 1
}

func (it *Item) Next() (string, error) {
	it.currentRepeat++

	var str = it.str
	var err error
	for i := 0; i < it.currentRepeat; i++ {
		it.child.Match(str, it.mode)
		str, err = it.child.Next()

		if err != nil {
			break
		}
	}

	return str, err
}

func (it *Item) AppendToProcessor(processor Processor) {
	for i := 0; i < it.currentRepeat; i++ {
		it.child.AppendToProcessor(processor)
	}
}

//func (i Item) ConsumeStack(str string, stack *stack.Stack) (string, int, error) {
//	return i.consumeStack(str, stack, i.repeatMin, i.repeatMax)
//}
//func (i Item) consumeStack(str string, stack *stack.Stack, min, max int) (string, int, error) {
//	if max == 0 {
//		return str, 0, nil
//	}
//
//	outStr, p, err := i.child.ConsumeStack(str, stack)
//
//	if err != nil {
//		if min <= 0 {
//			return str, p, nil
//		}
//
//		return "", p, err
//	}
//
//	out2, p2, err2 := i.consumeStack(outStr, stack, min - 1, max - 1)
//
//	return out2, p + p2, err2
//}
