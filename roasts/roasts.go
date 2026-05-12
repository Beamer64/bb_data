// Package roasts serves a random roast from the embedded roasts dataset.
// Call Load once at startup with the embed.FS from the parent package,
// then call Random from anywhere.
package roasts

import (
	"io/fs"

	"github.com/Beamer64/bb_data/internal/pick"
)

var pool []string

// Load reads datasets/roasts.txt from fsys into memory. Call once at
// startup. After Load returns nil, Random is safe for concurrent use.
func Load(fsys fs.FS) error {
	lines, err := pick.LoadLines(fsys, "datasets/roasts.txt")
	if err != nil {
		return err
	}
	pool = lines
	return nil
}

// Random returns a uniformly random roast. Returns "" if Load has not
// been called or the dataset is empty.
func Random() string {
	return pick.Random(pool)
}
