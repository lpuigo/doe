package model

type Troncon struct {
	Ref           string
	Pb            PT
	NbRacco       int
	NbFiber       int
	NeedSignature bool
	InstallDate   string
	MeasureDate   string
}

func MakeTroncon(ref string, pb PT, nbRacco, nbFiber int, needsign bool) Troncon {
	return Troncon{Ref: ref, Pb: pb, NbRacco: nbRacco, NbFiber: nbFiber, NeedSignature: needsign}
}
