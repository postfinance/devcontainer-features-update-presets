package renovate_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// repoRoot returns the absolute path to the repository root.
func repoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	return root
}

// startPresetAndDatasourceMockServer starts a single HTTP server that serves
// both the rewritten preset file and mocked datasource JSON fixtures.
func startPresetAndDatasourceMockServer(t *testing.T, repoRoot, fixtureDir string) string {
	t.Helper()

	rawPreset, err := os.ReadFile(filepath.Join(repoRoot, "renovate-preset.json"))
	if err != nil {
		t.Fatalf("read preset: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/renovate-preset.json", func(w http.ResponseWriter, r *http.Request) {
		baseURL := "http://" + r.Host
		rewritten := strings.ReplaceAll(string(rawPreset), "https://versionhistory.googleapis.com", baseURL)
		rewritten = strings.ReplaceAll(rewritten, "https://product-details.mozilla.org", baseURL)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(rewritten))
	})

	serveFixture := func(filename string) http.HandlerFunc {
		data, err := os.ReadFile(filepath.Join(fixtureDir, filename))
		if err != nil {
			t.Fatalf("read fixture %s: %v", filename, err)
		}
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(data)
		}
	}

	mux.HandleFunc("/v1/chrome/platforms/linux/channels/stable/versions", serveFixture("google-chrome-versions.json"))
	mux.HandleFunc("/1.0/firefox_history_stability_releases.json", serveFixture("firefox_history_stability_releases.json"))

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	srv := &http.Server{Handler: mux}
	go srv.Serve(ln) //nolint:errcheck
	t.Cleanup(func() { _ = srv.Close() })

	return fmt.Sprintf("http://127.0.0.1:%d", ln.Addr().(*net.TCPAddr).Port)
}

// renovateBinPath returns the path to the local renovate binary in testDir.
func renovateBinPath(testDir string) string {
	return filepath.Join(testDir, "node_modules", ".bin", "renovate")
}

// requireRenovateInstalled returns the renovate binary path and fails if it is missing.
func requireRenovateInstalled(t *testing.T, testDir string) string {
	t.Helper()
	bin := renovateBinPath(testDir)
	if _, err := os.Stat(bin); err != nil {
		t.Fatalf("renovate binary not found at %s: %v", bin, err)
	}
	return bin
}

// dep is a single dependency extracted or looked up by Renovate.
type dep struct {
	DepName      string `json:"depName"`
	CurrentValue string `json:"currentValue"`
	Datasource   string `json:"datasource"`
	SkipReason   string `json:"skipReason"`
	Warnings     []struct {
		Message string `json:"message"`
	} `json:"warnings"`
	Updates []struct {
		NewValue string `json:"newValue"`
	} `json:"updates"`
}

// parseDeps reads a Renovate JSONL log and returns all extracted deps
// grouped by manager name (e.g. "regex", "devcontainer").
func parseDeps(t *testing.T, logFile string) map[string][]dep {
	t.Helper()
	f, err := os.Open(logFile)
	if err != nil {
		t.Fatalf("open log %s: %v", logFile, err)
	}
	defer f.Close()

	// The packageFiles log line can be very large — use a 10 MB buffer.
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 10*1024*1024), 10*1024*1024)

	var result map[string][]dep
	for scanner.Scan() {
		line := scanner.Text()
		// extract mode: {"packageFiles":{"regex":[...]}}
		// lookup mode:  {"msg":"packageFiles with updates","config":{"regex":[...]}}
		if !strings.Contains(line, `"packageFiles"`) && !strings.Contains(line, `"packageFiles with updates"`) {
			continue
		}
		type pkgFileEntry struct {
			Deps []dep `json:"deps"`
		}
		var entry struct {
			PackageFiles map[string][]pkgFileEntry `json:"packageFiles"`
			Config       map[string][]pkgFileEntry `json:"config"`
		}
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		pf := entry.PackageFiles
		if pf == nil {
			pf = entry.Config
		}
		if pf == nil {
			continue
		}
		result = map[string][]dep{}
		for manager, files := range pf {
			for _, f := range files {
				result[manager] = append(result[manager], f.Deps...)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan log: %v", err)
	}
	if result == nil {
		t.Fatal("no packageFiles entry found in Renovate log")
	}
	return result
}

// runRenovate runs Renovate in --platform=local mode against caseDir using the
// preset served at baseURL. dryRunMode is passed to --dry-run (e.g. "extract" or
// "lookup"). Any additional environment variables can be supplied via extraEnv.
func runRenovate(t *testing.T, renovateBin, configFile, caseDir, baseURL, dryRunMode string, extraEnv ...string) string {
	t.Helper()
	tmp, err := os.CreateTemp("", "renovate-*.jsonl")
	if err != nil {
		t.Fatalf("create temp log: %v", err)
	}
	_ = tmp.Close()
	logFile := tmp.Name()
	// t.Cleanup(func() { _ = os.Remove(logFile) })

	cmd := exec.Command(renovateBin, "--platform=local", "--dry-run="+dryRunMode)
	cmd.Dir = caseDir
	cmd.Env = append(os.Environ(),
		"LOG_LEVEL=debug",
		"RENOVATE_CONFIG_FILE="+configFile,
		fmt.Sprintf("RENOVATE_EXTENDS=[%q]", baseURL+"/renovate-preset.json"),
		"RENOVATE_LOG_FILE="+logFile,
	)
	cmd.Env = append(cmd.Env, extraEnv...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("renovate output:\n%s", out)
		t.Fatalf("renovate exited with error: %v", err)
	}
	return logFile
}

// assertDep fails the test if want is not present in the got slice.
func assertDep(t *testing.T, got []dep, want dep) {
	t.Helper()
	for _, d := range got {
		if d.DepName == want.DepName && d.CurrentValue == want.CurrentValue && d.Datasource == want.Datasource {
			return
		}
	}
	t.Errorf("dep not found: %+v\ngot: %+v", want, got)
}

// assertDatasourceResolved fails if the named dep is missing or has skipReason=datasource-error.
func assertDatasourceResolved(t *testing.T, got []dep, depName string) {
	t.Helper()
	for _, d := range got {
		if d.DepName == depName {
			if d.SkipReason == "datasource-error" {
				t.Errorf("dep %q has skipReason=datasource-error (datasource lookup failed)", depName)
			}
			for _, w := range d.Warnings {
				if strings.Contains(strings.ToLower(w.Message), "failed to look up") {
					t.Errorf("dep %q has lookup warning: %s", depName, w.Message)
				}
			}
			return
		}
	}
	t.Errorf("dep %q not found in results", depName)
}

// assertDepHasUpdates ensures Renovate found at least one candidate update for a dep.
func assertDepHasUpdates(t *testing.T, got []dep, depName string) {
	t.Helper()
	for _, d := range got {
		if d.DepName == depName {
			if len(d.Updates) == 0 {
				t.Errorf("dep %q has no updates; expected updates from mocked datasource", depName)
			}
			return
		}
	}
	t.Errorf("dep %q not found in results", depName)
}

// assertDepProposesUpdate ensures Renovate proposes a specific newValue for a dep.
func assertDepProposesUpdate(t *testing.T, got []dep, depName, newValue string) {
	t.Helper()
	for _, d := range got {
		if d.DepName != depName {
			continue
		}
		for _, u := range d.Updates {
			if u.NewValue == newValue {
				return
			}
		}
		t.Errorf("dep %q does not propose update %q (updates=%+v)", depName, newValue, d.Updates)
		return
	}
	t.Errorf("dep %q not found in results", depName)
}

// testEnv holds the common test infrastructure shared across test functions.
type testEnv struct {
	root, testDir, baseURL, renovateBin, configFile string
}

// newTestEnv sets up the shared test infrastructure.
func newTestEnv(t *testing.T) testEnv {
	t.Helper()
	root := repoRoot(t)
	testDir := filepath.Join(root, "test")
	baseURL := startPresetAndDatasourceMockServer(t, root, filepath.Join(testDir, "testdata", "mock-datasources"))
	return testEnv{
		root:        root,
		testDir:     testDir,
		baseURL:     baseURL,
		renovateBin: requireRenovateInstalled(t, testDir),
		configFile:  filepath.Join(testDir, "renovate.json"),
	}
}
