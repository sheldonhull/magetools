// +build testfiles
package ci_test

import (
	"os"
	"testing"

	iz "github.com/matryer/is"
	"github.com/sheldonhull/magetools/ci"
)

func Test_DetectCI(t *testing.T) {
	tests := []struct {
		name   string
		envVar string
		want   bool
	}{
		{name: "no ci", envVar: "NOTIMPORTANT", want: false},
		{name: "github actions", envVar: "CI", want: true},
		{name: "azure devops", envVar: "AGENT_ID", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := iz.New(t)
			err := os.Setenv(tt.envVar, "not important")
			is.NoErr(err) // setting env variable shouldn't fail

			is.Equal(ci.IsCI(), tt.want) // should correctly detect IsCI
			_ = os.Unsetenv(tt.envVar)
		})
	}
}
