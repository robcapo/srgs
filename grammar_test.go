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

func TestDigitsPrefix(t *testing.T) {
	assert := assert.New(t)

	g := NewGrammar()
	err := g.LoadXml(digitsXml)

	if !assert.Nil(err) {
		return
	}

	assert.True(g.HasPrefix("one two"))
	assert.True(g.HasPrefix("one two three four "))
	assert.True(g.HasPrefix("two"))
	assert.True(g.HasPrefix("three five four one"))
	assert.True(g.HasPrefix("two three four five"))
	assert.False(g.HasPrefix("six five four three two two"))
	assert.True(g.HasPrefix("on"))
	assert.False(g.HasPrefix("fix"))
}

func TestDigitsMatch(t *testing.T) {
	assert := assert.New(t)

	g := NewGrammar()
	err := g.LoadXml(digitsXml)

	if !assert.Nil(err) {
		return
	}

	assert.True(g.HasMatch("one two three four five"))
}

func BenchmarkDigitsMatchOneTwo(b *testing.B) {
	g := NewGrammar()
	g.LoadXml(digitsXml)

	benchmarkMatch(b, g, "one two", ModeExact)
}
func BenchmarkDigitsPrefixOneTwo(b *testing.B) {
	g := NewGrammar()
	g.LoadXml(digitsXml)

	benchmarkMatch(b, g, "one two", ModePrefix)
}
func BenchmarkDigitsMatchFoo(b *testing.B) {
	g := NewGrammar()
	g.LoadXml(digitsXml)

	benchmarkMatch(b, g, "foo", ModeExact)
}
func BenchmarkDigitsPrefixFoo(b *testing.B) {
	g := NewGrammar()
	g.LoadXml(digitsXml)

	benchmarkMatch(b, g, "foo", ModePrefix)
}
func BenchmarkDigitsMatchOneTwoThreeFourFix(b *testing.B) {
	g := NewGrammar()
	g.LoadXml(digitsXml)

	benchmarkMatch(b, g, "one two three four fix", ModeExact)
}
func BenchmarkDigitsPrefixOneTwoThreeFourFix(b *testing.B) {
	g := NewGrammar()
	g.LoadXml(digitsXml)

	benchmarkMatch(b, g, "one two three four fix", ModePrefix)
}
func BenchmarkDigitsMatchOneTwoThreeFourFive(b *testing.B) {
	g := NewGrammar()
	g.LoadXml(digitsXml)

	benchmarkMatch(b, g, "one two three four five", ModeExact)
}
func BenchmarkDigitsPrefixOneTwoThreeFourFive(b *testing.B) {
	g := NewGrammar()
	g.LoadXml(digitsXml)

	benchmarkMatch(b, g, "one two three four five", ModePrefix)
}

var match bool

func benchmarkMatch(b *testing.B, g *Grammar, prefix string, mode MatchMode) {
	var out bool
	for i := 0; i < b.N; i++ {
		out = g.HasPrefix(prefix)
	}

	match = out
}

func TestSisr(t *testing.T) {
	assert := assert.New(t)

	xml := `<?xml version="1.0" encoding="UTF-8" ?>
<grammar xmlns="http://www.w3.org/2001/06/grammar" version="1.0" xml:lang="en-US" root="example" tag-format="swi-semantics/1.0">
	<rule scope="public" id="example">
		my age is
		<one-of>
			<item>
				ten
				<tag>out = 10;</tag>
			</item>
			<item>
				fifteen
				<tag>out = 15;</tag>
			</item>
			<item>
				twenty
				<tag>out = 20;</tag>
			</item>
		</one-of>
	</rule>
</grammar>
`
	g := NewGrammar()
	err := g.LoadXml(xml)

	if !assert.Nil(err) {
		return
	}

	p := new(SISRProcessor)

	err = g.GetMatch("my age is fifteen", p)

	if !assert.Nil(err) {
		return
	}

	assert.Equal("my age is fifteen", p.GetInterpretation())

	inst, err := p.GetInstance()

	if !assert.Nil(err) {
		return
	}

	assert.Equal("15", inst)
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := NewGrammar()
		g.LoadXml(digitsXml)
	}
}

func TestDigits(t *testing.T) {
	assert := assert.New(t)

	g := NewGrammar()
	if !assert.Nil(g.LoadXml(digitsXml)) {
		return
	}

	p := new(SISRProcessor)

	if !assert.Nil(g.GetMatch("one two three four five", p)) {
		return
	}

	assert.Equal("one two three four five", p.GetInterpretation())

	out, err := p.GetInstance()
	assert.Nil(err)
	assert.Equal("12345", out)

	p = new(SISRProcessor)
	if !assert.Nil(g.GetMatch("triple three four five", p)) {
		return
	}

	out, err = p.GetInstance()
	assert.Nil(err)

	assert.Equal("33345", out)
}
