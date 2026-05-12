// Package eightball serves a random Magic 8-Ball response from the
// embedded dataset. The package name is "eightball" because Go
// identifiers cannot start with a digit; consumers should still use
// the slash-command label "8ball" wherever it's user-facing.
package eightball

import (
	"io/fs"

	"github.com/Beamer64/bb_data/internal/pick"
)

var pool []string

// Load reads datasets/8ball.txt from fsys into memory. Call once at
// startup. After Load returns nil, Random is safe for concurrent use.
func Load(fsys fs.FS) error {
	lines, err := pick.LoadLines(fsys, "datasets/8ball.txt")
	if err != nil {
		return err
	}
	pool = lines
	return nil
}

// Random returns a uniformly random 8-Ball response. Returns "" if
// Load has not been called or the dataset is empty.
func Random() string {
	return pick.Random(pool)
}
