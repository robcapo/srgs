package srgs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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

func TestDigits(t *testing.T) {
	assert := assert.New(t)

	xml := `<?xml version="1.0" encoding="UTF-8" ?>
<grammar xmlns="http://www.w3.org/2001/06/grammar" version="1.0" xml:lang="en-US" root="example" tag-format="swi-semantics/1.0">
	<rule id="example">
		<item>my number is <ruleref uri="#triplet" /><tag>out = rules.triplet.out;</tag></item>
	</rule>

	<rule id="triplet">
		<one-of>
			<item>
				<ruleref uri="#doublet" />
				<ruleref uri="#digit" />
				<tag>out = rules.doublet.out + rules.digit.out;</tag>
			</item>
			<item>
				<ruleref uri="#digit" />
				<ruleref uri="#doublet" />
				<tag>out = rules.digit.out + rules.doublet.out;</tag>
			</item>
			<item>
				<ruleref uri="#digit" />
				<tag>out = rules.digit.out;</tag>
				<ruleref uri="#digit" />
				<tag>out += rules.digit.out;</tag>
				<ruleref uri="#digit" />
				<tag>out += rules.digit.out;</tag>
			</item>
			<item>
				triple
				<ruleref uri="#digit" />
				<tag>out = rules.digit.out + rules.digit.out + rules.digit.out;</tag>
			</item>
		</one-of>
	</rule>

	<rule id="doublet">
		<one-of>
			<item>ten <tag>out = '10';</tag></item>
			<item>eleven <tag>out = '10';</tag></item>
			<item>twelve <tag>out = '10';</tag></item>
			<item>thirteen <tag>out = '10';</tag></item>
			<item>fourteen <tag>out = '10';</tag></item>
			<item>fifteen <tag>out = '10';</tag></item>
			<item>sixteen <tag>out = '10';</tag></item>
			<item>seventeen <tag>out = '10';</tag></item>
			<item>eighteen <tag>out = '10';</tag></item>
			<item>nineteen <tag>out = '10';</tag></item>
			<item>
				<ruleref uri="#digit" />
				<tag>out = rules.digit.out;</tag>
				<ruleref uri="#digit" /> 
				<tag>out += rules.digit.out;</tag>
			</item>
		</one-of>
	</rule>

	<rule id="digit">
		<one-of>
			<item>one <tag>out = '1';</tag></item>
			<item>two <tag>out = '2';</tag></item>
			<item>three <tag>out = '3';</tag></item>
			<item>four <tag>out = '4';</tag></item>
			<item>five <tag>out = '5';</tag></item>
			<item>six <tag>out = '6';</tag></item>
			<item>seven <tag>out = '7';</tag></item>
			<item>eight <tag>out = '8';</tag></item>
			<item>nine <tag>out = '9';</tag></item>
			<item>zero <tag>out = '0';</tag></item>
			<item>oh <tag>out = '0';</tag></item>
		</one-of>
	</rule>
</grammar>
`
	g := NewGrammar()
	err := g.LoadXml(xml)

	if !assert.Nil(err) {
		return
	}

	_, seq, err := g.Root.Consume("my number is one two three")

	if !assert.Nil(err) {
		return
	}

	p := new(SISRProcessor)
	seq.AppendToProcessor(p)

	assert.Equal("my number is one two three", p.GetInterpretation())

	out, err := p.GetInstance()
	assert.Nil(err)
	assert.Equal("123", out)

	_, seq, err = g.Root.Consume("my number is triple three")

	p = new(SISRProcessor)
	seq.AppendToProcessor(p)

	assert.Nil(err)

	out, err = p.GetInstance()
	assert.Nil(err)

	assert.Equal("333", out)
}
