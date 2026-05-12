// Package kanye serves a random Kanye West quote from the embedded
// dataset. Backed by datasets/kanyequotes.json — a flat JSON array of
// strings.
package kanye

import (
	"io/fs"
	"strings"

	"github.com/Beamer64/bb_data/internal/pick"
)

var pool []string

// Load reads datasets/kanyequotes.json from fsys into memory. Call once
// at startup. After Load returns nil, Random is safe for concurrent use.
func Load(fsys fs.FS) error {
	var quotes []string
	if err := pick.LoadJSON(fsys, "datasets/kanyequotes.json", &quotes); err != nil {
		return err
	}
	p := make([]string, 0, len(quotes))
	for _, q := range quotes {
		q = strings.TrimSpace(q)
		if q == "" {
			continue
		}
		p = append(p, q)
	}
	pool = p
	return nil
}

// Random returns a uniformly random Kanye quote. Returns "" if Load
// has not been called or the dataset is empty.
func Random() string {
	return pick.Random(pool)
}
