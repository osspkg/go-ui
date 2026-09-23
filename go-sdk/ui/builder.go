/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package ui

import (
	"encoding/json"
	"fmt"
)

// App builds and stores declarative UI views for a plugin.
type App struct {
	manifest Manifest
	views    map[string]ViewSchema
}

// AppOption configures an App during construction.
type AppOption func(*App)

// AppID sets the plugin identifier in the application manifest.
func AppID(value string) AppOption {
	return func(app *App) {
		app.manifest.Plugin.ID = value
	}
}

// AppTitle sets the plugin title in the application manifest.
func AppTitle(value string) AppOption {
	return func(app *App) {
		app.manifest.Plugin.Title = value
	}
}

// AppVersion sets the plugin version in the application manifest.
func AppVersion(value string) AppOption {
	return func(app *App) {
		app.manifest.Plugin.Version = value
	}
}

// New creates an application with the supplied manifest options.
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

// AddView validates and registers a view in the application.
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

// Manifest returns a copy of the application manifest.
func (app *App) Manifest() Manifest {
	return cloneManifest(app.manifest)
}

// ManifestJSON serializes a copy of the application manifest.
func (app *App) ManifestJSON() ([]byte, error) {
	return json.Marshal(app.Manifest())
}

// View returns a copy of a registered view by ID.
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

// ViewBuilder incrementally constructs a ViewSchema.
type ViewBuilder struct {
	schema ViewSchema
}

// ViewOption configures a view during construction.
type ViewOption func(*ViewSchema)

// View creates a view builder with the supplied ID and options.
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

// Title sets the view title.
func Title(value string) ViewOption {
	return func(schema *ViewSchema) {
		schema.Title = value
	}
}

// State sets the initial view state.
func State(value map[string]any) ViewOption {
	return func(schema *ViewSchema) {
		schema.State = value
	}
}

// TopHeader places nodes in the top header region.
func TopHeader(nodes ...*NodeBuilder) ViewOption {
	return func(schema *ViewSchema) {
		schema.Regions.TopHeader = buildNodes(nodes)
	}
}

// LeftPanel places nodes in the left panel region.
func LeftPanel(nodes ...*NodeBuilder) ViewOption {
	return func(schema *ViewSchema) {
		schema.Regions.LeftPanel = buildNodes(nodes)
	}
}

// Content places nodes in the content region.
func Content(nodes ...*NodeBuilder) ViewOption {
	return func(schema *ViewSchema) {
		schema.Regions.Content = buildNodes(nodes)
	}
}

// RightPanel places nodes in the right panel region.
func RightPanel(nodes ...*NodeBuilder) ViewOption {
	return func(schema *ViewSchema) {
		schema.Regions.RightPanel = buildNodes(nodes)
	}
}

// Bottom places nodes in the bottom region.
func Bottom(nodes ...*NodeBuilder) ViewOption {
	return func(schema *ViewSchema) {
		schema.Regions.Bottom = buildNodes(nodes)
	}
}

// Action adds a named action to the view.
func (view *ViewBuilder) Action(name string, action *ActionBuilder) *ViewBuilder {
	view.schema.Actions[name] = action.action
	return view
}

// Schema returns the current view schema.
func (view *ViewBuilder) Schema() ViewSchema {
	return view.schema
}

// Revision sets the view revision.
func (view *ViewBuilder) Revision(value string) *ViewBuilder {
	view.schema.Revision = value
	return view
}

// Source adds a named data source to the view.
func (view *ViewBuilder) Source(name string, source *SourceBuilder) *ViewBuilder {
	view.schema.Sources[name] = source.source
	return view
}

// NodeBuilder incrementally constructs a UI node.
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

// Row sets the node's grid row.
func (node *NodeBuilder) Row(row int) *NodeBuilder {
	node.ensureLayout()
	node.node.Layout.Row = row
	return node
}

// Cols sets the node's grid column span.
func (node *NodeBuilder) Cols(cols int) *NodeBuilder {
	node.ensureLayout()
	node.node.Layout.Cols = cols
	return node
}

// Offset sets the node's grid column offset.
func (node *NodeBuilder) Offset(offset int) *NodeBuilder {
	node.ensureLayout()
	node.node.Layout.Offset = offset
	return node
}

// Prop sets a named node property.
func (node *NodeBuilder) Prop(name string, value any) *NodeBuilder {
	node.node.Props[name] = toValue(value)
	return node
}

// On attaches an event handler to the node.
func (node *NodeBuilder) On(name string, handler EventHandler) *NodeBuilder {
	node.node.Events[name] = handler
	return node
}

// Children sets the node's child nodes.
func (node *NodeBuilder) Children(children ...*NodeBuilder) *NodeBuilder {
	node.node.Children = buildNodes(children)
	return node
}

// Slots sets the nodes in a named slot.
func (node *NodeBuilder) Slots(name string, children ...*NodeBuilder) *NodeBuilder {
	if node.node.Slots == nil {
		node.node.Slots = map[string][]Node{}
	}

	node.node.Slots[name] = buildNodes(children)

	return node
}

// When sets the conditional expression for the node.
func (node *NodeBuilder) When(value Value) *NodeBuilder {
	node.node.When = &value
	return node
}

// Build returns the current node value.
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

// SourceBuilder incrementally constructs a data source.
type SourceBuilder struct{ source DataSource }

// Source adds a named data source to a view.
func Source(name string, source *SourceBuilder) ViewOption {
	return func(schema *ViewSchema) {
		if schema.Sources == nil {
			schema.Sources = map[string]DataSource{}
		}

		schema.Sources[name] = source.source
	}
}

// Tool creates a tool-backed data source builder.
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

// Input adds an input value to the data source.
func (source *SourceBuilder) Input(name string, value any) *SourceBuilder {
	source.source.Input[name] = toValue(value)
	return source
}

// OnMount configures the source to load when the view mounts.
func (source *SourceBuilder) OnMount() *SourceBuilder {
	source.source.Policy = SourceOnMount
	return source
}

// TTL sets the source cache lifetime in seconds.
func (source *SourceBuilder) TTL(value int64) *SourceBuilder {
	source.source.Cache = &CachePolicy{TTL: value}
	return source
}

// ActionBuilder incrementally constructs an action.
type ActionBuilder struct{ action Action }

// CallTool creates a tool-calling action builder.
func CallTool(name string) *ActionBuilder {
	return &ActionBuilder{
		action: Action{
			Type:  actionTypeTool,
			Tool:  name,
			Input: map[string]Value{},
		},
	}
}

// SetState creates an action that replaces a state path.
func SetState(path string, value any) *ActionBuilder {
	return &ActionBuilder{
		action: Action{
			Type:  actionTypeSetState,
			Path:  path,
			Value: ptrValue(value),
		},
	}
}

// MergeState creates an action that merges a value into a state path.
func MergeState(path string, value any) *ActionBuilder {
	return &ActionBuilder{
		action: Action{
			Type:  actionTypeMergeState,
			Path:  path,
			Value: ptrValue(value),
		},
	}
}

// Input adds an input value to the action.
func (action *ActionBuilder) Input(name string, value any) *ActionBuilder {
	action.action.Input[name] = toValue(value)
	return action
}

// Then appends effects to the action.
func (action *ActionBuilder) Then(effects ...Effect) *ActionBuilder {
	action.action.Effects = append(action.action.Effects, effects...)
	return action
}

// ActionRef creates an event handler that invokes a named action.
func ActionRef(name string) EventHandler {
	return EventHandler{Action: name}
}

// Invalidate creates an effect that invalidates a data source.
func Invalidate(source string) Effect {
	return Effect{
		Type:   actionTypeInvalidate,
		Source: source,
	}
}

// RefreshSource creates an effect that refreshes a data source.
func RefreshSource(source string) Effect {
	return Effect{
		Type:   actionTypeRefreshSource,
		Source: source,
	}
}

// ToastSuccess creates a successful toast effect.
func ToastSuccess(message string) Effect {
	return Effect{
		Type:    "toast",
		Variant: "success",
		Message: message,
	}
}

// SetStateEffect creates an effect that replaces a state path.
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
