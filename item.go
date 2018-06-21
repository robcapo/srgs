package srgs

type RepeatMode string

const (
	RepeatModeLazy   RepeatMode = "lazy"
	RepeatModeNormal RepeatMode = "normal"
	RepeatModeGreedy RepeatMode = "greedy"
)

type Item struct {
	children   []Expansion
	repeatMin  int
	repeatMax  int
	repeatMode RepeatMode

	saveString []string
	str        string
	mode       MatchMode

	nextInd int
	scanInd int
}

func (it *Item) Copy(refs RuleRefs) Expansion {
	children := make([]Expansion, len(it.children))
	for i := 0; i < len(it.children); i++ {
		children[i] = it.children[i].Copy(refs)
	}

	return &Item{
		children:   children,
		repeatMin:  it.repeatMin,
		repeatMax:  it.repeatMax,
		str:        it.str,
		mode:       it.mode,
		nextInd:    it.nextInd,
		repeatMode: it.repeatMode,
		saveString: it.saveString,
	}
}

func NewItem(child Expansion, repeatMode RepeatMode, repeatMin, repeatMax int, r RuleRefs) *Item {
	children := make([]Expansion, repeatMax+1)
	children[0] = NewToken("")
	for i := 1; i < len(children); i++ {
		children[i] = child.Copy(r)
	}

	return &Item{
		children:   children,
		repeatMin:  repeatMin,
		repeatMax:  repeatMax,
		repeatMode: repeatMode,
		saveString: make([]string, repeatMax+1),
	}
}

func (it *Item) Match(str string, mode MatchMode) {
	it.str = str
	it.mode = mode
	it.nextInd = 0
	it.saveString[0] = ""
	it.children[0].Match(str, mode)
}

func (it *Item) Next() (string, error) {
	if it.nextInd < 0 {
		return "", NoMatch
	}

	var str string
	var err error
	var lastSuccessfulString string
	var hadSuccess bool
	// loop all the way up to the child right before the min-repeat
	for {
		str, err = it.children[it.nextInd].Next()

		if err != nil {
			if hadSuccess && it.repeatMode == RepeatModeGreedy {
				return lastSuccessfulString, nil
			}
			it.scanInd--
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
		lastSuccessfulString = str
		if it.nextInd > 0 {
			hadSuccess = true
		}
		breaker := false
		if it.nextInd >= it.repeatMin && it.repeatMode != RepeatModeGreedy {
			breaker = true
		}

		it.scanInd = it.nextInd
		if it.nextInd+1 < len(it.children) && str != it.saveString[it.nextInd] {

			it.saveString[it.nextInd] = str
			it.nextInd++
			it.saveString[it.nextInd] = str
			it.children[it.nextInd].Match(str, it.mode)
		}

		if breaker {
			break
		}
	}

	return str, err
}

func (it *Item) Scan(processor Processor) {
	for i := 1; i <= it.scanInd; i++ {
		it.children[i].Scan(processor)
	}
}

func (it *Item) ScanIDAndMatch(s Scorer) {
	for i := 1; i <= it.scanInd; i++ {
		it.children[i].ScanIDAndMatch(s)
	}
}
