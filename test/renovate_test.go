package renovate_test

import (
	"path/filepath"
	"testing"
)

func TestCustomManagersRegexExtraction(t *testing.T) {
	type expect struct {
		dep dep
	}

	tests := []struct {
		name    string
		caseDir string
		expects []expect
	}{
		{
			name:    "browsers",
			caseDir: "features-custom-managers/browsers",
			expects: []expect{
				{dep: dep{DepName: "google-chrome", CurrentValue: "126.0.6478.182", Datasource: "custom.browsers-google-chrome"}},
				{dep: dep{DepName: "firefox", CurrentValue: "128.0.3", Datasource: "custom.browsers-firefox"}},
			},
		},
		{
			name:    "build-essential",
			caseDir: "features-custom-managers/build-essential",
			expects: []expect{
				{dep: dep{DepName: "build-essential", CurrentValue: "12.10", Datasource: "deb"}},
			},
		},
		{
			name:    "claude-code",
			caseDir: "features-custom-managers/claude-code",
			expects: []expect{
				{dep: dep{DepName: "anthropics/claude-code", CurrentValue: "0.7.5", Datasource: "github-tags"}},
			},
		},
		{
			name:    "docker-out",
			caseDir: "features-custom-managers/docker-out",
			expects: []expect{
				{dep: dep{DepName: "docker/cli", CurrentValue: "27.1.0", Datasource: "github-tags"}},
				{dep: dep{DepName: "docker/compose", CurrentValue: "2.29.2", Datasource: "github-tags"}},
				{dep: dep{DepName: "docker/buildx", CurrentValue: "0.16.2", Datasource: "github-tags"}},
			},
		},
		{
			name:    "dotnet",
			caseDir: "features-custom-managers/dotnet",
			expects: []expect{
				{dep: dep{DepName: "dotnet-sdk", CurrentValue: "8.0.7", Datasource: "dotnet-version"}},
			},
		},
		{
			name:    "git-lfs",
			caseDir: "features-custom-managers/git-lfs",
			expects: []expect{
				{dep: dep{DepName: "git-lfs/git-lfs", CurrentValue: "3.5.1", Datasource: "github-tags"}},
			},
		},
		{
			name:    "github-cli",
			caseDir: "features-custom-managers/github-cli",
			expects: []expect{
				{dep: dep{DepName: "cli/cli", CurrentValue: "2.55.0", Datasource: "github-tags"}},
			},
		},
		{
			name:    "github-copilot-cli",
			caseDir: "features-custom-managers/github-copilot-cli",
			expects: []expect{
				{dep: dep{DepName: "github/copilot-cli", CurrentValue: "0.0.31", Datasource: "github-tags"}},
			},
		},
		{
			name:    "gitlab-cli",
			caseDir: "features-custom-managers/gitlab-cli",
			expects: []expect{
				{dep: dep{DepName: "gitlab-org/cli", CurrentValue: "1.42.0", Datasource: "gitlab-tags"}},
			},
		},
		{
			name:    "go",
			caseDir: "features-custom-managers/go",
			expects: []expect{
				{dep: dep{DepName: "go", CurrentValue: "1.23.1", Datasource: "golang-version"}},
			},
		},
		{
			name:    "gonovate",
			caseDir: "features-custom-managers/gonovate",
			expects: []expect{
				{dep: dep{DepName: "roemer/gonovate", CurrentValue: "0.8.3", Datasource: "github-tags"}},
			},
		},
		{
			name:    "goreleaser",
			caseDir: "features-custom-managers/goreleaser",
			expects: []expect{
				{dep: dep{DepName: "goreleaser/goreleaser", CurrentValue: "2.1.0", Datasource: "github-tags"}},
			},
		},
		{
			name:    "jfrog-cli",
			caseDir: "features-custom-managers/jfrog-cli",
			expects: []expect{
				{dep: dep{DepName: "jfrog/jfrog-cli", CurrentValue: "2.59.0", Datasource: "github-tags"}},
			},
		},
		{
			name:    "kubectl",
			caseDir: "features-custom-managers/kubectl",
			expects: []expect{
				{dep: dep{DepName: "kubernetes/kubectl", CurrentValue: "1.30.2", Datasource: "github-tags"}},
				{dep: dep{DepName: "ahmetb/kubectx", CurrentValue: "0.9.5", Datasource: "github-tags"}},
				{dep: dep{DepName: "ahmetb/kubectx", CurrentValue: "0.9.4", Datasource: "github-tags"}},
				{dep: dep{DepName: "derailed/k9s", CurrentValue: "0.32.5", Datasource: "github-tags"}},
				{dep: dep{DepName: "helm/helm", CurrentValue: "3.16.1", Datasource: "github-tags"}},
				{dep: dep{DepName: "kubernetes-sigs/kustomize", CurrentValue: "5.4.2", Datasource: "github-tags"}},
				{dep: dep{DepName: "yannh/kubeconform", CurrentValue: "0.6.7", Datasource: "github-tags"}},
				{dep: dep{DepName: "zegl/kube-score", CurrentValue: "1.19.0", Datasource: "github-tags"}},
			},
		},
		{
			name:    "nginx",
			caseDir: "features-custom-managers/nginx",
			expects: []expect{
				{dep: dep{DepName: "nginx/nginx", CurrentValue: "release-1.27.1", Datasource: "github-tags"}},
			},
		},
		{
			name:    "node",
			caseDir: "features-custom-managers/node",
			expects: []expect{
				{dep: dep{DepName: "node", CurrentValue: "22.5.1", Datasource: "node-version"}},
				{dep: dep{DepName: "npm", CurrentValue: "10.8.2", Datasource: "npm"}},
				{dep: dep{DepName: "yarn", CurrentValue: "1.22.22", Datasource: "npm"}},
				{dep: dep{DepName: "pnpm", CurrentValue: "9.6.0", Datasource: "npm"}},
				{dep: dep{DepName: "corepack", CurrentValue: "0.29.3", Datasource: "npm"}},
			},
		},
		{
			name:    "opencode",
			caseDir: "features-custom-managers/opencode",
			expects: []expect{
				{dep: dep{DepName: "anomalyco/opencode", CurrentValue: "0.4.0", Datasource: "github-tags"}},
			},
		},
		{
			name:    "python",
			caseDir: "features-custom-managers/python",
			expects: []expect{
				{dep: dep{DepName: "python/cpython", CurrentValue: "v3.12.4", Datasource: "github-tags"}},
			},
		},
		{
			name:    "rust",
			caseDir: "features-custom-managers/rust",
			expects: []expect{
				{dep: dep{DepName: "rust-lang/rust", CurrentValue: "1.80.0", Datasource: "github-tags"}},
				{dep: dep{DepName: "rust-lang/rustup", CurrentValue: "1.27.1", Datasource: "github-tags"}},
			},
		},
		{
			name:    "sonar-scanner-cli",
			caseDir: "features-custom-managers/sonar-scanner-cli",
			expects: []expect{
				{dep: dep{DepName: "SonarSource/sonar-scanner-cli", CurrentValue: "6.1.0.4477", Datasource: "github-tags"}},
			},
		},
		{
			name:    "terraform",
			caseDir: "features-custom-managers/terraform",
			expects: []expect{
				{dep: dep{DepName: "hashicorp/terraform", CurrentValue: "1.9.3", Datasource: "github-tags"}},
			},
		},
		{
			name:    "vault-cli",
			caseDir: "features-custom-managers/vault-cli",
			expects: []expect{
				{dep: dep{DepName: "hashicorp/vault", CurrentValue: "1.17.2", Datasource: "github-tags"}},
			},
		},
		{
			name:    "zig",
			caseDir: "features-custom-managers/zig",
			expects: []expect{
				{dep: dep{DepName: "ziglang/zig", CurrentValue: "0.13.0", Datasource: "github-tags"}},
			},
		},
	}

	env := newTestEnv(t)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			caseDir := filepath.Join(env.testDir, "testdata", tc.caseDir)
			logFile := runRenovate(t, env.renovateBin, env.configFile, caseDir, env.baseURL, "extract")
			deps := parseDeps(t, logFile)

			for _, e := range tc.expects {
				t.Run(e.dep.DepName+"-"+e.dep.CurrentValue, func(t *testing.T) {
					assertDep(t, deps["regex"], e.dep)
				})
			}

			if got, want := len(deps["regex"]), len(tc.expects); got != want {
				t.Fatalf("unexpected number of regex deps extracted: got=%d want=%d deps=%+v", got, want, deps["regex"])
			}
		})
	}
}

func TestCustomManagersBrowsersDatasourceLookup(t *testing.T) {
	env := newTestEnv(t)
	caseDir := filepath.Join(env.testDir, "testdata", "features-custom-managers", "browsers")

	logFile := runRenovate(t, env.renovateBin, env.configFile, caseDir, env.baseURL, "lookup")
	deps := parseDeps(t, logFile)

	browserDeps := deps["regex"]
	assertDatasourceResolved(t, browserDeps, "google-chrome")
	assertDatasourceResolved(t, browserDeps, "firefox")
	assertDepHasUpdates(t, browserDeps, "google-chrome")
	assertDepHasUpdates(t, browserDeps, "firefox")
	assertDepProposesUpdate(t, browserDeps, "google-chrome", "150.0.7871.100")
	assertDepProposesUpdate(t, browserDeps, "firefox", "152.0.5")
}
