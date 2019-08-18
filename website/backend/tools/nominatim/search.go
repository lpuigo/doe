package nominatim

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	nominatimUrl        string = "https://nominatim.openstreetmap.org/search?q="
	nominatimFormatJson string = "&format=json"
)

type Geoloc struct {
	Lat         string  `json:"lat"`
	Long        string  `json:"lon"`
	FullAddress string  `json:"display_name"`
	Importance  float64 `json:"importance"`
}

func (g Geoloc) GetLatLong() (float64, float64, error) {
	lat, err := strconv.ParseFloat(g.Lat, 64)
	if err != nil {
		return 0, 0, err
	}
	long, err := strconv.ParseFloat(g.Long, 64)
	if err != nil {
		return 0, 0, err
	}
	return lat, long, nil
}

func GeolocSearch(address string) ([]Geoloc, error) {
	adr := strings.Replace(strings.Trim(address, " "), " ", "+", -1)
	uri := nominatimUrl + adr + nominatimFormatJson

	response, err := http.Get(uri)
	time.Sleep(time.Second)
	if err != nil {
		return []Geoloc{}, fmt.Errorf("get failed and returns: %s", err.Error())
	}

	defer response.Body.Close()

	res := []Geoloc{}
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		return []Geoloc{}, fmt.Errorf("could not unmarshall response: %s", err.Error())
	}
	return res, nil
}
