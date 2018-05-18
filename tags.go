package srgs

//
//type Tag string
//
//func (t Tag) Match(str string) {}
//func (t Tag) MatchPrefix(str string) (string, error) { return str, nil }
//func (t Tag) ConsumeStack(str string, stack *stack.Stack) (string, int, error) {
//	stack.Push(t)
//	return str, 1, nil
//}
//
//func (t Tag) Consume(str string) (string, Sequence, error) { return str, []Expansion{t}, nil }
//func (t Tag) AppendToProcessor(p Processor)                {
//	p.AppendTag(string(t))
//}
//
//
//type PreProcessTag string
//func (t PreProcessTag) MatchPrefix(str string) (string, error) { return str, nil }
//func (t PreProcessTag) Match(str string) (string, error) { return str, nil }
//func (t PreProcessTag) ConsumeStack(str string, stack *stack.Stack) (string, int, error) {
//	stack.Push(t)
//	return str, 1, nil
//}
//
//func (t PreProcessTag) Consume(str string) (string, Sequence, error) { return str, []Expansion{t}, nil }
//func (t PreProcessTag) AppendToProcessor(p Processor) {
//	p.AppendTag(fmt.Sprintf(`
//rules.%s = {};
//
//if (root == undefined) {
//	root = rules.%s;
//}
//
//(function() {
//	var out;
//`, string(t), string(t)))
//}
//
//type PostProcessTag string
//func (t PostProcessTag) Match(str string) (string, error) { return str, nil }
//func (t PostProcessTag) MatchPrefix(str string) (string, error) { return str, nil }
//func (t PostProcessTag) ConsumeStack(str string, stack *stack.Stack) (string, int, error) {
//	stack.Push(t)
//	return str, 1, nil
//}
//func (t PostProcessTag) Consume(str string) (string, Sequence, error) { return str, []Expansion{t}, nil }
//func (t PostProcessTag) AppendToProcessor(p Processor) {
//	p.AppendTag(fmt.Sprintf(`
//	rules.%s.out = out;
//})();
//`, string(t)))
//}
//
