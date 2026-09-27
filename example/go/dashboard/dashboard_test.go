package dashboard

import "testing"

func TestUnit_NewApp(t *testing.T) {
	app, err := NewApp()
	if err != nil {
		t.Fatal(err)
	}
	if err := app.Manifest().Validate(); err != nil {
		t.Fatal(err)
	}
	view, ok := app.View("metrics.dashboard")
	if !ok {
		t.Fatal("dashboard view is missing")
	}
	if len(view.Regions.Content) != 3 {
		t.Fatalf("content nodes = %d, want 3", len(view.Regions.Content))
	}
}
