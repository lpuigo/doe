package tronconstatustag

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ptedit"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

const (
	template string = `
<el-tag :type="Status" size="medium" style="width: 100%">{{StatusText}}</el-tag>
`
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func Register() {
	hvue.NewComponent("troncon-status-tag",
		ComponentOptions()...,
	)
}

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("troncon-status-tag", ComponentOptions()...)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		ptedit.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("value"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewTronconEditModel(vm)
		}),
		hvue.Computed("Status", func(vm *hvue.VM) interface{} {
			tst := &TronconStatusTagModel{Object: vm.Object}
			statusType, statusText := tst.SetStatus()
			tst.StatusText = statusText
			return statusType
		}),
		hvue.MethodsOf(&TronconStatusTagModel{}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type TronconStatusTagModel struct {
	*js.Object

	Troncon    *fm.Troncon `js:"value"`
	StatusText string      `js:"StatusText"`

	VM *hvue.VM `js:"VM"`
}

func NewTronconEditModel(vm *hvue.VM) *TronconStatusTagModel {
	tem := &TronconStatusTagModel{Object: tools.O()}
	tem.VM = vm
	tem.Troncon = nil
	tem.StatusText = ""
	return tem
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Actions

func (tst *TronconStatusTagModel) SetStatus() (statusType, statusText string) {
	tr := tst.Troncon

	switch {
	case tr.Ref == "" || !tr.Pb.IsFilledIn():
		return "warning", "A renseigner"
	case tr.Blockage && !tr.NeedSignature && tr.Comment == "":
		return "warning", "Saisir desc. bloquage"
	case tr.Blockage:
		return "info", "Bloqué"
	case !tools.Empty(tr.MeasureDate):
		return "success", "Terminé"
	case !tools.Empty(tr.InstallDate):
		return "", "Mesures à faire"
	case tr.Ref != "" && tr.Pb.IsFilledIn() && !tr.Blockage:
		return "", "A Réaliser"
	}

	return "danger", "Erreur"
}
