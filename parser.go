package tagram

import (
	"fmt"
	"reflect"

	"github.com/davecgh/go-spew/spew"
)

type Parser[G any] struct {
	parsers []*ruleParser[G]
}

func MustGenerate[G any]() *Parser[G] {
	pg, err := NewParser[G]()
	if err != nil {
		panic(err)
	}
	return pg
}

func NewParser[G any]() (*Parser[G], error) {
	structType, err := requireStruct(identity[G]())
	if err != nil {
		return nil, err
	}

	parser := &Parser[G]{}
	rules, err := parseGrammar(structType)
	if err != nil {
		return nil, err
	}
	spew.Dump(rules)
	parser.parsers = make([]*ruleParser[G], len(rules))
	for i, rule := range rules {
		parser.parsers[i] = newRuleParser[G](rule, i)
	}
	return parser, nil
}

func (p *Parser[G]) ParseField(taggedStruct any, i int) (*G, error) {
	structType, err := requireStruct(taggedStruct)
	if err != nil {
		return nil, err
	}
	return p.parseField(structType, i)
}

func (p *Parser[G]) ParseFieldByName(taggedStruct any, name string) (*G, error) {
	structType, err := requireStruct(taggedStruct)
	if err != nil {
		return nil, err
	}
	sf, ok := structType.FieldByName(name)
	if !ok {
		return nil, fmt.Errorf("field '%s' not found", name)
	}
	return p.parseField(structType, sf.Index[0])
}

func (p *Parser[G]) Parse(taggedStruct any) ([]*G, error) {
	structType, err := requireStruct(taggedStruct)
	if err != nil {
		return nil, err
	}
	results := make([]*G, structType.NumField())
	for i := 0; i < structType.NumField(); i++ {
		g, err := p.parseField(structType, i)
		if err != nil {
			return nil, err
		}
		results[i] = g
	}
	return results, nil
}

func (p *Parser[G]) parseField(structType reflect.Type, i int) (*G, error) {
	inputField := structType.Field(i)
	var g G
	for _, parser := range p.parsers {
		err := parser.parse(inputField, &g)
		if err != nil {
			return nil, err
		}
	}
	return &g, nil
}
