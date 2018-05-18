package srgs

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSequence_MatchSimplePrefixMode(t *testing.T) {
	assert := assert.New(t)
	seq := new(Sequence)
	seq.exps = []Expansion{
		NewToken("my"),
		NewToken("name is"),
	}

	seq.Match("my name", ModePrefix)
	str, err := seq.Next()

	assert.Nil(err)
	assert.Empty(str)


	seq.Match("my nam", ModePrefix)
	str, err = seq.Next()

	assert.Nil(err)
	assert.Empty(str)

	seq.Match("my name is rob", ModePrefix)
	str, err = seq.Next()

	assert.Nil(err)
	assert.Equal("rob", str)

	seq.Match("my names", ModePrefix)
	str, err = seq.Next()

	assert.Equal(NoMatch, err)

	seq.Match("your name is", ModePrefix)
	str, err = seq.Next()

	assert.Equal(NoMatch, err)
}
