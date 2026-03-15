package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/silee-tools/mydesk/internal/config"
	"github.com/silee-tools/mydesk/internal/shell"
	"github.com/silee-tools/mydesk/internal/ui"
)

func RunInstallShell(opts GlobalOpts) error {
	if opts.ConfigDir == "" && os.Getenv("MYDESK_CONFIG_DIR") == "" {
		return fmt.Errorf("install-shell requires --config-dir or MYDESK_CONFIG_DIR.\n" +
			"Example: mydesk --config-dir ~/my-dotfiles install-shell")
	}

	cfg, err := config.Load(opts.ConfigDir)
	if err != nil {
		return err
	}

	profilePath := shell.DetectProfile()

	// Resolve symlinks to edit the actual file
	realPath, err := filepath.EvalSymlinks(profilePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("cannot resolve %s: %w", profilePath, err)
	}
	if err == nil && realPath != profilePath {
		ui.Info("Resolved %s → %s (symlink)", filepath.Base(profilePath), realPath)
		profilePath = realPath
	}

	// Build block content
	var lines []string
	lines = append(lines, fmt.Sprintf("export MYDESK_CONFIG_DIR=%q", cfg.ConfigDir))

	gobin := resolveGoBinPath()
	if gobin != "" && !isInPath(gobin) {
		lines = append(lines, fmt.Sprintf("export PATH=\"$PATH:%s\"", gobin))
	}

	blockContent := strings.Join(lines, "\n")

	if opts.DryRun {
		ui.DryRun("would write to %s:", profilePath)
		ui.DryRun("  %s", shell.StartMarker)
		for _, line := range lines {
			ui.DryRun("  %s", line)
		}
		ui.DryRun("  %s", shell.EndMarker)
		return nil
	}

	// Read existing content
	existing, err := os.ReadFile(profilePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("cannot read %s: %w", profilePath, err)
	}

	newContent := shell.UpsertBlock(string(existing), shell.StartMarker, shell.EndMarker, blockContent)

	// Preserve file permissions
	mode := os.FileMode(0644)
	if info, err := os.Stat(profilePath); err == nil {
		mode = info.Mode()
	}

	if err := os.WriteFile(profilePath, []byte(newContent), mode); err != nil {
		return fmt.Errorf("cannot write %s: %w", profilePath, err)
	}

	ui.Success("Shell profile updated: %s", profilePath)
	ui.Info("Run: source %s", profilePath)

	return nil
}

// resolveGoBinPath returns the Go binary install directory.
func resolveGoBinPath() string {
	if gobin := os.Getenv("GOBIN"); gobin != "" {
		return gobin
	}
	if gopath := os.Getenv("GOPATH"); gopath != "" {
		return filepath.Join(gopath, "bin")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, "go", "bin")
}

// isInPath checks if dir is already in the PATH.
func isInPath(dir string) bool {
	for _, p := range filepath.SplitList(os.Getenv("PATH")) {
		if p == dir {
			return true
		}
	}
	return false
}
