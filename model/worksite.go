package model

type Worksite struct {
	Ref       string
	OrderDate string
	Pmz       PT
	Pa        PT
	Orders    []Order
}

func MakeWorksite(ref, orderdate string, pmz, pa PT, order ...Order) Worksite {
	return Worksite{Ref: ref, Pmz: pmz, Pa: pa, Orders: order, OrderDate: orderdate}
}
