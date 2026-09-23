package ui

import (
	"encoding/json"
	"testing"
)

func TestUnit_BuilderCreatesView(t *testing.T) {
	view := View(
		"users.list",
		Title("Users"),
		Source("users", Tool("users.list").OnMount()),
		Content(
			DataTable("users-table").Row(1).Cols(12).Prop(
				"rows",
				SourceRef("users.items"),
			),
		),
	).Action(
		"users.refresh",
		CallTool("users.refresh").Then(Invalidate("users")),
	)

	app := New(AppID("users"), AppTitle("Users"), AppVersion("1.0.0"))
	if err := app.AddView(view); err != nil {
		t.Fatal(err)
	}
	manifest := app.Manifest()
	if len(manifest.Views) != 1 || manifest.Views[0].ID != "users.list" {
		t.Fatalf("manifest = %#v", manifest)
	}
	payload, err := json.Marshal(view.Schema())
	if err != nil {
		t.Fatal(err)
	}
	var decoded struct {
		Regions map[string]json.RawMessage `json:"regions"`
	}
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatal(err)
	}
	for _, region := range []string{
		"top-header",
		"left-panel",
		"content",
		"right-panel",
		"bottom",
	} {
		if string(decoded.Regions[region]) == "null" {
			t.Fatalf("region %q serialized as null", region)
		}
	}
}

func TestUnit_AppSnapshotsViews(t *testing.T) {
	view := View("users.list", Title("Users")).Action(
		"initial",
		CallTool("initial"),
	)
	app := New(AppID("users"))
	if err := app.AddView(view); err != nil {
		t.Fatal(err)
	}

	manifest := app.Manifest()
	manifest.Views[0].ID = "changed"
	if got := app.Manifest().Views[0].ID; got != "users.list" {
		t.Fatalf("stored manifest view id = %q, want users.list", got)
	}

	view.Action("late", CallTool("late"))
	returned, ok := app.View("users.list")
	if !ok {
		t.Fatal("view is missing")
	}
	returned.Actions["returned"] = Action{Type: "tool", Tool: "returned"}

	stored, ok := app.View("users.list")
	if !ok {
		t.Fatal("view is missing after mutation")
	}
	if len(stored.Actions) != 1 {
		t.Fatalf("stored actions = %d, want 1", len(stored.Actions))
	}
	if _, ok := stored.Actions["initial"]; !ok {
		t.Fatal("initial action is missing")
	}
}
