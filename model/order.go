package model

type Order struct {
	Ref      string
	Comment  string
	Troncons []Troncon
}

func MakeOrder(ref string, troncon ...Troncon) Order {
	return Order{Ref: ref, Troncons: troncon}
}

func (o *Order) AddTroncon(troncon ...Troncon) {
	o.Troncons = append(o.Troncons, troncon...)
}
