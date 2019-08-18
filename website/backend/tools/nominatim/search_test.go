package nominatim

import (
	"testing"
	"time"
)

func TestGeolocSearch(t *testing.T) {
	addresses := []string{
		"62 rue de Verdun, Forbach",
		"142 Rue du Faubourg Saint Denis, 75010 Paris",
		"34 rue Principale, ARRIANCE",
		"8 rue de l'étang, HERNY",
		"4 chemin de Bonne House, MAINVILLERS",
	}

	for _, addr := range addresses {
		t.Logf("\nRequesting '%s'", addr)
		res, err := GeolocSearch(addr)
		if err != nil {
			t.Errorf("GeolocSearch('%s') returns unexpected: %s", addr, err.Error())
		}
		t.Log("result", res)
		time.Sleep(time.Second)
	}
}
