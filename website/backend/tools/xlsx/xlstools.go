package xlsx

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
)

func RcToAxis(row, col int) string {
	//res, err := excelize.CoordinatesToCellName(col, row)
	//if err != nil {
	//	res = "A1"
	//}
	//return res

	colname, err := excelize.ColumnNumberToName(col)
	if err != nil {
		colname = "A"
	}
	return colname + strconv.Itoa(row)
}
