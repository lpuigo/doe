package message

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
)

func confirmString(vm *hvue.VM, msg, msgtype string, confirm func()) {
	vm.Call("$confirm", msg, js.M{
		"confirmButtonText": "OK",
		"cancelButtonText":  "Retour",
		"type":              msgtype,
		"callback":          confirmCallBack(confirm),
	})
}

func ConfirmWarning(vm *hvue.VM, msg string, confirm func()) {
	confirmString(vm, msg, "warning", confirm)
}

func confirmCallBack(confirm func()) func(string) {
	return func(action string) {
		if action == "confirm" {
			confirm()
		}
	}
}
