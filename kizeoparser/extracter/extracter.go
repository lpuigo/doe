package extracter

import (
	"encoding/json"
	"fmt"
	"github.com/lpuig/ewin/doe/kizeoparser/api"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Extracter struct {
	ConfigPath   string
	ConfigFile   string
	Config       Config
	Progress     Progress
	KizeoContext *api.KizeoContext
}

func NewExtracterFromConfigFile(path, file string) (*Extracter, error) {
	configFile, err := os.Open(filepath.Join(path, file))
	if err != nil {
		return nil, fmt.Errorf("open returned unexpected: %s", err.Error())
	}
	defer configFile.Close()

	extracter := &Extracter{
		ConfigPath: path,
		ConfigFile: file,
		Config:     Config{},
		Progress:   NewProgress(),
	}
	err = json.NewDecoder(configFile).Decode(&extracter.Config)
	if err != nil {
		return nil, fmt.Errorf("decode returned unexpected: %s", err.Error())
	}
	extracter.KizeoContext = api.NewKizeoContext()
	return extracter, nil
}

func (e Extracter) SaveConfig() error {
	configFile, err := os.Create(filepath.Join(e.ConfigPath, e.ConfigFile))
	if err != nil {
		return fmt.Errorf("create returned unexpected: %s", err.Error())
	}
	defer configFile.Close()

	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(e.Config)
	if err != nil {
		return fmt.Errorf("encode returned unexpected: %s", err.Error())
	}
	return nil
}

func (e *Extracter) InitProgress() {
	e.Progress = Progress{Formulaires: make([]*FormProgress, len(e.Config.Forms))}
	for i, form := range e.Config.Forms {
		e.Progress.Formulaires[i] = &FormProgress{
			FormConfig:     form,
			ExtractionDate: date.Today().AddDays(-15).String(),
		}
	}
}

func (e *Extracter) LoadProgressFile() error {
	progressFile, err := os.Open(filepath.Join(e.ConfigPath, e.Config.ProgressFile))
	if err != nil {
		return fmt.Errorf("open returned unexpected: %s", err.Error())
	}
	defer progressFile.Close()

	err = json.NewDecoder(progressFile).Decode(&e.Progress)
	if err != nil {
		return fmt.Errorf("decode returned unexpected: %s", err.Error())
	}
	return nil
}

func (e Extracter) SaveProgressFile() error {
	progressFile, err := os.Create(filepath.Join(e.ConfigPath, e.Config.ProgressFile))
	if err != nil {
		return fmt.Errorf("create returned unexpected: %s", err.Error())
	}
	defer progressFile.Close()

	encoder := json.NewEncoder(progressFile)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(e.Progress)
	if err != nil {
		return fmt.Errorf("encode returned unexpected: %s", err.Error())
	}
	return nil
}

func (e *Extracter) GetKizeoProgress() error {
	if e.KizeoContext == nil {
		return fmt.Errorf("kizeo context is not set")
	}

	doSaveProgress := false
	defer func() {
		if !doSaveProgress {
			return
		}
		err := e.SaveProgressFile()
		if err != nil {
			fmt.Printf("failed to save progress: %s", err.Error())
			return
		}
		fmt.Printf("progress saved.")
	}()

	kc := api.NewKizeoContext()
	for i, form := range e.Progress.Formulaires {
		fmt.Printf("Kizeo search for '%s': ... ", form.FormName)

		t := time.Now()
		datas, err := kc.FormDatasSince(strconv.Itoa(form.FormId), form.ExtractionDate)
		if err != nil {
			fmt.Printf("failed after %s\n", time.Since(t).String())
			return fmt.Errorf("Kizeo Search for '%s' retured unexpected: %s", form.FormName, err.Error())
		}
		fmt.Printf("retreived %d record in %s\n", len(datas), time.Since(t).String())

		if len(datas) > 0 {
			e.Progress.forms[form.FormId] = datas
			nextDate := form.ExtractionDate
			for _, data := range datas {
				if data.UpdateTime > nextDate {
					nextDate = data.UpdateTime
				}
			}
			e.Progress.Formulaires[i].ExtractionDate = nextDate
			doSaveProgress = true
		}
	}
	return nil
}
