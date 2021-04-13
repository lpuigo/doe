package manager

type ManagerConfig struct {
	WorksitesDir      string
	IsWorksitesActive bool

	RipsitesDir      string
	IsRipsitesActive bool

	PolesitesDir      string
	IsPolesitesActive bool

	FoasitesDir      string
	IsFoasitesActive bool

	UsersDir      string
	ActorsDir     string
	ActorInfosDir string
	TimeSheetsDir string
	CalendarFile  string
	ClientsDir    string
	GroupsDir     string
	VehiculesDir  string
	TemplatesDir  string
	SessionKey    string

	SaveArchiveDir string
}
