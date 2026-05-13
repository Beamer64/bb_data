// Package tonguetwister serves a random tongue twister from the embedded
// dataset. Backed by datasets/tongue_twisters.txt — one twister per line.
package tonguetwister

import (
	"io/fs"

	"github.com/Beamer64/bb_data/internal/pick"
)

var pool []string

// Load reads datasets/tongue_twisters.txt from fsys into memory. Call once
// at startup. After Load returns nil, Random is safe for concurrent use.
func Load(fsys fs.FS) error {
	lines, err := pick.LoadLines(fsys, "datasets/tongue_twisters.txt")
	if err != nil {
		return err
	}
	pool = lines
	return nil
}

// Random returns a uniformly random tongue twister. Returns "" if Load has
// not been called or the dataset is empty.
func Random() string {
	return pick.Random(pool)
}
