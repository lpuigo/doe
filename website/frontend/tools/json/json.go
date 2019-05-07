// package JSON wraps the javascritp JSON api for GOPHERJS.
package json

import (
	"github.com/gopherjs/gopherjs/js"
)

func Stringify(obj interface{}) string {
	json := js.Global.Get("JSON")
	return json.Call("stringify", obj).String()
}

func Parse(jsonStr string) *js.Object {
	json := js.Global.Get("JSON")
	return json.Call("parse", jsonStr)
}
