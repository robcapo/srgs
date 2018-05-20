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

func (it *Item) Scan(processor Processor) {
	for i := 0; i < it.currentRepeat; i++ {
		it.child.Scan(processor)
	}
}
