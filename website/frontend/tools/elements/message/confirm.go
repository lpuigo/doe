package message

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
)

func confirmString(vm *hvue.VM, msg, msgtype string, confirm func()) {
	vm.Call("$confirm", msg, js.M{
		"confirmButtonText": "OK",
		"cancelButtonText":  "Non",
		"type":              msgtype,
		"callback":          confirmCallBack(confirm),
	})
}

func ConfirmWarning(vm *hvue.VM, msg string, confirm func()) {
	confirmString(vm, msg, "warning", confirm)
}

func ConfirmSuccess(vm *hvue.VM, msg string, confirm func()) {
	confirmString(vm, msg, "success", confirm)
}

func confirmCallBack(confirm func()) func(string) {
	return func(action string) {
		if action == "confirm" {
			confirm()
		}
	}
}

func ConfirmCancelWarning(vm *hvue.VM, msg string, confirm, cancel func()) {
	vm.Call("$confirm", msg, js.M{
		"confirmButtonText": "OK",
		"cancelButtonText":  "Non",
		"type":              "warning",
		"callback": func(action string) {
			switch action {
			case "confirm":
				confirm()
			case "cancel":
				cancel()
			default:
			}
		},
	})
}
