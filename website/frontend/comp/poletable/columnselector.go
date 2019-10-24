package poletable

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type Column struct {
	*js.Object
	Name string `js:"Name"`
	Show bool   `js:"Show"`
}

func NewColumn(name string) *Column {
	c := &Column{Object: tools.O()}
	c.Name = name
	c.Show = true
	return c
}

type ColumnSelector struct {
	*js.Object
	Mode    string                     `js:"Mode"`
	Columns []string                   `js:"Columns"`
	Show    map[string]bool            `js:"Show"`
	Refs    []string                   `js:"Refs"`
	Ref     map[string]map[string]bool `js:"Ref"`
}

func NewColumnSelector() *ColumnSelector {
	cs := &ColumnSelector{Object: tools.O()}
	cs.Mode = ""
	cs.Columns = []string{}
	cs.Show = map[string]bool{}
	cs.Refs = []string{}
	cs.Ref = map[string]map[string]bool{}
	return cs
}

func (cs *ColumnSelector) Apply(mode string) {
	cs.Mode = mode
	cs.Show = cs.Ref[mode]

	//config := map[string]bool{}
	//for key, value := range cs.Ref[mode] {
	//	config[key] = value
	//}
	//cs.Show = config

	//for key, value := range cs.Ref[mode] {
	//	cs.SetColumn(key, value)
	//}
}

func (cs *ColumnSelector) SaveAsRef(mode string) {
	cs.Refs = append(cs.Refs, mode)
	config := map[string]bool{}
	for key, value := range cs.Show {
		config[key] = value
	}
	cs.Get("Ref").Set(mode, config)
	//cs.Get("Ref").Set(mode, cs.Show)
}

func (cs *ColumnSelector) AddColumn(name string, show bool) {
	cs.Columns = append(cs.Columns, name)
	cs.SetColumn(name, show)
}

func (cs *ColumnSelector) SetColumn(name string, show bool) {
	cs.Get("Show").Set(name, show)
}

func DefaultColumnSelector() *ColumnSelector {
	cs := NewColumnSelector()

	cs.AddColumn("Appui", true)
	cs.AddColumn("Ville", false)
	cs.AddColumn("Adresse", true)
	cs.AddColumn("DT", true)
	cs.AddColumn("DICT", true)
	cs.AddColumn("Déb.Trx", true)
	cs.AddColumn("Info DICT", true)
	cs.AddColumn("Aspi.", false)
	cs.AddColumn("Type", true)
	cs.AddColumn("Produits", true)
	cs.AddColumn("Acteurs", false)
	cs.AddColumn("Statut", true)
	cs.AddColumn("Ref. Kizeo", false)
	cs.AddColumn("Date", false)
	cs.AddColumn("Attachement", false)
	cs.SaveAsRef("Création")

	cs.SetColumn("Appui", true)
	cs.SetColumn("Ville", false)
	cs.SetColumn("Adresse", true)
	cs.SetColumn("DT", false)
	cs.SetColumn("DICT", false)
	cs.SetColumn("Déb.Trx", false)
	cs.SetColumn("Info DICT", false)
	cs.SetColumn("Aspi.", true)
	cs.SetColumn("Type", true)
	cs.SetColumn("Produits", true)
	cs.SetColumn("Acteurs", true)
	cs.SetColumn("Statut", true)
	cs.SetColumn("Ref. Kizeo", true)
	cs.SetColumn("Date", true)
	cs.SetColumn("Attachement", true)
	cs.SaveAsRef("Suivi")

	cs.Apply("Suivi")

	return cs
}
