package srgs

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/golang-collections/collections/stack"
)

var digitsXml = `<?xml version="1.0" encoding="UTF-8" ?>
<grammar version="1.0" xml:lang="en-US" root="combined" mode="voice" tag-format="swi-semantics/1.0">
  <rule id="combined" scope="public">
	<item>
	  <ruleref uri="#quintet" />
	  <tag>out = rules.quintet.out;</tag>
	</item>
  </rule>
  
  <rule id="quintet">
	<one-of>
	  <item>
		<ruleref uri="#digit" />
		<tag>out = rules.digit.out</tag>
		<ruleref uri="#quartet" />
		<tag>out = out + rules.quartet.out</tag>
	  </item>
	  <item>
		<ruleref uri="#quartet" />
		<tag>out = rules.quartet.out</tag>
		<ruleref uri="#digit" />
		<tag>out = out + rules.digit.out</tag>
	  </item>
	  <item>
		<ruleref uri="#triplet" />
		<tag>out = rules.triplet.out</tag>
		<ruleref uri="#doublet" />
		<tag>out = out + rules.doublet.out</tag>
	  </item>
	  <item>
		<ruleref uri="#doublet" />
		<tag>out = rules.doublet.out</tag>
		<ruleref uri="#triplet" />
		<tag>out = out + rules.triplet.out</tag>
	  </item>
	</one-of>
  </rule>
  
  <rule id="quartet">
  	<one-of>
	<item>
	  <ruleref uri="#digit" />
	  <tag>out = rules.digit.out;</tag>
	  thousand
	  <tag>out = rules.digit.out + "000";</tag>
	</item>
	<!-- 1 + 3 -->
	<item>
	  <ruleref uri="#digit" />
	  <tag>out = rules.digit.out</tag>
	  <ruleref uri="#triplet" />
	  <tag>out = out + rules.triplet.out</tag>
	</item>
	<item>
	  <ruleref uri="#triplet" />
	  <tag>out = rules.triplet.out</tag>
	  <ruleref uri="#digit" />
	  <tag>out = out + rules.digit.out</tag>
	</item>
	<!-- 2 + 2 -->
	<item repeat="2">
	  <ruleref uri="#doublet" />
	  <tag>out = out ? out + rules.doublet.out : rules.doublet.out</tag>
	</item>
	<!-- quadruple 1 -->
	<item>
	  <ruleref uri="#four" />
	  <ruleref uri="#digit" />
	  <tag>out = "" + rules.digit.out + rules.digit.out + rules.digit.out + rules.digit.out</tag>
	</item>
	</one-of>
  </rule>
  
  <rule id="four" scope="public">
	 <one-of>
	  <item>quad</item>
	  <item>quadruple</item>
	 </one-of>
  </rule>
  
  <rule id="triplet">
	<one-of>
	  <item>
		<ruleref uri="#doublet"/>
		<tag>out = rules.doublet.out</tag>
		<ruleref uri="#digit" />
		<tag>out = out + rules.digit.out</tag>
	  </item>

	  <item>
		<ruleref uri="#digit"/>
		<tag>out = rules.digit.out</tag>
		<ruleref uri="#doublet" />
		<tag>out = out + rules.doublet.out</tag>
	  </item>

	  <item>
	  	triple <ruleref uri="#digit" /><tag>out = rules.digit.out + rules.digit.out + rules.digit.out</tag>
	  </item>
	</one-of>
  </rule>
  
  
  <rule id="doublet">
	<one-of>
	  <item>
		<ruleref uri="#digit" />
		<tag>out = rules.digit.out</tag>
		<ruleref uri="#digit" />
		<tag>out = out + rules.digit.out</tag>
	  </item>
	  
	  <item>
		double <ruleref uri="#digit" />
		<tag>out = rules.digit.out + rules.digit.out</tag>
	  </item>
	  
	  <item>
		<one-of>
		  <item>ten <tag>out = "10";</tag></item>
		  <item>eleven <tag>out = "11";</tag></item>
		  <item>twelve <tag>out = "12";</tag></item>
		  <item>thirteen <tag>out = "13";</tag></item>
		  <item>fourteen <tag>out = "14";</tag></item>
		  <item>fifteen <tag>out = "15";</tag></item>
		  <item>sixteen <tag>out = "16";</tag></item>
		  <item>seventeen <tag>out = "17";</tag></item>
		  <item>eighteen <tag>out = "18";</tag></item>
		  <item>nineteen <tag>out = "19";</tag></item>
		</one-of>
	  </item>
	  
	  <item>
		<ruleref uri="#tens" />
		<tag>out = rules.tens.tens;</tag>
		<ruleref uri="#ones" />
		<tag>out = out + rules.ones.ones;</tag>
	  </item>
	</one-of>
  </rule>
  
  <rule id="tens">
	<one-of>
	  <item>twenty <tag>out.tens = "2"; out="20";</tag></item>
	  <item>thirty <tag>out.tens = "3"; out="30";</tag></item>
	  <item>forty <tag>out.tens = "4"; out="40";</tag></item>
	  <item>fifty <tag>out.tens = "5"; out="50";</tag></item>
	  <item>sixty <tag>out.tens = "6"; out="60";</tag></item>
	  <item>seventy <tag>out.tens = "7"; out="70";</tag></item>
	  <item>eighty <tag>out.tens = "8"; out="80";</tag></item>
	  <item>ninety <tag>out.tens = "9"; out="90";</tag></item>
	</one-of>
  </rule>
  
  <rule id="ones">
	<one-of>
	  <item>one <tag>out.ones = "1";</tag></item>
	  <item>two <tag>out.ones = "2";</tag></item>
	  <item>three <tag>out.ones = "3";</tag></item>
	  <item>four <tag>out.ones = "4";</tag></item>
	  <item>five <tag>out.ones = "5";</tag></item>
	  <item>six <tag>out.ones = "6";</tag></item>
	  <item>seven <tag>out.ones = "7";</tag></item>
	  <item>eight <tag>out.ones = "8";</tag></item>
	  <item>nine <tag>out.ones = "9";</tag></item>
	</one-of>
  </rule>
  
  <rule id="digit">
	<one-of>
	  <item weight="0.1">
		oh
		<tag>out='0';</tag>
	  </item>
	  <item>
		zero
		<tag>out='0';</tag>
	  </item>
	  <item>
		one
		<tag>out='1';</tag>
	  </item>
	  <item>
		two
		<tag>out='2';</tag>
	  </item>
	  <item>
		three
		<tag>out='3';</tag>
	  </item>
	  <item>
		four
		<tag>out='4';</tag>
	  </item>
	  <item>
		five
		<tag>out='5';</tag>
	  </item>
	  <item>
		six
		<tag>out='6';</tag>
	  </item>
	  <item>
		seven
		<tag>out='7';</tag>
	  </item>
	  <item>
		eight
		<tag>out='8';</tag>
	  </item>
	  <item>
		nine
		<tag>out='9';</tag>
	  </item>
	</one-of>
  </rule>
</grammar>
`

func TestParseXml(t *testing.T) {
	assert := assert.New(t)

	xml := `<?xml version="1.0" encoding="UTF-8" ?>
<grammar xmlns="http://www.w3.org/2001/06/grammar" version="1.0" xml:lang="en-US" root="example" tag-format="swi-semantics/1.0">
	<rule scope="public" id="example">
		i am an <ruleref uri="#animal" />
	</rule>

	<rule id="animal">
		<one-of>
			<item>antler</item>
			<item>aardvark</item>
		</one-of>
	</rule>
</grammar>
`
	g := NewGrammar()
	err := g.LoadXml(xml)
	assert.Nil(err)

	prefix, seq, err := g.Root.Consume("i am an antler")

	processExpansion := func(exp Expansion) (string, error) {
		processor := new(SimpleProcessor)
		exp.AppendToProcessor(processor)

		return processor.GetInterpretation(), nil
	}

	assert.Nil(err)
	assert.Equal("", prefix)
	out, err := processExpansion(seq)
	assert.Nil(err)
	assert.Equal("i am an antler", out)

	prefix, seq, err = g.Root.Consume("i am an aardvark")

	assert.Nil(err)
	assert.Equal("", prefix)
	out, err = processExpansion(seq)
	assert.Nil(err)
	assert.Equal("i am an aardvark", out)

	_, _, err = g.Root.Consume("i am an")
	assert.Equal(PrefixOnly, err)

	_, _, err = g.Root.Consume("i am a")
	assert.Equal(PrefixOnly, err)

	_, _, err = g.Root.Consume("i am an ape")
	assert.Equal(NoMatch, err)
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

	_, seq, err := g.Root.Consume("my age is fifteen")

	if !assert.Nil(err) {
		return
	}

	p := new(SISRProcessor)

	seq.AppendToProcessor(p)

	assert.Equal("my age is fifteen", p.GetInterpretation())

	inst, err := p.GetInstance()

	if !assert.Nil(err) {
		return
	}

	assert.Equal("15", inst)
}

func BenchmarkDigitsOne(b *testing.B) {
	g := NewGrammar()
	g.LoadXml(digitsXml)

	benchmarkConsumeStack(b, g, "one")
}
func BenchmarkDigitsOneTwo(b *testing.B) {
	g := NewGrammar()
	g.LoadXml(digitsXml)

	benchmarkConsumeStack(b, g, "one two")
}
func BenchmarkDigitsOneTwoThree(b *testing.B) {
	g := NewGrammar()
	g.LoadXml(digitsXml)

	benchmarkConsumeStack(b, g, "one two three")
}
func BenchmarkDigitsOneTwoThreeFour(b *testing.B) {
	g := NewGrammar()
	g.LoadXml(digitsXml)

	benchmarkConsumeStack(b, g, "one two three four")
}
func BenchmarkDigitsOneTwoThreeFourFive(b *testing.B) {
	g := NewGrammar()
	g.LoadXml(digitsXml)

	benchmarkConsumeStack(b, g, "one two three four five")
}


func benchmarkConsumeStack(b *testing.B, g *Grammar, prefix string) {
	for i := 0; i < b.N; i++ {
		stk := stack.New()
		g.Root.ConsumeStack(prefix, stk)
	}
}

func BenchmarkDigitsMatchOneTwo(b *testing.B) {
	g := NewGrammar()
	g.LoadXml(digitsXml)

	benchmarkMatch(b, g, "one two")
}
func BenchmarkDigitsMatchOneTwoThreeFourFive(b *testing.B) {
	g := NewGrammar()
	g.LoadXml(digitsXml)

	benchmarkMatch(b, g, "one two three four five")
}

var (
	outStr string
	outErr error
)
func benchmarkMatch(b *testing.B, g *Grammar, prefix string) {
	var str string
	var err error
	for i := 0; i < b.N; i++ {
		str, err = g.Root.Match(prefix)
	}
	outStr = str
	outErr = err
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
	err := g.LoadXml(digitsXml)

	if !assert.Nil(err) {
		return
	}

	stk := stack.New()
	_, _, err = g.Root.ConsumeStack("one two three four five", stk)

	if !assert.Nil(err) {
		return
	}

	p := new(SISRProcessor)
	p.ProcessStack(stk)

	assert.Equal("one two three four five", p.GetInterpretation())

	out, err := p.GetInstance()
	assert.Nil(err)
	assert.Equal("12345", out)

	stk = stack.New()
	_, _, err = g.Root.ConsumeStack("triple three four five", stk)

	p = new(SISRProcessor)
	p.ProcessStack(stk)

	assert.Nil(err)

	out, err = p.GetInstance()
	assert.Nil(err)

	assert.Equal("33345", out)
}
