package model

type Defect struct {
	PT             string
	SubmissionDate string
	Description    string
	NbOK           int
	NbKO           int
	ToBeFixed      bool
	FixDate        string
	FixActor       string
}

type Rework struct {
	ControlDate    string
	SubmissionDate string
	CompletionDate string
	Defects        []Defect
}

func NewRework() *Rework {
	return &Rework{Defects: []Defect{}}
}
