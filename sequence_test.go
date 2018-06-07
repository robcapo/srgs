package srgs

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
	_, err = seq.Next()

	assert.Equal(NoMatch, err)
}

func TestSequenceWithGarbage(t *testing.T) {
	assert := assert.New(t)

	seq := new(Sequence)
	seq.exps = []Expansion{
		new(Garbage),
		NewToken("ten"),
		new(Garbage),
	}

	seq.Match("i am ten years old", ModeExact)

	var str string
	var err error

	str, err = seq.Next()
	assert.Nil(err)
	assert.Equal("years old", str)

	str, err = seq.Next()
	assert.Nil(err)
	assert.Equal("old", str)

	str, err = seq.Next()
	assert.Nil(err)
	assert.Empty(str)

	_, err = seq.Next()
	assert.Equal(NoMatch, err)
}
