package message

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/prjptf/src/client/goel/message"
	"honnef.co/go/js/xhr"
	"strconv"
)

var (
	duration int    = 2000
	msgclass string = "message"
)

func SetDuration(msec int) {
	duration = msec
}

func messageString(vm *hvue.VM, msgtype, msg string, close bool) {
	vm.Call("$message", js.M{
		"showClose":   close,
		"message":     msg,
		"type":        msgtype,
		"duration":    duration,
		"customClass": msgclass,
	})
}

func InfoStr(vm *hvue.VM, msg string, closeButton bool) {
	pdur := duration
	if closeButton {
		duration = 0
	}
	messageString(vm, "info", msg, closeButton)
	duration = pdur
}

func SuccesStr(vm *hvue.VM, msg string) {
	messageString(vm, "success", msg, false)
}

func WarningStr(vm *hvue.VM, msg string) {
	messageString(vm, "warning", msg, false)
}

func ErrorStr(vm *hvue.VM, msg string, closeButton bool) {
	pdur := duration
	if closeButton {
		duration = 0
	}
	messageString(vm, "error", msg, closeButton)
	duration = pdur
}

func ErrorMsgStr(vm *hvue.VM, msg string, o *js.Object, closeButton bool) {
	msg += ErrorMsgFromJS(o).Error
	ErrorStr(vm, msg, closeButton)
}

type ErrorMsg struct {
	*js.Object
	Error string `js:"Error"`
}

func ErrorMsgFromJS(o *js.Object) *ErrorMsg {
	return &ErrorMsg{Object: o}
}

func ErrorMessage(vm *hvue.VM, req *xhr.Request) {
	message.SetDuration(tools.WarningMsgDuration)
	msg := "Quelquechose c'est mal passé !\n"
	msg += "Le server retourne un code " + strconv.Itoa(req.Status) + "\n"
	ErrorMsgStr(vm, msg, req.Response, true)
}
