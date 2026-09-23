/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package ui defines the transport-independent declarative UI protocol.
package ui

//go:generate easyjson

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// ProtocolVersion is the supported declarative UI protocol version.
	ProtocolVersion = "1.0"
	// UIMIMEType is the MIME type used for declarative UI resources.
	UIMIMEType = "application/vnd.osspkg.ui+json"
)

var (
	// ErrInvalidSchema reports an invalid manifest or view schema.
	ErrInvalidSchema = errors.New("invalid ui schema")
	// ErrInvalidLayout reports an invalid node layout.
	ErrInvalidLayout = errors.New("invalid ui layout")
	// ErrInvalidValue reports an invalid declarative value.
	ErrInvalidValue = errors.New("invalid ui value")
)

// Manifest describes the application entry points exposed to a host.
//
//easyjson:json
//nolint:recvcheck // easyjson generates value marshal and pointer unmarshal methods.
type Manifest struct {
	ProtocolVersion    string             `json:"protocolVersion"`
	Plugin             PluginManifest     `json:"plugin"`
	Views              []UIViewDescriptor `json:"views"`
	RequiredComponents []string           `json:"requiredComponents,omitempty"`
}

// PluginManifest identifies the plugin that owns the UI.
//
//easyjson:json
//nolint:recvcheck // easyjson generates value marshal and pointer unmarshal methods.
type PluginManifest struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Version     string `json:"version"`
	Description string `json:"description,omitempty"`
}

// UIViewDescriptor identifies a view schema resource in a manifest.
//
//easyjson:json
//nolint:recvcheck // easyjson generates value marshal and pointer unmarshal methods.
type UIViewDescriptor struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Route  string `json:"route,omitempty"`
	Schema string `json:"schema"`
}

// ViewSchema is the declarative layout, state, source, and action model for a view.
//
//easyjson:json
//nolint:recvcheck // easyjson generates value marshal and pointer unmarshal methods.
type ViewSchema struct {
	ProtocolVersion string                `json:"protocolVersion"`
	ID              string                `json:"id"`
	Revision        string                `json:"revision,omitempty"`
	Title           string                `json:"title,omitempty"`
	State           map[string]any        `json:"state,omitempty"`
	Sources         map[string]DataSource `json:"sources,omitempty"`
	Actions         map[string]Action     `json:"actions,omitempty"`
	Regions         Regions               `json:"regions"`
}

// Regions groups view nodes into the host-supported layout regions.
//
//easyjson:json
//nolint:recvcheck // easyjson generates value marshal and pointer unmarshal methods.
type Regions struct {
	TopHeader  []Node `json:"top-header"`
	LeftPanel  []Node `json:"left-panel"`
	Content    []Node `json:"content"`
	RightPanel []Node `json:"right-panel"`
	Bottom     []Node `json:"bottom"`
}

// Layout places a node in a twelve-column row layout.
//
//easyjson:json
//nolint:recvcheck // easyjson generates value marshal and pointer unmarshal methods.
type Layout struct {
	Row    int `json:"row"`
	Cols   int `json:"cols,omitempty"`
	Offset int `json:"offset,omitempty"`
}

// Node describes one declarative UI component and its descendants.
//
//easyjson:json
//nolint:recvcheck // easyjson generates value marshal and pointer unmarshal methods.
type Node struct {
	ID        string                  `json:"id"`
	Component string                  `json:"component"`
	Layout    *Layout                 `json:"layout,omitempty"`
	Props     map[string]Value        `json:"props,omitempty"`
	Events    map[string]EventHandler `json:"events,omitempty"`
	Children  []Node                  `json:"children,omitempty"`
	Slots     map[string][]Node       `json:"slots,omitempty"`
	When      *Value                  `json:"when,omitempty"`
}

// EventHandler maps a UI event to an action or sequence of action steps.
//
//easyjson:json
//nolint:recvcheck // easyjson generates value marshal and pointer unmarshal methods.
type EventHandler struct {
	Action string       `json:"action,omitempty"`
	Steps  []ActionStep `json:"steps,omitempty"`
}

// DataSource describes a host tool used to populate view data.
//
//easyjson:json
//nolint:recvcheck // easyjson generates value marshal and pointer unmarshal methods.
type DataSource struct {
	Type      string           `json:"type"`
	Tool      string           `json:"tool"`
	Input     map[string]Value `json:"input,omitempty"`
	Policy    SourcePolicy     `json:"policy,omitempty"`
	Cache     *CachePolicy     `json:"cache,omitempty"`
	RefreshOn []string         `json:"refreshOn,omitempty"`
}

// SourcePolicy controls when a data source is fetched.
type SourcePolicy string

const (
	// SourceManual fetches the source only when explicitly requested.
	SourceManual SourcePolicy = "manual"
	// SourceOnMount fetches the source when the view is mounted.
	SourceOnMount SourcePolicy = "on-mount"
)

// CachePolicy controls how long fetched source data may be reused.
//
//easyjson:json
//nolint:recvcheck // easyjson generates value marshal and pointer unmarshal methods.
type CachePolicy struct {
	TTL int64 `json:"ttl,omitempty"`
}

// Action describes a named operation available to the view.
//
//easyjson:json
//nolint:recvcheck // easyjson generates value marshal and pointer unmarshal methods.
type Action struct {
	Type    string           `json:"type"`
	Tool    string           `json:"tool,omitempty"`
	Input   map[string]Value `json:"input,omitempty"`
	Path    string           `json:"path,omitempty"`
	Value   *Value           `json:"value,omitempty"`
	Source  string           `json:"source,omitempty"`
	To      string           `json:"to,omitempty"`
	Variant string           `json:"variant,omitempty"`
	Message string           `json:"message,omitempty"`
	Effects []Effect         `json:"effects,omitempty"`
}

// ActionStep is one step in an event action sequence.
//
//easyjson:json
//nolint:recvcheck // easyjson generates value marshal and pointer unmarshal methods.
type ActionStep struct {
	Type   string           `json:"type"`
	Action string           `json:"action,omitempty"`
	Input  map[string]Value `json:"input,omitempty"`
}

// Effect describes a client-side consequence of an action.
//
//easyjson:json
//nolint:recvcheck // easyjson generates value marshal and pointer unmarshal methods.
type Effect struct {
	Type         string `json:"type"`
	Path         string `json:"path,omitempty"`
	Source       string `json:"source,omitempty"`
	Value        *Value `json:"value,omitempty"`
	BaseRevision string `json:"baseRevision,omitempty"`
	Revision     string `json:"revision,omitempty"`
	To           string `json:"to,omitempty"`
	Variant      string `json:"variant,omitempty"`
	Message      string `json:"message,omitempty"`
}

// Validate checks whether the manifest contains valid plugin and view entries.
func (m Manifest) Validate() error {
	if m.ProtocolVersion != ProtocolVersion {
		return fmt.Errorf("%w: unsupported manifest protocol version %q", ErrInvalidSchema, m.ProtocolVersion)
	}

	if m.Plugin.ID == "" {
		return fmt.Errorf("%w: plugin id is required", ErrInvalidSchema)
	}

	seen := make(map[string]struct{}, len(m.Views))
	for _, view := range m.Views {
		if view.ID == "" || view.Schema == "" {
			return fmt.Errorf("%w: view id and schema are required", ErrInvalidSchema)
		}

		if _, ok := seen[view.ID]; ok {
			return fmt.Errorf("%w: duplicate view %q", ErrInvalidSchema, view.ID)
		}

		seen[view.ID] = struct{}{}
	}

	return nil
}

func validName(name string) bool {
	return name != "" && !strings.ContainsAny(name, "\r\n")
}
