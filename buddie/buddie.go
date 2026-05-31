// Package buddie serves a random image of the dog BuddieBot was named
// after, embedded in the binary alongside the text datasets. Call Load
// once at startup with the embed.FS from the parent package, then call
// Random freely from anywhere.
package buddie

import (
	"fmt"
	"io/fs"
	"path"

	"github.com/Beamer64/bb_data/internal/pick"
)

// Image is one picked entry: the file's basename + raw bytes. Files in
// datasets/buddie are JPEG today, but the loader doesn't care about the
// extension — anything in the directory is offered as-is.
type Image struct {
	Filename string
	Data     []byte
}

var pool []Image

// Load reads every file under datasets/buddie/ into memory. Call once at
// startup. After Load returns nil, Random is safe for concurrent use.
// Eager-loads to match the rest of bb_data's "load once, lock-free reads"
// pattern; the embedded folder is small (~2 MiB), so the memory cost is
// trivial vs the per-call simplicity it buys.
func Load(fsys fs.FS) error {
	const dir = "datasets/buddie"
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return fmt.Errorf("buddie: read dir: %w", err)
	}
	imgs := make([]Image, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		p := path.Join(dir, e.Name())
		b, readErr := fs.ReadFile(fsys, p)
		if readErr != nil {
			return fmt.Errorf("buddie: read %s: %w", p, readErr)
		}
		imgs = append(imgs, Image{Filename: e.Name(), Data: b})
	}
	pool = imgs
	return nil
}

// Random returns a uniformly random Image. Returns the zero Image
// (empty Filename, nil Data) when Load has not been called or the
// dataset is empty — callers can guard on Filename == "".
func Random() Image {
	return pick.Random(pool)
}

// Count returns the number of loaded images. 0 means Load wasn't called
// or the dataset directory was empty.
func Count() int {
	return len(pool)
}
