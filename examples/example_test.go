package examples

import (
	"fmt"
	"testing"

	"github.com/connormckelvey/stag"
	"github.com/stretchr/testify/assert"
)

type FlagGrammar struct {
	Name    string   `grammar:"flag,0,name"`
	Aliases []string `grammar:"flag,1,aliases"`
	Usage   string   `grammar:"flag,2,usage"`
}

var FlagParser = stag.MustGenerate(FlagGrammar{})

type MyFlags struct {
	Props   any `flag:"props,p,Load the props file"`
	Include any `flag:"include,i,Specify file or glob"`
}

func TestGenerate(t *testing.T) {
	results, err := FlagParser.Parse(MyFlags{})
	assert.NoError(t, err)

	for _, tagValues := range results {
		fmt.Printf("%+v\n", tagValues)
	}
}
