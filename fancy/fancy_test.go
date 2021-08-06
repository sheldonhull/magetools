// fancy uses pterm to provide some nice output that's not really critical but makes the experience nicer with summary
// +build testfiles

package fancy_test

import (
	"testing"

	"github.com/sheldonhull/magetools/fancy"
)

func TestIntroScreen(t *testing.T) {
	tests := []struct {
		name           string
		disableStyling bool
	}{
		{name: "intro screen with color", disableStyling: false},
		{name: "intro screen without color", disableStyling: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fancy.IntroScreen(tt.disableStyling)
		})
	}
}
