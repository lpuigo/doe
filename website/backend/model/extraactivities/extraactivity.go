package extraactivities

type ExtraActivity struct {
	Name           string
	State          string
	NbPoints       float64
	Income         float64
	Date           string
	AttachmentDate string
	Actors         []string
	Comment        string
}
