package srgs

type Alternative struct {
	items []Expansion

	str        string
	currentInd int
}

func (a *Alternative) Copy(g *Grammar) Expansion {
	out := new(Alternative)
	out.items = make([]Expansion, len(a.items))

	for ind, e := range a.items {
		out.items[ind] = e.Copy(g)
	}

	return out
}

func NewAlternative(items ...Expansion) *Alternative {
	return &Alternative{items: items}
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

func (a *Alternative) Scan(p Processor) {
	a.items[a.currentInd].Scan(p)
}
