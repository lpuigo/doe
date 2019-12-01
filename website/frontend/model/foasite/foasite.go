package foasite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/model/foasite/foaconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
)

// type Foa reflects backend/model/foasites.foa struct
type FoaSite struct {
	*js.Object
	Id         int    `js:"Id"`
	Client     string `js:"Client"`
	Ref        string `js:"Ref"`
	Manager    string `js:"Manager"`
	OrderDate  string `js:"OrderDate"`
	UpdateDate string `js:"UpdateDate"`
	Status     string `js:"Status"`
	Comment    string `js:"Comment"`
	Foas       []*Foa `js:"Foas"`
}

func NewFoaSite() *FoaSite {
	return &FoaSite{Object: tools.O()}
}

func FoaSiteFromJS(o *js.Object) *FoaSite {
	return &FoaSite{Object: o}
}

func (ps *FoaSite) SearchInString() string {
	return json.Stringify(ps)
}

func (fs *FoaSite) getNextId() int {
	// naive algorithm ... something smarter must be possible
	maxid := -1
	for _, foa := range fs.Foas {
		if foa.Id >= maxid {
			maxid = foa.Id + 1
		}
	}
	return maxid
}

func (fs *FoaSite) Copy(ofs *FoaSite) {
	fs.Id = ofs.Id
	fs.Client = ofs.Client
	fs.Ref = ofs.Ref
	fs.Manager = ofs.Manager
	fs.OrderDate = ofs.OrderDate
	fs.UpdateDate = ofs.UpdateDate
	fs.Status = ofs.Status
	fs.Comment = ofs.Comment
	foas := make([]*Foa, len(ofs.Foas))
	for id, foa := range ofs.Foas {
		fs.Foas[id] = foa.Clone()
	}
	fs.Foas = foas
}

func (fs *FoaSite) Clone() *FoaSite {
	return &FoaSite{Object: json.Parse(json.Stringify(fs))}
}

func FoaSiteStatusLabel(status string) string {
	switch status {
	case foaconst.FsStatusNew:
		return foaconst.FsStatusLabelNew
	case foaconst.FsStatusInProgress:
		return foaconst.FsStatusLabelInProgress
	case foaconst.FsStatusBlocked:
		return foaconst.FsStatusLabelBlocked
	case foaconst.FsStatusCancelled:
		return foaconst.FsStatusLabelCancelled
	case foaconst.FsStatusDone:
		return foaconst.FsStatusLabelDone
	default:
		return "<" + status + ">"
	}
}

func FoaSiteRowClassName(status string) string {
	var res string = ""
	switch status {
	case foaconst.FsStatusNew:
		return "worksite-row-new"
	case foaconst.FsStatusInProgress:
		return "worksite-row-inprogress"
	case foaconst.FsStatusBlocked:
		return "worksite-row-blocked"
	case foaconst.FsStatusCancelled:
		return "worksite-row-canceled"
	case foaconst.FsStatusDone:
		return "worksite-row-done"
	default:
		res = "worksite-row-error"
	}
	return res
}
