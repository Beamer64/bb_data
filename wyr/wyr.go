// Package wyr serves Would You Rather polls from the embedded WYR dataset.
// Each poll has two options with historical vote counts from the source —
// useful for showing "what other people picked" stats next to the choice.
// Call Load once at startup with the embed.FS from the parent package,
// then call Random from anywhere.
package wyr

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"strconv"

	"github.com/Beamer64/bb_data/internal/pick"
)

// Poll is a single Would You Rather entry.
type Poll struct {
	OptionA string
	VotesA  int
	OptionB string
	VotesB  int
}

var pool []Poll

// Load reads datasets/WYR.csv from fsys into memory. Call once at startup.
// After Load returns nil, Random and Count are safe for concurrent use.
func Load(fsys fs.FS) error {
	f, err := fsys.Open("datasets/WYR.csv")
	if err != nil {
		return err
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return fmt.Errorf("wyr: read csv: %w", err)
	}

	polls, err := parseRecords(records)
	if err != nil {
		return err
	}
	pool = polls
	return nil
}

// Random returns a uniformly random Poll. Returns the zero-value Poll
// (empty options) when Load has not been called or the dataset is
// empty, so callers can guard on OptionA == "".
func Random() Poll {
	return pick.Random(pool)
}

// Count returns the number of loaded polls. 0 means Load wasn't called
// or the dataset was empty.
func Count() int {
	return len(pool)
}

// parseRecords converts CSV records (header + rows) into polls.
// Malformed rows (too few columns, non-numeric vote counts) are logged
// and skipped rather than aborting the whole load. Returns an error
// only when no usable rows remain.
//
// Pure function — no I/O, suitable for unit tests.
func parseRecords(records [][]string) ([]Poll, error) {
	if len(records) <= 1 {
		return nil, errors.New("no data rows found in WYR CSV")
	}

	polls := make([]Poll, 0, len(records)-1)
	for rowNum, row := range records[1:] {
		if len(row) < 4 {
			log.Printf("WYR CSV row %d skipped (only %d columns)", rowNum+2, len(row))
			continue
		}
		votesA, errA := strconv.Atoi(row[1])
		votesB, errB := strconv.Atoi(row[3])
		if errA != nil || errB != nil {
			log.Printf("WYR CSV row %d skipped (non-numeric vote: %v / %v)", rowNum+2, errA, errB)
			continue
		}
		polls = append(polls, Poll{
			OptionA: row[0],
			VotesA:  votesA,
			OptionB: row[2],
			VotesB:  votesB,
		})
	}

	if len(polls) == 0 {
		return nil, errors.New("no valid WYR rows after parsing")
	}
	return polls, nil
}
