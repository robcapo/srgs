package srgs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestItem_MatchRepeatRange(t *testing.T) {
	assert := assert.New(t)

	tok := NewToken("rob")
	item := NewItem(tok, 3, 5)

	item.Match("rob", ModeExact)
	_, err := item.Next()
	assert.Equal(PrefixOnly, err)

	item.Match("rob", ModePrefix)
	_, err = item.Next()
	assert.Nil(err)

	item.Match("rob rob rob", ModeExact)
	_, err = item.Next()
	assert.Nil(err)

	item.Match("rob rob rob rob", ModeExact)
	_, err = item.Next()
	assert.Nil(err)

	item.Match("rob rob rob rob rob", ModeExact)
	str, err := item.Next()
	assert.Nil(err)
	assert.Equal("rob rob", str)
	str, _ = item.Next()
	assert.Equal("rob", str)
	str, err = item.Next()
	assert.Empty(str)
	assert.Nil(err)
	_, err = item.Next()
	assert.Equal(PrefixOnly, err)

}
