/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// ValueKind identifies how a Value obtains its data.
type ValueKind string

// Value kind constants identify the supported value representations.
const (
	// ValueLiteral stores data directly in the value.
	ValueLiteral ValueKind = "literal"
	// ValueState references view state.
	ValueState ValueKind = "state"
	// ValueContext references host context.
	ValueContext ValueKind = "context"
	// ValueSource references source data.
	ValueSource ValueKind = "source"
	ValueEvent  ValueKind = "event"
	ValueResult ValueKind = "result"
	ValueExpr   ValueKind = "expression"
)

// Value is a literal, reference, or expression used by the declarative UI model.
//
//nolint:recvcheck // JSON values intentionally marshal by value and unmarshal by pointer.
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

// Literal creates a value containing data directly.
func Literal(data any) Value {
	return Value{Kind: ValueLiteral, Data: data}
}

// StateRef creates a value referencing view state at path.
func StateRef(path string) Value {
	return Value{Kind: ValueState, Path: path}
}

// ContextRef creates a value referencing host context at path.
func ContextRef(path string) Value {
	return Value{Kind: ValueContext, Path: path}
}

// SourceRef creates a value referencing source data at path.
func SourceRef(path string) Value {
	return Value{Kind: ValueSource, Path: path}
}

// EventRef creates a value referencing the current event at path.
func EventRef(path string) Value {
	return Value{Kind: ValueEvent, Path: path}
}

// ResultRef creates a value referencing the previous action result at path.
func ResultRef(path string) Value {
	return Value{Kind: ValueResult, Path: path}
}

// Expr creates a value representing an expression operator and its arguments.
func Expr(operator string, args ...Value) Value {
	return Value{Kind: ValueExpr, Path: operator, Data: append([]Value(nil), args...)}
}

// Validate checks that the value kind, path, and expression arguments are valid.
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

// MarshalJSON encodes the value in the declarative UI wire format.
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

func (v *Value) unmarshalSpecial(data []byte) (bool, error) {
	var object map[string]json.RawMessage
	if err := json.Unmarshal(data, &object); err != nil {
		return false, fmt.Errorf("%w: %s", ErrInvalidValue, err.Error())
	}

	if len(object) != 1 {
		return false, fmt.Errorf("%w: expected 1 object, got %d", ErrInvalidValue, len(object))
	}

	for key, raw := range object {
		return v.unmarshalSpecialEntry(key, raw)
	}

	return false, nil
}

func (v *Value) unmarshalSpecialEntry(key string, raw json.RawMessage) (bool, error) {
	switch key {
	case "$state", "$context", "$source", "$event", "$result":
		var path string
		if err := json.Unmarshal(raw, &path); err != nil {
			return true, fmt.Errorf("%w: reference path: %w", ErrInvalidValue, err)
		}

		v.Kind = ValueKind(key[1:])
		v.Path = path
		v.Data = nil
		return true, v.Validate()
	}

	if len(key) == 0 || key[0] != '$' {
		return false, fmt.Errorf("%w: reference path %q", ErrInvalidValue, key)
	}

	var args []Value
	if err := json.Unmarshal(raw, &args); err != nil {
		return false, fmt.Errorf("%w: reference arguments have invalid type", ErrInvalidValue)
	}

	v.Kind = ValueExpr
	v.Path = key
	v.Data = args

	return true, v.Validate()
}

// UnmarshalJSON decodes a literal, reference, or expression value.
func (v *Value) UnmarshalJSON(data []byte) error {
	if v == nil {
		return fmt.Errorf("%w: nil destination", ErrInvalidValue)
	}

	handled, err := v.unmarshalSpecial(data)
	if handled {
		return err
	}

	var literal any
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err = decoder.Decode(&literal); err != nil {
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

	for _, part := range strings.Split(path, ".") {
		if part == "__proto__" || part == "prototype" || part == "constructor" || part == "" {
			return false
		}
	}

	return true
}
