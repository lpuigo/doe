package bpu

import (
	"fmt"
	"sort"
	"strings"
)

type CategoryArticles map[string][]*Article

func NewCategoryChapters() CategoryArticles {
	return make(CategoryArticles)
}

// SortChapters sort CategoryChapters by ascending size order
func (ca CategoryArticles) SortChapters() {
	for cat, chapters := range ca {
		sort.Slice(chapters, func(i, j int) bool {
			return chapters[i].Unit < chapters[j].Unit
		})
		ca[cat] = chapters
	}
}

func (ca CategoryArticles) GetChapterForUnit(category string, unit int) (*Article, error) {
	chapters := ca[strings.ToUpper(category)]
	if len(chapters) == 0 {
		return nil, fmt.Errorf("unknown category '%s'", category)
	}
	for _, p := range chapters {
		if unit <= p.Unit {
			return p, nil
		}
	}
	return chapters[len(chapters)-1], nil
}
