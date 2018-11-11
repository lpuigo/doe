package model

type PT struct {
	Ref     string
	RefPt   string
	Address string
}

func MakePT(ref, refpt, address string) PT {
	return PT{Ref: ref, RefPt: refpt, Address: address}
}
