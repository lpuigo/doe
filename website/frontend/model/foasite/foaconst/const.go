package foaconst

const (
	FsStatusNew        string = "00 New"
	FsStatusInProgress string = "20 InProgress"
	FsStatusBlocked    string = "90 Blocked"
	FsStatusCancelled  string = "98 Canceled"
	FsStatusDone       string = "99 Done"

	FsStatusLabelNew        string = "Nouveau"
	FsStatusLabelInProgress string = "En cours"
	FsStatusLabelBlocked    string = "Bloqué"
	FsStatusLabelCancelled  string = "Annulé"
	FsStatusLabelDone       string = "Terminé"
)

const (
	StateToDo       string = "10 To Do"
	StateIncident   string = "25 Incident"
	StateDone       string = "90 Done"
	StateAttachment string = "95 Attachment"
	StateCancelled  string = "99 Cancelled"

	LabelToDo       string = "A faire"
	LabelIncident   string = "Incident"
	LabelDone       string = "Fait"
	LabelAttachment string = "Attach. fait"
	LabelCancelled  string = "Annulé"
)

const (
	FilterValueAll     string = ""
	FilterValueRef     string = "REF:"
	FilterValueInsee   string = "INS:"
	FilterValueType    string = "TYP:"
	FilterValueComment string = "CMT:"

	FilterLabelAll     string = "Tout"
	FilterLabelRef     string = "Référence"
	FilterLabelInsee   string = "Insee"
	FilterLabelType    string = "Type"
	FilterLabelComment string = "Commentaire"
)
