package polesite

func GetCenterAndBounds(poles []*Pole) (clat, clong, blat1, blong1, blat2, blong2 float64) {

	min := func(pole *Pole) {
		if pole.Lat < blat1 {
			blat1 = pole.Lat
		}
		if pole.Long < blong1 {
			blong1 = pole.Long
		}
	}

	max := func(pole *Pole) {
		if pole.Lat > blat2 {
			blat2 = pole.Lat
		}
		if pole.Long > blong2 {
			blong2 = pole.Long
		}
	}

	blat1, blong1 = 500, 500
	blat2, blong2 = -500, -500
	var nbPole int = 0
	for _, pole := range poles {
		if pole.Deleted() {
			continue
		}
		nbPole++
		clat += pole.Lat
		clong += pole.Long
		min(pole)
		max(pole)
	}

	if nbPole == 0 {
		return 47, 5, 46.5, 4.5, 47.5, 5.5
	}

	nb := float64(nbPole)
	clat /= nb
	clong /= nb
	return
}
