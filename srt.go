package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	initialEntriesCapacity = 500
)

var (
	timestampRegexp = regexp.MustCompile(`(\d{2}):(\d{2}):(\d{2}),(\d{3})\s*-->\s*(\d{2}):(\d{2}):(\d{2}),(\d{3})`)
)

type SRTParser struct{}

// Read the given stream until exhausted, and parse it as an SRT subtitle file.
func (p *SRTParser) Read(reader io.Reader) (*SubtitleFile, error) {
	entries := make([]*SubtitleEntry, 0, initialEntriesCapacity)
	scanner := bufio.NewScanner(reader)
	for {
		entry, err := readEntry(scanner)
		if err != nil {
			return nil, err
		}

		if entry == nil {
			break
		}

		entries = append(entries, entry)
	}

	return &SubtitleFile{
		Entries: entries,
	}, nil
}

// Write the given subtitle file into the given stream, in SRT format.
func (p *SRTParser) Write(subtitle *SubtitleFile, writer io.Writer) error {
	buffer := bufio.NewWriter(writer)

	var err error
	for i, entry := range subtitle.Entries {
		_, err = fmt.Fprintf(buffer, "%d\n", entry.Index)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(buffer, "%s --> %s\n", timestampString(entry.Start), timestampString(entry.End))
		if err != nil {
			return err
		}

		for _, line := range entry.Text {
			_, err = fmt.Fprintf(buffer, "%s\n", line)
			if err != nil {
				return err
			}
		}

		if i < (len(subtitle.Entries) - 1) {
			_, err = fmt.Fprintf(buffer, "\n")
			if err != nil {
				return err
			}
		}
	}

	return buffer.Flush()
}

// readEntry consumes the scanner until a full SRT entry is read.
func readEntry(scanner *bufio.Scanner) (*SubtitleEntry, error) {

	// Skip whitespace
	for isWhitespace(scanner.Text()) {
		next := scanner.Scan()
		if !next {
			if scanner.Err() != nil {
				return nil, scanner.Err()
			}
			// EOF
			return nil, nil
		}
	}

	// Parse index
	index, err := parseIndex(scanner.Text())
	if err != nil {
		return nil, err
	}

	// Advance scanner
	if !scanner.Scan() {
		err := scanner.Err()
		if err == nil {
			return nil, fmt.Errorf("Incomplete SRT entry at index %d", index)
		}
		return nil, err
	}

	// Parse timestamps
	start, end, err := parseTimestamps(scanner.Text())
	if err != nil {
		return nil, err
	}

	// Parse text
	text := make([]string, 0, 2)
	for scanner.Scan() {
		line := scanner.Text()
		if isWhitespace(line) {
			break
		}
		text = append(text, line)
	}

	// Final error check
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	// Done
	return &SubtitleEntry{
		Index: index,
		Start: start,
		End:   end,
		Text:  text,
	}, nil
}

// isWhitespace determines whether the given string is comprised of whitespace only.
func isWhitespace(s string) bool {
	return strings.TrimSpace(s) == ""
}

// parseIndex parses the given string as an SRT entry index, i.e., an integer.
func parseIndex(s string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(s))
}

// parseTimestamps parses the given string as an SRT entry start/end time specification,
// i.e. "hh:mm:ss,iii --> hh:mm:ss,iii".
func parseTimestamps(s string) (time.Duration, time.Duration, error) {
	g := timestampRegexp.FindStringSubmatch(s)
	if g == nil {
		return 0, 0, fmt.Errorf("Invalid subtitle timestamp: %s", s)
	}

	return timestamp(g[1], g[2], g[3], g[4]),
		timestamp(g[5], g[6], g[7], g[8]),
		nil
}

// timestamp constructs a Duration out of the given hours, minutes, seconds and millies,
// given in the form of integer strings.
func timestamp(hours, minutes, seconds, millis string) time.Duration {
	return durationFrom(hours)*time.Hour +
		durationFrom(minutes)*time.Minute +
		durationFrom(seconds)*time.Second +
		durationFrom(millis)*time.Millisecond
}

// timestampString converts the given duration into an SRT style timestamp, i.e. "hh:mm:ss,iii".
func timestampString(d time.Duration) string {
	fmt.Println(d)
	return fmt.Sprintf("%02d:%02d:%02d,%03d",
		d/time.Hour,
		(d%time.Hour)/time.Minute,
		(d%time.Minute)/time.Second,
		(d%time.Second)/time.Millisecond)
}

// durationFrom parses the given string as an integer, and converts it to a Duration.
func durationFrom(s string) time.Duration {
	i, _ := strconv.Atoi(s)
	return time.Duration(i)
}
