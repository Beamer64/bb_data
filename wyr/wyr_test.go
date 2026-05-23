package wyr

import (
	"testing"
)

func TestParseRecords(t *testing.T) {
	t.Run("valid CSV produces all rows", func(t *testing.T) {
		records := [][]string{
			{"option_a", "votes_a", "option_b", "votes_b"}, // header
			{"Pizza", "10", "Tacos", "20"},
			{"Cats", "5", "Dogs", "5"},
		}
		polls, err := parseRecords(records)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(polls) != 2 {
			t.Fatalf("got %d polls, want 2", len(polls))
		}
		if polls[0].OptionA != "Pizza" || polls[0].VotesA != 10 ||
			polls[0].OptionB != "Tacos" || polls[0].VotesB != 20 {
			t.Errorf("first poll mismatch: %+v", polls[0])
		}
	})

	t.Run("empty records returns error", func(t *testing.T) {
		_, err := parseRecords(nil)
		if err == nil {
			t.Fatal("expected error for empty records, got nil")
		}
	})

	t.Run("only header returns error", func(t *testing.T) {
		records := [][]string{{"option_a", "votes_a", "option_b", "votes_b"}}
		_, err := parseRecords(records)
		if err == nil {
			t.Fatal("expected error for header-only CSV, got nil")
		}
	})

	t.Run("rows with too few columns are skipped", func(t *testing.T) {
		records := [][]string{
			{"option_a", "votes_a", "option_b", "votes_b"},
			{"Pizza", "10"}, // truncated — should be skipped
			{"Cats", "5", "Dogs", "5"},
		}
		polls, err := parseRecords(records)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(polls) != 1 {
			t.Fatalf("got %d polls, want 1 (truncated row should be skipped)", len(polls))
		}
		if polls[0].OptionA != "Cats" {
			t.Errorf("expected Cats poll, got %+v", polls[0])
		}
	})

	t.Run("rows with non-numeric votes are skipped", func(t *testing.T) {
		records := [][]string{
			{"option_a", "votes_a", "option_b", "votes_b"},
			{"Pizza", "ten", "Tacos", "20"}, // bad parse — should be skipped
			{"Cats", "5", "Dogs", "abc"},    // bad parse — should be skipped
			{"Yes", "1", "No", "2"},
		}
		polls, err := parseRecords(records)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(polls) != 1 {
			t.Fatalf("got %d polls, want 1 (bad-parse rows should be skipped)", len(polls))
		}
		if polls[0].OptionA != "Yes" {
			t.Errorf("expected Yes poll, got %+v", polls[0])
		}
	})

	t.Run("all rows malformed returns error", func(t *testing.T) {
		records := [][]string{
			{"option_a", "votes_a", "option_b", "votes_b"},
			{"Pizza", "ten", "Tacos", "twenty"},
			{"Cats", "5"}, // truncated
		}
		_, err := parseRecords(records)
		if err == nil {
			t.Fatal("expected error when all rows are unusable, got nil")
		}
	})
}
