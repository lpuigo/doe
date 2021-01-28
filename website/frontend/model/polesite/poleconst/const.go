package poleconst

const (
	DictValidityDuration int = 90
)

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
	StateNotSubmitted      string = "00 Not Submitted"
	StateNoGo              string = "05 NoGo"
	StateDictToDo          string = "08 DICT To Do"
	StateDaToDo            string = "081 DA To Do"
	StateDaExpected        string = "083 DA Expected"
	StatePermissionPending string = "09 Permission Pending"
	StateToDo              string = "10 To Do"
	StateMarked            string = "11 Marked"
	StateNoAccess          string = "12 No Access"
	StateDenseNetwork      string = "14 Dense Network"
	StateHoleDone          string = "20 Hole Done"
	StateIncident          string = "25 Incident"
	StateDone              string = "90 Done"
	StateAttachment        string = "95 Attachment"
	StateCancelled         string = "99 Cancelled"
	StateDeleted           string = "999 Deleted"

	LabelNotSubmitted      string = "Non soumis"
	LabelNoGo              string = "NoGo Client"
	LabelDictToDo          string = "DICT à faire"
	LabelDaToDo            string = "Demande AC à faire"
	LabelDaExpected        string = "Attente retour AC"
	LabelPermissionPending string = "Attente Permission"
	LabelToDo              string = "A faire"
	LabelMarked            string = "Marqué"
	LabelNoAccess          string = "Inaccessible"
	LabelDenseNetwork      string = "Réseaux denses"
	LabelHoleDone          string = "Trou fait"
	LabelIncident          string = "Incident"
	LabelDone              string = "Fait"
	LabelAttachment        string = "Attach. fait"
	LabelCancelled         string = "Annulé"
	LabelDeleted           string = "Supprimé"

	FilterValueAll      string = ""
	FilterValueRef      string = "REF:"
	FilterValueCity     string = "CTY:"
	FilterValueAddr     string = "ADD:"
	FilterValueComment  string = "CMT:"
	FilterValueMaterial string = "MAT:"
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
	FilterLabelMaterial string = "Matériau"
	FilterLabelHeigth   string = "Hauteur"
	FilterLabelProduct  string = "Produit"
	FilterLabelDt       string = "DT"
	FilterLabelDict     string = "DICT"
	FilterLabelDictInfo string = "DICT Info"

	MaterialWood          string = "Bois"
	MaterialMetal         string = "Métal"
	MaterialEnforcedMetal string = "Métal Renforcé"
	MaterialComp          string = "Composite"
	MaterialEnforcedComp  string = "Composite Renforcé"

	ProductCreation         string = "Création"
	ProductCoated           string = "Enrobé"
	ProductInRow            string = "Enfilade"
	ProductHandDigging      string = "Implantation manuelle"
	ProductMechDigging      string = "Terrassement mécanique"
	ProductPruning          string = "Elagage"
	ProductMoise            string = "Moisé"
	ProductCouple           string = "Couplé"
	ProductHauban           string = "Haubané"
	ProductReplace          string = "Remplacement"
	ProductTrickyReplace    string = "Remplacement complexe"
	ProductRemove           string = "Retrait"
	ProductStraighten       string = "Redressement"
	ProductNoAccess         string = "Inaccessible"
	ProductDenseNetwork     string = "Réseaux denses"
	ProductReplenishment    string = "Réappro."
	ProductFarReplenishment string = "Réappro. (>50km)"
)

const (
	OpacityBlur     float64 = 0.25
	OpacityNormal   float64 = 0.6
	OpacityFiltered float64 = 1.0
	OpacitySelected float64 = 0.8

	ZoomLevelOnPole int = 20
)
