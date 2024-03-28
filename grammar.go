package tagram

import (
	"reflect"
	"strings"
)

const (
	tagName = "grammar"
	sep     = ","
)

type positionType int

const (
	positionTypeUnknown positionType = iota
	positionTypeEmpty
	positionTypeExplicit
	positionTypeInferred
)

type parseContext struct {
	Type  reflect.Type
	Field reflect.StructField
}

type parseFunc func(pctx *parseContext, literal string, target *rule) error

func parseNamespace(pctx *parseContext, literal string, target *rule) error {
	return parsePrimitive(literal, &target.Namespace)
}

func parsePosition(pctx *parseContext, literal string, target *rule) error {
	switch literal {
	case "":
		target.PositionType = positionTypeEmpty
	case "_":
		target.PositionType = positionTypeInferred
		target.Position = pctx.Field.Index[0]
	default:
		target.PositionType = positionTypeExplicit
		return parsePrimitive(literal, &target.Position)
	}
	return nil
}

func parseKeys(pctx *parseContext, literal string, target *rule) error {
	target.Keys = append(target.Keys, literal)
	return nil
}

func parseGrammar(structType reflect.Type) (results []*rule, err error) {
	parsers := []parseFunc{parseNamespace, parsePosition, parseKeys}
	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)
		rawTagValue, ok := fieldType.Tag.Lookup(tagName)
		rawTagValue = strings.TrimSpace(rawTagValue)
		if !ok || rawTagValue == "" {
			return nil, nil
		}
		tagTokens := strings.Split(rawTagValue, sep)
		if len(tagTokens) < 1 {
			return nil, nil
		}

		var target rule
		for j, token := range tagTokens {
			var parse parseFunc
			// repeatedly use last parse if num tokens > num parsers
			if j >= len(parsers) {
				parse = parsers[len(parsers)-1]
			} else {
				parse = parsers[j]
			}
			ctx := &parseContext{
				Type:  structType,
				Field: fieldType,
			}
			err := parse(ctx, token, &target)
			if err != nil {
				return nil, err
			}
		}
		results = append(results, &target)
	}
	return results, nil
}
