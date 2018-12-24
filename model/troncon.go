package model

type Troncon struct {
	Ref           string
	Pb            PT
	NbRacco       int
	NbFiber       int
	NeedSignature bool
	Signed        bool
	InstallDate   string
	MeasureDate   string
	Comment       string
}

func MakeTroncon(ref string, pb PT, nbRacco, nbFiber int, needsign bool) Troncon {
	return Troncon{Ref: ref, Pb: pb, NbRacco: nbRacco, NbFiber: nbFiber, NeedSignature: needsign}
}
