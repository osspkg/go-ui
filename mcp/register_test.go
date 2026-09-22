package mcp

import (
	"testing"

	goMCP "go.osspkg.com/mcp"

	"go.osspkg.com/ui"
)

func TestUnit_Register(t *testing.T) {
	server, err := goMCP.New(
		goMCP.ServerInfo{Name: "users", Version: "1.0.0"},
		Capabilities(),
	)
	if err != nil {
		t.Fatal(err)
	}
	app := ui.New(
		ui.AppID("users"),
		ui.AppTitle("Users"),
		ui.AppVersion("1.0.0"),
	)
	if err := app.AddView(
		ui.View(
			"users.list",
			ui.Title("Users"),
		),
	); err != nil {
		t.Fatal(err)
	}
	if err := Register(server, app); err != nil {
		t.Fatal(err)
	}
}

func TestUnit_RegisterViewWithoutTitle(t *testing.T) {
	server, err := goMCP.New(
		goMCP.ServerInfo{Name: "users", Version: "1.0.0"},
		Capabilities(),
	)
	if err != nil {
		t.Fatal(err)
	}
	app := ui.New(ui.AppID("users"))
	if err := app.AddView(ui.View("users.list")); err != nil {
		t.Fatal(err)
	}
	if err := Register(server, app); err != nil {
		t.Fatal(err)
	}
}
