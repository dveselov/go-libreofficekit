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
