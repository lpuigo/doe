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
	datas, err := kc.FormDatas("630190")
	if err != nil {
		t.Fatalf("Forms retured unexpected: %s", err.Error())
	}
	for i, data := range datas {
		t.Logf("form %d: %s", i, data)
	}
}

func TestKizeoContext_FormDatasSince(t *testing.T) {
	kc := NewKizeoContext()
	datas, err := kc.FormDatasSince("630190", "2021-04-09 13:58:11")
	if err != nil {
		t.Fatalf("Forms retured unexpected: %s", err.Error())
	}
	for i, data := range datas {
		t.Logf("form %d: %s", i, data)
	}
}

func TestKizeoContext_FormUnreadDatas(t *testing.T) {
	kc := NewKizeoContext()
	datas, err := kc.FormUnreadDatas("630190")
	if err != nil {
		t.Fatalf("Forms retured unexpected: %s", err.Error())
	}
	for i, data := range datas {
		t.Logf("form %d: %s", i, data)
	}
}

func TestKizeoContext_FormData(t *testing.T) {
	kc := NewKizeoContext()
	formData, err := kc.FormData("664879", "102550347")
	if err != nil {
		t.Fatalf("Forms retured unexpected: %s", err.Error())
	}

	t.Log(formData.String())
}

func TestKizeoContext_FormDataPicture(t *testing.T) {
	kc := NewKizeoContext()
	_, info, err := kc.FormDataPicture("664879", "102550347", "c55532f664879pu447378_20210408132938_c15ef482-abe7-4120-a27e-c23d95999ef2")
	if err != nil {
		t.Fatalf("Forms retured unexpected: %s", err.Error())
	}
	t.Log(info)
}
