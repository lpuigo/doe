package tools

import "github.com/gopherjs/gopherjs/js"

func O() *js.Object {
	return js.Global.Get("Object").New()
}

func Empty(s string) bool {
	return !(s != "" && s != "null")
}
