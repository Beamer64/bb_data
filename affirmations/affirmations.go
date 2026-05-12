// Package affirmations serves a random affirmation from the embedded
// dataset. Backed by datasets/affirmations.jsonl — each line is a
// {prompt, completion} object; only the completion text is exposed
// (the prompt encodes category metadata not needed for random pick).
package affirmations

import (
	"io/fs"
	"strings"

	"github.com/Beamer64/bb_data/internal/pick"
)

type affirmationEntry struct {
	Completion string `json:"completion"`
}

var pool []string

// Load reads datasets/affirmations.jsonl from fsys into memory. Call
// once at startup. After Load returns nil, Random is safe for
// concurrent use.
func Load(fsys fs.FS) error {
	entries, err := pick.LoadJSONL[affirmationEntry](fsys, "datasets/affirmations.jsonl")
	if err != nil {
		return err
	}
	p := make([]string, 0, len(entries))
	for _, e := range entries {
		c := strings.TrimSpace(e.Completion)
		if c == "" {
			continue
		}
		p = append(p, c)
	}
	pool = p
	return nil
}

// Random returns a uniformly random affirmation. Returns "" if Load
// has not been called or the dataset is empty.
func Random() string {
	return pick.Random(pool)
}
