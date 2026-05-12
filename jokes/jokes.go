// Package jokes serves a random general-purpose joke from the embedded
// dataset. Backed by datasets/shortjokes.json — entries are flat
// {id, joke} objects; only the joke text is exposed.
package jokes

import (
	"io/fs"
	"strings"

	"github.com/Beamer64/bb_data/internal/pick"
)

type jokesFile struct {
	Data []struct {
		Joke string `json:"joke"`
	} `json:"data"`
}

var pool []string

// Load reads datasets/shortjokes.json from fsys into memory. Call once
// at startup. After Load returns nil, Random is safe for concurrent use.
func Load(fsys fs.FS) error {
	var f jokesFile
	if err := pick.LoadJSON(fsys, "datasets/shortjokes.json", &f); err != nil {
		return err
	}
	p := make([]string, 0, len(f.Data))
	for _, e := range f.Data {
		j := strings.TrimSpace(e.Joke)
		if j == "" {
			continue
		}
		p = append(p, j)
	}
	pool = p
	return nil
}

// Random returns a uniformly random joke. Returns "" if Load has not
// been called or the dataset is empty.
func Random() string {
	return pick.Random(pool)
}
