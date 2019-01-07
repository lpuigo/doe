package progressbar

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"strconv"
)

const template = `
<div>
    <div v-if="showProgressBar" :class="progressStatus">
        <el-progress 
                 :show-text="false"
                 :stroke-width="5"
                 :percentage="progressPct"
        ></el-progress>
        <span class="small-font">{{ progressText }}</span>
    </div>
    <span v-else>{{ progressText }}</span>
</div>

`

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("worksiteinfo-progress-bar", ComponentOptions()...)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("value", "measure"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewWorksiteProgressBarModel(vm)
		}),
		hvue.MethodsOf(&WorksiteProgressBarModel{}),
		hvue.Computed("progressPct", func(vm *hvue.VM) interface{} {
			wspb := &WorksiteProgressBarModel{Object: vm.Object}
			return wspb.ProgressPct()
		}),
		hvue.Computed("showProgressBar", func(vm *hvue.VM) interface{} {
			wspb := &WorksiteProgressBarModel{Object: vm.Object}
			return wspb.mustShow()
		}),
		hvue.Computed("progressStatus", func(vm *hvue.VM) interface{} {
			wspb := &WorksiteProgressBarModel{Object: vm.Object}
			return wspb.ProgressStatus()
		}),
		hvue.Computed("progressText", func(vm *hvue.VM) interface{} {
			wspb := &WorksiteProgressBarModel{Object: vm.Object}
			return wspb.Format()
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type WorksiteProgressBarModel struct {
	*js.Object

	WorksiteInfo *fm.WorksiteInfo `js:"value"`
	Progress     float64          `js:"progress"`
	Measure      bool             `js:"measure"`

	VM *hvue.VM `js:"VM"`
}

func NewWorksiteProgressBarModel(vm *hvue.VM) *WorksiteProgressBarModel {
	wpbm := &WorksiteProgressBarModel{Object: tools.O()}
	wpbm.WorksiteInfo = nil
	wpbm.Progress = 0
	wpbm.Measure = false
	wpbm.VM = vm
	return wpbm
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Formatting Related Methods

func (wspbm *WorksiteProgressBarModel) mustShow() bool {
	return wspbm.WorksiteInfo.NbElTotal > 0 && wspbm.WorksiteInfo.NbElBlocked < wspbm.WorksiteInfo.NbElTotal
}

func (wspbm *WorksiteProgressBarModel) ProgressPct() (pct float64) {
	wsi := wspbm.WorksiteInfo
	effNb := float64(wsi.NbElTotal - wsi.NbElBlocked)
	val := float64(wsi.NbElInstalled)
	if wspbm.Measure {
		val = float64(wsi.NbElMeasured)
	}
	pct = val / effNb * 100.0
	return
}

func (wspbm *WorksiteProgressBarModel) ProgressStatus() (res string) {
	res = "progress-bar"
	return
}

// Filter related Funcs
//

func (wspbm *WorksiteProgressBarModel) Format() (res string) {
	wsi := wspbm.WorksiteInfo
	tot := wsi.NbElTotal - wsi.NbElBlocked
	val := wsi.NbElInstalled
	if wspbm.Measure {
		val = wsi.NbElMeasured
	}
	if wspbm.mustShow() {
		res = strconv.Itoa(val) + " / " + strconv.Itoa(tot)
		return
	}
	res = "-"
	return
}
