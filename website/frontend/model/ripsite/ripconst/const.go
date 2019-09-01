package ripconst

const (
	StateToDo       string = "00 A faire"
	StateInProgress string = "10 En cours"
	StateBlocked    string = "20 Bloqué"
	StateRedo       string = "30 A reprendre"
	StateWarning2   string = "40 Warning2"
	StateWarning1   string = "45 Warning1"
	StateDone       string = "90 Fait"
	StateCanceled   string = "99 Annulé"
)

const (
	RsStatusNew        string = "00 New"
	RsStatusInProgress string = "20 InProgress"
	RsStatusBlocked    string = "90 Blocked"
	RsStatusCancelled  string = "98 Canceled"
	RsStatusDone       string = "99 Done"
)

const (
	FilterValueAll     string = ""
	FilterValuePtRef   string = "PT:"
	FilterValueTrRef   string = "TR:"
	FilterValueOpe     string = "OPE:"
	FilterValueComment string = "CMT:"

	FilterLabelAll     string = "Tous"
	FilterLabelPtRef   string = "PT"
	FilterLabelTrRef   string = "TR"
	FilterLabelOpe     string = "Opération"
	FilterLabelComment string = "Commentaire"
)
