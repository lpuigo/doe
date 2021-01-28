package api

import (
	"fmt"
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
		fmt.Printf("form %5d:\n%s\n", i, form)
	}
}

func TestKizeoContext_FormDatas(t *testing.T) {
	kc := NewKizeoContext()
	datas, err := kc.FormDatas("640312")
	if err != nil {
		t.Fatalf("Forms retured unexpected: %s", err.Error())
	}
	for i, data := range datas {
		t.Logf("form %d: %s", i, data)
	}
}

func TestKizeoContext_FormUnreadDatas(t *testing.T) {
	kc := NewKizeoContext()
	datas, err := kc.FormUnreadDatas("640312")
	if err != nil {
		t.Fatalf("Forms retured unexpected: %s", err.Error())
	}
	for i, data := range datas {
		t.Logf("form %d: %s", i, data)
	}
}

func TestKizeoContext_FormData(t *testing.T) {
	kc := NewKizeoContext()
	formData, err := kc.FormData("640312", "97775188")
	if err != nil {
		t.Fatalf("Forms retured unexpected: %s", err.Error())
	}

	t.Log(formData.String())
}

func TestKizeoContext_FormDataPicture(t *testing.T) {
	kc := NewKizeoContext()
	_, info, err := kc.FormDataPicture("640312", "97775188", "c55532f640312pu412860_20210127095039_433bf31c-c7f5-4327-8b99-166cc6e3626a")
	if err != nil {
		t.Fatalf("Forms retured unexpected: %s", err.Error())
	}
	t.Log(info)
}
