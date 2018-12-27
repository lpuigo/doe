package model

type Rework struct {
	ControlDate    string
	SubmissionDate string
	CompletionDate string
	Defects        []string
}

func NewRework() *Rework {
	return &Rework{Defects: []string{}}
}
