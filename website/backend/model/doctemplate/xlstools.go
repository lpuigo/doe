package doctemplate

import (
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func RcToAxis(row, col int) string {
	return excelize.ToAlphaString(col) + strconv.Itoa(row+1)
}
