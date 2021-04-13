package api

import (
	"encoding/json"
	"fmt"
	"github.com/disintegration/imaging"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type SearchResult struct {
	RecordsTotal    int           `json:"recordsTotal"`
	RecordsFiltered int           `json:"recordsFiltered"`
	Data            []*SearchData `json:"data"`
	Status          string        `json:"status"`
}

type SearchData struct {
	ID               string `json:"_id"`
	RecordNumber     string `json:"_record_number"`
	FormID           string `json:"_form_id"`
	UserID           string `json:"_user_id"`
	UpdateUserID     string `json:"_update_user_id"`
	CreateTime       string `json:"_create_time"`
	UpdateTime       string `json:"_update_time"`
	AnswerTime       string `json:"_answer_time"`
	UpdateAnswerTime string `json:"_update_answer_time"`
	History          string `json:"_history"`
	FormUniqueID     string `json:"_form_unique_id"`
	UserName         string `json:"_user_name"`
	UpdateUserName   string `json:"_update_user_name"`
	SummaryTitle     string `json:"_summary_title"`
	SummarySubtitle  string `json:"_summary_subtitle"`
	Comment          string `json:"commentaire"`
	Geoloc           struct {
		Lat  string `json:"lat"`
		Long string `json:"long"`
	} `json:"_geoloc"`
	Pictures    map[string]string
	ExtractData bool
}

func (sd *SearchData) UnmarshalJSON(data []byte) error {
	var dat map[string]interface{}
	if err := json.Unmarshal(data, &dat); err != nil {
		return err
	}

	picts := make(map[string][]string)

	getImageUIDs := func(chapt, val string) {
		if val == "" {
			return
		}
		rawUids := strings.Split(val, ",")
		res := make([]string, len(rawUids))
		for i, rawUid := range rawUids {
			res[i] = strings.TrimSuffix(rawUid, ".jpg")
		}
		picts[chapt] = append(picts[chapt], res...)
	}

	sd.ID, _ = dat["_id"].(string)
	delete(dat, "_id")
	sd.RecordNumber, _ = dat["_record_number"].(string)
	delete(dat, "_record_number")
	sd.FormID, _ = dat["_form_id"].(string)
	delete(dat, "_form_id")
	sd.UserID, _ = dat["_user_id"].(string)
	delete(dat, "_user_id")
	sd.UpdateTime, _ = dat["_update_time"].(string)
	delete(dat, "_update_time")
	sd.History, _ = dat["_history"].(string)
	delete(dat, "_history")
	sd.FormUniqueID, _ = dat["_form_unique_id"].(string)
	delete(dat, "_form_unique_id")
	sd.SummarySubtitle, _ = dat["_summary_subtitle"].(string)
	delete(dat, "_summary_subtitle")
	sd.Comment, _ = dat["commentaire"].(string)
	delete(dat, "commentaire")

	for field, itf := range dat {
		switch field {
		case "numero_poteau", "reference_poteau":
			if sd.SummarySubtitle == "" {
				sd.SummarySubtitle, _ = itf.(string)
			}
		case "localisation_gps_poteau_latitude":
			sd.Geoloc.Lat, _ = itf.(string)
		case "localisation_gps_poteau_longitude":
			sd.Geoloc.Long, _ = itf.(string)
		case "pendant_commentaire_travaux":
			comment, _ := itf.(string)
			sd.Comment += comment
		// Pictures infos
		case "avant_vue_d_ensemble":
			uids, _ := itf.(string)
			getImageUIDs("A Avant Vue d Ensemble", uids)
		case "etiquettes_avant", "etiquettes_avant1":
			uids, _ := itf.(string)
			getImageUIDs("B Avant Etiquette", uids)
		case "tete_poteau_avant", "tete_poteau_avant1":
			uids, _ := itf.(string)
			getImageUIDs("C Avant Tete Poteau", uids)
		case "information_terrassement", "pendant_photos_trou":
			uids, _ := itf.(string)
			getImageUIDs("D Pendant Information Terrassement", uids)
		case "fond_de_trou":
			uids, _ := itf.(string)
			getImageUIDs("E Pendant Fond de Trou", uids)
		case "profondeur_trou":
			uids, _ := itf.(string)
			getImageUIDs("E Pendant Profondeur Trou", uids)
		case "etiquette", "apres_etiquette":
			uids, _ := itf.(string)
			getImageUIDs("F Apres Etiquette", uids)
		case "apres_pied_du_poteau", "pied_de_poteau_apres":
			uids, _ := itf.(string)
			getImageUIDs("G Apres Pied Poteau", uids)
		case "vue_d_ensemble_apres", "apres_vue_d_ensemble":
			uids, _ := itf.(string)
			getImageUIDs("H Apres Vue d Ensemble", uids)
		case "tete_poteau_apres", "apres_tete_du_poteau":
			uids, _ := itf.(string)
			getImageUIDs("I Apres Tete Poteau", uids)
		case "pied_poteau":
			uids, _ := itf.(string)
			switch sd.FormID {
			case "664879":
				getImageUIDs("A Avant Vue d Ensemble", uids)
			default:
				getImageUIDs("G Apres Pied Poteau", uids)
			}
		case "vue_d_ensemble":
			uids, _ := itf.(string)
			switch sd.FormID {
			case "664879":
				getImageUIDs("H Apres Vue d Ensemble", uids)
			default:
				getImageUIDs("A Avant Vue d Ensemble", uids)
			}
		}
	}

	sd.Pictures = make(map[string]string)
	for chapter, uuids := range picts {
		for i, uuid := range uuids {
			sd.Pictures[fmt.Sprintf("%s %d", chapter, i+1)] = uuid
		}
	}
	sd.ExtractData = true

	return nil
}

func (sd *SearchData) WriteComment(path string) error {
	if sd.Comment == "" {
		return nil
	}

	commentName := fmt.Sprintf("%s %s.txt", sd.GetSafeRef(), "Commentaire")
	commentFullName := filepath.Join(path, commentName)
	commentFile, err := os.Create(commentFullName)
	if err != nil {
		return fmt.Errorf("could not create comment file '%s': %s\n", commentFullName, err.Error())
	}
	defer commentFile.Close()
	_, err = fmt.Fprint(commentFile, sd.Comment)
	if err != nil {
		return fmt.Errorf("could not write to comment file '%s': %s\n", commentFullName, err.Error())
	}
	return nil
}

func (sd *SearchData) GetSafeRef() string {
	_, ref := sd.GetSroRef()
	return safeName(ref)
}

func safeName(name string) string {
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, ":", "")
	name = strings.ReplaceAll(name, "  ", " ")
	name = strings.Trim(name, " \t")
	return name
}

func (sd *SearchData) GetDateHour() (recDate, recHour string) {
	recDate, recHour = sd.UpdateTime[0:10], sd.UpdateTime[11:16]
	return
}

func (sd *SearchData) GetSroRef() (sro, ref string) {
	refs := strings.Split(sd.SummarySubtitle, "|")
	sro = sd.SummarySubtitle
	if len(refs) < 2 {
		return
	}
	sro = strings.Trim(refs[0], "  ")
	ref = strings.Trim(refs[1], "  ")
	return
}

func (sd *SearchData) GetSafeSroRef() (sro, ref string) {
	sro, ref = sd.GetSroRef()
	sro, ref = safeName(sro), safeName(ref)
	return
}

const (
	kizeoImgMaxSize int = 1600
	kizeoImgQuality int = 75
)

func (sd *SearchData) WriteAllPictures(path string, kc *KizeoContext, parallel int) error {
	workers := make(chan struct{}, parallel)
	defer func() {
		close(workers)
	}()
	//t := time.Now()
	for imgLabel, _ := range sd.Pictures {
		workers <- struct{}{}
		go func(dir, label string) {
			err := sd.WritePicture(dir, label, kc)
			if err != nil {
				fmt.Printf("%s\n", err.Error())
			}
			<-workers
		}(path, imgLabel)
	}

	// wait for goroutine completion
	for i := 0; i < parallel; i++ {
		workers <- struct{}{}
	}
	//fmt.Printf(" took %s\n", time.Since(t).String())
	return nil
}

func (sd *SearchData) WritePicture(dir, label string, kc *KizeoContext) error {
	uuid, found := sd.Pictures[label]
	if !found {
		return fmt.Errorf("could not find image with label '%s'\n", label)
	}
	//t := time.Now()
	rc, err := kc.FormDataPictureStream(sd.FormID, sd.ID, uuid)
	if err != nil {
		return err
	}
	defer rc.Close()

	imgName := fmt.Sprintf("%s %s.jpg", sd.GetSafeRef(), label)

	imgFullName := filepath.Join(dir, imgName)
	imgFile, err := os.Create(imgFullName)
	if err != nil {
		return fmt.Errorf("could not create image file '%s': %s\n", imgFullName, err.Error())
	}
	defer imgFile.Close()
	// process image
	err = processImage(rc, imgFile)
	if err != nil {
		_ = os.Remove(imgFullName)
		return fmt.Errorf("could not process image file '%s': %s\n", imgFullName, err.Error())
	}
	//fmt.Printf(".")
	//fmt.Printf("Wrote '%s' in %s\n", imgName, time.Since(t).String())
	return nil
}

// processImage ensure resize / sharpen and quality downgrade on given img body
func processImage(img io.Reader, target io.Writer) error {
	im, err := imaging.Decode(img, imaging.AutoOrientation(true))
	if err != nil {
		return fmt.Errorf("image decode returned: %s", err.Error())
	}
	// Resize image if too large
	if im.Bounds().Size().X > im.Bounds().Size().Y && im.Bounds().Size().X > kizeoImgMaxSize {
		im = imaging.Resize(im, kizeoImgMaxSize, 0, imaging.CatmullRom)
	}
	if im.Bounds().Size().Y > im.Bounds().Size().X && im.Bounds().Size().Y > kizeoImgMaxSize {
		im = imaging.Resize(im, 0, kizeoImgMaxSize, imaging.CatmullRom)
	}
	err = jpeg.Encode(target, im, &jpeg.Options{Quality: kizeoImgQuality})
	if err != nil {
		return fmt.Errorf("image jpeg encode returned: %s", err.Error())
	}
	return nil
}
