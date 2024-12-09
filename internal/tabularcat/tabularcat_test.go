package tabularcat_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
	"md0.org/reporoot"
	"md0.org/tcat/internal/tabularcat"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"tcat": tabularcat.Main,
	}))
}

func TestScript(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: filepath.Join(reporoot.MustRoot(), "testdata/script"),
	})
}
