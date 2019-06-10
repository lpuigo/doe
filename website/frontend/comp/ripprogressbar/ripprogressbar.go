package ripprogressbar

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/tvprogressbar"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"strconv"
)

const template = `
<div>
    <div v-if="showProgressBar" class="small-font">
		<twovalues-progressbar 
			height="5px"
			:pct1="progressPct"
			:pct2="progress2"
		></twovalues-progressbar>
<!--        <el-progress -->
<!--                 :show-text="false"-->
<!--                 :stroke-width="5"-->
<!--                 :percentage="progressPct"-->
<!--				 :status="progressStatus"-->
<!--        ></el-progress>-->
        <span>{{ progressText }}</span>
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
		tvprogressbar.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("total", "done", "blocked"),
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

	Total     int     `js:"total"`
	Done      int     `js:"done"`
	Blocked   int     `js:"blocked"`
	Progress1 float64 `js:"progress1"`
	Progress2 float64 `js:"progress2"`

	VM *hvue.VM `js:"VM"`
}

func NewWorksiteProgressBarModel(vm *hvue.VM) *RipsiteProgressBarModel {
	rpbm := &RipsiteProgressBarModel{Object: tools.O()}
	rpbm.Total = 0
	rpbm.Done = 0
	rpbm.Blocked = 0
	rpbm.Progress1 = 0
	rpbm.Progress2 = 0
	rpbm.VM = vm
	return rpbm
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Formatting Related Methods

func (rpbm *RipsiteProgressBarModel) mustShow() bool {
	return rpbm.Total > 0
}

func (rpbm *RipsiteProgressBarModel) ProgressPct() (pctOK float64) {
	effNb := float64(rpbm.Total)
	valOK := float64(rpbm.Done)
	pctOK = valOK * 100.0 / effNb
	rpbm.Progress1 = pctOK
	rpbm.Progress2 = float64(rpbm.Blocked) * 100.0 / effNb
	return
}

// Filter related Funcs
//

func (rpbm *RipsiteProgressBarModel) Format() (res string) {
	if !rpbm.mustShow() {
		res = "-"
		return
	}
	res = strconv.Itoa(rpbm.Done)
	if rpbm.Blocked > 0 {
		res += " + " + strconv.Itoa(rpbm.Blocked)
	}
	res += " / " + strconv.Itoa(rpbm.Total)
	pct := 0.0
	if rpbm.Total-rpbm.Blocked > 0 {
		pct = float64(rpbm.Done) * 100.0 / float64(rpbm.Total-rpbm.Blocked)
	} else {
		pct = 0
	}
	res += " ( " + strconv.FormatFloat(pct, 'f', 0, 64) + "% )"
	return
}
