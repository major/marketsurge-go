package marketsurge_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestImportSmoke(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	dir := t.TempDir()
	goMod := fmt.Sprintf(
		"module marketsurge_import_smoke\n\ngo %s\n\nrequire github.com/major/marketsurge-go v0.0.0\n\nreplace github.com/major/marketsurge-go => %s\n",
		moduleGoDirective(t),
		root,
	)
	goTest := "package smoke\n\nimport (\n\t\"testing\"\n\n\t\"github.com/major/marketsurge-go/marketsurge\"\n)\n\nfunc TestCoreExplicitSessionImport(t *testing.T) {\n\t_ = t\n\t_ = marketsurge.NewClient\n\t_ = marketsurge.NewSession\n}\n"
	mustWriteFile(t, filepath.Join(dir, "go.mod"), goMod)
	mustWriteFile(t, filepath.Join(dir, "smoke_test.go"), goTest)

	cmd := exec.Command("go", "test", "./...")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("external module import smoke test failed:\n%s", string(output))
	}
}

func TestRootPackageDependencyIsolation(t *testing.T) {
	t.Parallel()

	deps := listPackageDeps(t, "github.com/major/marketsurge-go/marketsurge")
	err := forbidDependencySubstrings(deps, []string{
		"kooky",
		"sqlite",
		"keyring",
		"dbus",
		"browserutils",
	})
	if err != nil {
		t.Fatalf("forbidDependencySubstrings() = %v, want nil", err)
	}
}

func TestRootPackageDependencyIsolationFailsClosed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		deps      []string
		wantError bool
	}{
		{
			name: "clean dependency list",
			deps: []string{
				"github.com/major/marketsurge-go/marketsurge",
				"net/http",
			},
			wantError: false,
		},
		{
			name: "forbidden dependency is rejected",
			deps: []string{
				"github.com/major/marketsurge-go/marketsurge",
				"github.com/browserutils/kooky",
			},
			wantError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := forbidDependencySubstrings(tc.deps, []string{
				"kooky",
				"sqlite",
				"keyring",
				"dbus",
				"browserutils",
			})
			if tc.wantError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), "github.com/browserutils/kooky") {
					t.Errorf("error %q does not contain %q", err.Error(), "github.com/browserutils/kooky")
				}
				return
			}

			if err != nil {
				t.Fatalf("forbidDependencySubstrings() = %v, want nil", err)
			}
		})
	}
}

func moduleRoot(t *testing.T) string {
	t.Helper()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() = %v", err)
	}

	return filepath.Dir(cwd)
}

func moduleGoDirective(t *testing.T) string {
	t.Helper()

	goMod, err := os.ReadFile(filepath.Join(moduleRoot(t), "go.mod"))
	if err != nil {
		t.Fatalf("ReadFile(go.mod) = %v", err)
	}

	for line := range strings.Lines(string(goMod)) {
		if version, ok := strings.CutPrefix(strings.TrimSpace(line), "go "); ok {
			return strings.TrimSpace(version)
		}
	}

	t.Fatal("root go.mod is missing a go directive")
	return ""
}

func mustWriteFile(t *testing.T, path string, contents string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("WriteFile(%q) = %v", path, err)
	}
}

func listPackageDeps(t *testing.T, packagePath string) []string {
	t.Helper()

	root := moduleRoot(t)
	cmd := exec.Command("go", "list", "-deps", "-f", "{{.ImportPath}}", packagePath)
	cmd.Dir = root
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go list -deps failed:\n%s", string(output))
	}

	deps := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(deps) == 1 && deps[0] == "" {
		return nil
	}

	return deps
}

func forbidDependencySubstrings(deps []string, forbidden []string) error {
	for _, dep := range deps {
		for _, bad := range forbidden {
			if strings.Contains(dep, bad) {
				return fmt.Errorf("forbidden dependency %q found in %q", bad, dep)
			}
		}
	}

	return nil
}
