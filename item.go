package srgs

type Item struct {
	child     Expansion
	repeatMin int
	repeatMax int

	str  string
	mode MatchMode

	state      ItemState
	trackState bool

	deferToChild bool

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
	it.deferToChild = false
}

func (it *Item) Next() (string, error) {
	if it.deferToChild {
		str, err := it.child.Next()

		if err == nil {
			if it.trackState && it.currentRepeat > 0 {
				it.state[it.currentRepeat-1] = it.child.GetState()
			}
			return str, err
		}
	}

	it.deferToChild = false
	it.currentRepeat++

	if it.currentRepeat > it.repeatMax {
		it.currentRepeat--
		return "", NoMatch
	}

	if it.trackState {
		it.state = make(ItemState, it.currentRepeat)
	}

	var str = it.str
	var err error
	for i := 0; i < it.currentRepeat; i++ {
		it.child.Match(str, it.mode)
		str, err = it.child.Next()

		if err != nil {
			break
		}

		if it.trackState {
			it.state[i] = it.child.GetState()
		}
	}

	if it.currentRepeat > 0 && err == nil {
		it.deferToChild = true
	}

	return str, err
}

func (it *Item) Scan(processor Processor) {
	for i := 0; i < len(it.state); i++ {
		if it.trackState {
			it.child.SetState(it.state[i])
		}
		it.child.Scan(processor)
	}
}

func (it *Item) SetState(s State) {
	state, ok := s.(ItemState)

	if !ok {
		panic("Got invalid state. Expecting ItemState")
	}

	it.currentRepeat = len(state)

	it.state = state
}

func (it *Item) GetState() State {
	return it.state
}

func (it *Item) TrackState(t bool) {
	it.trackState = t
	it.child.TrackState(t)
}
