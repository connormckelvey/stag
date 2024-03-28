package stag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserGrammar(t *testing.T) {
	type grammar struct {
		Namespace string `grammar:"grammar,_,namespace"`
		Position  string `grammar:"grammar,_,position"`
	}

	type tests struct {
		Keyed any `grammar.namespace:"foo" grammar.position:"bar"`
		Args  any `grammar:"foo,bar"`
		Mixed any `grammar:"foo" grammar.position:"bar"`
	}

	p, err := New(grammar{})
	assert.NoError(t, err)

	tags, err := p.Parse(tests{})
	assert.NoError(t, err)

	assert.Greater(t, len(tags), 0)
	for _, tag := range tags {
		assert.Equal(t, "foo", tag.Namespace)
		assert.Equal(t, "bar", tag.Position)
	}
}

func TestParserPointers(t *testing.T) {
	type grammar struct {
		First *string `grammar:"name,_,first"`
		Last  *string `grammar:"name,_,last"`
	}

	pp, err := New(grammar{})
	assert.NoError(t, err)

	type tests struct {
		Test any `name:"robert,paulson"`
	}

	nameTag, err := pp.ParseFieldByName(tests{}, "Test")
	assert.NoError(t, err)

	assert.Equal(t, "robert", *nameTag.First)
	assert.Equal(t, "paulson", *nameTag.Last)

}

func TestParserMatrix(t *testing.T) {
	type grammar struct {
		First string `grammar:"name,_,first"`
		Last  string `grammar:"name,_,last"`
	}

	np, err := New(grammar{})
	assert.NoError(t, err)

	tests := []struct {
		name  string
		st    any
		first string
		last  string
	}{
		{
			"Test0",
			struct {
				Test any `name:""`
			}{},
			"", "",
		},
		{
			"Test1",
			struct {
				Test any `name:","`
			}{},
			"", "",
		},
		{
			"Test2",
			struct {
				Test any `name:"robert"`
			}{},
			"robert", "",
		},
		{
			"Test3",
			struct {
				Test any `name:",paulson"`
			}{},
			"", "paulson",
		},
		{
			"Test4",
			struct {
				Test any `name:"robert,paulson"`
			}{},
			"robert", "paulson",
		},
		{
			"Test5",
			struct {
				Test any `name.first:"robert"`
			}{},
			"robert", "",
		},
		{
			"Test6",
			struct {
				Test any `name.last:"paulson"`
			}{},
			"", "paulson",
		},
		{
			"Test7",
			struct {
				Test any `name:"" name.first:"robert"`
			}{},
			"robert", "",
		},
		{
			"Test8",
			struct {
				Test any `name:"" name.last:"paulson"`
			}{},
			"", "paulson",
		},
		{
			"Test9",
			struct {
				Test any `name.first:"robert" name:""`
			}{},
			"robert", "",
		},
		{
			"Test10",
			struct {
				Test any `name.last:"paulson" name:""`
			}{},
			"", "paulson",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := np.ParseFieldByName(test.st, "Test")
			assert.NoError(t, err)

			assert.Equal(t, test.first, actual.First)
			assert.Equal(t, test.last, actual.Last)
		})
	}
}

func TestParserMultipleGrammars(t *testing.T) {
	type grammar struct {
		First string `grammar:"name,_,first"`
		Last  string `grammar:"name,_,last"`
	}

	np, err := New(grammar{})
	assert.NoError(t, err)

	type testExpectTag struct {
		First string `grammar:"expect,,first"`
		Last  string `grammar:"expect,,last"`
	}

	ep, err := New(testExpectTag{})
	assert.NoError(t, err)

	type tests struct {
		Test0  any `expect.first:"" expect.last:"" name:""`
		Test1  any `expect.first:"" expect.last:"" name:","`
		Test2  any `expect.first:"robert" expect.last:"" name:"robert"`
		Test3  any `expect.first:"" expect.last:"paulson" name:",paulson"`
		Test4  any `expect.first:"robert" expect.last:"paulson" name:"robert,paulson"`
		Test5  any `expect.first:"robert" expect.last:"" name.first:"robert"`
		Test6  any `expect.first:"" expect.last:"paulson" name.last:"paulson"`
		Test7  any `expect.first:"robert" expect.last:"" name:"" name.first:"robert"`
		Test8  any `expect.first:"" expect.last:"paulson" name:"" name.last:"paulson"`
		Test9  any `expect.first:"robert" expect.last:"" name.first:"robert" name:""`
		Test10 any `expect.first:"" expect.last:"paulson" name.last:"paulson" name:""`
	}

	tags, err := np.Parse(tests{})
	assert.NoError(t, err)
	assert.Greater(t, len(tags), 0)

	for i, actual := range tags {
		expect, err := ep.ParseField(tests{}, i)
		assert.NoError(t, err)

		assert.Equal(t, expect.First, actual.First)
		assert.Equal(t, expect.Last, actual.Last)
	}
}
