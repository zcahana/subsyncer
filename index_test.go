package main

import "testing"

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

func TestIndexedSubtitleSearchProximity(t *testing.T) {
	indexedSub, err := NewIndexedSubtitle(testSubtitle)
	if err != nil {
		t.Fatalf("Expected no error to occur while indexing subtitle, got error: %v", err)
	}

	resEntry, err := indexedSub.Search("Something bad happend")
	if err != nil {
		t.Fatalf("Got error while searching: %v", err)
	}

	if resEntry == nil {
		t.Fatalf("Got nil entry while searching")
	}

	if !equalSubtitleEntries(testSubtitle.Entries[1], resEntry) {
		t.Errorf("Got wrong entry while searching: Expected %v, got %v", testSubtitle.Entries[1], resEntry)
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
