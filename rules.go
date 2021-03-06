package srgs

import (
	"fmt"
)

type Rules map[string]Expansion
type RuleRefs map[string][]*RuleRef

type RuleRef struct {
	rule   Expansion
	ruleId string
}

func (r *RuleRef) Match(str string, mode MatchMode) {
	r.rule.Match(str, mode)
}
func (r *RuleRef) Next() (string, error) {
	return r.rule.Next()
}
func (r *RuleRef) Copy(rr RuleRefs) Expansion {
	ref := new(RuleRef)
	ref.ruleId = r.ruleId
	if r.rule != nil {
		ref.rule = r.rule.Copy(rr)
	}

	rr[ref.ruleId] = append(rr[ref.ruleId], ref)

	return ref
}

func (r *RuleRef) Scan(p Processor) {
	p.AppendTag("scopes.push({'rules':{}, 'out':undefined, 'raw':undefined});")
	r.rule.Scan(p)
	p.AppendTag(fmt.Sprintf(`var last = scopes.pop();
scopes[scopes.length-1]['rules']['%s'] = {'out': last.out, 'raw': last.raw};
scopes[scopes.length-1]['raw'] = scopes[scopes.length-1]['raw'] ? scopes[scopes.length-1]['raw'] + ' ' + last.raw : last.raw;
`, r.ruleId))
}
