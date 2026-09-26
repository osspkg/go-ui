/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package ui

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	defaultMaxBytes           = 1 << 20
	defaultMaxNodes           = 1000
	defaultMaxDepth           = 32
	defaultMaxPropsPerNode    = 64
	defaultMaxActions         = 128
	defaultMaxSources         = 64
	defaultMaxStrings         = 10000
	defaultMaxStringLength    = 16 << 10
	defaultMaxExpressionDepth = 32
	defaultMaxArrayItems      = 1000
	defaultMaxObjectKeys      = 256
	defaultMaxPathLength      = 1024

	actionTypeTool          = "tool"
	actionTypeSetState      = "set-state"
	actionTypeMergeState    = "merge-state"
	actionTypeRefreshSource = "refresh-source"
	actionTypeInvalidate    = "invalidate"
	actionStepTypeAction    = "action"
)

// Limits bounds the size and complexity of a view schema.
type Limits struct {
	MaxBytes           int
	MaxNodes           int
	MaxDepth           int
	MaxPropsPerNode    int
	MaxActions         int
	MaxSources         int
	MaxStrings         int
	MaxStringLength    int
	MaxExpressionDepth int
	MaxArrayItems      int
	MaxObjectKeys      int
	MaxPathLength      int
}

// DefaultLimits returns the default schema validation limits.
func DefaultLimits() Limits {
	return Limits{
		MaxBytes:           defaultMaxBytes,
		MaxNodes:           defaultMaxNodes,
		MaxDepth:           defaultMaxDepth,
		MaxPropsPerNode:    defaultMaxPropsPerNode,
		MaxActions:         defaultMaxActions,
		MaxSources:         defaultMaxSources,
		MaxStrings:         defaultMaxStrings,
		MaxStringLength:    defaultMaxStringLength,
		MaxExpressionDepth: defaultMaxExpressionDepth,
		MaxArrayItems:      defaultMaxArrayItems,
		MaxObjectKeys:      defaultMaxObjectKeys,
		MaxPathLength:      defaultMaxPathLength,
	}
}

// Validate checks the schema against the default validation limits.
func (v ViewSchema) Validate() error {
	return v.ValidateWithLimits(DefaultLimits())
}

// ValidateWithLimits checks the schema against caller-provided limits.
func (v ViewSchema) ValidateWithLimits(limits Limits) error {
	limits = withDefaultLimits(limits)

	if v.ProtocolVersion != ProtocolVersion {
		return fmt.Errorf("%w: unsupported view protocol version %q", ErrInvalidSchema, v.ProtocolVersion)
	}

	if !validName(v.ID) {
		return fmt.Errorf("%w: view id is required", ErrInvalidSchema)
	}

	if len(v.Actions) > limits.MaxActions || len(v.Sources) > limits.MaxSources {
		return fmt.Errorf("%w: source or action limit exceeded", ErrInvalidSchema)
	}

	if v.Regions.TopHeader == nil ||
		v.Regions.LeftPanel == nil ||
		v.Regions.Content == nil ||
		v.Regions.RightPanel == nil ||
		v.Regions.Bottom == nil {
		return fmt.Errorf("%w: all regions must be arrays", ErrInvalidSchema)
	}

	if err := validateSerializedLimits(v, limits); err != nil {
		return err
	}

	for name, source := range v.Sources {
		if !validName(name) || source.Type != actionTypeTool || !validName(source.Tool) {
			return fmt.Errorf("%w: invalid source %q", ErrInvalidSchema, name)
		}

		if source.Policy != "" && source.Policy != SourceManual && source.Policy != SourceOnMount {
			return fmt.Errorf("%w: source %q has unsupported policy", ErrInvalidSchema, name)
		}

		if source.Cache != nil && source.Cache.TTL < 0 {
			return fmt.Errorf("%w: source %q has negative cache ttl", ErrInvalidSchema, name)
		}

		if err := validateValues(source.Input, limits, v); err != nil {
			return fmt.Errorf("%w: source %q input: %w", ErrInvalidSchema, name, err)
		}

		for _, action := range source.RefreshOn {
			if _, ok := v.Actions[action]; !ok {
				return fmt.Errorf("%w: source %q refresh action %q is not defined", ErrInvalidSchema, name, action)
			}
		}
	}

	for name, action := range v.Actions {
		if !validName(name) || !validActionType(action.Type) {
			return fmt.Errorf("%w: invalid action %q", ErrInvalidSchema, name)
		}

		if err := validateAction(action, limits, v); err != nil {
			return fmt.Errorf("%w: action %q: %w", ErrInvalidSchema, name, err)
		}
	}

	context := validationContext{
		limits: limits,
		view:   v,
		seen:   make(map[string]struct{}),
	}

	for region, list := range map[string][]Node{
		"top-header":  v.Regions.TopHeader,
		"left-panel":  v.Regions.LeftPanel,
		"content":     v.Regions.Content,
		"right-panel": v.Regions.RightPanel,
		"bottom":      v.Regions.Bottom,
	} {
		for _, n := range list {
			if err := context.validateNode(n, region, 1); err != nil {
				return err
			}
		}
	}

	return nil
}

func withDefaultLimits(limits Limits) Limits {
	defaults := DefaultLimits()

	if limits.MaxBytes <= 0 {
		limits.MaxBytes = defaults.MaxBytes
	}

	if limits.MaxNodes <= 0 {
		limits.MaxNodes = defaults.MaxNodes
	}

	if limits.MaxDepth <= 0 {
		limits.MaxDepth = defaults.MaxDepth
	}

	if limits.MaxPropsPerNode <= 0 {
		limits.MaxPropsPerNode = defaults.MaxPropsPerNode
	}

	if limits.MaxActions <= 0 {
		limits.MaxActions = defaults.MaxActions
	}

	if limits.MaxSources <= 0 {
		limits.MaxSources = defaults.MaxSources
	}

	if limits.MaxStrings <= 0 {
		limits.MaxStrings = defaults.MaxStrings
	}

	if limits.MaxStringLength <= 0 {
		limits.MaxStringLength = defaults.MaxStringLength
	}

	if limits.MaxExpressionDepth <= 0 {
		limits.MaxExpressionDepth = defaults.MaxExpressionDepth
	}

	if limits.MaxArrayItems <= 0 {
		limits.MaxArrayItems = defaults.MaxArrayItems
	}

	if limits.MaxObjectKeys <= 0 {
		limits.MaxObjectKeys = defaults.MaxObjectKeys
	}

	if limits.MaxPathLength <= 0 {
		limits.MaxPathLength = defaults.MaxPathLength
	}

	return limits
}

func validateSerializedLimits(schema any, limits Limits) error {
	encoded, err := json.Marshal(schema)
	if err != nil {
		return fmt.Errorf("%w: serialize schema for limits: %w", ErrInvalidSchema, err)
	}

	if len(encoded) > limits.MaxBytes {
		return fmt.Errorf("%w: schema byte limit exceeded", ErrInvalidSchema)
	}

	decoder := json.NewDecoder(bytes.NewReader(encoded))
	decoder.UseNumber()
	var document any

	if err := decoder.Decode(&document); err != nil {
		return fmt.Errorf("%w: inspect schema limits: %w", ErrInvalidSchema, err)
	}

	count := 0
	if err := walkLimitValue(document, 0, limits, &count); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidSchema, err)
	}

	return nil
}

func walkLimitValue(value any, depth int, limits Limits, count *int) error {
	if depth > limits.MaxDepth {
		return errors.New("schema depth limit exceeded")
	}

	switch typed := value.(type) {
	case string:
		if len(typed) > limits.MaxStringLength {
			return errors.New("string length limit exceeded")
		}

		*count++
		if *count > limits.MaxStrings {
			return errors.New("string count limit exceeded")
		}

	case []any:
		if len(typed) > limits.MaxArrayItems {
			return errors.New("array item limit exceeded")
		}

		for _, item := range typed {
			if err := walkLimitValue(item, depth+1, limits, count); err != nil {
				return err
			}
		}

	case map[string]any:
		if len(typed) > limits.MaxObjectKeys {
			return errors.New("object key limit exceeded")
		}

		for key, item := range typed {
			if key == "__proto__" || key == "prototype" || key == "constructor" {
				return errors.New("unsafe object key")
			}

			if err := walkLimitValue(key, depth+1, limits, count); err != nil {
				return err
			}

			if err := walkLimitValue(item, depth+1, limits, count); err != nil {
				return err
			}
		}

	default:
	}

	return nil
}

type validationContext struct {
	limits Limits
	view   ViewSchema
	seen   map[string]struct{}
	nodes  int
}

func (context *validationContext) validateNode(node Node, region string, depth int) error {
	context.nodes++
	if context.nodes > context.limits.MaxNodes || depth > context.limits.MaxDepth {
		return fmt.Errorf("%w: node limit exceeded", ErrInvalidSchema)
	}

	if !validName(node.ID) || !validComponentName(node.Component) {
		return fmt.Errorf("%w: node id and component are required", ErrInvalidSchema)
	}

	if _, exists := context.seen[node.ID]; exists {
		return fmt.Errorf("%w: duplicate node id %q", ErrInvalidSchema, node.ID)
	}

	context.seen[node.ID] = struct{}{}
	if node.Layout != nil {
		if err := node.Layout.Validate(region); err != nil {
			return err
		}
	}

	if len(node.Props) > context.limits.MaxPropsPerNode {
		return fmt.Errorf("%w: too many props on node %q", ErrInvalidSchema, node.ID)
	}

	if err := validateValues(node.Props, context.limits, context.view); err != nil {
		return fmt.Errorf("%w: node %q props: %w", ErrInvalidSchema, node.ID, err)
	}

	for event, handler := range node.Events {
		if !validName(event) || handler.Action == "" && len(handler.Steps) == 0 {
			return fmt.Errorf("%w: invalid event on node %q", ErrInvalidSchema, node.ID)
		}

		if handler.Action != "" {
			if _, ok := context.view.Actions[handler.Action]; !ok {
				return fmt.Errorf("%w: action %q is not defined", ErrInvalidSchema, handler.Action)
			}
		}

		for _, step := range handler.Steps {
			if step.Type != actionStepTypeAction || step.Action == "" {
				return fmt.Errorf("%w: invalid event step on node %q", ErrInvalidSchema, node.ID)
			}

			if _, ok := context.view.Actions[step.Action]; !ok {
				return fmt.Errorf("%w: action %q is not defined", ErrInvalidSchema, step.Action)
			}

			if err := validateValues(step.Input, context.limits, context.view); err != nil {
				return fmt.Errorf("%w: event step input: %w", ErrInvalidSchema, err)
			}
		}
	}

	if node.When != nil {
		if err := validateValue(*node.When, context.limits, context.view, 0); err != nil {
			return fmt.Errorf("%w: node %q condition: %w", ErrInvalidSchema, node.ID, err)
		}
	}

	for _, child := range node.Children {
		if err := context.validateNode(child, region, depth+1); err != nil {
			return err
		}
	}

	for _, children := range node.Slots {
		for _, child := range children {
			if err := context.validateNode(child, region, depth+1); err != nil {
				return err
			}
		}
	}

	return nil
}

// Validate checks whether the layout fits the specified region.
func (l Layout) Validate(region string) error {
	if l.Row < 1 {
		return fmt.Errorf("%w: row must be positive", ErrInvalidLayout)
	}

	if region == "left-panel" || region == "right-panel" {
		if l.Cols != 0 || l.Offset != 0 {
			return fmt.Errorf("%w: side panels do not accept cols or offset", ErrInvalidLayout)
		}
		return nil
	}

	if l.Cols < 1 || l.Cols > 12 || l.Offset < 0 || l.Offset > 11 || l.Cols+l.Offset > 12 {
		return fmt.Errorf("%w: cols and offset must fit 12 columns", ErrInvalidLayout)
	}

	return nil
}

func validateValues(values map[string]Value, limits Limits, view ViewSchema) error {
	for key, value := range values {
		if len(key) > limits.MaxStringLength {
			return errors.New("property name is too long")
		}

		if err := validateValue(value, limits, view, 0); err != nil {
			return err
		}
	}

	return nil
}

func validateAction(action Action, limits Limits, view ViewSchema) error {
	if action.Type == actionTypeTool && !validName(action.Tool) {
		return errors.New("tool action requires a tool")
	}

	if action.Type == actionTypeSetState || action.Type == actionTypeMergeState {
		if !safePathWithLimit(action.Path, limits.MaxPathLength) {
			return errors.New("state action requires a safe path")
		}
	}

	if action.Type == actionTypeRefreshSource || action.Type == actionTypeInvalidate {
		if _, ok := view.Sources[action.Source]; !ok {
			return fmt.Errorf("source %q is not defined", action.Source)
		}
	}

	if action.Path != "" && !safePathWithLimit(action.Path, limits.MaxPathLength) {
		return errors.New("unsafe action path")
	}

	if err := validateValues(action.Input, limits, view); err != nil {
		return err
	}

	if action.Value != nil {
		if err := validateValue(*action.Value, limits, view, 0); err != nil {
			return err
		}
	}

	for _, effect := range action.Effects {
		if err := validateEffect(effect, limits, view); err != nil {
			return err
		}
	}

	return nil
}

func validateEffect(effect Effect, limits Limits, view ViewSchema) error {
	if effect.Type == "" {
		return errors.New("effect type is required")
	}

	if !validEffectType(effect.Type) {
		return fmt.Errorf("unsupported effect type %q", effect.Type)
	}

	if effect.Path != "" && !safePathWithLimit(effect.Path, limits.MaxPathLength) {
		return errors.New("unsafe effect path")
	}

	if effect.Type == actionTypeSetState || effect.Type == actionTypeMergeState {
		if !safePathWithLimit(effect.Path, limits.MaxPathLength) {
			return errors.New("state effect requires a safe path")
		}
	}

	if effect.Type == actionTypeInvalidate || effect.Type == actionTypeRefreshSource {
		if _, ok := view.Sources[effect.Source]; !ok {
			return fmt.Errorf("source %q is not defined", effect.Source)
		}
	}

	if effect.Type == "toast" && effect.Message == "" {
		return errors.New("toast effect requires a message")
	}

	if effect.Type == "navigate" && effect.To == "" {
		return errors.New("navigate effect requires a target")
	}

	if effect.Value != nil {
		if err := validateValue(*effect.Value, limits, view, 0); err != nil {
			return err
		}
	}

	return nil
}

func validateValue(value Value, limits Limits, view ViewSchema, depth int) error {
	if depth > limits.MaxExpressionDepth {
		return errors.New("expression depth limit exceeded")
	}

	switch value.Kind {
	case ValueLiteral:
		return validateLiteralData(value.Data, limits, depth)

	case ValueState, ValueContext, ValueEvent, ValueResult:
		if !safePathWithLimit(value.Path, limits.MaxPathLength) {
			return fmt.Errorf("%w: unsafe reference path %q", ErrInvalidValue, value.Path)
		}

		return nil

	case ValueSource:
		if !safePathWithLimit(value.Path, limits.MaxPathLength) {
			return fmt.Errorf("%w: unsafe reference path %q", ErrInvalidValue, value.Path)
		}

		if _, ok := sourceForPath(value.Path, view.Sources); !ok {
			return fmt.Errorf("source %q is not defined", value.Path)
		}

		return nil

	case ValueExpr:
		args, ok := value.Data.([]Value)
		if !ok {
			return errors.New("expression arguments have invalid type")
		}

		arity, ok := expressionArity[value.Path]
		if !ok {
			return fmt.Errorf("%w: unsupported expression operator %q", ErrInvalidValue, value.Path)
		}

		if len(args) < arity[0] || (arity[1] >= 0 && len(args) > arity[1]) {
			return fmt.Errorf(
				"%w: operator %s expects %d..%d arguments",
				ErrInvalidValue,
				value.Path,
				arity[0],
				arity[1],
			)
		}

		for _, arg := range args {
			if err := validateValue(arg, limits, view, depth+1); err != nil {
				return err
			}
		}

		return nil

	default:
		return fmt.Errorf("%w: unknown value kind %q", ErrInvalidValue, value.Kind)
	}
}

func sourceForPath(path string, sources map[string]DataSource) (string, bool) {
	for name := path; name != ""; {
		if _, ok := sources[name]; ok {
			return name, true
		}

		index := strings.LastIndexByte(name, '.')
		if index < 0 {
			break
		}

		name = name[:index]
	}

	return "", false
}

func validActionType(value string) bool {
	switch value {
	case actionTypeTool, actionTypeSetState, actionTypeMergeState, actionTypeRefreshSource, actionTypeInvalidate:
		return true
	default:
		return false
	}
}

func validEffectType(value string) bool {
	switch value {
	case actionTypeSetState, actionTypeMergeState, actionTypeInvalidate, actionTypeRefreshSource,
		"refresh-view", "patch-view", "navigate", "toast",
		"dialog", "close-dialog":
		return true
	default:
		return false
	}
}

func validComponentName(value string) bool {
	if value == "" {
		return false
	}

	for index, char := range value {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9' && index > 0) || char == '-' {
			continue
		}

		return false
	}

	return true
}
