/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type ValueKind string

const (
	ValueLiteral ValueKind = "literal"
	ValueState   ValueKind = "state"
	ValueContext ValueKind = "context"
	ValueSource  ValueKind = "source"
	ValueEvent   ValueKind = "event"
	ValueResult  ValueKind = "result"
	ValueExpr    ValueKind = "expression"
)

type Value struct {
	Kind ValueKind
	Path string
	Data any
}

var expressionArity = map[string][2]int{
	"$eq":       {2, 2},
	"$ne":       {2, 2},
	"$gt":       {2, 2},
	"$gte":      {2, 2},
	"$lt":       {2, 2},
	"$lte":      {2, 2},
	"$and":      {1, -1},
	"$or":       {1, -1},
	"$not":      {1, 1},
	"$in":       {2, 2},
	"$exists":   {1, 1},
	"$concat":   {0, -1},
	"$coalesce": {1, -1},
	"$if":       {3, 3},
}

func Literal(data any) Value {
	return Value{Kind: ValueLiteral, Data: data}
}

func StateRef(path string) Value {
	return Value{Kind: ValueState, Path: path}
}

func ContextRef(path string) Value {
	return Value{Kind: ValueContext, Path: path}
}

func SourceRef(path string) Value {
	return Value{Kind: ValueSource, Path: path}
}

func EventRef(path string) Value {
	return Value{Kind: ValueEvent, Path: path}
}

func ResultRef(path string) Value {
	return Value{Kind: ValueResult, Path: path}
}

func Expr(operator string, args ...Value) Value {
	return Value{Kind: ValueExpr, Path: operator, Data: append([]Value(nil), args...)}
}

func (v Value) Validate() error {
	switch v.Kind {
	case ValueLiteral:
		return nil

	case ValueState, ValueContext, ValueSource, ValueEvent, ValueResult:
		if !safePath(v.Path) {
			return fmt.Errorf("%w: unsafe reference path %q", ErrInvalidValue, v.Path)
		}

		return nil

	case ValueExpr:
		arity, ok := expressionArity[v.Path]
		if !ok {
			return fmt.Errorf("%w: unsupported expression operator %q", ErrInvalidValue, v.Path)
		}

		if v.Path == "" {
			return fmt.Errorf("%w: expression operator is required", ErrInvalidValue)
		}

		args, ok := v.Data.([]Value)
		if !ok {
			return fmt.Errorf("%w: expression arguments have invalid type", ErrInvalidValue)
		}

		if len(args) < arity[0] || (arity[1] >= 0 && len(args) > arity[1]) {
			return fmt.Errorf("%w: operator %s expects %d..%d arguments", ErrInvalidValue, v.Path, arity[0], arity[1])
		}

		for _, arg := range args {
			if err := arg.Validate(); err != nil {
				return err
			}
		}

		return nil

	default:
		return fmt.Errorf("%w: unknown value kind %q", ErrInvalidValue, v.Kind)
	}
}

func (v Value) MarshalJSON() ([]byte, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}

	if v.Kind == ValueLiteral {
		return json.Marshal(v.Data)
	}

	if v.Kind == ValueExpr {
		args, ok := v.Data.([]Value)
		if !ok {
			return nil, fmt.Errorf("%w: expression arguments have invalid type", ErrInvalidValue)
		}

		return json.Marshal(map[string]any{v.Path: args})
	}

	return json.Marshal(map[string]string{"$" + string(v.Kind): v.Path})
}

func (v *Value) UnmarshalJSON(data []byte) error {
	if v == nil {
		return fmt.Errorf("%w: nil destination", ErrInvalidValue)
	}

	var object map[string]json.RawMessage
	if err := json.Unmarshal(data, &object); err == nil && object != nil {
		if len(object) == 1 {
			for key, raw := range object {
				if key == "$state" || key == "$context" || key == "$source" || key == "$event" || key == "$result" {
					var path string
					if err := json.Unmarshal(raw, &path); err != nil {
						return fmt.Errorf("%w: reference path: %w", ErrInvalidValue, err)
					}

					v.Kind = ValueKind(key[1:])
					v.Path = path
					v.Data = nil
					return v.Validate()
				}

				var args []Value
				if len(key) > 0 && key[0] == '$' && json.Unmarshal(raw, &args) == nil {
					v.Kind = ValueExpr
					v.Path = key
					v.Data = args
					return v.Validate()
				}
			}
		}
	}

	var literal any
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&literal); err != nil {
		return fmt.Errorf("%w: literal: %w", ErrInvalidValue, err)
	}

	v.Kind = ValueLiteral
	v.Path = ""
	v.Data = literal

	return nil
}

func safePath(path string) bool {
	if path == "" {
		return false
	}

	for _, part := range splitPath(path) {
		if part == "__proto__" || part == "prototype" || part == "constructor" || part == "" {
			return false
		}
	}

	return true
}

func splitPath(path string) []string {
	parts := make([]string, 0, 4)
	start := 0

	for i := 0; i <= len(path); i++ {
		if i == len(path) || path[i] == '.' {
			parts = append(parts, path[start:i])
			start = i + 1
		}
	}

	return parts
}
