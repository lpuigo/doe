package actorconst

const (
	StateCandidate string = "00 Candidate"
	StateActive    string = "10 Active"
	StateOnHoliday string = "20 On Holiday"
	StateGone      string = "90 Fait"

	StateLabelCandidate string = "Candidat"
	StateLabelActive    string = "Employé"
	StateLabelOnHoliday string = "En Congés"
	StateLabelGone      string = "Parti"
)

const (
	RolePuller     string = "Tireur"
	RoleJuncter    string = "Racordeur"
	RoleDriver     string = "Chauffeur"
	RoleTeamleader string = "Chef d'Equipe"
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
