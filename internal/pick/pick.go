// Package pick holds tiny helpers shared by every dataset subpackage:
// uniform random selection plus filesystem readers for JSON, JSONL,
// and line-delimited text files. Internal by design — not part of
// the public API.
package pick

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"math/rand/v2"
	"strings"
)

// Random returns a uniformly random element of xs. Returns the zero
// value of T when xs is empty so callers don't have to special-case
// the unloaded state.
func Random[T any](xs []T) T {
	var zero T
	if len(xs) == 0 {
		return zero
	}
	return xs[rand.IntN(len(xs))]
}

// LoadJSON reads path from fsys and decodes the JSON into out.
func LoadJSON(fsys fs.FS, path string, out any) error {
	b, err := fs.ReadFile(fsys, path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}

// LoadJSONL reads path from fsys as a stream of JSON values and decodes
// each into a value of type T. Tolerates JSON Lines (one value per
// line), concatenated values with no separator, and arbitrary whitespace
// between values — useful because hand-curated JSONL files occasionally
// have a missing newline between records.
func LoadJSONL[T any](fsys fs.FS, path string) ([]T, error) {
	f, err := fsys.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []T
	dec := json.NewDecoder(f)
	for {
		var v T
		err := dec.Decode(&v)
		if errors.Is(err, io.EOF) {
			return out, nil
		}
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
}

// LoadLines reads path from fsys and returns each non-empty line with
// trailing \r stripped. Handles both LF and CRLF input.
func LoadLines(fsys fs.FS, path string) ([]string, error) {
	b, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, err
	}
	raw := strings.Split(string(b), "\n")
	out := make([]string, 0, len(raw))
	for _, line := range raw {
		line = strings.TrimRight(line, "\r")
		if line == "" {
			continue
		}
		out = append(out, line)
	}
	return out, nil
}
