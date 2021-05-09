package ripsites

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func siteRecordFromFile(t *testing.T, dir, file string) *SiteRecord {
	site, err := NewSiteRecordFromFile(filepath.Join(dir, file))
	if err != nil {
		t.Fatalf("could not create site from file '%s': %s", file, err.Error())
	}
	return site
}

func Test_FixSite(t *testing.T) {
	dir := `C:\Users\Laurent\OneDrive\Documents\TEMPORAIRE\Sogetrel\Chantier Fibre Aube\2020-11-30 SRO 10_018_098`
	origFile := "000057_orig.json"
	newFile := "000057_fixed.json"
	resFile := "000057.json"

	origSite := siteRecordFromFile(t, dir, origFile)
	newSite := siteRecordFromFile(t, dir, newFile)

	fixMeasDict := make(map[string]*Measurement)
	for _, measurement := range newSite.Measurements {
		fixMeasDict[measurement.DestNodeName] = measurement
	}

	for _, measurement := range origSite.Site.Measurements {
		fixMeas, found := fixMeasDict[measurement.DestNodeName]
		if !found {
			t.Fatalf("could not find fixed measurement for '%s'", measurement.DestNodeName)
		}
		measurement.NbFiber = fixMeas.NbFiber
	}

	resF, err := os.Create(filepath.Join(dir, resFile))
	if err != nil {
		t.Fatalf("could not create file: %s", err.Error())
	}
	defer resF.Close()
	err = json.NewEncoder(resF).Encode(origSite)
	if err != nil {
		t.Fatalf("could not encode origSite: %s", err.Error())
	}
}

func TestSite_AuditMeasurementEvents(t *testing.T) {
	dir := `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Ripsites`
	siteFile := "000062.json"

	site := siteRecordFromFile(t, dir, siteFile)

	for _, auditMsg := range site.AuditMeasurementEvents() {
		t.Logf("%s: %s %s", auditMsg.Msg, auditMsg.NodeName, auditMsg.Info)
	}
}
