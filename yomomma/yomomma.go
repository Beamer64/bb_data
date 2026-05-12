// Package yomomma serves a random "yo mama" joke from the embedded
// dataset. Call Load once at startup with the embed.FS from the parent
// package, then call Random from anywhere.
package yomomma

import (
	"io/fs"

	"github.com/Beamer64/bb_data/internal/pick"
)

type yomommaFile struct {
	Data []struct {
		Description string `json:"description"`
	} `json:"data"`
}

var pool []string

// Load reads datasets/yomomma.json from fsys into memory. Call once at
// startup. After Load returns nil, Random is safe for concurrent use.
func Load(fsys fs.FS) error {
	var f yomommaFile
	if err := pick.LoadJSON(fsys, "datasets/yomomma.json", &f); err != nil {
		return err
	}
	p := make([]string, 0, len(f.Data))
	for _, j := range f.Data {
		p = append(p, j.Description)
	}
	pool = p
	return nil
}

// Random returns a uniformly random yo-mama joke. Returns "" if Load
// has not been called or the dataset is empty.
func Random() string {
	return pick.Random(pool)
}
