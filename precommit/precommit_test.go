package precommit_test

import (
	"os"
	"testing"

	"github.com/magefile/mage/sh"
	iz "github.com/matryer/is"
	"github.com/sheldonhull/magetools/precommit"
)

func Test_ProvideDefaultConfig(t *testing.T) {
	is := iz.New(t)
	executedOutputFile := ".pre-commit-config.yaml"
	defer sh.Rm(executedOutputFile) //nolint:errcheck // not required for this
	config, err := precommit.ProvideDefaultConfig()
	is.NoErr(err)                        // running ProvideDefaultConfig shouldn't fail.
	is.Equal(config, executedOutputFile) // expected output file should exist
	f, err := os.Stat(executedOutputFile)
	is.NoErr(err)                          // .pre-commit-config.yaml should be found.
	is.Equal(f.Name(), executedOutputFile) // .pre-commit-config.yaml should be the filename.
}
