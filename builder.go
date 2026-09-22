/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package ui

import (
	"encoding/json"
	"fmt"
)

type App struct {
	manifest Manifest
	views    map[string]ViewSchema
}

type AppOption func(*App)

func AppID(value string) AppOption {
	return func(app *App) {
		app.manifest.Plugin.ID = value
	}
}

func AppTitle(value string) AppOption {
	return func(app *App) {
		app.manifest.Plugin.Title = value
	}
}

func AppVersion(value string) AppOption {
	return func(app *App) {
		app.manifest.Plugin.Version = value
	}
}

func New(options ...AppOption) *App {
	app := &App{
		manifest: Manifest{
			ProtocolVersion: ProtocolVersion,
			Views:           []UIViewDescriptor{},
		},
		views: make(map[string]ViewSchema),
	}

	for _, option := range options {
		option(app)
	}

	return app
}

func (app *App) AddView(view *ViewBuilder) error {
	schema := view.Schema()
	if err := schema.Validate(); err != nil {
		return err
	}

	if _, exists := app.views[schema.ID]; exists {
		return ErrInvalidSchema
	}

	snapshot, err := cloneViewSchema(schema)
	if err != nil {
		return fmt.Errorf("%w: copy view schema: %w", ErrInvalidSchema, err)
	}

	app.views[schema.ID] = snapshot
	app.manifest.Views = append(
		app.manifest.Views,
		UIViewDescriptor{
			ID:     schema.ID,
			Title:  schema.Title,
			Schema: "ui://views/" + schema.ID,
		},
	)

	return nil
}

func (app *App) Manifest() Manifest {
	return cloneManifest(app.manifest)
}

func (app *App) ManifestJSON() ([]byte, error) {
	return json.Marshal(app.Manifest())
}

func (app *App) View(id string) (ViewSchema, bool) {
	view, ok := app.views[id]
	if !ok {
		return ViewSchema{}, false
	}

	snapshot, err := cloneViewSchema(view)
	if err != nil {
		return ViewSchema{}, false
	}

	return snapshot, true
}

func cloneManifest(manifest Manifest) Manifest {
	if manifest.Views != nil {
		manifest.Views = append(make([]UIViewDescriptor, 0, len(manifest.Views)), manifest.Views...)
	}

	if manifest.RequiredComponents != nil {
		manifest.RequiredComponents = append(
			make([]string, 0, len(manifest.RequiredComponents)),
			manifest.RequiredComponents...,
		)
	}

	return manifest
}

func cloneViewSchema(schema ViewSchema) (ViewSchema, error) {
	payload, err := json.Marshal(schema)
	if err != nil {
		return ViewSchema{}, err
	}

	var clone ViewSchema
	if err = json.Unmarshal(payload, &clone); err != nil {
		return ViewSchema{}, err
	}

	return clone, nil
}

type ViewBuilder struct {
	schema ViewSchema
}

type ViewOption func(*ViewSchema)

func View(id string, options ...ViewOption) *ViewBuilder {
	builder := &ViewBuilder{
		schema: ViewSchema{
			ProtocolVersion: ProtocolVersion,
			ID:              id,
			Regions:         emptyRegions(),
			Sources:         map[string]DataSource{},
			Actions:         map[string]Action{},
		},
	}

	for _, option := range options {
		option(&builder.schema)
	}

	return builder
}

func emptyRegions() Regions {
	return Regions{
		TopHeader:  []Node{},
		LeftPanel:  []Node{},
		Content:    []Node{},
		RightPanel: []Node{},
		Bottom:     []Node{},
	}
}

func Title(value string) ViewOption {
	return func(schema *ViewSchema) {
		schema.Title = value
	}
}

func State(value map[string]any) ViewOption {
	return func(schema *ViewSchema) {
		schema.State = value
	}
}

func TopHeader(nodes ...*NodeBuilder) ViewOption {
	return func(schema *ViewSchema) {
		schema.Regions.TopHeader = buildNodes(nodes)
	}
}

func LeftPanel(nodes ...*NodeBuilder) ViewOption {
	return func(schema *ViewSchema) {
		schema.Regions.LeftPanel = buildNodes(nodes)
	}
}

func Content(nodes ...*NodeBuilder) ViewOption {
	return func(schema *ViewSchema) {
		schema.Regions.Content = buildNodes(nodes)
	}
}

func RightPanel(nodes ...*NodeBuilder) ViewOption {
	return func(schema *ViewSchema) {
		schema.Regions.RightPanel = buildNodes(nodes)
	}
}

func Bottom(nodes ...*NodeBuilder) ViewOption {
	return func(schema *ViewSchema) {
		schema.Regions.Bottom = buildNodes(nodes)
	}
}

func (view *ViewBuilder) Action(name string, action *ActionBuilder) *ViewBuilder {
	view.schema.Actions[name] = action.action
	return view
}

func (view *ViewBuilder) Schema() ViewSchema {
	return view.schema
}

func (view *ViewBuilder) Revision(value string) *ViewBuilder {
	view.schema.Revision = value
	return view
}

func (view *ViewBuilder) Source(name string, source *SourceBuilder) *ViewBuilder {
	view.schema.Sources[name] = source.source
	return view
}

type NodeBuilder struct{ node Node }

func node(component, id string) *NodeBuilder {
	return &NodeBuilder{
		node: Node{
			Component: component,
			ID:        id,
			Props:     map[string]Value{},
			Events:    map[string]EventHandler{},
		},
	}
}

func (node *NodeBuilder) Row(row int) *NodeBuilder {
	node.ensureLayout()
	node.node.Layout.Row = row
	return node
}

func (node *NodeBuilder) Cols(cols int) *NodeBuilder {
	node.ensureLayout()
	node.node.Layout.Cols = cols
	return node
}

func (node *NodeBuilder) Offset(offset int) *NodeBuilder {
	node.ensureLayout()
	node.node.Layout.Offset = offset
	return node
}

func (node *NodeBuilder) Prop(name string, value any) *NodeBuilder {
	node.node.Props[name] = toValue(value)
	return node
}

func (node *NodeBuilder) On(name string, handler EventHandler) *NodeBuilder {
	node.node.Events[name] = handler
	return node
}

func (node *NodeBuilder) Children(children ...*NodeBuilder) *NodeBuilder {
	node.node.Children = buildNodes(children)
	return node
}

func (node *NodeBuilder) Slots(name string, children ...*NodeBuilder) *NodeBuilder {
	if node.node.Slots == nil {
		node.node.Slots = map[string][]Node{}
	}

	node.node.Slots[name] = buildNodes(children)

	return node
}

func (node *NodeBuilder) When(value Value) *NodeBuilder {
	node.node.When = &value
	return node
}

func (node *NodeBuilder) Build() Node {
	return node.node
}

func (node *NodeBuilder) ensureLayout() {
	if node.node.Layout == nil {
		node.node.Layout = &Layout{}
	}
}

func buildNodes(nodes []*NodeBuilder) []Node {
	result := make([]Node, 0, len(nodes))
	for _, n := range nodes {
		result = append(result, n.Build())
	}

	return result
}

type SourceBuilder struct{ source DataSource }

func Source(name string, source *SourceBuilder) ViewOption {
	return func(schema *ViewSchema) {
		if schema.Sources == nil {
			schema.Sources = map[string]DataSource{}
		}

		schema.Sources[name] = source.source
	}
}

func Tool(name string) *SourceBuilder {
	return &SourceBuilder{
		source: DataSource{
			Type:   "tool",
			Tool:   name,
			Policy: SourceManual,
			Input:  map[string]Value{},
		},
	}
}

func (source *SourceBuilder) Input(name string, value any) *SourceBuilder {
	source.source.Input[name] = toValue(value)
	return source
}

func (source *SourceBuilder) OnMount() *SourceBuilder {
	source.source.Policy = SourceOnMount
	return source
}

func (source *SourceBuilder) TTL(value int64) *SourceBuilder {
	source.source.Cache = &CachePolicy{TTL: value}
	return source
}

type ActionBuilder struct{ action Action }

func CallTool(name string) *ActionBuilder {
	return &ActionBuilder{
		action: Action{
			Type:  "tool",
			Tool:  name,
			Input: map[string]Value{},
		},
	}
}

func SetState(path string, value any) *ActionBuilder {
	return &ActionBuilder{
		action: Action{
			Type:  "set-state",
			Path:  path,
			Value: ptrValue(value),
		},
	}
}

func MergeState(path string, value any) *ActionBuilder {
	return &ActionBuilder{
		action: Action{
			Type:  "merge-state",
			Path:  path,
			Value: ptrValue(value),
		},
	}
}

func (action *ActionBuilder) Input(name string, value any) *ActionBuilder {
	action.action.Input[name] = toValue(value)
	return action
}

func (action *ActionBuilder) Then(effects ...Effect) *ActionBuilder {
	action.action.Effects = append(action.action.Effects, effects...)
	return action
}

func ActionRef(name string) EventHandler {
	return EventHandler{Action: name}
}

func Invalidate(source string) Effect {
	return Effect{
		Type:   "invalidate",
		Source: source,
	}
}

func RefreshSource(source string) Effect {
	return Effect{
		Type:   "refresh-source",
		Source: source,
	}
}

func ToastSuccess(message string) Effect {
	return Effect{
		Type:    "toast",
		Variant: "success",
		Message: message,
	}
}

func SetStateEffect(path string, value any) Effect {
	return Effect{
		Type:  "set-state",
		Path:  path,
		Value: ptrValue(value),
	}
}

func toValue(value any) Value {
	if typed, ok := value.(Value); ok {
		return typed
	}

	return Literal(value)
}

func ptrValue(value any) *Value {
	result := toValue(value)
	return &result
}
