package ui

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestUnit_SecuritySchemaBoundaries(t *testing.T) {
	t.Run(
		"unsafe value path", func(t *testing.T) {
			var value Value
			if err := json.Unmarshal([]byte(`{"$state":"constructor.name"}`), &value); !errors.Is(
				err,
				ErrInvalidValue,
			) {
				t.Fatalf("Unmarshal() error = %v, want ErrInvalidValue", err)
			}
		},
	)
	t.Run(
		"oversized serialized schema", func(t *testing.T) {
			view := ViewSchema{
				ProtocolVersion: ProtocolVersion,
				ID:              "safe",
				Title:           strings.Repeat("x", 128),
				Regions:         Regions{},
			}
			if err := view.ValidateWithLimits(Limits{MaxBytes: 64}); !errors.Is(err, ErrInvalidSchema) {
				t.Fatalf("ValidateWithLimits() error = %v, want ErrInvalidSchema", err)
			}
		},
	)
	t.Run(
		"unsafe action path", func(t *testing.T) {
			view := ViewSchema{
				ProtocolVersion: ProtocolVersion,
				ID:              "safe",
				Actions: map[string]Action{
					"save": {
						Type:  "set-state",
						Path:  "__proto__.polluted",
						Value: ptrValue("x"),
					},
				},
				Regions: Regions{},
			}
			if err := view.Validate(); !errors.Is(err, ErrInvalidSchema) {
				t.Fatalf("Validate() error = %v, want ErrInvalidSchema", err)
			}
		},
	)
}
