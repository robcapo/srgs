package srgs

type Item struct {
	children  []Expansion
	repeatMin int
	repeatMax int

	str  string
	mode MatchMode

	nextInd int
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

	str, err = it.children[it.nextInd].Next()

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

	if it.nextInd+1 < len(it.children) {
		it.nextInd++
		it.children[it.nextInd].Match(str, it.mode)
	}

	return str, err
}

func (it *Item) Scan(processor Processor) {
	for i := 0; i < it.nextInd+1; i++ {
		it.children[i].Scan(processor)
	}
}
