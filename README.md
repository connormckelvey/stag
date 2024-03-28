# TaGram - Declarative Struct Tag Parser Generator

TaGram is an open-source Go library designed to simplify the parsing of struct tags by generating a parser based on a declarative grammar. TaGram allows developers to define the structure and rules for parsing directly within their Go structs, streamlining the development process.

## Get Started

### Install
```bash
go get github.com/connormckelvey/tagram
```

### Import
```go
import "github.com/connormckelvey/tagram"
```

### Define grammar
```go
type CliFlagGrammar struct {
	Name    string   `grammar:"flag,0,name"`
	Aliases []string `grammar:"flag,1,aliases"`
	Usage   string   `grammar:"flag,2,usage"`
}
```

### Generate parser
```go
var FlagParser = tagram.MustGenerate[CliFlagGrammar]()
```

### Parse
```go
type MyFlags struct {
    Props   any `flag:"props,p,Load the props file"`
    Include any `flag:"include,i,Specify file or glob"`
    Beep    any `flag.usage:"wow" flag.aliases:"foo;bar"`
}

result, err := FlagParser.Parse(MyFlags{})
if err != nil {
    panic(err)
}

spew.Dump(result)
```