package actorconst

const (
	StateCandidate string = "00 Candidate"
	StateActive    string = "10 Active"
	StateOnHoliday string = "20 On Holiday"
	StateGone      string = "90 Gone"
	StateDefection string = "91 Defection"

	StateLabelCandidate string = "Candidat"
	StateLabelActive    string = "Actif"
	StateLabelOnHoliday string = "En Congés"
	StateLabelGone      string = "Parti"
	StateLabelDefection string = "Défection"
)

const (
	RolePuller     string = "Tireur"
	RoleJuncter    string = "Raccordeur"
	RoleDriver     string = "Chauffeur"
	RoleTeamleader string = "Chef d'Equipe"

	ContractTemp string = "Intérim"
	ContractCDD  string = "CDD"
	ContractCDI  string = "CDI"
)

const (
	FilterValueAll     string = ""
	FilterValueCompany string = "CMPY:"
	FilterValueName    string = "NAM:"
	FilterValueClient  string = "CLT:"
	FilterValueComment string = "CMT:"

	FilterLabelAll     string = "Tout"
	FilterLabelCompany string = "Compagnie"
	FilterLabelName    string = "Nom"
	FilterLabelClient  string = "Client"
	FilterLabelComment string = "Commentaire"
)

const (
	InfoTrainingBirth         string = "Date Naissance"
	InfoTrainingMedic         string = "Visite Med."
	InfoTrainingProBTP        string = "ProBTP"
	InfoTrainingAIPR          string = "AIPR"
	InfoTrainingElec          string = "Prox. Elec."
	InfoTrainingCacesChantier string = "Sign. Chantier"
	InfoTrainingCacesGrue     string = "Caces Grue"
	InfoTrainingCacesHauteur  string = "Caces Hauteur"
	InfoTrainingCacesNacelle  string = "Caces Nacelle"
)
