// Package textfonts maps ASCII characters to stylized Unicode variants
// (cursive, bubble, leet, flipped, cursed). Call Load once at startup
// with the embed.FS from the parent package, then call Convert from
// anywhere to transform a string under a named group.
//
// Some characters in the source dataset map to multiple stylized variants;
// Convert picks one uniformly at random per character, so repeated calls
// with the same input may differ. Characters with no mapping are passed
// through unchanged.
package textfonts

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/Beamer64/bb_data/internal/pick"
)

// groups holds the parsed font tables keyed by group name. The source
// JSON shape is `{group: [{char: [variants]}]}` — the outer array always
// has exactly one element, so we flatten to `{group: {char: [variants]}}`
// at Load time for simpler lookups.
var groups map[string]map[string][]string

// Load reads datasets/text_fonts.json from fsys into memory. Call once
// at startup. After Load returns nil, Convert and Groups are safe for
// concurrent use.
func Load(fsys fs.FS) error {
	b, err := fs.ReadFile(fsys, "datasets/text_fonts.json")
	if err != nil {
		return err
	}

	var raw map[string][]map[string][]string
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("textfonts: parse json: %w", err)
	}

	flat := make(map[string]map[string][]string, len(raw))
	for group, layers := range raw {
		if len(layers) == 0 {
			continue
		}
		// The source JSON wraps each group's char-map in a one-element
		// array. If a future dataset adds more layers, merge them; first
		// occurrence of a char wins.
		merged := make(map[string][]string)
		for _, layer := range layers {
			for char, variants := range layer {
				if _, exists := merged[char]; !exists {
					merged[char] = variants
				}
			}
		}
		flat[group] = merged
	}
	groups = flat
	return nil
}

// Convert returns text with each character substituted according to the
// group's mapping. Characters with no mapping in the group, and all
// characters when the group is unknown, are passed through unchanged.
// Characters that map to multiple variants get a uniformly-random pick
// per call.
func Convert(text, group string) string {
	chars, ok := groups[group]
	if !ok {
		return text
	}
	var b strings.Builder
	b.Grow(len(text))
	for _, r := range text {
		if variants := chars[string(r)]; len(variants) > 0 {
			b.WriteString(pick.Random(variants))
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// Groups returns the loaded font group names sorted alphabetically.
// Returns nil if Load has not been called.
func Groups() []string {
	if len(groups) == 0 {
		return nil
	}
	names := make([]string, 0, len(groups))
	for g := range groups {
		names = append(names, g)
	}
	sort.Strings(names)
	return names
}
