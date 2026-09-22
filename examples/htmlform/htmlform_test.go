package htmlform

import "testing"

func TestUnit_NewApp(t *testing.T) {
	app, err := NewApp()
	if err != nil {
		t.Fatal(err)
	}
	if err := app.Manifest().Validate(); err != nil {
		t.Fatal(err)
	}
	if _, ok := app.View("profile.edit"); !ok {
		t.Fatal("profile view is missing")
	}
}
