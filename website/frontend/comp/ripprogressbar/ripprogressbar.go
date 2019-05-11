package ripprogressbar

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"strconv"
)

const template = `
<div>
    <div v-if="showProgressBar" >
        <el-progress 
                 :show-text="false"
                 :stroke-width="5"
                 :percentage="progressPct"
				 :status="progressStatus"
        ></el-progress>
        <span class="small-font">{{ progressText }}</span>
    </div>
    <span v-else>{{ progressText }}</span>
</div>

`

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("ripsiteinfo-progress-bar", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("total", "done"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewWorksiteProgressBarModel(vm)
		}),
		hvue.MethodsOf(&RipsiteProgressBarModel{}),
		hvue.Computed("progressPct", func(vm *hvue.VM) interface{} {
			wspb := &RipsiteProgressBarModel{Object: vm.Object}
			return wspb.ProgressPct()
		}),
		hvue.Computed("showProgressBar", func(vm *hvue.VM) interface{} {
			wspb := &RipsiteProgressBarModel{Object: vm.Object}
			return wspb.mustShow()
		}),
		hvue.Computed("progressStatus", func(vm *hvue.VM) interface{} {
			wspb := &RipsiteProgressBarModel{Object: vm.Object}
			return wspb.ProgressStatus()
		}),
		hvue.Computed("progressText", func(vm *hvue.VM) interface{} {
			wspb := &RipsiteProgressBarModel{Object: vm.Object}
			return wspb.Format()
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type RipsiteProgressBarModel struct {
	*js.Object

	Total    int     `js:"total"`
	Done     int     `js:"done"`
	Progress float64 `js:"progress"`

	VM *hvue.VM `js:"VM"`
}

func NewWorksiteProgressBarModel(vm *hvue.VM) *RipsiteProgressBarModel {
	rpbm := &RipsiteProgressBarModel{Object: tools.O()}
	rpbm.Total = 0
	rpbm.Done = 0
	rpbm.Progress = 0
	rpbm.VM = vm
	return rpbm
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Formatting Related Methods

func (rpbm *RipsiteProgressBarModel) mustShow() bool {
	return rpbm.Total > 0
}

func (rpbm *RipsiteProgressBarModel) ProgressPct() (pct float64) {
	effNb := float64(rpbm.Total)
	val := float64(rpbm.Done)
	pct = val / effNb * 100.0
	rpbm.Progress = pct
	return
}

func (rpbm *RipsiteProgressBarModel) ProgressStatus() (res string) {
	res = ""
	if rpbm.Progress >= 100 {
		res = "success"
	}
	return
}

// Filter related Funcs
//

func (rpbm *RipsiteProgressBarModel) Format() (res string) {
	if rpbm.mustShow() {
		res = strconv.Itoa(rpbm.Done) + " / " + strconv.Itoa(rpbm.Total)
		pct := strconv.FormatFloat(rpbm.ProgressPct(), 'f', 0, 64)
		res += " ( " + pct + "% )"
		return
	}
	res = "-"
	return
}
