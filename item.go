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

	var str = it.str
	var err error

	// loop all the way up to the child right before the min-repeat
	for i := it.nextInd; i < it.repeatMin-1; i++ {
		str, err = it.children[i].Next()

		if err != nil {
			it.nextInd--

			str2, err2 := it.Next()

			if err2 == nil {
				return str2, err2
			}

			if err == PrefixOnly {
				return str, err
			}

			return str2, err2
		}

		if i+1 < len(it.children) {
			it.nextInd = i + 1
			it.children[it.nextInd].Match(str, it.mode)
		}
	}

	if it.currentRepeat > 0 {
		str, err = it.children[it.currentRepeat-1].Next()
	}

	if err != nil {
		it.currentRepeat--
		it.nextInd = it.currentRepeat

		if it.currentRepeat < it.repeatMin {
			it.currentRepeat = it.repeatMin
		}
	}

	if it.currentRepeat < len(it.children) {
		it.children[it.currentRepeat].Match(str, it.mode)
		it.currentRepeat++
	}

	return str, err
}

func (it *Item) Scan(processor Processor) {
	for i := 0; i < it.currentRepeat; i++ {
		it.children[i].Scan(processor)
	}
}
