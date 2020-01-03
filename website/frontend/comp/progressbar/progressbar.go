package progressbar

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/tvprogressbar"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"strconv"
)

const template = `
<div>
    <div v-if="showProgressBar" class="small-font">
		<twovalues-progressbar 
			height="5px"
			pct1=0
			:pct2="progressPct"
			:pct3="progressKo"
		></twovalues-progressbar>
        <span>{{ progressText }}</span>
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
		tvprogressbar.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("value", "measure"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewWorksiteProgressBarModel(vm)
		}),
		hvue.MethodsOf(&WorksiteProgressBarModel{}),
		hvue.Computed("showProgressBar", func(vm *hvue.VM) interface{} {
			wspb := &WorksiteProgressBarModel{Object: vm.Object}
			return wspb.mustShow()
		}),
		hvue.Computed("progressPct", func(vm *hvue.VM) interface{} {
			wspb := &WorksiteProgressBarModel{Object: vm.Object}
			return wspb.ProgressPct()
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
	ProgressOK   float64          `js:"progressOk"`
	ProgressKO   float64          `js:"progressKo"`
	Measure      bool             `js:"measure"`

	VM *hvue.VM `js:"VM"`
}

func NewWorksiteProgressBarModel(vm *hvue.VM) *WorksiteProgressBarModel {
	wpbm := &WorksiteProgressBarModel{Object: tools.O()}
	wpbm.WorksiteInfo = nil
	wpbm.ProgressOK = 0
	wpbm.ProgressKO = 0
	wpbm.Measure = false
	wpbm.VM = vm
	return wpbm
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Formatting Related Methods

func (wspbm *WorksiteProgressBarModel) mustShow() bool {
	return wspbm.WorksiteInfo.NbElTotal > 0
}

func (wspbm *WorksiteProgressBarModel) ProgressPct() (pctOK float64) {
	wsi := wspbm.WorksiteInfo
	totNb := float64(wsi.NbElTotal)
	//totNb := float64(wsi.NbElTotal - wsi.NbElBlocked)
	valOK := float64(wsi.NbElInstalled)
	valKO := float64(wsi.NbElBlocked)
	if wspbm.Measure {
		valOK = float64(wsi.NbElMeasured)
		totNb -= valKO
		valKO = 0
	}
	pctOK = valOK * 100.0 / totNb
	wspbm.ProgressOK = pctOK
	wspbm.ProgressKO = valKO * 100.0 / totNb
	return
}

// Filter related Funcs
//

func (wspbm *WorksiteProgressBarModel) Format() (res string) {
	if !wspbm.mustShow() {
		res = "-"
		return
	}
	wsi := wspbm.WorksiteInfo
	pct, val := 0.0, 0.0
	if !wspbm.Measure {
		res = strconv.Itoa(wsi.NbElInstalled)
		if wsi.NbElBlocked > 0 {
			res += " + " + strconv.Itoa(wsi.NbElBlocked)
		}
		res += " / " + strconv.Itoa(wsi.NbElTotal)
		val = float64(wsi.NbElInstalled)
	} else {
		res = strconv.Itoa(wsi.NbElMeasured) + " / " + strconv.Itoa(wsi.NbElTotal-wsi.NbElBlocked)
		val = float64(wsi.NbElMeasured)
	}
	if wsi.NbElTotal-wsi.NbElBlocked > 0 {
		pct = val * 100.0 / float64(wsi.NbElTotal-wsi.NbElBlocked)
	} else {
		pct = 0.0
	}
	res += " ( " + strconv.FormatFloat(pct, 'f', 0, 64) + "% )"
	return
}
