package tools

import "github.com/gopherjs/gopherjs/js"

type DomElement struct {
	*js.Object
}

func GetElementById(id string) *DomElement {
	el := js.Global.Get("document").Call("getElementById", id)
	return &DomElement{Object: el}
}

func (de *DomElement) SetInnerText(text string) {
	de.Set("innerText", text)
}

func (de *DomElement) SetInnerHTML(text string) {
	de.Set("innerHTML", text)
}

func (de *DomElement) SetValue(text string) {
	de.Set("value", text)
}

func (de *DomElement) Select() {
	de.Call("select")
	de.Call("setSelectionRange", 0, 99999) // For mobile devices
}

func CopyToClipBoard(text string) {
	document := js.Global.Get("document")
	textArea := document.Call("createElement", "textarea")
	styleTa := textArea.Get("style")
	styleTa.Set("position", "fixed")
	styleTa.Set("top", 0)
	styleTa.Set("left", 0)
	styleTa.Set("width", "2em")
	styleTa.Set("height", "2em")
	styleTa.Set("padding", 0)
	styleTa.Set("border", "none")
	styleTa.Set("outline", "none")
	styleTa.Set("boxShadow", "none")
	styleTa.Set("background", "transparent")
	textArea.Set("value", text)

	document.Get("body").Call("appendChild", textArea)
	textArea.Call("focus")
	textArea.Call("select")

	res := document.Call("execCommand", "copy")
	if !res.Bool() {
		print("DocumentCopy failed")
	}
	document.Get("body").Call("removeChild", textArea)
}
