package srgs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToken_Consume(t *testing.T) {
	utterance := "i want to go to the mall"

	tok := Token("i want to")
	str, seq, err := tok.Consume(utterance)

	assert := assert.New(t)

	assert.Equal("go to the mall", str)
	assert.Nil(err, "Expected no error.")
	assert.Len(seq, 1)

	tok = Token("i am an ant")

	_, _, err = tok.Consume("i am an antler")
	assert.Equal(NoMatch, err)

	tok = Token("i am an antler")

	str, seq, err = tok.Consume("i am an antler")
	assert.Empty(str)
	assert.Len(seq, 1)
	assert.Equal(seq[0], tok)
	assert.Nil(err)

	str, _, err = tok.Consume("i am an")
	assert.Empty(str)
	assert.Equal(PrefixOnly, err)

	str, _, err = tok.Consume("i am an an")
	assert.Empty(str)
	assert.Equal(PrefixOnly, err)

	_, _, err = tok.Consume("i am an owl")
	assert.Equal(NoMatch, err)
}

func TestTag_Consume(t *testing.T) {
	assert := assert.New(t)

	tag := Tag("out = 1;")
	str, seq, err := tag.Consume("two three four")

	assert.Equal("two three four", str)
	assert.Len(seq, 1)
	assert.Equal(tag, seq[0])
	assert.Nil(err)
}

func TestSequence_Consume(t *testing.T) {
	assert := assert.New(t)

	iam := Token("i am")
	tag := Tag("out.person = 'me';")
	an := Token("an")
	antler := Token("antler")
	tag2 := Tag("out.animal = 'antler';")

	sequence := Sequence{iam, tag, an, antler, tag2}

	str, seq, err := sequence.Consume("i am an antler")

	assert.Empty(str)
	assert.Equal(nil, err)

	if assert.Len(seq, 5) {
		assert.Equal(iam, seq[0])
		assert.Equal(tag, seq[1])
		assert.Equal(an, seq[2])
		assert.Equal(antler, seq[3])
		assert.Equal(tag2, seq[4])
	}

	str, seq, err = sequence.Consume("i am an antler eater")

	assert.Equal("eater", str)
	assert.Nil(err)

	if assert.Len(seq, 5) {
		assert.Equal(iam, seq[0])
		assert.Equal(tag, seq[1])
		assert.Equal(an, seq[2])
		assert.Equal(antler, seq[3])
		assert.Equal(tag2, seq[4])
	}

	str, seq, err = sequence.Consume("i am an")

	assert.Equal(PrefixOnly, err)
}

func TestItem_Consume(t *testing.T) {
	assert := assert.New(t)

	iam := Token("i am")
	an := Token("an")
	antler := Token("antler")
	sequence := Sequence{iam, an, antler}

	item := Item{sequence, 0, 2}

	str, seq, err := item.Consume("i am an antler i am an antler")

	assert.Empty(str)
	assert.Nil(err)

	if assert.Len(seq, 6) {
		assert.Equal(iam, seq[0])
		assert.Equal(an, seq[1])
		assert.Equal(antler, seq[2])
		assert.Equal(iam, seq[3])
		assert.Equal(an, seq[4])
		assert.Equal(antler, seq[5])
	}

	str, seq, err = item.Consume("i am an antler i am an antler i am an antler")

	assert.Equal("i am an antler", str)
	assert.Nil(err)

	if assert.Len(seq, 6) {
		assert.Equal(iam, seq[0])
		assert.Equal(an, seq[1])
		assert.Equal(antler, seq[2])
		assert.Equal(iam, seq[3])
		assert.Equal(an, seq[4])
		assert.Equal(antler, seq[5])
	}

	str, seq, err = item.Consume("i am")
	assert.Equal("i am", str)
	assert.Nil(err)
	assert.Empty(seq)

	item = Item{sequence, 2, 2}

	_, _, err = item.Consume("i am an antler")
	assert.Equal(PrefixOnly, err)

	_, _, err = item.Consume("i am an antler i am")
	assert.Equal(PrefixOnly, err)

	_, _, err = item.Consume("i am an aardvark")
	assert.Equal(NoMatch, err)

	_, _, err = item.Consume("i am an antler i am an aard")
	assert.Equal(NoMatch, err)
}

func TestAlternative_Consume(t *testing.T) {
	assert := assert.New(t)

	iAmAn := Item{Sequence{Token("i am an")}, 1, 1}
	antler := Item{Sequence{Token("antler")}, 1, 1}
	aardvark := Item{Sequence{Token("aardvark")}, 1, 1}

	all := Sequence{iAmAn, Alternative{antler, aardvark}}

	str, seq, err := all.Consume("i am an antler")

	assert.Empty(str)
	assert.Nil(err)
	if assert.Len(seq, 2) {
		assert.Equal("i am an", string(seq[0].(Token)))
		assert.Equal("antler", string(seq[1].(Token)))
	}

	str, seq, err = all.Consume("i am an aardvark")

	assert.Empty(str)
	assert.Nil(err)
	if assert.Len(seq, 2) {
		assert.Equal("i am an", string(seq[0].(Token)))
		assert.Equal("aardvark", string(seq[1].(Token)))
	}

	_, _, err = all.Consume("i am an")
	assert.Equal(PrefixOnly, err)

	_, _, err = all.Consume("i am an ant")
	assert.Equal(PrefixOnly, err)

	_, _, err = all.Consume("i am an ape")
	assert.Equal(NoMatch, err)
}
