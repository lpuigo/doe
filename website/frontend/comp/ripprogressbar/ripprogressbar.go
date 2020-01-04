package ripprogressbar

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/tvprogressbar"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"strconv"
)

const template = `
<div style="width: 100%">
    <div v-if="showProgressBar" class="small-font">
		<twovalues-progressbar 
			:height="height"
			:pct1="progressPct"
			:pct2="progress2"
			:pct3="progress3"
		></twovalues-progressbar>
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
		hvue.Props("total", "billed", "done", "blocked", "height"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewWorksiteProgressBarModel(vm)
		}),
		hvue.MethodsOf(&RipsiteProgressBarModel{}),
		hvue.Computed("showProgressBar", func(vm *hvue.VM) interface{} {
			wspb := &RipsiteProgressBarModel{Object: vm.Object}
			return wspb.mustShow()
		}),
		hvue.Computed("progressPct", func(vm *hvue.VM) interface{} {
			wspb := &RipsiteProgressBarModel{Object: vm.Object}
			return wspb.ProgressPct()
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
	Billed    int     `js:"billed"`
	Done      int     `js:"done"`
	Blocked   int     `js:"blocked"`
	Progress1 float64 `js:"progress1"`
	Progress2 float64 `js:"progress2"`
	Progress3 float64 `js:"progress3"`
	Height    string  `js:"height"`

	VM *hvue.VM `js:"VM"`
}

func NewWorksiteProgressBarModel(vm *hvue.VM) *RipsiteProgressBarModel {
	rpbm := &RipsiteProgressBarModel{Object: tools.O()}
	rpbm.Total = 0
	rpbm.Billed = 0
	rpbm.Done = 0
	rpbm.Blocked = 0
	rpbm.Progress1 = .0
	rpbm.Progress2 = .0
	rpbm.Progress3 = .0
	rpbm.Height = "5px"
	rpbm.VM = vm
	return rpbm
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Formatting Related Methods

func (rpbm *RipsiteProgressBarModel) mustShow() bool {
	return rpbm.Total > 0
}

func (rpbm *RipsiteProgressBarModel) ProgressPct() float64 {
	effNb := float64(rpbm.Total)
	rpbm.Progress1 = float64(rpbm.Billed) * 100.0 / effNb
	rpbm.Progress2 = float64(rpbm.Done) * 100.0 / effNb
	rpbm.Progress3 = float64(rpbm.Blocked) * 100.0 / effNb
	return rpbm.Progress1
}

// Filter related Funcs
//

func (rpbm *RipsiteProgressBarModel) Format() (res string) {
	if !rpbm.mustShow() {
		res = "-"
		return
	}
	if rpbm.Billed+rpbm.Done+rpbm.Blocked == 0 {
		res = "0"
	} else {
		if rpbm.Billed > 0 {
			res = strconv.Itoa(rpbm.Billed)
		}
		if rpbm.Done > 0 {
			if res != "" {
				res += " + "
			}
			res += strconv.Itoa(rpbm.Done)
		}
		if rpbm.Blocked > 0 {
			if res != "" {
				res += " + "
			}
			res += strconv.Itoa(rpbm.Blocked)
		}
	}
	res += " / " + strconv.Itoa(rpbm.Total)
	pct := 0.0
	todo := rpbm.Total - rpbm.Blocked
	if todo > 0 {
		pct = float64(rpbm.Billed+rpbm.Done) * 100.0 / float64(todo)
	} else {
		pct = 0
	}
	res += " ( " + strconv.FormatFloat(pct, 'f', 0, 64) + "% )"
	return
}
