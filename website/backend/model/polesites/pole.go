package polesites

import (
	"fmt"
	"strconv"
	"strings"
)

type Pole struct {
	Ref      string
	City     string
	Address  string
	Lat      float64
	Long     float64
	State    string
	DtRef    string
	DictRef  string
	Height   int
	Product  map[string]int
	DictInfo map[string]string
}

func (p *Pole) SearchString() string {
	var searchBuilder strings.Builder
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "Ref", strings.ToUpper(p.Ref))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "City", strings.ToUpper(p.City))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "Address", strings.ToUpper(p.Address))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "DtRef", strings.ToUpper(p.DtRef))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "DictRef", strings.ToUpper(p.DictRef))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "Height", strconv.Itoa(p.Height)+"M")
	for key, _ := range p.Product {
		fmt.Fprintf(&searchBuilder, "poleProduct:%s,", strings.ToUpper(key))
	}
	for key, value := range p.DictInfo {
		fmt.Fprintf(&searchBuilder, "poleDict%s:%s,", strings.ToUpper(key), strings.ToUpper(value))
	}
	return searchBuilder.String()
}
