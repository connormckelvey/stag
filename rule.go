package stag

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

func (rp *ruleParser[G]) parseNamespace(fieldTag *reflect.StructTag, target *G) error {
	namespaceTagLit, ok := fieldTag.Lookup(rp.rule.Namespace)
	if !ok {
		return nil
	}

	targetValue := reflect.ValueOf(target).Elem()
	targetFieldValue := targetValue.Field(rp.ruleIdx)

	tokens := strings.Split(namespaceTagLit, ",")
	var pos int = -1
	switch rp.rule.PositionType {
	case positionTypeExplicit, positionTypeInferred:
		pos = rp.rule.Position
	}
	if pos < 0 || pos >= len(tokens) {
		return nil
	}
	token := tokens[pos]
	return parseInto(token, targetFieldValue)
}
func (rp *ruleParser[G]) parseKeys(fieldTag *reflect.StructTag, target *G) error {
	targetValue := reflect.ValueOf(target).Elem()
	targetFieldValue := targetValue.Field(rp.ruleIdx)

	for _, key := range rp.rule.Keys {
		keyed := strings.Join([]string{rp.rule.Namespace, key}, ".")
		keyedTagLit, ok := fieldTag.Lookup(keyed)
		if !ok {
			continue
		}
		if err := parseInto(keyedTagLit, targetFieldValue); err != nil {
			return err
		}
	}
	return nil
}

func (rp *ruleParser[G]) parse(fieldTag *reflect.StructTag, target *G) error {
	if rp.rule == nil {
		return nil
	}
	if err := rp.parseNamespace(fieldTag, target); err != nil {
		return err
	}
	if err := rp.parseKeys(fieldTag, target); err != nil {
		return err
	}
	return nil
}
