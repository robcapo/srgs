package srgs

type Rules map[string]Expansion

type RuleRef struct {
	rule   *Expansion
	ruleId string
}

func (r *RuleRef) Match(str string, mode MatchMode) {
	(*r.rule).Match(str, mode)
}
func (r *RuleRef) Next() (string, error) {
	return (*r.rule).Next()
}

//func (r RuleRef) ConsumeStack(str string, stack *stack.Stack) (string, int, error) {
//	PreProcessTag(r.ruleId).ConsumeStack(str, stack)
//	out, p, err := (*r.rule).ConsumeStack(str, stack)
//	PostProcessTag(r.ruleId).ConsumeStack(str, stack)
//	return out, p + 2, err
//}

func (r *RuleRef) AppendToProcessor(p Processor) {}
