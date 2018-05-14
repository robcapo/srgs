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
		i am an
		<one-of>
			<item>antler</item>
			<item>aardvark</item>
		</one-of>
	</rule>
</grammar>
`
	root, err := ParseXml(xml)

	assert.Nil(err)

	prefix, seq, err := root.Consume("i am an antler")

	assert.Nil(err)
	assert.Equal("", prefix)

	processor := new(SimpleProcessor)

	seq.AppendToProcessor(processor)

	assert.Equal("i am an antler", processor.GetString())

}
