package manager

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/items"
	"io"
)

func (m Manager) GetItemizableSite(siteType string) (site items.ItemizableContainer, err error) {
	switch siteType {
	case "ripsites":
		site = m.Ripsites
	case "polesites":
		site = m.Polesites
	case "foasites":
		site = m.Foasites
	default:
		err = fmt.Errorf("'%s' site type not handled", siteType)
	}
	return
}

func (m Manager) GetItemizableSiteXLSAttachementName(site items.ItemizableSite) string {
	return fmt.Sprintf("%s ATTACHEMENT %s.xlsx", date.Today().String(), site.GetRef())
}

func (m Manager) GetItemizableSiteXLSAttachement(writer io.Writer, site items.ItemizableSite) error {
	return m.TemplateEngine.GetItemizableSiteXLSAttachement(writer, site, m.genGetClient(), m.genActorById())
}
