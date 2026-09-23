/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package mcp registers declarative UI resources on a go-mcp server.
package mcp

import (
	"encoding/json"
	"errors"
	"strings"

	ossmcp "go.osspkg.com/mcp"

	"go.osspkg.com/ui"
)

// Capabilities returns the MCP option advertising the UI capability.
func Capabilities() ossmcp.Option {
	return ossmcp.WithCapabilities(
		map[string]any{
			"experimental": map[string]any{
				"osspkg.ui": map[string]any{
					"version":  ui.ProtocolVersion,
					"manifest": "ui://manifest",
				},
			},
		},
	)
}

// Register publishes an app manifest and its view resources on server.
func Register(server *ossmcp.Server, app *ui.App) error {
	if server == nil || app == nil {
		return errors.New("ui/mcp: server and app are required")
	}

	manifest, err := app.ManifestJSON()
	if err != nil {
		return err
	}

	if err = app.Manifest().Validate(); err != nil {
		return err
	}

	if err = server.RegisterResource(
		ossmcp.Resource{
			URI:         "ui://manifest",
			Name:        "UI Manifest",
			Description: "UI entry points exposed by the application",
			MIMEType:    ui.UIMIMEType,
			Text:        string(manifest),
		},
	); err != nil {
		return err
	}

	for _, descriptor := range app.Manifest().Views {
		view, ok := app.View(descriptor.ID)
		if !ok {
			return errors.New("ui/mcp: manifest view is missing")
		}

		payload, err := json.Marshal(view)
		if err != nil {
			return err
		}

		name := descriptor.Title
		if strings.TrimSpace(name) == "" {
			name = descriptor.ID
		}

		if err := server.RegisterResource(
			ossmcp.Resource{
				URI:      descriptor.Schema,
				Name:     name,
				MIMEType: ui.UIMIMEType,
				Text:     string(payload),
			},
		); err != nil {
			return err
		}
	}

	return nil
}
