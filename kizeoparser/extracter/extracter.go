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

// GetKizeoProgress retieves data from Kizeo API for all reciever.Progress.Formulaires.
//
// For each Formulaire with new Kizeo Form found, receiver.Progress is populated and ExtractionDate is updated.
//
// If Kizeo Form's Pole Réference is invalid, form is set as not to be extracted (form.ExtractData = false)
func (e *Extracter) GetKizeoProgress() error {
	if e.KizeoContext == nil {
		return fmt.Errorf("kizeo context is not set")
	}

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

		if len(datas) == 0 {
			continue
		}

		e.Progress.forms[form.FormId] = datas
		nextDate := form.ExtractionDate
		for _, data := range datas {
			data.ExtractData = data.CheckSroRef()
			if !data.ExtractData {
				fmt.Printf("\twarning in %s: misformatted pole ref : %s\n", form.FormName, data.GetRawSroRef())
			}
			if data.UpdateTime > nextDate {
				nextDate = data.UpdateTime
			}
		}
		e.Progress.Formulaires[i].ExtractionDate = nextDate
	}
	return nil
}

func (e Extracter) WriteXLSForms() error {
	timeinfo := date.Now().TimeStampShort()
	for _, formulaire := range e.Progress.Formulaires {
		forms, found := e.Progress.forms[formulaire.FormId]
		if !found {
			continue
		}
		if len(forms) == 0 {
			continue
		}
		file := filepath.Join(e.ConfigPath, timeinfo+"_"+formulaire.FormName+".xlsx")
		err := writeXlsFormsFile(file, formulaire.FormName, forms)
		if err != nil {
			return fmt.Errorf("writeXlsFormsFile returned unexpected: %s", err.Error())
		}
	}
	return nil
}

func (e *Extracter) ReadXLSForms(timestamp string) error {
	files, err := filepath.Glob(filepath.Join(e.ConfigPath, timestamp) + "*.xlsx")
	if err != nil {
		return fmt.Errorf("glob returned unexpected: %s", err.Error())
	}
	for _, file := range files {
		err = readXlsFormsFile(file, e.Progress.forms)
		if err != nil {
			return fmt.Errorf("readXlsFormsFile from '%s' returned: %s", filepath.Base(file), err.Error())
		}
	}
	return nil
}

func (e *Extracter) ReadXLSFormsFromFile(file string) error {
	actualFile := filepath.Join(e.ConfigPath, file)
	if !fileExists(actualFile) {
		return fmt.Errorf("file '%s' not found", actualFile)
	}
	err := readXlsFormsFile(actualFile, e.Progress.forms)
	if err != nil {
		return fmt.Errorf("readXlsFormsFile from '%s' returned: %s", actualFile, err.Error())
	}
	return nil
}

func (e *Extracter) ExtractFormsData(timestamp string) error {
	for formId, records := range e.Progress.forms {
		formName := e.GetFormNameById(formId)
		if formName == "" {
			fmt.Printf("skipping unknown forulaire Id %d\n", formId)
			continue
		}
		path := filepath.Join(e.ConfigPath, "Extract_"+timestamp, formName)
		err := e.ExtractRecords(path, records)
		if err != nil {
			return fmt.Errorf("extractRecords from '%s' returned: %s", formName, err.Error())
		}
	}
	return nil
}

func (e Extracter) GetFormNameById(formId int) string {
	for _, formConf := range e.Config.Forms {
		if formConf.FormId == formId {
			return formConf.FormName
		}
	}
	return ""
}

func (e Extracter) ExtractRecords(path string, records []*api.SearchData) error {
	parallel := 4
	for _, record := range records {
		if !record.ExtractData {
			continue
		}
		formId, _ := strconv.Atoi(record.FormID)
		formName := e.GetFormNameById(formId)
		sro, ref := record.GetSafeSroRef()
		recordPath := filepath.Join(path, sro, ref)
		err := ensure(recordPath)
		if err != nil {
			return err
		}
		err = record.WriteAllPictures(recordPath, e.KizeoContext, parallel)
		if err != nil {
			fmt.Printf("\tWarning: WriteAllImage failed on '%s'.'%s': %s\n", formName, record.SummarySubtitle, err.Error())
		}
		err = record.WriteComment(recordPath)
		if err != nil {
			fmt.Printf("\tWarning: WriteComment failed on '%s'.'%s': %s\n", formName, record.SummarySubtitle, err.Error())
		}
	}
	return nil
}

func ensure(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("%s: creating '%s'\n", time.Now().Format("2006-01-02 15:04:05.0"), dir)
		err := os.MkdirAll(dir, os.ModeDir)
		if err != nil {
			return err
		}
	}
	return nil
}
