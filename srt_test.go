package main

import (
	"bytes"
	"testing"
	"time"
)

func TestSRTParserReadSingleEntry(t *testing.T) {
	content := `1
00:01:15,760 --> 00:01:17,479
Entry 1 line 1
`

	sub, err := (&SRTParser{}).Read(bytes.NewReader([]byte(content)))
	if err != nil {
		t.Fatalf("Expected no error to occur while reading subtitle, got error: %v", err)
	}

	if sub == nil {
		t.Fatalf("Expected returned subtitle to be non-nil")
	}

	if len(sub.Entries) != 1 {
		t.Errorf("Expected 1 entry in subtitle, got %d", len(sub.Entries))
	}

	assertEntry(t, sub.Entries[0], 1, "1m15s760ms", "1m17s479ms", "Entry 1 line 1", "Entry 1 line 2")
}

func TestSRTParserReadSingleEntrySingleLine(t *testing.T) {
	content := `1
00:01:15,760 --> 00:01:17,479
Entry 1 line 1
`

	sub, err := (&SRTParser{}).Read(bytes.NewReader([]byte(content)))
	if err != nil {
		t.Fatalf("Expected no error to occur while reading subtitle, got error: %v", err)
	}

	if sub == nil {
		t.Fatalf("Expected returned subtitle to be non-nil")
	}

	if len(sub.Entries) != 1 {
		t.Errorf("Expected 1 entry in subtitle, got %d", len(sub.Entries))
	}

	assertEntry(t, sub.Entries[0], 1, "1m15s760ms", "1m17s479ms", "Entry 1 line 1")
}

func TestSRTParserReadMultiEntries(t *testing.T) {
	content := `1
00:01:15,760 --> 00:01:17,479
Entry 1 line 1
Entry 1 line 2

2
00:01:20,150 --> 00:01:22,204
Entry 2 line 1
Entry 2 line 2

3
01:01:25,250 --> 01:01:30,000
Entry 3 line 1
Entry 3 line 2
`

	sub, err := (&SRTParser{}).Read(bytes.NewReader([]byte(content)))
	if err != nil {
		t.Fatalf("Expected no error to occur while reading subtitle, got error: %v", err)
	}

	if sub == nil {
		t.Fatalf("Expected returned subtitle to be non-nil")
	}

	if len(sub.Entries) != 3 {
		t.Errorf("Expected 3 entries in subtitle, got %d", len(sub.Entries))
	}

	assertEntry(t, sub.Entries[0], 1, "1m15s760ms", "1m17s479ms", "Entry 1 line 1", "Entry 1 line 2")
	assertEntry(t, sub.Entries[1], 2, "1m20s150ms", "1m22s204ms", "Entry 2 line 1", "Entry 2 line 2")
	assertEntry(t, sub.Entries[2], 3, "1h1m25s250ms", "1h1m30s", "Entry 3 line 1", "Entry 3 line 2")
}

func TestSRTParserWrite(t *testing.T) {
	expected := `1
00:01:15,760 --> 00:01:17,479
Entry 1 line 1
Entry 1 line 2

2
00:01:20,150 --> 00:01:22,204
Entry 2 line 1
Entry 2 line 2

3
01:01:25,250 --> 01:01:30,000
Entry 3 line 1
Entry 3 line 2
`

	subtitle := &SubtitleFile{
		Entries: []*SubtitleEntry{
			{
				Index: 1,
				Start: mustParseDuration("1m15s760ms"),
				End:   mustParseDuration("1m17s479ms"),
				Text:  []string{"Entry 1 line 1", "Entry 1 line 2"},
			},
			{
				Index: 2,
				Start: mustParseDuration("1m20s150ms"),
				End:   mustParseDuration("1m22s204ms"),
				Text:  []string{"Entry 2 line 1", "Entry 2 line 2"},
			},
			{
				Index: 3,
				Start: mustParseDuration("1h1m25s250ms"),
				End:   mustParseDuration("1h1m30s"),
				Text:  []string{"Entry 3 line 1", "Entry 3 line 2"},
			},
		},
	}

	buffer := new(bytes.Buffer)
	err := (&SRTParser{}).Write(subtitle, buffer)
	if err != nil {
		t.Fatalf("Expected no error to occur while writing subtitle, got error: %v", err)
	}

	actual := buffer.String()
	if expected != actual {
		t.Errorf("Expected written subtitle content to be:\n%s\n\ngot:\n%s\n", expected, actual)
	}
}

func assertEntry(t *testing.T, entry *SubtitleEntry, index int, start, end string, text ...string) {
	if entry.Index != index {
		t.Errorf("Expected entry index to be %d, got %d", index, entry.Index)
	}

	startDuration := mustParseDuration(start)
	if entry.Start != startDuration {
		t.Errorf("Expected start timestamp to be at %v, got %v", startDuration, entry.Start)
	}

	endDuration := mustParseDuration(end)
	if entry.End != endDuration {
		t.Errorf("Expected end timestamp to be at %v, got %v", endDuration, entry.End)
	}

	if len(entry.Text) != len(text) {
		t.Errorf("Expected %d text lines, got %d", len(text), len(entry.Text))
	}

	for i, line := range text {
		if entry.Text[i] != line {
			t.Errorf("Expected line %d to be '%s', got '%s'", line, entry.Text[i])
		}
	}
}

func mustParseDuration(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}
