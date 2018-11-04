package model

import (
	"encoding/json"
	"os"
	"testing"
)

func genCommande(t *testing.T) Commande {
	c := NewCommande()
	c.Ref = "F12345JJMMAA"

	pa := NewPA()
	pa.Ref = "PA 12345ABCD"
	pa.Ville = "Roubaix"
	c.Pa = pa

	tr := NewTroncon()
	tr.Ref = "TR 12 3456"
	tr.NbFiber = 6

	pt := NewPT()
	pt.Ref = "PT 123456"
	pt.NbFiber = 6
	pt.NbELRacco = 5
	pt.RefPB = "PB 123456"
	pt.Type = PB_Facade
	pt.Address = "10, rue de la Pompe"

	tr.AddPT(pt)
	c.AddTroncons(tr)

	return c
}

func TestNewCommande(t *testing.T) {
	c := genCommande(t)
	t.Logf("Commande : %v", c)
}

func TestNewCommandeJson(t *testing.T) {
	c := genCommande(t)
	je := json.NewEncoder(os.Stdout)
	je.SetIndent("", "\t")
	err := je.Encode(c)
	if err != nil {
		t.Error("Json encode returns", err)
	}
}
