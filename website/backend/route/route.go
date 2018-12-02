package route

import (
	"github.com/lpuig/ewin/doe/website/backend/logger"
	mgr "github.com/lpuig/ewin/doe/website/backend/manager"
	"net/http"
)

type MgrHandlerFunc func(*mgr.Manager, http.ResponseWriter, *http.Request)

func addError(w http.ResponseWriter, errmsg string, code int) string {
	res := logger.LogResponse(code)
	res += logger.LogInfo(errmsg)
	http.Error(w, errmsg, code)
	return res
}
