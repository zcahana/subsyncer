package main

import (
	"strings"
	"testing"
)

var testSubtitle = &SubtitleFile{
	Entries: []*SubtitleEntry{
		{
			Index: 1,
			Start: mustParseDuration("1m15s760ms"),
			End:   mustParseDuration("1m17s479ms"),
			Text:  []string{"Once upon a time", "in a far away land"},
		},
		{
			Index: 2,
			Start: mustParseDuration("1m20s150ms"),
			End:   mustParseDuration("1m22s204ms"),
			Text:  []string{"Something horrible happend,", "but then it was somehow solved"},
		},
		{
			Index: 3,
			Start: mustParseDuration("1h1m25s250ms"),
			End:   mustParseDuration("1h1m30s"),
			Text:  []string{"And they all lived", "happily ever after..."},
		},
	},
}

func TestIndexedSubtitleSearchExact(t *testing.T) {
	indexedSub, err := NewIndexedSubtitle(testSubtitle)
	if err != nil {
		t.Fatalf("Expected no error to occur while indexing subtitle, got error: %v", err)
	}

	for _, entry := range testSubtitle.Entries {
		for i, text := range entry.Text {
			resEntry, err := indexedSub.Search(text)
			if err != nil {
				t.Fatalf("Got error while searching entry %d (line %d): %v", entry.Index, i+1, err)
			}

			if resEntry == nil {
				t.Fatalf("Got nil entry while searching entry %d (line %d)", entry.Index, i+1)
			}

			if !equalSubtitleEntries(entry, resEntry) {
				t.Errorf("Got wrong result while searching entry %d (line %d): Expected %v, got %v",
					entry.Index, i+1, entry, resEntry)
			}
		}
	}
}

func TestIndexedSubtitleSearchJoinedLines(t *testing.T) {
	indexedSub, err := NewIndexedSubtitle(testSubtitle)
	if err != nil {
		t.Fatalf("Expected no error to occur while indexing subtitle, got error: %v", err)
	}

	for _, entry := range testSubtitle.Entries {
		text := strings.Join(entry.Text, " ")
		resEntry, err := indexedSub.Search(text)
		if err != nil {
			t.Fatalf("Got error while searching entry %d: %v", entry.Index, err)
		}

		if resEntry == nil {
			t.Fatalf("Got nil entry while searching entry %d", entry.Index)
		}

		if !equalSubtitleEntries(entry, resEntry) {
			t.Errorf("Got wrong result while searching entry %d: Expected %v, got %v",
				entry.Index, entry, resEntry)
		}
	}
}

func TestIndexedSubtitleSearchProximity(t *testing.T) {
	indexedSub, err := NewIndexedSubtitle(testSubtitle)
	if err != nil {
		t.Fatalf("Expected no error to occur while indexing subtitle, got error: %v", err)
	}

	cases := []struct {
		text  string
		entry *SubtitleEntry
	}{
		{"in a far land", testSubtitle.Entries[0]},
		{"in a very far land", testSubtitle.Entries[0]},
		{"Something bad happend", testSubtitle.Entries[1]},
		{"Something horrible occurred", testSubtitle.Entries[1]},
		{"They lived", testSubtitle.Entries[2]},
		{"All they lived", testSubtitle.Entries[2]},
	}

	for i, c := range cases {
		resEntry, err := indexedSub.Search(c.text)
		if err != nil {
			t.Fatalf("Got error while searching (case %d): %v", i, err)
		}

		if resEntry == nil {
			t.Fatalf("Got nil entry while searching (case %d)", i)
		}

		if !equalSubtitleEntries(c.entry, resEntry) {
			t.Errorf("Got wrong entry while searching (case %d): Expected %v, got %v", i, c.entry, resEntry)
		}
	}
}

func TestIndexedSubtitleSearchNoMatch(t *testing.T) {
	indexedSub, err := NewIndexedSubtitle(testSubtitle)
	if err != nil {
		t.Fatalf("Expected no error to occur while indexing subtitle, got error: %v", err)
	}

	cases := []struct {
		text string
	}{
		{"This phrases are completely unrelated"},
		{"And must return no match"},
		{"Even if they have are somehow similar terms"},
		{"Or some common ideas from time to time"},
		{"Still it's far enough to be matched"},
	}

	for i, c := range cases {
		resEntry, err := indexedSub.Search(c.text)
		if err != nil {
			t.Fatalf("Got error while searching (case %d): %v", i, err)
		}

		if resEntry != nil {
			t.Fatalf("Expected nil entry while searching (case %d), got entry %d", i, resEntry.Index)
		}
	}
}

func equalSubtitleEntries(e1, e2 *SubtitleEntry) bool {
	if e1.Index != e2.Index {
		return false
	}

	if e1.Start != e2.Start {
		return false
	}

	if e1.End != e2.End {
		return false
	}

	if len(e1.Text) != len(e2.Text) {
		return false
	}

	for i, line := range e1.Text {
		if line != e2.Text[i] {
			return false
		}
	}

	return true
}
