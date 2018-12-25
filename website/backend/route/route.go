package route

import (
	"encoding/json"
	"github.com/lpuig/ewin/doe/website/backend/logger"
	mgr "github.com/lpuig/ewin/doe/website/backend/manager"
	"net/http"
)

type ErrorMsg struct {
	Error string
}

type MgrHandlerFunc func(*mgr.Manager, http.ResponseWriter, *http.Request)

func addError(w http.ResponseWriter, errmsg string, code int) string {
	res := logger.LogResponse(code)
	res += logger.LogInfo(errmsg)
	em := ErrorMsg{Error: errmsg}
	sem, _ := json.Marshal(em)
	http.Error(w, string(sem), code)
	return res
}
