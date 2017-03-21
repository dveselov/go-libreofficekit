package libreofficekit

import (
	"testing"
)

const (
	DefaultLibreOfficePath  = "/usr/lib/libreoffice/program/"
	DocumentThatDoesntExist = "testdata/kittens.docx"
	SampleDocument          = "testdata/sample.docx"
)

func TestInvalidOfficePath(t *testing.T) {
	_, err := NewOffice("/etc/passwd")
	if err == nil {
		t.Fail()
	}
}

func TestValidOfficePath(t *testing.T) {
	_, err := NewOffice(DefaultLibreOfficePath)
	if err != nil {
		t.Fail()
	}
}

func TestGetOfficeErrorMessage(t *testing.T) {
	office, _ := NewOffice(DefaultLibreOfficePath)
	office.LoadDocument(DocumentThatDoesntExist)
	message := office.GetError()
	if len(message) == 0 {
		t.Fail()
	}
}

func TestLoadDocumentThatDoesntExist(t *testing.T) {
	office, _ := NewOffice(DefaultLibreOfficePath)
	_, err := office.LoadDocument(DocumentThatDoesntExist)
	if err == nil {
		t.Fail()
	}
}

func TestSuccessLoadDocument(t *testing.T) {
	office, _ := NewOffice(DefaultLibreOfficePath)
	_, err := office.LoadDocument(SampleDocument)
	if err != nil {
		t.Fail()
	}
}

func TestGetPartPageRectangles(t *testing.T) {
	office, _ := NewOffice(DefaultLibreOfficePath)
	document, _ := office.LoadDocument(SampleDocument)
	rectangles := document.GetPartPageRectangles()
	if len(rectangles) != 2 {
		t.Fail()
	}
}

func TestGetParts(t *testing.T) {
	office, _ := NewOffice(DefaultLibreOfficePath)
	document, _ := office.LoadDocument(SampleDocument)
	parts := document.GetParts()
	if parts != 2 {
		t.Fail()
	}
}

func TestGetTileMode(t *testing.T) {
	office, _ := NewOffice(DefaultLibreOfficePath)
	document, _ := office.LoadDocument(SampleDocument)
	mode := document.GetTileMode()
	if mode != RGBATilemode && mode != BGRATilemode {
		t.Fail()
	}
}

func TestGetViewsCount(t *testing.T) {
	office, _ := NewOffice(DefaultLibreOfficePath)
	document, _ := office.LoadDocument(SampleDocument)
	views := document.GetViewsCount()
	if views != 1 {
		t.Fail()
	}
}

func TestGetType(t *testing.T) {
	office, _ := NewOffice(DefaultLibreOfficePath)
	document, _ := office.LoadDocument(SampleDocument)
	documentType := document.GetType()
	if documentType != TextDocument {
		t.Fail()
	}
}

func TestTextSelection(t *testing.T) {
	office, _ := NewOffice(DefaultLibreOfficePath)
	document, _ := office.LoadDocument(SampleDocument)
	rectangle := document.GetPartPageRectangles()[0]
	document.SetTextSelection(SetGraphicSelectionStart, rectangle.Min.X, rectangle.Min.Y)
	document.SetTextSelection(SetGraphicSelectionEnd, rectangle.Max.X, rectangle.Max.Y)
	plaintext := document.GetTextSelection("text/plain;charset=utf-8")
	if len(plaintext) < 1000 {
		t.Fail()
	}
	document.ResetTextSelection()
	plaintext = document.GetTextSelection("text/plain;charset=utf-8")
	if len(plaintext) != 0 {
		t.Fail()
	}
}
