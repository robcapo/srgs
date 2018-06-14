package srgs

type Item struct {
	children  []Expansion
	repeatMin int
	repeatMax int

	str  string
	mode MatchMode

	currentRepeat int
	nextInd       int
}

func (it *Item) Copy(refs RuleRefs) Expansion {
	children := make([]Expansion, len(it.children))
	for i := 0; i < len(it.children); i++ {
		children[i] = it.children[i].Copy(refs)
	}

	return &Item{
		children:      children,
		repeatMin:     it.repeatMin,
		repeatMax:     it.repeatMax,
		str:           it.str,
		mode:          it.mode,
		currentRepeat: it.currentRepeat,
		nextInd:       it.nextInd,
	}
}

func NewItem(child Expansion, repeatMin, repeatMax int, r RuleRefs) *Item {
	children := make([]Expansion, repeatMax)
	for i := 0; i < len(children); i++ {
		children[i] = child.Copy(r)
	}

	return &Item{
		children:  children,
		repeatMin: repeatMin,
		repeatMax: repeatMax,
	}
}

func (it *Item) Match(str string, mode MatchMode) {
	it.str = str
	it.mode = mode

	it.currentRepeat = it.repeatMin
	it.nextInd = 0

	it.children[0].Match(str, mode)
}

func (it *Item) Next() (string, error) {
	if it.nextInd < 0 {
		return "", NoMatch
	}

	var str string
	var err error

	for i := it.nextInd; i < it.currentRepeat; i++ {
		str, err = it.children[i].Next()

		if err != nil {
			it.nextInd--
			return it.Next()
		}

		if i+1 < len(it.children) {
			it.nextInd = i + 1
			it.children[it.nextInd].Match(str, it.mode)
		}
	}

	str, err = it.children[it.currentRepeat].Next()

	if err != nil {
		it.currentRepeat--
		it.nextInd = it.currentRepeat

		if it.currentRepeat < it.repeatMin {
			it.currentRepeat = it.repeatMin
		}
	}

	if it.currentRepeat+1 < len(it.children) {
		it.currentRepeat++
		it.children[it.currentRepeat].Match(str, it.mode)
	}

	return str, err
}

func (it *Item) Scan(processor Processor) {
	for i := 0; i < it.currentRepeat; i++ {
		it.children[i].Scan(processor)
	}
}
