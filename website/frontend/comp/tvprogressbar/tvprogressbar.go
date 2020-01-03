package tvprogressbar

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"strconv"
)

const template = `
<div class="twovalues-progressbar">
	<div class="outer" :style="Style0">
		<div v-if="pct1>0" :class="Class1" :style="Style1"></div>
		<div v-if="pct2>0" :class="Class2" :style="Style2"></div>
		<div v-if="pct3>0" :class="Class3" :style="Style3"></div>
	</div>	
</div>
`

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

// RegisterComponent registers "twovalues-progressbar" component
func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("twovalues-progressbar", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("pct1", "pct2", "pct3", "height"),
		hvue.Computed("Style0", func(vm *hvue.VM) interface{} {
			tvpbm := &TwoValuesProgressBarModel{Object: vm.Object}
			tvpbm.SetParam()
			return "height:" + tvpbm.Height
		}),
		//hvue.Computed("Class1", func(vm *hvue.VM) interface{} {
		//	tvpbm := &TwoValuesProgressBarModel{Object: vm.Object}
		//	return tvpbm.SetParam()
		//}),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewTwoValuesProgressBarModel(vm)
		}),
		hvue.MethodsOf(&TwoValuesProgressBarModel{}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type TwoValuesProgressBarModel struct {
	*js.Object

	Pct1   float64 `js:"pct1"`
	Pct2   float64 `js:"pct2"`
	Pct3   float64 `js:"pct3"`
	Height string  `js:"height"`

	Class1 string `js:"Class1"`
	Class2 string `js:"Class2"`
	Class3 string `js:"Class3"`
	Style1 string `js:"Style1"`
	Style2 string `js:"Style2"`
	Style3 string `js:"Style3"`

	VM *hvue.VM `js:"VM"`
}

func NewTwoValuesProgressBarModel(vm *hvue.VM) *TwoValuesProgressBarModel {
	tvpbm := &TwoValuesProgressBarModel{Object: tools.O()}

	tvpbm.Pct1 = 0.0
	tvpbm.Pct2 = 0.0
	tvpbm.Pct3 = 0.0
	tvpbm.Height = "6px"

	tvpbm.Class1 = ""
	tvpbm.Class2 = ""
	tvpbm.Class3 = ""
	tvpbm.Style1 = ""
	tvpbm.Style2 = ""
	tvpbm.Style3 = ""

	tvpbm.VM = vm
	return tvpbm
}

func (tvpbm *TwoValuesProgressBarModel) SetParam() {
	fpct := func(pct float64) string {
		if pct > 100 {
			pct = 100
		}
		if pct < 0 {
			pct = 0
		}
		return strconv.FormatFloat(pct, 'f', 1, 64) + "%"
	}

	// class 1
	tvpbm.Class1 = "inner begin inprogress2"
	if tvpbm.Pct1 >= 100 {
		tvpbm.Class1 += " end"
	}
	// class 2
	tvpbm.Class2 = "inner"
	if tvpbm.Pct1 <= 0 {
		tvpbm.Class2 += " begin"
	}
	if tvpbm.Pct2 >= 100-tvpbm.Pct1-tvpbm.Pct3 {
		tvpbm.Class2 += " complete"
	} else {
		tvpbm.Class2 += " inprogress"
	}
	if tvpbm.Pct3 <= 0 {
		tvpbm.Class2 += " end"
	}
	// class 3
	tvpbm.Class3 = "inner ko"
	if tvpbm.Pct1+tvpbm.Pct2 <= 0.0 {
		tvpbm.Class3 += " begin"
	}
	tvpbm.Class3 += " end"

	pct1 := fpct(tvpbm.Pct1)
	pct2 := fpct(tvpbm.Pct2 + .1) // add 0.1% to avoid display aliasing issue (white bar)
	pct12 := fpct(tvpbm.Pct1 + tvpbm.Pct2)
	pct3 := fpct(tvpbm.Pct3)
	tvpbm.Style1 = "width: " + pct1
	tvpbm.Style2 = "width: " + pct2 + "; left: " + pct1
	tvpbm.Style3 = "width: " + pct3 + "; left: " + pct12
}
