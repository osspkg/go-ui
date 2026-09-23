/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package dashboard demonstrates shadcn components, grid placement, and data sources.
package dashboard

import "go.osspkg.com/ui"

const (
	metricsCacheTTL   = 30
	headerTitleCols   = 9
	refreshButtonCols = 3
	refreshButtonPos  = 9
	metricCardCols    = 6
	latencyCardOffset = 6
	eventTableCols    = 12
	eventTableRow     = 2
)

// NewApp builds a dashboard whose metrics source is loaded when the view opens.
func NewApp() (*ui.App, error) {
	app := ui.New(
		ui.AppID("metrics"),
		ui.AppTitle("Metrics"),
		ui.AppVersion("1.0.0"),
	)
	view := ui.View(
		"metrics.dashboard",
		ui.Title("Metrics dashboard"),
		ui.Source("metrics", ui.Tool("metrics.summary").OnMount().TTL(metricsCacheTTL)),
		ui.TopHeader(
			ui.HTML(ui.HTMLH1, "title").Row(1).Cols(headerTitleCols).Prop(
				"text",
				"Metrics dashboard",
			),
			ui.Shadcn(
				ui.ShadcnButton,
				"refresh",
			).Row(1).Cols(refreshButtonCols).Offset(refreshButtonPos).Prop("text", "Refresh").On(
				"click",
				ui.ActionRef("metrics.refresh"),
			),
		),
		ui.Content(
			ui.Shadcn(ui.ShadcnCard, "requests-card").Row(1).Cols(metricCardCols).Children(
				ui.HTML(ui.HTMLH2, "requests-heading").Prop("text", "Requests"),
				ui.Shadcn(ui.ShadcnTypography, "requests-value").Prop(
					"text",
					ui.SourceRef("metrics.requests"),
				),
			),
			ui.Shadcn(
				ui.ShadcnCard,
				"latency-card",
			).Row(1).Cols(metricCardCols).Offset(latencyCardOffset).Children(
				ui.HTML(ui.HTMLH2, "latency-heading").Prop("text", "Latency"),
				ui.Shadcn(ui.ShadcnTypography, "latency-value").Prop(
					"text",
					ui.SourceRef("metrics.latency"),
				),
			),
			ui.DataTable("recent-events").Row(eventTableRow).Cols(eventTableCols).Prop(
				"rows",
				ui.SourceRef("metrics.events"),
			).Prop("loading", ui.SourceRef("metrics.$loading")).Prop(
				"error",
				ui.SourceRef("metrics.$error"),
			),
		),
	).Action(
		"metrics.refresh",
		ui.CallTool("metrics.refresh").Then(
			ui.RefreshSource("metrics"),
			ui.ToastSuccess("Metrics refreshed"),
		),
	)
	if err := app.AddView(view); err != nil {
		return nil, err
	}
	return app, nil
}
