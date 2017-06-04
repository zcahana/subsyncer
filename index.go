package main

import (
	"strconv"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
)

type IndexedSubtitle interface {
	Search(text string) (*SubtitleEntry, error)
}

func NewIndexedSubtitle(subtitle *SubtitleFile) (IndexedSubtitle, error) {
	// TODO: Tune the IndexMapping as needed
	index, err := bleve.NewMemOnly(bleve.NewIndexMapping())
	if err != nil {
		return nil, err
	}

	bis := &bleveIndexedSubtitle{
		subtitle: subtitle,
		index:    index,
	}

	err = bis.initialize()
	if err != nil {
		return nil, err
	}

	return bis, nil
}

type bleveIndexedSubtitle struct {
	subtitle *SubtitleFile
	index    bleve.Index
}

func (bis *bleveIndexedSubtitle) initialize() error {
	for i, entry := range bis.subtitle.Entries {
		docId := strconv.Itoa(i)
		docContent := strings.Join(entry.Text, " ")
		err := bis.index.Index(docId, docContent)
		if err != nil {
			return err
		}
	}

	return nil
}

func (bis *bleveIndexedSubtitle) Search(text string) (*SubtitleEntry, error) {
	q := query.NewQueryStringQuery(text)
	req := bleve.NewSearchRequestOptions(q, 1, 0, false)

	res, err := bis.index.Search(req)
	if err != nil {
		return nil, err
	}

	if len(res.Hits) < 1 {
		// TODO: should that be an error?
		return nil, nil
	}

	hit := res.Hits[0]

	i, err := strconv.Atoi(hit.ID)
	if err != nil {
		return nil, err
	}

	return bis.subtitle.Entries[i], nil
}
