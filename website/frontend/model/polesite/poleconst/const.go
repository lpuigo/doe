package poleconst

const (
	PsStatusNew        string = "00 New"
	PsStatusInProgress string = "20 InProgress"
	PsStatusBlocked    string = "90 Blocked"
	PsStatusCancelled  string = "98 Canceled"
	PsStatusDone       string = "99 Done"

	PsStatusLabelNew        string = "Nouveau"
	PsStatusLabelInProgress string = "En cours"
	PsStatusLabelBlocked    string = "Bloqué"
	PsStatusLabelCancelled  string = "Annulé"
	PsStatusLabelDone       string = "Terminé"
)

const (
	StateNotSubmitted string = "00 Not Submitted"
	StateNoGo         string = "05 NoGo"
	StateDictToDo     string = "08 DICT To Do"
	StateToDo         string = "10 To Do"
	StateHoleDone     string = "20 Hole Done"
	StateIncident     string = "25 Incident"
	StateDone         string = "90 Done"
	StateAttachment   string = "95 Attachment"
	StateCancelled    string = "99 Cancelled"

	LabelNotSubmitted string = "Non soumis"
	LabelNoGo         string = "NoGo Client"
	LabelDictToDo     string = "DICT à faire"
	LabelToDo         string = "A faire"
	LabelHoleDone     string = "Trou fait"
	LabelIncident     string = "Incident"
	LabelDone         string = "Fait"
	LabelAttachment   string = "Attach. fait"
	LabelCancelled    string = "Annulé"

	FilterValueAll      string = ""
	FilterValueRef      string = "REF:"
	FilterValueCity     string = "CTY:"
	FilterValueAddr     string = "ADD:"
	FilterValueComment  string = "CMT:"
	FilterValueHeigth   string = "HGT:"
	FilterValueProduct  string = "PRD:"
	FilterValueDt       string = "DT:"
	FilterValueDict     string = "DCT:"
	FilterValueDictInfo string = "DCI:"

	FilterLabelAll      string = "Tout"
	FilterLabelRef      string = "Référence"
	FilterLabelCity     string = "Ville"
	FilterLabelAddr     string = "Adresse"
	FilterLabelComment  string = "Commentaire"
	FilterLabelHeigth   string = "Hauteur"
	FilterLabelProduct  string = "Produit"
	FilterLabelDt       string = "DT"
	FilterLabelDict     string = "DICT"
	FilterLabelDictInfo string = "DICT Info"

	MaterialWood  string = "Bois"
	MaterialMetal string = "Métal"
	MaterialComp  string = "Composite"

	ProductCoated  string = "Enrobé"
	ProductMoise   string = "Moisé"
	ProductCouple  string = "Couplé"
	ProductReplace string = "Remplacement"
	ProductRemove  string = "Retrait"

	OpacityBlur     float64 = 0.2
	OpacityNormal   float64 = 0.5
	OpacityFiltered float64 = 0.8
	OpacitySelected float64 = 0.9
)
