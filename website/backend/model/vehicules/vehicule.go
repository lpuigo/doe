package vehicules

type ActorHistory struct {
	Date    string
	ActorId int
}

type Event struct {
	StartDate string
	EndDate   string
	Type      string
	Comment   string
}

type Vehicule struct {
	Id             int
	Type           string
	Model          string
	Company        string
	Immat          string
	InCharge       []ActorHistory
	ServiceDate    string
	EndServiceDate string
	Comment        string

	Inventories []Inventory
	Events      []Event
}

func NewVehicule(vType, immat string) *Vehicule {
	return &Vehicule{
		Id:          0,
		Type:        vType,
		Immat:       immat,
		InCharge:    []ActorHistory{},
		Inventories: []Inventory{},
	}
}

// VehiculeByImmat is a getter function to retrieve Vehicule by Immat. returns nil if Vehicule's Immat not found
type VehiculeByImmat func(immat string) *Vehicule

// VehiculeById is a getter function to retrieve Vehicule by id. returns nil if Vehicule's id not found
type VehiculeById func(VehiculeId int) *Vehicule
