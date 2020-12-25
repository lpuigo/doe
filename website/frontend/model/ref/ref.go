package ref

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type Ref struct {
	*js.Object
	Reference    string        `js:"Reference"`
	Dirty        bool          `js:"Dirty"`
	GetReference func() string `js:"GetReference"`
}

func NewRef(getRefFunc func() string) *Ref {
	ref := &Ref{Object: tools.O()}
	ref.Reference = ""
	ref.Dirty = false
	ref.GetReference = getRefFunc
	return ref
}

func (r *Ref) SetReference() {
	r.Reference = r.GetReference()
	r.Dirty = false
}

// IsDirty sets Dirty flag to true if Reference if no longer up to date, false otherwise, and return its new value
func (r *Ref) IsDirty() bool {
	r.Dirty = r.Reference != r.GetReference()
	return r.Dirty
}
