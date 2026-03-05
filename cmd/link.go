package cmd

import (
	"fmt"

	"github.com/silee-tools/mydesk/internal/config"
	"github.com/silee-tools/mydesk/internal/linker"
	"github.com/silee-tools/mydesk/internal/native"
	"github.com/silee-tools/mydesk/internal/ui"
)

func RunLink(opts GlobalOpts) error {
	cfg, err := config.Load(opts.ConfigDir)
	if err != nil {
		return err
	}

	entries, err := collectEntries(cfg)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		ui.Warning("No link entries found")
		return nil
	}

	ui.Info("Linking %d entries from %s", len(entries), cfg.ConfigDir)
	fmt.Println()

	r := linker.New(entries, opts.DryRun).Link()

	fmt.Println()
	ui.Info("Done: %d linked, %d skipped, %d backed up, %d failed",
		r.Linked, r.Skipped, r.BackedUp, r.Failed)

	if r.Failed > 0 {
		return fmt.Errorf("%d link(s) failed", r.Failed)
	}
	return nil
}

func collectEntries(cfg *config.Config) ([]linker.LinkEntry, error) {
	nativeEntries, err := native.Scan(cfg)
	if err != nil {
		return nil, fmt.Errorf("scanning native directories: %w", err)
	}

	cfgAdapter := &configAdapter{cfg}
	customEntries, err := linker.ParseLinksConf(cfg.LinksConfPath(), cfgAdapter)
	if err != nil {
		return nil, fmt.Errorf("parsing links.conf: %w", err)
	}

	return linker.MergeEntries(nativeEntries, customEntries), nil
}

// configAdapter adapts config.Config to linker.ConfigDirProvider
type configAdapter struct {
	cfg *config.Config
}

func (a *configAdapter) ExpandVars(path string) string {
	return a.cfg.ExpandVars(path)
}

func (a *configAdapter) GetConfigDir() string {
	return a.cfg.ConfigDir
}
