package vehiculeconst

const (
	TypeTariere     string = "Tarière"
	TypeNacelle     string = "Nacelle"
	TypeFourgon     string = "Fourgon"
	TypeCar         string = "Voiture"
	TypePorteTouret string = "Porte-Touret"

	InChargeNotAffected string = "non affecté"
	InventoryNotFound   string = "pas d'inventaire"
)

const (
	StatusInUse    string = "Disponible"
	StatusInRepair string = "En réparation"
	StatusReturned string = "Rendu"
)

const (
	EventTypeIncident string = "Accident"
	EventTypeRepair   string = "Reparation"
	EventTypeCheck    string = "Contrôle"
	EventTypeMisc     string = "Divers"
)

const (
	FilterValueAll     string = ""
	FilterValueCompany string = "CMPY:"
	FilterValueImmat   string = "IMAT:"
	FilterValueType    string = "TYP:"
	FilterValueComment string = "CMT:"

	FilterLabelAll     string = "Tout"
	FilterLabelCompany string = "Compagnie"
	FilterLabelImmat   string = "Immat"
	FilterLabelType    string = "Type"
	FilterLabelComment string = "Commentaire"
)
