/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD-3-Clause license that can be found in the LICENSE file.
 */

package ui

// Host RPC method names used by the React runtime. The MCP adapter publishes
// UI resources; a host transport is responsible for exposing these methods.
const (
	RPCMethodGetView  = "ui.get"
	RPCMethodDataCall = "data.call"
	RPCMethodAction   = "ui.action"
)

// GetViewRequest requests a view schema from a host UI backend.
type GetViewRequest struct {
	Plugin  string `json:"plugin,omitempty"`
	View    string `json:"view"`
	Context any    `json:"context,omitempty"`
}

// GetViewResponse contains a view schema returned by a host UI backend.
type GetViewResponse struct {
	Schema ViewSchema `json:"schema"`
}

// DataCallRequest requests a data source operation from a host UI backend.
type DataCallRequest struct {
	Plugin    string         `json:"plugin,omitempty"`
	View      string         `json:"view"`
	Source    string         `json:"source"`
	Operation string         `json:"operation"`
	Input     map[string]any `json:"input,omitempty"`
}

// ActionRequest requests a declared view action from a host UI backend.
type ActionRequest struct {
	Plugin string         `json:"plugin,omitempty"`
	View   string         `json:"view"`
	Action string         `json:"action"`
	Input  map[string]any `json:"input,omitempty"`
}
