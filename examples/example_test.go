package examples

import (
	"encoding/json"
	"testing"

	"github.com/connormckelvey/tagram"
	"github.com/stretchr/testify/assert"
)

type FlagGrammar struct {
	Name    string   `grammar:"flag,0,name"`
	Aliases []string `grammar:"flag,1,aliases"`
	Usage   string   `grammar:"flag,2,usage"`
}

var FlagParser = tagram.MustGenerate[FlagGrammar]()

func TestGenerate(t *testing.T) {
	type MyFlags struct {
		Props   any `flag:"props,p,Load the props file"`
		Include any `flag:"include,i,Specify file or glob"`
		Beep    any `flag.usage:"wow" flag.aliases:"foo;bar"`
	}

	result, err := FlagParser.Parse(MyFlags{})
	assert.NoError(t, err)

	// spew.Dump(result)
	jj, _ := json.MarshalIndent(result, "", "\t")
	t.Log(string(jj))

}
