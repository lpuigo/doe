package latlong

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

func DecToDeg(val float64) string {
	main, remain := math.Modf(val)
	res := strconv.Itoa(int(main)) + "°"

	min, sec := math.Modf(remain * 60.0)
	res += strconv.Itoa(int(min)) + "'" + strconv.FormatFloat(sec*60.0, 'f', 3, 64)

	return res
}

func DegToDec(val string) (float64, error) {
	vals := strings.FieldsFunc(strings.ReplaceAll(val, ",", "."), func(r rune) bool {
		switch r {
		case '°', '\'', '"':
			return true
		default:
			return false
		}
	})
	if len(vals) < 3 {
		return 0, errors.New("'" + val + "' is not a proper DMS value")
	}
	deg, err := strconv.Atoi(strings.Trim(vals[0], " "))
	if err != nil {
		return 0, err
	}
	min, err := strconv.Atoi(strings.Trim(vals[1], " "))
	if err != nil {
		return 0, err
	}
	sec, err := strconv.ParseFloat(strings.Trim(vals[2], " "), 64)
	if err != nil {
		return 0, err
	}
	return float64(deg) + float64(min)/60 + float64(sec)/3600, nil
}
