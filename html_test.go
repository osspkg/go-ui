package ui

import "testing"

func TestUnit_HTMLBuilderProducesTagName(t *testing.T) {
	for _, tag := range []HTMLTag{HTMLArticle, HTMLDiv, HTMLH1, HTMLInput, HTMLTable, HTMLVideo} {
		t.Run(string(tag), func(t *testing.T) {
			if got := HTML(tag, "node").Build().Component; got != string(tag) {
				t.Fatalf("component = %q, want %q", got, tag)
			}
		})
	}
}
