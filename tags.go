package srgs

type Tag struct {
	text string

	match  string
	called bool
}

func NewTag(str string) *Tag {
	return &Tag{text: str}
}

func (t *Tag) Match(str string, mode MatchMode) {
	t.match = str
	t.called = false
}

func (t *Tag) Next() (string, error) {
	if t.called == true {
		return "", NoMatch
	}

	t.called = true

	return t.match, nil
}

func (t *Tag) Scan(p Processor) {
	p.AppendTag(`
(function () {
	var rules = scopes[scopes.length-1]['rules'];
	var out = scopes[scopes.length-1]['out'];
	var raw = scopes[scopes.length-1]['raw'];
	var GARBAGE = scopes[scopes.length-1]['GARBAGE'];
	var SLM = scopes[scopes.length-1]['SLM'];
`)
	p.AppendTag(t.text)

	p.AppendTag(`
	scopes[scopes.length-1]['out'] = out;
})();`)
}
func (t *Tag) Copy(r RuleRefs) Expansion {
	return &Tag{
		text:   t.text,
		match:  t.match,
		called: t.called,
	}
}
