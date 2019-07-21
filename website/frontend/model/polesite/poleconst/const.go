package poleconst

const (
	StateNotSubmitted string = "00 Not Submitted"
	StateNoGo         string = "05 NoGo"
	StateToDo         string = "10 To Do"
	StateHoleDone     string = "20 Hole Done"
	StateIncident     string = "25 Incident"
	StateDone         string = "90 Done"
	StateCancelled    string = "99 Cancelled"

	LabelNotSubmitted string = "Non soumis"
	LabelNoGo         string = "NoGo Client"
	LabelToDo         string = "A faire"
	LabelHoleDone     string = "Trou fait"
	LabelIncident     string = "Incident"
	LabelDone         string = "Fait"
	LabelCancelled    string = "Annulé"

	MaterialWood  string = "Bois"
	MaterialMetal string = "Métal"
	MaterialComp  string = "Composite"

	ProductCoated  string = "Enrobé"
	ProductMoise   string = "Moisé"
	ProductReplace string = "Remplacement"
	ProductRemove  string = "Retrait"
)
