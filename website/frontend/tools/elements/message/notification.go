package message

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
)

func notify(vm *hvue.VM, notifType, title, msg string, dur int) {
	vm.Call("$notify", js.M{
		"title":    title,
		"message":  msg,
		"type":     notifType,
		"duration": dur,
		"offset":   0,
	})
}

func NotifySuccess(vm *hvue.VM, title, msg string) {
	notify(vm, "success", title, msg, 3500)
}

func NotifyWarning(vm *hvue.VM, title, msg string) {
	notify(vm, "warning", title, msg, 0)
}

func NotifyError(vm *hvue.VM, title, msg string) {
	notify(vm, "error", title, msg, 5500)
}
