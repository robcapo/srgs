package srgs

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAlternative_Match(t *testing.T) {
	assert := assert.New(t)

	alt := NewAlternative(
		NewItem(NewToken("rob"), 1, 1),
		NewItem(NewToken("rob"), 1, 1),
		NewItem(NewToken("ram"), 1, 1),
		NewItem(NewToken("ram malav"), 1, 1),
		NewItem(NewToken("kaustav"),1,1),
		NewItem(NewToken("kaustav datta"), 1,1),
	)

	alt.Match("kaustav", ModeExact)
	_, err := alt.Next()
	assert.Nil(err)
	_, err = alt.Next()
	assert.Equal(PrefixOnly, err)

}
