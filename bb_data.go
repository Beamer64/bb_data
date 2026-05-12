// Package bb_data provides in-process access to curated datasets
// (jokes, facts, etc.) embedded in the binary. Consumers import a
// subpackage per dataset and call its Random or other accessors after
// invoking Load once at startup.
package bb_data

import (
	"embed"

	"github.com/Beamer64/bb_data/affirmations"
	"github.com/Beamer64/bb_data/eightball"
	"github.com/Beamer64/bb_data/facts"
	"github.com/Beamer64/bb_data/jokes"
	"github.com/Beamer64/bb_data/kanye"
	"github.com/Beamer64/bb_data/pickuplines"
	"github.com/Beamer64/bb_data/roasts"
	"github.com/Beamer64/bb_data/yomomma"
)

//go:embed datasets
var FS embed.FS

// Load loads every supported dataset into memory. Call exactly once at
// program startup. Returns the first error encountered so callers can
// fail fast before serving traffic. After Load returns nil, subpackage
// accessors are safe for concurrent use.
func Load() error {
	if err := yomomma.Load(FS); err != nil {
		return err
	}
	if err := roasts.Load(FS); err != nil {
		return err
	}
	if err := eightball.Load(FS); err != nil {
		return err
	}
	if err := pickuplines.Load(FS); err != nil {
		return err
	}
	if err := jokes.Load(FS); err != nil {
		return err
	}
	if err := facts.Load(FS); err != nil {
		return err
	}
	if err := affirmations.Load(FS); err != nil {
		return err
	}
	if err := kanye.Load(FS); err != nil {
		return err
	}
	return nil
}
