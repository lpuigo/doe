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

	FilterValueAll      string = ""
	FilterValueRef      string = "REF:"
	FilterValueCity     string = "CTY:"
	FilterValueComment  string = "CMT:"
	FilterValueHeigth   string = "HGT:"
	FilterValueProduct  string = "PRD:"
	FilterValueDt       string = "DT:"
	FilterValueDict     string = "DCT:"
	FilterValueDictInfo string = "DCI:"

	FilterLabelAll      string = "Tout"
	FilterLabelRef      string = "Référence"
	FilterLabelCity     string = "Ville"
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
	ProductReplace string = "Remplacement"
	ProductRemove  string = "Retrait"

	OpacityBlur     float64 = 0.2
	OpacityNormal   float64 = 0.5
	OpacityFiltered float64 = 0.8
	OpacitySelected float64 = 0.9
)
