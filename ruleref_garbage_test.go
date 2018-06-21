package srgs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGarbage(t *testing.T) {
	assert := assert.New(t)
	g := new(Garbage)

	g.Match("my name is rob", ModeExact)

	var str string
	var err error

	str, err = g.Next()
	assert.Nil(err)
	assert.Equal("my name is rob", str)

	str, err = g.Next()
	assert.Nil(err)
	assert.Equal("name is rob", str)

	str, err = g.Next()
	assert.Nil(err)
	assert.Equal("is rob", str)

	str, err = g.Next()
	assert.Nil(err)
	assert.Equal("rob", str)

	str, err = g.Next()
	assert.Nil(err)
	assert.Empty(str)

	_, err = g.Next()
	assert.Equal(NoMatch, err)
}
