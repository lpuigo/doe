package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"strings"
)

const (
	authToken string = "perma_restv3_laurentpuig_2d0b22a927b9e3bdcc9d69855eadbd970a610082"
	urlV3     string = "https://www.kizeoforms.com/rest/v3/"
)

type KizeoContext struct {
	URL      string
	User     string
	Password string
	Company  string
	Auth     string
}

func NewKizeoContext() *KizeoContext {
	return &KizeoContext{
		URL:      urlV3,
		User:     "",
		Password: "",
		Company:  "",
		Auth:     authToken,
	}
	//return &KizeoContext{
	//	URL:      urlV3,
	//	User:     "laurentpuig",
	//	Password: "1sc0m1ng",
	//	Company:  "EWINSE",
	//	Auth:     authToken,
	//}
}

// Login calls login api if KizeoContext has no Auth token already set.
// If login call succeed, receiver KizeoContext is updated with retrieved Authorisation token
func (kc *KizeoContext) Login() error {
	if kc.Auth != "" {
		return nil
	}
	loginReqBody, err := json.Marshal(map[string]string{
		"user":     kc.User,
		"password": kc.Password,
		"company":  kc.Company,
	})
	if err != nil {
		return fmt.Errorf("could not marshal request body: %s", err.Error())
	}

	resp, err := http.Post(kc.URL+"login", "application/json", bytes.NewBuffer(loginReqBody))
	if err != nil {
		return fmt.Errorf("sending post failed: %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response has non ok HTTP status: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var loginResp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Token string `js:"token"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	if err != nil {
		return fmt.Errorf("decoding response body failed: %s\nRecieved Response: %v", err.Error(), resp.Body)
	}
	if loginResp.Status != "ok" {
		return fmt.Errorf("unexpected response: %+v", loginResp)
	}
	kc.Auth = loginResp.Data.Token
	return nil
}

func (kc *KizeoContext) addAuth(req *http.Request) {
	req.Header.Set("Authorization", kc.Auth)
}

func (kc *KizeoContext) Forms() ([]*FormMin, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", kc.URL+"forms", nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %s", err.Error())
	}
	kc.addAuth(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP call failed: %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response has non ok HTTP status: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var formsResp struct {
		Status string     `json:"status"`
		Forms  []*FormMin `json:"forms"`
	}
	err = json.NewDecoder(resp.Body).Decode(&formsResp)
	if err != nil {
		return nil, fmt.Errorf("decoding response body failed: %s\nRecieved Response: %v", err.Error(), resp.Body)
	}
	if formsResp.Status != "ok" {
		return nil, fmt.Errorf("unexpected response: %+v", formsResp)
	}
	return formsResp.Forms, nil
}

func (kc *KizeoContext) FormDatas(formId string) ([]*DataMin, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%sforms/%s/data", kc.URL, formId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %s", err.Error())
	}
	kc.addAuth(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP call failed: %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response has non ok HTTP status: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var formDataResp struct {
		Status string     `json:"status"`
		Data   []*DataMin `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&formDataResp)
	if err != nil {
		return nil, fmt.Errorf("decoding response body failed: %s\nRecieved Response: %v", err.Error(), resp.Body)
	}
	if formDataResp.Status != "ok" {
		return nil, fmt.Errorf("unexpected response: %+v", formDataResp)
	}
	return formDataResp.Data, nil
}

func (kc *KizeoContext) FormDatasSince(formId, date string) ([]*SearchData, error) {
	as := NewAdvancedSearch().
		SetFilters(NewAdvancedSearchFilter("update_time", ">", date)).
		SetOrder(NewAdvancedSearchOrder("id", false))
	return kc.FormDatasAdvanced(formId, as)
}

func (kc *KizeoContext) FormDatasAdvanced(formId string, advSearch *AdvancedSearch) ([]*SearchData, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%sforms/%s/data/advanced", kc.URL, formId)
	advBody, err := json.Marshal(advSearch)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request body: %s", err.Error())
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(advBody))
	if err != nil {
		return nil, fmt.Errorf("could not create request: %s", err.Error())
	}
	kc.addAuth(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP call failed: %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response has non ok HTTP status: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var searchResult SearchResult
	// Debug
	//buf := new(strings.Builder)
	//tee := io.TeeReader(resp.Body, buf)
	//fmt.Printf("%s", buf.String())
	//err = json.NewDecoder(tee).Decode(&searchResult)
	err = json.NewDecoder(resp.Body).Decode(&searchResult)
	if err != nil {
		return nil, fmt.Errorf("decoding response body failed: %s\nRecieved Response: %v", err.Error(), resp.Body)
	}
	//fmt.Printf("%s", buf.String())
	return searchResult.Data, nil
}

func (kc *KizeoContext) FormUnreadDatas(formId string) ([]*DataMin, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%sforms/%s/data/readnew", kc.URL, formId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %s", err.Error())
	}
	kc.addAuth(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP call failed: %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response has non ok HTTP status: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var formDataResp struct {
		Status string     `json:"status"`
		Data   []*DataMin `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&formDataResp)
	if err != nil {
		return nil, fmt.Errorf("decoding response body failed: %s\nRecieved Response: %v", err.Error(), resp.Body)
	}
	if formDataResp.Status != "ok" {
		return nil, fmt.Errorf("unexpected response: %+v", formDataResp)
	}
	return formDataResp.Data, nil
}

func (kc *KizeoContext) FormData(formId, dataId string) (*FormulaireData, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%sforms/%s/data/%s", kc.URL, formId, dataId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %s", err.Error())
	}
	kc.addAuth(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP call failed: %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response has non ok HTTP status: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var formDataResp struct {
		Status  string          `json:"status"`
		Message string          `json:"message"`
		Data    *FormulaireData `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&formDataResp)
	if err != nil {
		return nil, fmt.Errorf("decoding response body failed: %s\nRecieved Response: %v", err.Error(), resp.Body)
	}
	if formDataResp.Status != "ok" {
		return nil, fmt.Errorf("unexpected response: %+v", formDataResp)
	}
	return formDataResp.Data, nil
}

func (kc *KizeoContext) FormDataPicture(formId, dataId, picId string) (image.Image, string, error) {
	rc, err := kc.FormDataPictureStream(formId, dataId, picId)
	if err != nil {
		return nil, "", err
	}
	defer rc.Close()

	img, err := jpeg.Decode(rc)
	if err != nil {
		return nil, "", fmt.Errorf("decoding image failed : %s", err.Error())
	}
	return img, "", nil
}

func (kc *KizeoContext) FormDataPictureStream(formId, dataId, picId string) (io.ReadCloser, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%sforms/%s/data/%s/medias/%s", kc.URL, formId, dataId, picId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %s", err.Error())
	}
	kc.addAuth(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP call failed: %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response has non ok HTTP status: %d", resp.StatusCode)
	}

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "image") {
		return nil, fmt.Errorf("response is not an image")
	}
	return resp.Body, nil
}
