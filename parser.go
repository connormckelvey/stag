package stag

import (
	"fmt"
	"reflect"
)

// A Stag represents a parser for the provided grammar.
type Stag[G any] struct {
	parsers []*ruleParser[G]
}

// MustGenerate is like New but panics if the grammar cannot be parsed.
func MustGenerate[G any](grammar G) *Stag[G] {
	pg, err := New(grammar)
	if err != nil {
		panic(err)
	}
	return pg
}

// New parses the input grammar `G` into rules ready for parsing.
func New[G any](grammar G) (*Stag[G], error) {
	parser := &Stag[G]{}
	rules, err := parseGrammar(grammar)
	if err != nil {
		return nil, err
	}
	parser.parsers = make([]*ruleParser[G], len(rules))
	for i, rule := range rules {
		parser.parsers[i] = newRuleParser[G](rule, i)
	}
	return parser, nil
}

// ParseField parses the struct tags from `taggedStructs` i'th field using the configured grammar.
func (s *Stag[G]) ParseField(taggedStruct any, i int) (g G, err error) {
	structType, err := requireStruct(taggedStruct)
	if err != nil {
		return g, err
	}
	return s.parseField(structType, i)
}

// ParseFieldByName is similar to ParseField but takes a field name instead of an index.
func (s *Stag[G]) ParseFieldByName(taggedStruct any, name string) (g G, err error) {
	structType, err := requireStruct(taggedStruct)
	if err != nil {
		return g, err
	}
	sf, ok := structType.FieldByName(name)
	if !ok {
		return g, fmt.Errorf("field '%s' not found", name)
	}
	return s.parseField(structType, sf.Index[0])
}

// Parse parses every field in the struct returning a slice containing output for each field in order.
func (s *Stag[G]) Parse(taggedStruct any) ([]G, error) {
	structType, err := requireStruct(taggedStruct)
	if err != nil {
		return nil, err
	}
	results := make([]G, structType.NumField())
	for i := 0; i < structType.NumField(); i++ {
		g, err := s.parseField(structType, i)
		if err != nil {
			return nil, err
		}
		results[i] = g
	}
	return results, nil
}

func (s *Stag[G]) parseField(structType reflect.Type, i int) (g G, err error) {
	inputField := structType.Field(i)
	fieldTag := inputField.Tag
	for _, parser := range s.parsers {
		err := parser.parse(&fieldTag, &g)
		if err != nil {
			return g, err
		}
	}
	return g, nil
}
