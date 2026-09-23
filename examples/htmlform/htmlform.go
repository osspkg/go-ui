/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package htmlform demonstrates a declarative form built from safe HTML tags.
package htmlform

import "go.osspkg.com/ui"

const (
	profileCardCols   = 8
	profileCardOffset = 2
)

// NewApp builds a profile form with native HTML elements and a local action.
func NewApp() (*ui.App, error) {
	app := ui.New(
		ui.AppID("profile"),
		ui.AppTitle("Profile"),
		ui.AppVersion("1.0.0"),
	)
	view := ui.View(
		"profile.edit",
		ui.Title("Profile"),
		ui.State(
			map[string]any{
				"profile":   map[string]any{"name": ""},
				"submitted": false,
			},
		),
		ui.Content(
			ui.HTML(
				ui.HTMLArticle,
				"profile-card",
			).Row(1).Cols(profileCardCols).Offset(profileCardOffset).Children(
				ui.HTML(ui.HTMLH1, "title").Prop("text", "Edit profile"),
				ui.HTML(ui.HTMLP, "hint").Prop(
					"text",
					"Changes are saved when the form is submitted.",
				),
				ui.HTML(ui.HTMLForm, "profile-form").On(
					"submit",
					ui.ActionRef("profile.submit"),
				).Children(
					ui.HTML(ui.HTMLLabel, "name-label").Prop("text", "Name"),
					ui.HTML(ui.HTMLInput, "name").Prop(
						"name",
						"name",
					).Prop(
						"value",
						ui.StateRef("profile.name"),
					).Prop("placeholder", "Ada Lovelace"),
					ui.HTML(ui.HTMLButton, "save").Prop(
						"type",
						"submit",
					).Prop("text", "Save profile"),
				),
			),
		),
	).Action(
		"profile.submit",
		ui.SetState("submitted", true).Then(ui.ToastSuccess("Profile saved")),
	)
	if err := app.AddView(view); err != nil {
		return nil, err
	}
	return app, nil
}
