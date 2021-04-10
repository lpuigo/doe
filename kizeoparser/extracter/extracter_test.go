package extracter

import (
	"fmt"
	"path/filepath"
	"testing"
)

const (
	ExtracterConfigDir string = `C:\Users\Laurent\OneDrive\Documents\TEMPORAIRE\KizeoExtract`
	ConfigFile         string = `config.json`
)

func TestExtracter_SaveConfig(t *testing.T) {
	extracter := Extracter{
		ConfigPath: filepath.Join(ExtracterConfigDir, ConfigFile),
		Config: Config{
			ProgressFile:    "progress.json",
			ProgressXLSFile: "Synthese.xlsx",
			Forms: []FormConfig{
				{
					FormId:   664879,
					FormName: "Poteau Eiffage Signes",
				},
				{
					FormId:   630190,
					FormName: "Poteau Fiitelcom",
				},
				{
					FormId:   640312,
					FormName: "Poteau Sogetrel",
				},
			},
		},
	}
	err := extracter.SaveConfig()
	if err != nil {
		t.Fatalf("SaveConfig returned unexpected %s", err.Error())
	}
}

func TestNewExtracterFromConfigFile(t *testing.T) {
	_, err := NewExtracterFromConfigFile(ExtracterConfigDir, ConfigFile)
	if err != nil {
		t.Fatalf("NewExtracterFromConfigFile returned unexpected %s", err.Error())
	}
}

func TestExtracter_SaveProgressFile(t *testing.T) {
	extracter, err := NewExtracterFromConfigFile(ExtracterConfigDir, ConfigFile)
	if err != nil {
		t.Fatalf("NewExtracterFromConfigFile returned unexpected %s", err.Error())
	}
	extracter.InitProgress()
	err = extracter.SaveProgressFile()
	if err != nil {
		t.Fatalf("SaveProgressFile returned unexpected %s", err.Error())
	}
}

func TestExtracter_LoadProgressFile(t *testing.T) {
	extracter, err := NewExtracterFromConfigFile(ExtracterConfigDir, ConfigFile)
	if err != nil {
		t.Fatalf("NewExtracterFromConfigFile returned unexpected %s", err.Error())
	}
	err = extracter.LoadProgressFile()
	if err != nil {
		t.Fatalf("LoadProgressFile returned unexpected %s", err.Error())
	}
	for i, formulaire := range extracter.Progress.Formulaires {
		fmt.Printf("Formulaire %d: %s > %s\n", i, formulaire.FormName, formulaire.ExtractionDate)
	}
}

func TestExtracter_GetKizeoProgress(t *testing.T) {
	extracter, err := NewExtracterFromConfigFile(ExtracterConfigDir, ConfigFile)
	if err != nil {
		t.Fatalf("NewExtracterFromConfigFile returned unexpected %s", err.Error())
	}
	err = extracter.LoadProgressFile()
	if err != nil {
		t.Fatalf("LoadProgressFile returned unexpected %s", err.Error())
	}
	err = extracter.GetKizeoProgress()
	if err != nil {
		t.Fatalf("LoadProgressFile returned unexpected %s", err.Error())
	}
}
