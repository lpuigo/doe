package main

import (
	"fmt"
	"github.com/lpuig/ewin/doe/model"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
)

func main() {
	walk.FocusEffect, _ = walk.NewBorderGlowEffect(walk.RGB(0, 63, 255))
	tr := model.NewTroncon()
	tr.Ref = "TR123456"
	tr.NbFiber = 12

	var mainwindow *walk.MainWindow
	var db *walk.DataBinder
	var teTr *walk.TextEdit
	var comp *walk.Composite

	refreshTr := func() {
		teTr.SetText(fmt.Sprintf("%+v", *tr))
	}

	submit := func() {
		err := db.Submit()
		if err != nil {
			log.Println("Can not Submit change", err)
		}
		refreshTr()
	}

	gbTroncon := GroupBox{
		Layout: Grid{Columns: 2},
		Title:  "Tronçon",
		Children: []Widget{
			Label{Text: "Tronçon"},
			LineEdit{
				Text:          Bind("Ref"),
				OnTextChanged: submit,
			},
			Label{Text: "Nb Fibre"},
			ComboBox{
				Value:                 Bind("NbFiber", SelRequired{}),
				BindingMember:         "Nb",
				DisplayMember:         "Number",
				Model:                 ListNbFibers(),
				OnCurrentIndexChanged: submit,
			},
		},
	}

	attach := func(w Widget, c walk.Container) {
		b := NewBuilder(c)
		err := w.Create(b)
		if err != nil {
			log.Println("Error while create:", err)
		}
	}

	mw := MainWindow{
		Title:    "EWIN Suivi des dossiers",
		AssignTo: &mainwindow,
		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "troncon",
			DataSource:     tr,
			ErrorPresenter: ToolTipErrorPresenter{},
			//AutoSubmitDelay: 200 * time.Millisecond,
			//AutoSubmit:      false,
		},
		MinSize: Size{640, 480},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				AssignTo: &comp,
				Layout:   VBox{},
				Children: []Widget{
					gbTroncon,
					PushButton{
						Text:      "ajouter un Tronçon",
						OnClicked: func() { attach(gbTroncon, comp) },
					},
				},
			},
			Label{
				Text: "Tronçon:",
			},
			TextEdit{
				AssignTo: &teTr,
				ReadOnly: true,
				MinSize:  Size{10, 100},
				Text:     fmt.Sprintf("%+v", tr),
			},
		},
	}

	_, err := mw.Run()
	if err != nil {
		log.Fatal("MainWindow returns", err)
	}
}

type ChoiceNbFibre struct {
	Nb     int
	Number string
}

func ListNbFibers() []*ChoiceNbFibre {
	return []*ChoiceNbFibre{
		{6, "Six"},
		{12, "Douze"},
		{24, "Vingt-Quatre"},
		{36, "Trente-Six"},
		{48, "Quarante-Huit"},
	}
}
