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
	Date     string
	Actors   []string
	DtRef    string
	DictRef  string
	DictInfo string
	Height   int
	Material string
	Product  []string
}

func (p *Pole) SearchString() string {
	var searchBuilder strings.Builder
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "Ref", strings.ToUpper(p.Ref))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "City", strings.ToUpper(p.City))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "Address", strings.ToUpper(p.Address))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "DtRef", strings.ToUpper(p.DtRef))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "DictRef", strings.ToUpper(p.DictRef))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "Height", strconv.Itoa(p.Height)+"M")
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "Material", strings.ToUpper(p.Material))
	for _, key := range p.Product {
		fmt.Fprintf(&searchBuilder, "poleProduct:%s,", strings.ToUpper(key))
	}
	for _, actor := range p.Actors {
		fmt.Fprintf(&searchBuilder, "poleActor:%s,", strings.ToUpper(actor))
	}
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "DictInfo", strings.ToUpper(p.DictInfo))
	return searchBuilder.String()
}
