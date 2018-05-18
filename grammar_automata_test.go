package srgs

//
//func TestConsume(t *testing.T) {
//	assert := assert.New(t)
//
//	g := NewGrammar()
//	err := g.LoadXml(`<?xml version="1.0" encoding="UTF-8" ?>
//<grammar xmlns="http://www.w3.org/2001/06/grammar" version="1.0" xml:lang="en-US" root="example" tag-format="swi-semantics/1.0">
//	<rule id="example">
//		<one-of>
//			<item>my name is</item>
//			<item>my name</item>
//		</one-of>
//		is rob
//	</rule>
//</grammar>
//`)
//
//	if !assert.Nil(err) {
//		return
//	}
//
//	g.Root.Match("my name is rob", ModePrefix)
//	str, err := g.Root.Next()
//
//	assert.Nil(err)
//	assert.Empty(str)
//}
