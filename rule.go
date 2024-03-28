package tagram

import (
	"reflect"
	"strings"
)

type rule struct {
	Namespace    string
	PositionType positionType
	Position     int
	Keys         []string
}

type ruleParser[G any] struct {
	rule    *rule
	ruleIdx int
}

func newRuleParser[G any](rule *rule, ruleIdx int) *ruleParser[G] {
	return &ruleParser[G]{
		rule:    rule,
		ruleIdx: ruleIdx,
	}
}

func (rp *ruleParser[G]) parseNamespace(inputField reflect.StructField, target *G) error {
	// dereferenced value
	targetValue := reflect.ValueOf(target).Elem()
	targetFieldValue := targetValue.Field(rp.ruleIdx)

	namespaceTagLit, ok := inputField.Tag.Lookup(rp.rule.Namespace)
	if !ok {
		return nil
	}

	tokens := strings.Split(namespaceTagLit, ",")
	var pos int = -1
	switch rp.rule.PositionType {
	case positionTypeInferred, positionTypeExplicit:
		pos = rp.rule.Position
	}
	if pos < 0 || pos >= len(tokens) {
		return nil
	}
	token := tokens[pos]
	return parseInto(token, targetFieldValue)
}
func (rp *ruleParser[G]) parseKeys(inputField reflect.StructField, target *G) error {
	targetValue := reflect.ValueOf(target).Elem()
	targetFieldValue := targetValue.Field(rp.ruleIdx)

	// Find tag strings for keyed tag names (namespace.<key>)
	for _, key := range rp.rule.Keys {
		keyed := strings.Join([]string{rp.rule.Namespace, key}, ".")
		keyedTagLit, ok := inputField.Tag.Lookup(keyed)
		if !ok {
			continue
		}
		if err := parseInto(keyedTagLit, targetFieldValue); err != nil {
			return err
		}
	}
	return nil
}

func (rp *ruleParser[G]) parse(inputField reflect.StructField, target *G) error {
	if rp.rule == nil {
		return nil
	}
	if err := rp.parseNamespace(inputField, target); err != nil {
		return err
	}
	if err := rp.parseKeys(inputField, target); err != nil {
		return err
	}
	return nil
}
