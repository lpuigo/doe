package bpu

import (
	"fmt"
	xlstool "github.com/lpuig/ewin/doe/website/backend/tools/xlsx"
	"github.com/tealeg/xlsx"
	"sort"
	"strings"
)

type Bpu struct {
	Activities map[string]CategoryArticles // map[Activity]map[Category][]*Article
	Boxes      map[string]map[string]*Box  // map[Category]map[Name]*Box
}

func NewBpu() *Bpu {
	return &Bpu{
		Activities: make(map[string]CategoryArticles),
		Boxes:      make(map[string]map[string]*Box),
	}
}

func (bpu *Bpu) GetCategoryArticles(activity string) CategoryArticles {
	return bpu.Activities[strings.ToUpper(activity)]
}

func (bpu *Bpu) IsBoxDefined(cat, boxName string) bool {
	catBox, found := bpu.Boxes[strings.ToUpper(cat)]
	if !found {
		return false
	}
	_, found = catBox[strings.ToUpper(boxName)]
	return found
}

func (bpu *Bpu) GetBox(cat, boxName string) *Box {
	catBox, found := bpu.Boxes[strings.ToUpper(cat)]
	if !found {
		return nil
	}
	return catBox[strings.ToUpper(boxName)]
}

func (bpu *Bpu) GetArticleNames(activity string) []string {
	res := []string{}
	for _, catChapters := range bpu.GetCategoryArticles(activity) {
		for _, chapter := range catChapters {
			res = append(res, chapter.Name)
		}
	}
	sort.Strings(res)
	return res
}

const (
	bpuPriceSheetName = "Prices"
	bpuBoxeSheetName  = "Boxes"
)

const (
	colPricesActivity int = iota
	colPricesCategory
	colPricesName
	colPricesSize
	colPricesPrice
	colPricesWork
	colPricesEnd
)

func NewBpuFromXLS(file string) (bpu *Bpu, err error) {
	xf, err := xlsx.OpenFile(file)
	if err != nil {
		return
	}
	priceSheet := xf.Sheet[bpuPriceSheetName]
	if priceSheet == nil {
		err = fmt.Errorf("could not find '%s' sheet in '%s'", bpuPriceSheetName, file)
		return
	}
	boxSheet := xf.Sheet[bpuBoxeSheetName]
	if boxSheet == nil {
		err = fmt.Errorf("could not find '%s' sheet in '%s'", bpuBoxeSheetName, file)
		return
	}
	bpu = NewBpu()
	err = bpu.parseActivities(priceSheet)
	if err != nil {
		return
	}
	err = bpu.parseBoxes(boxSheet)
	return
}

func (bpu *Bpu) parseActivities(sheet *xlsx.Sheet) (err error) {
	entryFound := true
	for row := 1; entryFound; row++ {
		activity := sheet.Cell(row, colPricesActivity).Value
		if activity == "" { // check for data ending (first column is empty => we are done)
			entryFound = false
			continue
		}
		cat := sheet.Cell(row, colPricesCategory).Value
		nChapter, err := parseActivityRow(sheet, row)
		if err != nil {
			return err
		}
		cat = strings.ToUpper(cat)
		activity = strings.ToUpper(activity)
		currentActivityCatChapters, found := bpu.Activities[activity]
		if !found {
			currentActivityCatChapters = NewCategoryChapters()
			bpu.Activities[activity] = currentActivityCatChapters
		}
		currentActivityCatChapters[cat] = append(currentActivityCatChapters[cat], nChapter)
	}

	// sort price categories by ascending size
	//for _, actCatChapt := range bpu.Chapters {
	//	actCatChapt.SortChapters()
	//}
	return
}

func parseActivityRow(sh *xlsx.Sheet, row int) (p *Article, err error) {
	cSize := sh.Cell(row, colPricesSize)
	size, e := cSize.Int()
	if e != nil {
		err = fmt.Errorf("could not get size info '%s' in sheet '%s!%s'", cSize.Value, bpuPriceSheetName, xlstool.RcToAxis(row, colPricesSize))
		return
	}
	cPrice := sh.Cell(row, colPricesPrice)
	price, e := cPrice.Float()
	if e != nil {
		err = fmt.Errorf("could not get price info '%s' in sheet '%s!%s'", cSize.Value, bpuPriceSheetName, xlstool.RcToAxis(row, colPricesPrice))
		return
	}
	cWork := sh.Cell(row, colPricesWork)
	work, e := cWork.Float()
	if e != nil {
		err = fmt.Errorf("could not get work info '%s' in sheet '%s!%s'", cSize.Value, bpuPriceSheetName, xlstool.RcToAxis(row, colPricesWork))
		return
	}
	p = NewArticle()
	p.Name = sh.Cell(row, colPricesName).Value
	p.Unit = size
	p.Price = price
	p.Work = work
	return
}

const (
	colBoxesCategory int = iota
	colBoxesName
	colBoxesSize
	colBoxesUsage
	colBoxesEnd
)

func (bpu *Bpu) parseBoxes(sheet *xlsx.Sheet) (err error) {
	entryFound := true
	for row := 1; entryFound; row++ {
		cat := strings.ToUpper(sheet.Cell(row, colBoxesCategory).Value)
		boxName := strings.ToUpper(sheet.Cell(row, colBoxesName).Value)
		if boxName == "" { // check for data ending (first column is empty => we are done)
			entryFound = false
			continue
		}
		sizeCell := sheet.Cell(row, colBoxesSize)
		size, e := sizeCell.Int()
		if e != nil {
			err = fmt.Errorf("could not get size info '%s' in sheet '%s!%s'", sizeCell.Value, bpuBoxeSheetName, xlstool.RcToAxis(row, colBoxesSize))
			return
		}
		nBox := NewBox()
		nBox.Name = boxName
		nBox.Size = size
		nBox.Usage = sheet.Cell(row, colBoxesUsage).Value
		catBox, found := bpu.Boxes[cat]
		if !found {
			catBox = make(map[string]*Box)
			bpu.Boxes[cat] = catBox
		}
		catBox[nBox.Name] = nBox
	}
	return
}

//func (bpu *Bpu) String() string {
//	res := ""
//	for _, p := range bpu.BpePrices {
//		res += fmt.Sprintf("%3d : %v\n", p.Size, p.Article)
//	}
//	return res
//}
