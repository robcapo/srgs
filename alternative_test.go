package srgs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlternative_Match(t *testing.T) {
	assert := assert.New(t)

	alt := NewAlternative(
		NewItem(NewToken("rob"), 1, 1, make(RuleRefs)),
		NewItem(NewToken("rob"), 1, 1, make(RuleRefs)),
		NewItem(NewToken("ram"), 1, 1, make(RuleRefs)),
		NewItem(NewToken("ram malav"), 1, 1, make(RuleRefs)),
		NewItem(NewToken("kaustav"), 1, 1, make(RuleRefs)),
		NewItem(NewToken("kaustav datta"), 1, 1, make(RuleRefs)),
	)

	alt.Match("kaustav", ModeExact)
	_, err := alt.Next()
	assert.Nil(err)
	_, err = alt.Next()
	assert.Equal(PrefixOnly, err)

}
