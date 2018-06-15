package srgs

type Item struct {
	children  []Expansion
	repeatMin int
	repeatMax int

	str  string
	mode MatchMode

	nextInd int
	scanInd int
}

func (it *Item) Copy(refs RuleRefs) Expansion {
	children := make([]Expansion, len(it.children))
	for i := 0; i < len(it.children); i++ {
		children[i] = it.children[i].Copy(refs)
	}

	return &Item{
		children:  children,
		repeatMin: it.repeatMin,
		repeatMax: it.repeatMax,
		str:       it.str,
		mode:      it.mode,
		nextInd:   it.nextInd,
	}
}

func NewItem(child Expansion, repeatMin, repeatMax int, r RuleRefs) *Item {
	children := make([]Expansion, repeatMax+1)
	children[0] = NewToken("")
	for i := 1; i < len(children); i++ {
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

	it.nextInd = 0

	it.children[0].Match(str, mode)
}

func (it *Item) Next() (string, error) {
	if it.nextInd < 0 {
		return "", NoMatch
	}

	// loop all the way up to the child right before the min-repeat
	for {
		if it.nextInd >= it.repeatMin {
			break
		}

		str, err := it.next()

		if err != nil {
			return str, err
		}
	}

	return it.next()
}

func (it *Item) next() (string, error) {
	str, err := it.children[it.nextInd].Next()

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

	it.scanInd = it.nextInd
	if it.nextInd+1 < len(it.children) {
		it.nextInd++
		it.children[it.nextInd].Match(str, it.mode)
	}

	return str, err
}

func (it *Item) Scan(processor Processor) {
	for i := 1; i <= it.scanInd; i++ {
		it.children[i].Scan(processor)
	}
}
