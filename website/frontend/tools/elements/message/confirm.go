package message

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
)

func confirmString(vm *hvue.VM, msg, msgtype string) {
	vm.Call("$confirm", msg, js.M{
		"confirmButtonText": "OK",
		"cancelButtonText":  "Abandon",
		"type":              msgtype,
		"callback":          confirmCallBack,
	})
}

func ConfirmWarning(vm *hvue.VM, msg string) {
	confirmString(vm, msg, "warning")
}

func confirmCallBack(action string) {
	print("confirmCallBack", action)
}
