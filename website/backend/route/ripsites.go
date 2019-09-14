package route

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lpuig/ewin/doe/website/backend/logger"
	mgr "github.com/lpuig/ewin/doe/website/backend/manager"
	"github.com/lpuig/ewin/doe/website/backend/model/ripsites"
	"github.com/lpuig/ewin/doe/website/backend/model/ripsites/measurementreport"
	"net/http"
	"strconv"
	"strings"
)

func GetRipsitesInfo(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetRipsitesInfo").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")

	err := mgr.GetRipsitesInfo(w)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
}

func GetRipSite(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetRipSite").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	wsrid, err := strconv.Atoi(vars["rsid"])
	if err != nil {
		AddError(w, logmsg, "mis-formatted RipSite id '"+vars["rsid"]+"'", http.StatusBadRequest)
		return
	}
	rsr := mgr.Ripsites.GetById(wsrid)
	if rsr == nil {
		AddError(w, logmsg, fmt.Sprintf("ripSite with id %d does not exist", wsrid), http.StatusNotFound)
		return
	}
	err = rsr.Marshall(w)
	if err != nil {
		AddError(w, logmsg, "could not marshall RipSite. "+err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("ripSite Id %d (%s) returned", rsr.Id, rsr.Ref), http.StatusOK)
}

func UpdateRipSite(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("UpdateRipSite").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	if r.Body == nil {
		AddError(w, logmsg, "request RipSite missing", http.StatusBadRequest)
		return
	}
	reqRipSiteId := mux.Vars(r)["rsid"]
	rsrid, err := strconv.Atoi(reqRipSiteId)
	if err != nil {
		AddError(w, logmsg, "mis-formatted RipSite id '"+reqRipSiteId+"'", http.StatusBadRequest)
		return
	}
	rsr := mgr.Ripsites.GetById(rsrid)
	if rsr == nil {
		AddError(w, logmsg, fmt.Sprintf("ripSite with id %d does not exist", rsrid), http.StatusNotFound)
		return
	}
	rsr, err = ripsites.NewSiteRecordFrom(r.Body)
	if err != nil {
		AddError(w, logmsg, "malformed RipSite: "+err.Error(), http.StatusBadRequest)
		return
	}
	if rsr.Id != rsrid {
		AddError(w, logmsg, fmt.Sprintf("inconsitent RipSite id between request (%d) and body (%d)", rsrid, rsr.Id), http.StatusBadRequest)
		return
	}
	err = mgr.Ripsites.Update(rsr)
	if err != nil {
		AddError(w, logmsg, fmt.Sprintf("could not update RipSite with id %d: %v", rsrid, err), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("RipSite with id %d (%s) updated", rsrid, rsr.Ref), http.StatusOK)
}

func DeleteRipSite(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("DeleteRipSite").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	reqRipSiteId := mux.Vars(r)["rsid"]
	rsrid, err := strconv.Atoi(reqRipSiteId)
	if err != nil {
		AddError(w, logmsg, "mis-formatted RipSite id '"+reqRipSiteId+"'", http.StatusBadRequest)
		return
	}
	rsr := mgr.Ripsites.GetById(rsrid)
	if rsr == nil {
		AddError(w, logmsg, fmt.Sprintf("ripSite with id %d does not exist", rsrid), http.StatusNoContent)
		return
	}
	err = mgr.Ripsites.Remove(rsr)
	if err != nil {
		AddError(w, logmsg, fmt.Sprintf("could not delete RipSite with id %d: %v", rsrid, err), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("RipSite with id %d (%s) deleted", rsrid, rsr.Ref), http.StatusOK)
}

// GetRipSiteAttachement
func GetRipSiteAttachement(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetRipSiteAttachement").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	reqRipSiteId := mux.Vars(r)["rsid"]
	rsrid, err := strconv.Atoi(reqRipSiteId)
	if err != nil {
		AddError(w, logmsg, "mis-formatted RipSite id '"+reqRipSiteId+"'", http.StatusBadRequest)
		return
	}
	rsr := mgr.Ripsites.GetById(rsrid)
	if rsr == nil {
		AddError(w, logmsg, fmt.Sprintf("ripSite with id %d does not exist", rsrid), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", mgr.TemplateEngine.GetRipsiteXLSAttachementName(rsr.Site)))
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	err = mgr.GetRipsiteXLSAttachement(w, rsr.Site)
	if err != nil {
		AddError(w, logmsg, "could not generate RipSite XLS Attachment file. "+err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("Attachment XLS produced for RipSite id %d (%s)", rsrid, rsr.Site.Ref), http.StatusOK)
}

func GetRipsitesStats(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetRipsitesStats").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	groupBy := vars["groupby"]
	freq := vars["freq"]
	err := mgr.GetRipsitesStats(w, freq, groupBy)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("%s ripsites %s stats produced", freq, groupBy), http.StatusOK)
}

func GetRipsitesArchive(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetRipsitesArchive").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", mgr.RipsitesArchiveName()))
	w.Header().Set("Content-Type", "application/zip")

	err := mgr.CreateRipsitesArchive(w)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
}

func MeasurementRipSite(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("MeasurementRipSite").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")

	// Parse our multipart form, 30 << 20 specifies a maximum
	// upload of 30 MB files.
	if r.ParseMultipartForm(30<<20) != nil {
		AddError(w, logmsg, "file info missing", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		AddError(w, logmsg, "error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if !strings.HasSuffix(strings.ToUpper(handler.Filename), ".ZIP") {
		AddError(w, logmsg, "uploaded file is not a Zip Archive", http.StatusBadRequest)
		return
	}

	mr, err := measurementreport.ParseZipMeasurementFiles(file, handler.Size)
	if err != nil {
		AddError(w, logmsg, fmt.Sprintf("error processing measurement Zip archive : %s", err.Error()), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(mr)

	logmsg.Response = http.StatusOK
}
