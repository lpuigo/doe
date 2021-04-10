package extracter

import "github.com/lpuig/ewin/doe/kizeoparser/api"

type FormProgress struct {
	FormConfig
	ExtractionDate string `json:"extraction_date"`
}

type Progress struct {
	Formulaires []*FormProgress `json:"formulaires"`
	forms       map[int][]*api.SearchData
}

func NewProgress() Progress {
	return Progress{
		Formulaires: []*FormProgress{},
		forms:       make(map[int][]*api.SearchData),
	}
}
