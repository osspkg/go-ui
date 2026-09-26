package ui

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestUnit_ValueRoundTrip(t *testing.T) {
	value := Expr("$eq", StateRef("status"), Literal("active"))
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	var decoded Value
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Kind != ValueExpr || decoded.Path != "$eq" {
		t.Fatalf("decoded value = %#v", decoded)
	}
}

func TestUnit_ValueLiteralObjects(t *testing.T) {
	for _, data := range []string{
		`{"rows":[{"id":1}]}`,
		`{"disabled":null}`,
	} {
		t.Run(data, func(t *testing.T) {
			var value Value
			if err := json.Unmarshal([]byte(data), &value); err != nil {
				t.Fatal(err)
			}
			if value.Kind != ValueLiteral {
				t.Fatalf("value kind = %q, want %q", value.Kind, ValueLiteral)
			}
			encoded, err := json.Marshal(value)
			if err != nil {
				t.Fatal(err)
			}
			var want, got any
			if err := json.Unmarshal([]byte(data), &want); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal(encoded, &got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("round-trip JSON differs: got %s, want %s", encoded, data)
			}
		})
	}
}

func TestUnit_LayoutValidation(t *testing.T) {
	tests := []struct {
		name   string
		layout Layout
		region string
		want   error
	}{
		{name: "valid", layout: Layout{Row: 1, Cols: 8, Offset: 2}, region: "content"},
		{name: "overflow", layout: Layout{Row: 1, Cols: 8, Offset: 5}, region: "content", want: ErrInvalidLayout},
		{name: "side panel", layout: Layout{Row: 1, Cols: 1}, region: "left-panel", want: ErrInvalidLayout},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.layout.Validate(test.region)
			if !errors.Is(err, test.want) {
				t.Fatalf("Validate() error = %v, want %v", err, test.want)
			}
		})
	}
}

func TestUnit_SchemaLimits(t *testing.T) {
	view := ViewSchema{
		ProtocolVersion: ProtocolVersion,
		ID:              "users.list",
		Regions: func() Regions {
			regions := emptyRegions()
			regions.Content = []Node{{ID: "title", Component: "text", Props: map[string]Value{"text": Literal("Users")}}}
			return regions
		}(),
	}
	tests := []struct {
		name   string
		limits Limits
	}{
		{name: "bytes", limits: Limits{MaxBytes: 64}},
		{name: "strings", limits: Limits{MaxStrings: 1}},
		{name: "string length", limits: Limits{MaxStringLength: 3}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := view.ValidateWithLimits(test.limits); !errors.Is(err, ErrInvalidSchema) {
				t.Fatalf("ValidateWithLimits() error = %v, want ErrInvalidSchema", err)
			}
		})
	}
}

func TestUnit_SemanticValidation(t *testing.T) {
	tests := []struct {
		name  string
		value func() error
		want  error
	}{
		{name: "unsupported operator", value: func() error { return Expr("$unknown", Literal(true)).Validate() }, want: ErrInvalidValue},
		{name: "operator arity", value: func() error { return Expr("$eq", Literal(true)).Validate() }, want: ErrInvalidValue},
		{name: "unknown source reference", value: func() error {
			regions := emptyRegions()
			regions.Content = []Node{{ID: "n", Component: "text", Props: map[string]Value{"value": SourceRef("missing.items")}}}
			return (ViewSchema{ProtocolVersion: ProtocolVersion, ID: "users.list", Regions: regions}).Validate()
		}, want: ErrInvalidSchema},
		{name: "unknown effect source", value: func() error {
			return (ViewSchema{ProtocolVersion: ProtocolVersion, ID: "users.list", Actions: map[string]Action{"save": {Type: "tool", Tool: "users.save", Effects: []Effect{{Type: "invalidate", Source: "missing"}}}}, Regions: emptyRegions()}).Validate()
		}, want: ErrInvalidSchema},
		{name: "unknown action value source", value: func() error {
			return (ViewSchema{ProtocolVersion: ProtocolVersion, ID: "users.list", Actions: map[string]Action{"save": {Type: "set-state", Path: "form", Value: ptrValue(SourceRef("missing.items"))}}, Regions: emptyRegions()}).Validate()
		}, want: ErrInvalidSchema},
		{name: "module component name", value: func() error {
			regions := emptyRegions()
			regions.Content = []Node{{ID: "n", Component: "@scope/button"}}
			return (ViewSchema{ProtocolVersion: ProtocolVersion, ID: "users.list", Regions: regions}).Validate()
		}, want: ErrInvalidSchema},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.value(); !errors.Is(err, test.want) {
				t.Fatalf("validation error = %v, want %v", err, test.want)
			}
		})
	}
}

func TestUnit_DottedSourceReference(t *testing.T) {
	regions := emptyRegions()
	regions.Content = []Node{{
		ID:        "orders-table",
		Component: "table",
		Props:     map[string]Value{"rows": SourceRef("orders.list.items")},
	}}
	view := ViewSchema{
		ProtocolVersion: ProtocolVersion,
		ID:              "orders",
		Sources: map[string]DataSource{
			"orders.list": {Type: "tool", Tool: "orders.list"},
		},
		Regions: regions,
	}
	if err := view.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestUnit_SchemaRequiresAllRegionArrays(t *testing.T) {
	view := ViewSchema{ProtocolVersion: ProtocolVersion, ID: "users.list", Regions: Regions{Content: []Node{}}}
	if err := view.Validate(); !errors.Is(err, ErrInvalidSchema) {
		t.Fatalf("Validate() error = %v, want ErrInvalidSchema", err)
	}
}

func TestUnit_SchemaExpressionDepthLimit(t *testing.T) {
	value := Literal(true)
	for range 3 {
		value = Expr("$not", value)
	}
	regions := emptyRegions()
	regions.Content = []Node{{ID: "condition", Component: "text", When: &value}}
	view := ViewSchema{ProtocolVersion: ProtocolVersion, ID: "users.list", Regions: regions}
	if err := view.ValidateWithLimits(Limits{MaxExpressionDepth: 2}); !errors.Is(err, ErrInvalidSchema) {
		t.Fatalf("ValidateWithLimits() error = %v, want ErrInvalidSchema", err)
	}
}

func TestUnit_JSONFixtures(t *testing.T) {
	valid, err := os.ReadFile("../../fixtures/valid-view.json")
	if err != nil {
		t.Fatal(err)
	}
	var view ViewSchema
	if err := json.Unmarshal(valid, &view); err != nil {
		t.Fatal(err)
	}
	if err := view.Validate(); err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(view)
	if err != nil {
		t.Fatal(err)
	}
	var want, got any
	if err := json.Unmarshal(valid, &want); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(encoded, &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("round-trip JSON differs:\n got: %s\nwant: %s", encoded, valid)
	}

	invalid, err := os.ReadFile("../../fixtures/invalid-layout.json")
	if err != nil {
		t.Fatal(err)
	}
	var broken ViewSchema
	if err := json.Unmarshal(invalid, &broken); err != nil {
		t.Fatal(err)
	}
	if err := broken.Validate(); !errors.Is(err, ErrInvalidLayout) {
		t.Fatalf("invalid fixture error = %v, want ErrInvalidLayout", err)
	}
}

func TestUnit_ManifestValidationParity(t *testing.T) {
	base := Manifest{
		ProtocolVersion: ProtocolVersion,
		Plugin:          PluginManifest{ID: "users", Title: "Users", Version: "1.0"},
		Views:           []UIViewDescriptor{},
	}

	if err := (Manifest{ProtocolVersion: ProtocolVersion, Plugin: base.Plugin}).Validate(); !errors.Is(err, ErrInvalidSchema) {
		t.Fatalf("Validate() error = %v, want ErrInvalidSchema for nil views", err)
	}

	base.RequiredComponents = []string{"/button"}
	if err := base.Validate(); !errors.Is(err, ErrInvalidSchema) {
		t.Fatalf("Validate() error = %v, want ErrInvalidSchema for invalid component", err)
	}

	base.RequiredComponents = nil
	if err := base.ValidateWithLimits(Limits{MaxBytes: 1}); !errors.Is(err, ErrInvalidSchema) {
		t.Fatalf("ValidateWithLimits() error = %v, want ErrInvalidSchema for byte limit", err)
	}
}

func TestUnit_ValueValidationParity(t *testing.T) {
	viewWith := func(value Value) ViewSchema {
		regions := emptyRegions()
		regions.Content = []Node{{ID: "value", Component: "text", Props: map[string]Value{"value": value}}}
		return ViewSchema{ProtocolVersion: ProtocolVersion, ID: "values", Regions: regions}
	}

	tooManyItems := make([]any, DefaultLimits().MaxArrayItems+1)
	if err := viewWith(Literal(tooManyItems)).Validate(); !errors.Is(err, ErrInvalidSchema) {
		t.Fatalf("Validate() error = %v, want ErrInvalidSchema for array limit", err)
	}

	tooManyKeys := make(map[string]any, DefaultLimits().MaxObjectKeys+1)
	for index := 0; index <= DefaultLimits().MaxObjectKeys; index++ {
		tooManyKeys[fmt.Sprintf("key-%d", index)] = true
	}
	if err := viewWith(Literal(tooManyKeys)).Validate(); !errors.Is(err, ErrInvalidSchema) {
		t.Fatalf("Validate() error = %v, want ErrInvalidSchema for object limit", err)
	}

	if err := viewWith(Literal(map[string]any{"__proto__": true})).Validate(); !errors.Is(err, ErrInvalidSchema) {
		t.Fatalf("Validate() error = %v, want ErrInvalidSchema for unsafe key", err)
	}

	longPath := strings.Repeat("a", DefaultLimits().MaxPathLength+1)
	if err := viewWith(StateRef(longPath)).Validate(); !errors.Is(err, ErrInvalidSchema) {
		t.Fatalf("Validate() error = %v, want ErrInvalidSchema for path limit", err)
	}

	if err := viewWith(Literal([]any{1, 2})).ValidateWithLimits(Limits{MaxArrayItems: 1}); !errors.Is(err, ErrInvalidSchema) {
		t.Fatalf("ValidateWithLimits() error = %v, want ErrInvalidSchema for custom array limit", err)
	}
}

func FuzzValueUnmarshal(f *testing.F) {
	f.Add([]byte(`{"$state":"form.email"}`))
	f.Add([]byte(`{"$eq":[1,2]}`))
	f.Fuzz(func(t *testing.T, data []byte) {
		var value Value
		_ = json.Unmarshal(data, &value)
	})
}
