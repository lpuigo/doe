package worksitesummary

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

const template = `
<el-card shadow="hover" class="wssum">
    <div slot="header">
        <H2 class="spread text">
            <span>{{worksite.Ref}}</span>
            <span class="smaller">{{worksite.OrderDate}}</span>
        </H2>
        <!--<el-button style="float: right; padding: 3px 0" type="text">Operation button</el-button>-->
    </div>
    <div class="spread">
        <div class="pt spread text">
            <span><strong>PA :</strong></span>
            <span>{{worksite.Pa.Ref}}</span>
            <span>({{worksite.Pa.RefPt}})</span>
            <span class="smaller">{{worksite.Pa.Address}}</span>
        </div>
        <div class="pt spread text">
            <span><strong>PMZ :</strong></span>
            <span>{{worksite.Pmz.Ref}}</span>
            <span>({{worksite.Pmz.RefPt}})</span>
            <span class="smaller">{{worksite.Pmz.Address}}</span>
        </div>
    </div>
    <div class="pt spread text">
        <span><strong>Commentaire :</strong></span>
        <span>{{worksite.Comment}}</span>
    </div>
    <!-- Orders Here-->
</el-card>
`

func Register() {
	hvue.NewComponent("worksite-summary",
		ComponentOptions()...,
	)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		//hvue.Component("project-progress-bar", wl_progress_bar.ComponentOptions()...),
		hvue.Props("worksite"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewWorksiteSummaryModel(vm)
		}),
		hvue.MethodsOf(&WorksiteSummaryModel{}),
		//hvue.Computed("filteredProjects", func(vm *hvue.VM) interface{} {
		//	ptm := &ProjectTableModel{Object: vm.Object}
		//	if ptm.Filter == "" {
		//		return ptm.Projects
		//	}
		//	res := []*fm.Project{}
		//	for _, p := range ptm.Projects {
		//		if fm.TextFiltered(p, ptm.Filter) {
		//			res = append(res, p)
		//		}
		//	}
		//	return res
		//}),
		//hvue.Filter("DateFormat", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
		//	return fm.DateString(value.String())
		//}),
	}
}

type WorksiteSummaryModel struct {
	*js.Object
	VM *hvue.VM `js:"VM"`

	Worksite *fm.Worksite `js:"worksite"`
}

func NewWorksiteSummaryModel(vm *hvue.VM) *WorksiteSummaryModel {
	wssm := &WorksiteSummaryModel{Object: tools.O()}
	wssm.VM = vm

	wssm.Worksite = nil

	return wssm
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Event Related Methods
