package items

import "sort"

type ElemSet map[string]int

func (es ElemSet) SortedKeys(keep func(string) bool) []string {
	res := []string{}
	for key, _ := range es {
		if keep(key) {
			res = append(res, key)
		}
	}
	sort.Strings(res)
	return res
}

var KeepAll func(string) bool = func(string) bool {
	return true
}

func SortedSetKeys(set ElemSet) []string {
	return set.SortedKeys(KeepAll)
}
