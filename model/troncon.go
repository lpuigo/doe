package model

type Troncon struct {
	Ref           string
	Pb            PT
	NbRacco       int
	NbFiber       int
	Article       string
	Blockage      bool
	NeedSignature bool
	Signed        bool
	InstallDate   string
	InstallActor  string
	MeasureDate   string
	MeasureActor  string
	Comment       string
}

func MakeTroncon(ref string, pb PT, nbRacco, nbFiber int, needsign bool) Troncon {
	return Troncon{Ref: ref, Pb: pb, NbRacco: nbRacco, NbFiber: nbFiber, NeedSignature: needsign}
}
