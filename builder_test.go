package ui

import (
	"encoding/json"
	"testing"
)

func TestUnit_BuilderCreatesView(t *testing.T) {
	view := View("users.list", Title("Users"), Source("users", Tool("users.list").OnMount()), Content(
		DataTable("users-table").Row(1).Cols(12).Prop("rows", SourceRef("users.items")),
	)).Action("users.refresh", CallTool("users.refresh").Then(Invalidate("users")))

	app := New(PluginID("users"), PluginTitle("Users"), PluginVersion("1.0.0"))
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
	for _, region := range []string{"top-header", "left-panel", "content", "right-panel", "bottom"} {
		if string(decoded.Regions[region]) == "null" {
			t.Fatalf("region %q serialized as null", region)
		}
	}
}
