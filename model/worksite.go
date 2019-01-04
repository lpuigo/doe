package model

type Worksite struct {
	Id             int
	Client         string
	Ref            string
	OrderDate      string
	DoeDate        string
	AttachmentDate string
	PaymentDate    string
	City           string
	Status         string
	Pmz            PT
	Pa             PT
	Comment        string
	Orders         []Order
	Rework         *Rework
}

func MakeWorksite(ref, orderdate string, pmz, pa PT, order ...Order) Worksite {
	return Worksite{Ref: ref, Pmz: pmz, Pa: pa, Orders: order, OrderDate: orderdate}
}

func (w Worksite) FileName() string {
	return w.OrderDate + "_" + w.Ref
}
