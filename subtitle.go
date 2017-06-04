package main

import (
	"time"
	"io"
)

type SubtitleReader interface {
	Read(reader io.Reader) (*SubtitleFile, error)
}

type SubtitleWriter interface {
	Write(subtitle *SubtitleFile, writer io.Writer) error
}

type SubtitleReaderWriter interface {
	SubtitleReader
	SubtitleWriter
}

type SubtitleFile struct {
	Entries []*SubtitleEntry
}

type SubtitleEntry struct {
	Index int
	Start time.Duration
	End   time.Duration
	Text  []string
}

func (f *SubtitleFile) Shift(duration time.Duration) error {
	// TODO implement
	return nil
}

func (f *SubtitleFile) Scale(factor float32) error {
	// TODO implement
	return nil
}