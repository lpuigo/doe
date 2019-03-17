package worksiteedit

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/orderedit"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ptedit"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksitestatustag"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/autocomplete"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func Register() {
	hvue.NewComponent("worksite-edit",
		ComponentOptions()...,
	)
}

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("worksite-edit", ComponentOptions()...)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		ptedit.RegisterComponent(),
		orderedit.RegisterComponent(),
		worksitestatustag.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("worksite", "user", "readonly"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewWorksiteDetailModel(vm)
		}),
		hvue.Computed("HasChanged", func(vm *hvue.VM) interface{} {
			wdm := &WorksiteDetailModel{Object: vm.Object}
			if wdm.ReferenceWorksite.Object == nil {
				wdm.ReferenceWorksite = wdm.Worksite.Clone()
				return wdm.Worksite.Dirty
			}
			s1 := wdm.Worksite.SearchInString()
			s2 := wdm.ReferenceWorksite.SearchInString()
			wdm.Worksite.Dirty = s1 != s2
			return wdm.Worksite.Dirty
		}),
		hvue.MethodsOf(&WorksiteDetailModel{}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type WorksiteDetailModel struct {
	*js.Object

	Worksite          *fm.Worksite `js:"worksite"`
	User              *fm.User     `js:"user"`
	ReferenceWorksite *fm.Worksite `js:"refWorksite"`
	ReadOnly          bool         `js:"readonly"`

	VM *hvue.VM `js:"VM"`
}

func NewWorksiteDetailModel(vm *hvue.VM) *WorksiteDetailModel {
	wdm := &WorksiteDetailModel{Object: tools.O()}
	wdm.VM = vm
	wdm.Worksite = nil
	wdm.User = fm.NewUser()
	wdm.ReferenceWorksite = nil
	wdm.ReadOnly = false
	return wdm
}

func WorksiteDetailModelFromJS(o *js.Object) *WorksiteDetailModel {
	return &WorksiteDetailModel{Object: o}
}

func (wdm *WorksiteDetailModel) ClientSearch(vm *hvue.VM, query string, callback *js.Object) {
	wdm = WorksiteDetailModelFromJS(vm.Object)
	res := []*autocomplete.Result{}

	q := strings.ToLower(query)
	for _, client := range wdm.User.Clients {
		if q == "" || strings.Contains(strings.ToLower(client.Name), q) {
			res = append(res, autocomplete.NewResult(client.Name))
		}
	}
	callback.Invoke(res)
}

func (wdm *WorksiteDetailModel) IsDisabled(vm *hvue.VM, info string) bool {
	wdm = &WorksiteDetailModel{Object: vm.Object}
	return wdm.Worksite.IsInfoDisabled(info)
}

func (wdm *WorksiteDetailModel) DeleteOrder(vm *hvue.VM, i int) {
	wdm = &WorksiteDetailModel{Object: vm.Object}
	wdm.Worksite.DeleteOrder(i)
}

func (wdm *WorksiteDetailModel) AddOrder(vm *hvue.VM) {
	wdm = &WorksiteDetailModel{Object: vm.Object}
	wdm.Worksite.AddOrder()
}

func (wdm *WorksiteDetailModel) Save(vm *hvue.VM) {
	wdm = &WorksiteDetailModel{Object: vm.Object}
	vm.Emit("save_worksite", wdm.Worksite)
}

func (wdm *WorksiteDetailModel) Undo(vm *hvue.VM) {
	wdm = &WorksiteDetailModel{Object: vm.Object}
	wdm.Worksite.Copy(wdm.ReferenceWorksite)
}
