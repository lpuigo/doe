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
		<div :class="Class1" :style="Style1"></div>
		<div v-if="pct2>0" :class="Class2" :style="Style2"></div>
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
		hvue.Props("pct1", "pct2", "height"),
		hvue.Computed("Style0", func(vm *hvue.VM) interface{} {
			tvpbm := &TwoValuesProgressBarModel{Object: vm.Object}
			return "height:" + tvpbm.Height
		}),
		hvue.Computed("Class1", func(vm *hvue.VM) interface{} {
			tvpbm := &TwoValuesProgressBarModel{Object: vm.Object}
			return tvpbm.SetParam()
		}),
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
	Height string  `js:"height"`

	Class2 string `js:"Class2"`
	Style1 string `js:"Style1"`
	Style2 string `js:"Style2"`

	VM *hvue.VM `js:"VM"`
}

func NewTwoValuesProgressBarModel(vm *hvue.VM) *TwoValuesProgressBarModel {
	tvpbm := &TwoValuesProgressBarModel{Object: tools.O()}

	tvpbm.Pct1 = 0.0
	tvpbm.Pct2 = 0.0
	tvpbm.Height = "6px"

	tvpbm.Class2 = ""
	tvpbm.Style1 = ""
	tvpbm.Style2 = ""

	tvpbm.VM = vm
	return tvpbm
}

func (tvpbm *TwoValuesProgressBarModel) SetParam() string {
	class1 := "inner begin"
	if tvpbm.Pct1 >= 100-tvpbm.Pct2 {
		class1 += " complete"
	} else {
		class1 += " inprogress"
	}
	if tvpbm.Pct2 <= 0.0 {
		class1 += " end"
	}
	tvpbm.Class2 = "inner ko"
	if tvpbm.Pct1 <= 0.0 {
		tvpbm.Class2 += " begin"
	}
	tvpbm.Class2 += " end"
	pct1 := strconv.FormatFloat(tvpbm.Pct1, 'f', 3, 64) + "%"
	tvpbm.Style1 = "width: " + pct1
	tvpbm.Style2 = "width: " + strconv.FormatFloat(tvpbm.Pct2, 'f', 3, 64) + "%; left: " + pct1
	return class1
}
