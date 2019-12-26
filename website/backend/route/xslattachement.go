package route

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lpuig/ewin/doe/website/backend/logger"
	mgr "github.com/lpuig/ewin/doe/website/backend/manager"
	"net/http"
	"strconv"
)

func GetItemizableSiteAttachement(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetItemizableSiteAttachement").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	siteType := mux.Vars(r)["sitetype"]
	reqSiteId := mux.Vars(r)["id"]
	rsid, err := strconv.Atoi(reqSiteId)
	if err != nil {
		AddError(w, logmsg, "mis-formatted FoaSite id '"+reqSiteId+"'", http.StatusBadRequest)
		return
	}

	site, err := mgr.GetItemizableSite(siteType)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusBadRequest)
		return
	}

	psr := site.GetItemizableSiteById(rsid)
	if psr == nil {
		AddError(w, logmsg, fmt.Sprintf("%s with id %d does not exist", siteType, rsid), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", mgr.GetItemizableSiteXLSAttachementName(psr)))
	w.Header().Set("Content-Type", "vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	err = mgr.GetItemizableSiteXLSAttachement(w, psr)
	if err != nil {
		AddError(w, logmsg, fmt.Sprintf("could not generate %s XLS Attachment file. %s", siteType, err.Error()), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("Attachment XLS produced for %s id %d (%s)", siteType, rsid, psr.GetRef()), http.StatusOK)
}
