package api

type AdvancedSearchFilter struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Type     string `json:"type"`
	Value    string `json:"val"`
}

func NewAdvancedSearchFilter(field, operator, value string) *AdvancedSearchFilter {
	return &AdvancedSearchFilter{
		Field:    "_" + field,
		Operator: operator,
		Type:     "simple",
		Value:    value,
	}
}

type AdvancedSearchOrder struct {
	Col  string `json:"col"`
	Type string `json:"type"`
}

func NewAdvancedSearchOrder(col string, asc bool) *AdvancedSearchOrder {
	t := "desc"
	if asc {
		t = "asc"
	}
	return &AdvancedSearchOrder{
		Col:  "_" + col,
		Type: t,
	}
}

type AdvancedSearch struct {
	Format        string                  `json:"format"`
	GlobalFilters string                  `json:"global_filters"`
	Filters       []*AdvancedSearchFilter `json:"filters"`
	Order         []*AdvancedSearchOrder  `json:"order"`
}

func NewAdvancedSearch() *AdvancedSearch {
	return &AdvancedSearch{
		Format:        "simple",
		GlobalFilters: "",
		Filters:       []*AdvancedSearchFilter{},
		Order:         []*AdvancedSearchOrder{},
	}
}

func (as *AdvancedSearch) SetFilters(filters ...*AdvancedSearchFilter) *AdvancedSearch {
	as.Filters = filters
	return as
}

func (as *AdvancedSearch) SetOrder(orders ...*AdvancedSearchOrder) *AdvancedSearch {
	as.Order = orders
	return as
}
