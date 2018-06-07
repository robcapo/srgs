package srgs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToken_MatchPrefixMode(t *testing.T) {
	assert := assert.New(t)

	tok := Token{token: "my name is"}

	// Test covered utterance
	tok.Match("my name is rob", ModePrefix)
	str, err := tok.Next()

	assert.Nil(err)
	assert.Equal("rob", str)

	// Ensure next can only be called once
	str, err = tok.Next()
	assert.Equal(NoMatch, err)

	// Test prefix utterance
	tok.Match("my name", ModePrefix)
	str, err = tok.Next()

	assert.Nil(err)
	assert.Equal("", str)

	tok.Match("my name is", ModePrefix)
	str, err = tok.Next()

	assert.Nil(err)
	assert.Empty(str)
}

func TestToken_MatchExactMode(t *testing.T) {
	assert := assert.New(t)

	tok := Token{token: "my name is"}

	// Test covered utterance
	tok.Match("my name is rob", ModeExact)
	str, err := tok.Next()

	assert.Nil(err)
	assert.Equal("rob", str)

	// Ensure next can only be called once
	str, err = tok.Next()
	assert.Equal(NoMatch, err)

	// Test prefix utterance
	tok.Match("my name", ModeExact)
	str, err = tok.Next()

	assert.Equal(PrefixOnly, err)
	assert.Equal("", str)

	tok.Match("my name is", ModeExact)
	str, err = tok.Next()

	assert.Nil(err)
	assert.Empty(str)
}
