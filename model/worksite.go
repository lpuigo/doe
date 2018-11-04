package model

type Commande struct {
	Ref      string
	Pa       PA
	Troncons []*Troncon
}

func NewCommande() Commande {
	c := Commande{}
	return c
}

func (c *Commande) AddTroncons(t ...*Troncon) {
	c.Troncons = append(c.Troncons, t...)
}

type PA struct {
	Ref   string
	Ville string
}

func NewPA() PA {
	pa := PA{}
	return pa
}

type Troncon struct {
	Ref     string
	NbFiber int
	Pts     []*PT
}

func NewTroncon() *Troncon {
	tr := &Troncon{}
	tr.Pts = []*PT{}
	return tr
}

func (tr *Troncon) AddPT(pt ...*PT) {
	tr.Pts = append(tr.Pts, pt...)
}

type TypePB string

const (
	PB_Facade       TypePB = "Façade"
	PB_PoteauFT     TypePB = "Poteau FT"
	PB_PoteauEnedis TypePB = "Poteau Enedis"
	PB_Chambre      TypePB = "Chambre"
)

type PT struct {
	Ref       string
	NbFiber   int
	NbELRacco int
	RefPB     string
	Type      TypePB
	Address   string
}

func NewPT() *PT {
	pt := &PT{}
	return pt
}
