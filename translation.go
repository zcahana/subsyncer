package main

import (
	"github.com/kkdai/mstranslator"
	"strings"
)

type Translator interface {
	Translate(subtitle *SubtitleFile, from, to string) (*SubtitleFile, error)
}

func NewMicrosoftTranslator(clientID, clientSecret string) Translator {
	return &microsoftTranslator{
		client: mstranslator.NewClient(clientID, clientSecret),
	}
}

type microsoftTranslator struct {
	client *mstranslator.Client
}

func (t *microsoftTranslator) Translate(subtitle *SubtitleFile, from, to string) (*SubtitleFile, error) {
	// TODO: Implement bulk translation using TranslateArray
	tSubtitle := &SubtitleFile{
		Entries: make([]*SubtitleEntry, len(subtitle.Entries)),
	}

	for i, entry := range subtitle.Entries {
		tEntry := &SubtitleEntry{
			Index: entry.Index,
			Start: entry.Start,
			End: entry.End,
			Text: make([]string, 0, 1),
		}

		text := strings.Join(entry.Text, " ")
		tText, err := t.client.Translate(text, from, to)

		if err != nil {
			// TODO: better error handling, e.g. skip entries
			// until a threshold is reached
			return nil, err
		}

		entry.Text = append(entry.Text, tText)
		tSubtitle.Entries[i] = tEntry
	}

	return tSubtitle, nil
}
