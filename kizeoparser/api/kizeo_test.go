package api

import (
	"testing"
)

func TestLogin(t *testing.T) {
	kc := NewKizeoContext()
	kc.Auth = ""
	err := kc.Login()
	if err != nil {
		t.Fatalf("Login returned unexpected: %s", err.Error())
	}
	t.Logf("Auth Token: %s", kc.Auth)
}

func TestKizeoContext_Forms(t *testing.T) {
	kc := NewKizeoContext()
	forms, err := kc.Forms()
	if err != nil {
		t.Fatalf("Forms retured unexpected: %s", err.Error())
	}
	for i, form := range forms {
		t.Logf("form %d: %+v", i, form)
	}
}

func TestKizeoContext_FormDatas(t *testing.T) {
	kc := NewKizeoContext()
	datas, err := kc.FormDatas("630190")
	if err != nil {
		t.Fatalf("Forms retured unexpected: %s", err.Error())
	}
	for i, data := range datas {
		t.Logf("form %d: %+v", i, data)
	}
}

func TestKizeoContext_FormData(t *testing.T) {
	kc := NewKizeoContext()
	formData, err := kc.FormData("630190", "95370695")
	if err != nil {
		t.Fatalf("Forms retured unexpected: %s", err.Error())
	}

	t.Log(formData.String())
}
