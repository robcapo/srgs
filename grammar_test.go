package srgs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrefix(t *testing.T) {
	assert := assert.New(t)

	g := NewGrammar()
	err := g.LoadXml(animalXml)
	assert.Nil(err)

	assert.True(g.HasPrefix("i am an antler"))
	assert.True(g.HasPrefix("i am an aardvark"))
	assert.True(g.HasPrefix("i am"))
	assert.True(g.HasPrefix("i am an"))
	assert.True(g.HasPrefix("i am an an"))
	assert.False(g.HasPrefix("i am an antler eater"))
	assert.False(g.HasPrefix("i am an ape"))
	assert.False(g.HasPrefix("i an"))
}

func TestMatch(t *testing.T) {
	assert := assert.New(t)

	g := NewGrammar()
	err := g.LoadXml(animalXml)
	assert.Nil(err)

	assert.True(g.HasMatch("i am an antler"))
	assert.True(g.HasMatch("i am an aardvark"))

	assert.False(g.HasMatch("i am an antler eater"))
	assert.False(g.HasMatch("i am an"))
	assert.False(g.HasMatch("i am an ape"))
}

func TestSequenceWithMultipleAlternatives(t *testing.T) {
	assert := assert.New(t)

	g := NewGrammar()
	err := g.LoadXml(nameXml)
	if !assert.Nil(err) {
		return
	}

	assert.True(g.HasMatch("my name is rob"))
	assert.True(g.HasMatch("my name is ram"))
	assert.True(g.HasMatch("my name is kaustav"))
}

//func TestSisr(t *testing.T) {
//	assert := assert.New(t)
//
//	xml := `<?xml version="1.0" encoding="UTF-8" ?>
//<grammar xmlns="http://www.w3.org/2001/06/grammar" version="1.0" xml:lang="en-US" root="example" tag-format="swi-semantics/1.0">
//	<rule scope="public" id="example">
//		my age is
//		<one-of>
//			<item>
//				ten
//				<tag>out = 10;</tag>
//			</item>
//			<item>
//				fifteen
//				<tag>out = 15;</tag>
//			</item>
//			<item>
//				twenty
//				<tag>out = 20;</tag>
//			</item>
//		</one-of>
//	</rule>
//</grammar>
//`
//	g := NewGrammar()
//	err := g.LoadXml(xml)
//
//	if !assert.Nil(err) {
//		return
//	}
//
//	_, seq, err := g.Root.Consume("my age is fifteen")
//
//	if !assert.Nil(err) {
//		return
//	}
//
//	p := new(SISRProcessor)
//
//	seq.AppendToProcessor(p)
//
//	assert.Equal("my age is fifteen", p.GetInterpretation())
//
//	inst, err := p.GetInstance()
//
//	if !assert.Nil(err) {
//		return
//	}
//
//	assert.Equal("15", inst)
//}

//func BenchmarkDigitsOne(b *testing.B) {
//	g := NewGrammar()
//	g.LoadXml(digitsXml)
//
//	benchmarkConsumeStack(b, g, "one")
//}
//func BenchmarkDigitsOneTwo(b *testing.B) {
//	g := NewGrammar()
//	g.LoadXml(digitsXml)
//
//	benchmarkConsumeStack(b, g, "one two")
//}
//func BenchmarkDigitsOneTwoThree(b *testing.B) {
//	g := NewGrammar()
//	g.LoadXml(digitsXml)
//
//	benchmarkConsumeStack(b, g, "one two three")
//}
//func BenchmarkDigitsOneTwoThreeFour(b *testing.B) {
//	g := NewGrammar()
//	g.LoadXml(digitsXml)
//
//	benchmarkConsumeStack(b, g, "one two three four")
//}
//func BenchmarkDigitsOneTwoThreeFourFive(b *testing.B) {
//	g := NewGrammar()
//	g.LoadXml(digitsXml)
//
//	benchmarkConsumeStack(b, g, "one two three four five")
//}
//
//
//func benchmarkConsumeStack(b *testing.B, g *Grammar, prefix string) {
//	for i := 0; i < b.N; i++ {
//		stk := stack.New()
//		g.Root.ConsumeStack(prefix, stk)
//	}
//}

//func BenchmarkDigitsMatchOneTwo(b *testing.B) {
//	g := NewGrammar()
//	g.LoadXml(digitsXml)
//
//	benchmarkMatch(b, g, "one two")
//}
//func BenchmarkDigitsMatchOneTwoThreeFourFive(b *testing.B) {
//	g := NewGrammar()
//	g.LoadXml(digitsXml)
//
//	benchmarkMatch(b, g, "one two three four five")
//}
//
//var (
//	outStr string
//	outErr error
//)
//func benchmarkMatch(b *testing.B, g *Grammar, prefix string) {
//	var str string
//	var err error
//	for i := 0; i < b.N; i++ {
//		str, err = g.Root.Match(prefix)
//	}
//	outStr = str
//	outErr = err
//}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := NewGrammar()
		g.LoadXml(digitsXml)
	}
}

//func TestDigits(t *testing.T) {
//	assert := assert.New(t)
//
//	g := NewGrammar()
//	err := g.LoadXml(digitsXml)
//
//	if !assert.Nil(err) {
//		return
//	}
//
//	stk := stack.New()
//	_, _, err = g.Root.ConsumeStack("one two three four five", stk)
//
//	if !assert.Nil(err) {
//		return
//	}
//
//	p := new(SISRProcessor)
//	p.ProcessStack(stk)
//
//	assert.Equal("one two three four five", p.GetInterpretation())
//
//	out, err := p.GetInstance()
//	assert.Nil(err)
//	assert.Equal("12345", out)
//
//	stk = stack.New()
//	_, _, err = g.Root.ConsumeStack("triple three four five", stk)
//
//	p = new(SISRProcessor)
//	p.ProcessStack(stk)
//
//	assert.Nil(err)
//
//	out, err = p.GetInstance()
//	assert.Nil(err)
//
//	assert.Equal("33345", out)
//}
