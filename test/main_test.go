package renovate_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	if err := preflightRenovate(); err != nil {
		fmt.Fprintf(os.Stderr, "renovate preflight failed: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func preflightRenovate() error {
	root, err := filepath.Abs("..")
	if err != nil {
		return fmt.Errorf("resolve repo root: %w", err)
	}
	testDir := filepath.Join(root, "test")
	renovateBin := renovateBinPath(testDir)

	if _, err := os.Stat(renovateBin); os.IsNotExist(err) {
		cmd := exec.Command("npm", "ci", "--prefer-offline", "--no-audit", "--no-fund")
		cmd.Dir = testDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("npm ci: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("stat renovate binary: %w", err)
	}

	out, err := exec.Command(renovateBin, "--version").CombinedOutput()
	if err != nil {
		return fmt.Errorf("renovate --version failed: %w (output: %s)", err, strings.TrimSpace(string(out)))
	}
	if strings.TrimSpace(string(out)) == "" {
		return fmt.Errorf("renovate --version produced empty output")
	}

	return nil
}
