package tools

import "github.com/gopherjs/gopherjs/js"

// BeforeUnloadConfirmation activate confirm leave alert if askBeforeLeave func return true
func BeforeUnloadConfirmation(askBeforeLeave func() bool) {
	js.Global.Get("window").Call(
		"addEventListener",
		"beforeunload",
		func(event *js.Object) {
			if !askBeforeLeave() {
				return
			}
			event.Call("preventDefault")
			event.Set("returnValue", "")
			//js.Global.Call("confirm", "Sur ?")

		},
		false)
}

func GetLocationQueryString() string {
	location := js.Global.Get("location")
	return location.Get("search").String()
}

func GetURLSearchParam(param string) *js.Object {
	usp := js.Global.Get("URLSearchParams").New(GetLocationQueryString())
	return usp.Call("get", param)
}

func WindowHasHover() bool {
	res := js.Global.Get("window").Call("matchMedia", "(any-hover: hover)")
	return res.Get("matches").Bool()
}
