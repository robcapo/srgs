package main

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

func processExpansion(exp Expansion) (string, error) {
	processor := new(SimpleProcessor)
	exp.AppendToProcessor(processor)

	return processor.GetString()
}
