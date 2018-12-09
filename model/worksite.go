package model

type Worksite struct {
	Id        int
	Ref       string
	OrderDate string
	City      string
	Status    string
	Pmz       PT
	Pa        PT
	Comment   string
	Orders    []Order
}

func MakeWorksite(ref, orderdate string, pmz, pa PT, order ...Order) Worksite {
	return Worksite{Ref: ref, Pmz: pmz, Pa: pa, Orders: order, OrderDate: orderdate}
}

func (w Worksite) FileName() string {
	return w.OrderDate + "_" + w.Ref
}
