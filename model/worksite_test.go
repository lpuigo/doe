package model

import (
	"encoding/json"
	"os"
	"testing"
)

func genWorksite(t *testing.T) Worksite {
	return MakeWorksite(
		"PA-59163-003T",
		"2018-11-06",
		MakePT("PMZ-38467", "PT-007605", "02, Rue Kléber, CROIX"),
		MakePT("PA-50071", "PT-008020", "02, Rue Jean Jaurès, CROIX"),
		MakeOrder("F53349061118",
			MakeTroncon("TR-18-0502",
				MakePT("PB-63742", "PT-008881", "07, Rue Pasteur"),
				6, 6, false),
			MakeTroncon("TR-18-0505",
				MakePT("PB-63744", "PT-008882", "19, Rue Pasteur"),
				6, 6, false),
			MakeTroncon("TR-18-0506",
				MakePT("PB-63746", "PT-008883", "31, Rue Pasteur"),
				6, 6, false),
			MakeTroncon("TR-18-0507",
				MakePT("PB-63749", "PT-008884", "44, Rue de Bapaume"),
				5, 6, false),
			MakeTroncon("TR-18-0508",
				MakePT("PB-63751", "PT-008885", "36, Rue de Bapaume"),
				6, 6, false),
		),
		MakeOrder("F53361061118",
			MakeTroncon("TR-18-0514",
				MakePT("PB-63753", "PT-008895", "08, Rue Pasteur"),
				6, 6, false),
			MakeTroncon("TR-18-0515",
				MakePT("PB-63754", "PT-008896", "20, Rue Pasteur"),
				6, 6, false),
			MakeTroncon("TR-18-0517",
				MakePT("PB-63756", "PT-008897", "30, Rue Pasteur"),
				6, 6, false),
			MakeTroncon("TR-18-0519",
				MakePT("PB-63758", "PT-008898", "40, Rue Pasteur"),
				3, 6, true),
		),
		MakeOrder("F53370061118",
			MakeTroncon("TR-18-0520",
				MakePT("PB-63760", "PT-008899", "332, Rue Jean Jaurès"),
				3, 6, true),
			MakeTroncon("TR-18-0521",
				MakePT("PB-63762", "PT-008900", "326, Rue Jean Jaurès"),
				6, 6, false),
		),
		MakeOrder("F53393061118",
			MakeTroncon("TR-18-0522",
				MakePT("PB-63763", "PT-008901", "06, Rue du Professeur Langevin"),
				6, 6, false),
		),
	)
}

func TestMakeWorksite(t *testing.T) {
	ws := genWorksite(t)

	t.Logf("File : %s\n", ws.FileName())
	jm := json.NewEncoder(os.Stdout)
	jm.SetIndent("", "\t")
	jm.Encode(ws)

}
