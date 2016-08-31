package libreofficekit

/*
#cgo CFLAGS: -I ./ -D LOK_USE_UNSTABLE_API
#cgo LDFLAGS: -ldl
#include <lokbridge.h>
*/
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

type Office struct {
	handle *C.struct__LibreOfficeKit
	mutex  *sync.Mutex
}

func NewOffice(path string) (*Office, error) {
	office := new(Office)

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	lokit := C.lok_init(c_path)
	if lokit == nil {
		return nil, fmt.Errorf("Failed to initialize LibreOfficeKit with path: '%s'", path)
	}

	office.handle = lokit
	office.mutex = &sync.Mutex{}

	return office, nil

}

func (self *Office) Close() {
	C.destroy_office(self.handle)
}

func (self *Office) GetError() string {
	message := C.get_error(self.handle)
	return C.GoString(message)
}

func (self *Office) LoadDocument(path string) (*Document, error) {
	document := new(Document)
	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))
	handle := C.document_load(self.handle, c_path)
	if handle == nil {
		return nil, fmt.Errorf("Failed to load document")
	}
	document.handle = handle
	return document, nil
}

// Types of documents returned by Document.GetType function
const (
	TextDocument = iota
	SpreadsheetDocument
	PresentationDocument
	DrawingDocument
	OtherDocument
)

type Document struct {
	handle *C.struct__LibreOfficeKitDocument
}

func (self *Document) Close() {
	C.destroy_document(self.handle)
}

// Returns type of loaded document
func (self *Document) GetType() int {
	return int(C.get_document_type(self.handle))
}

// Returns count of slides (for presentations) or pages (for text documents)
func (self *Document) GetParts() int {
	return int(C.get_document_parts(self.handle))
}

// GetPart returns current part of document, e.g.
// if document was just loaded it's current part will be 0
func (self *Document) GetPart() int {
	return int(C.get_document_part(self.handle))
}

// SetPart updates current part of document
func (self *Document) SetPart(part int) {
	C.set_document_part(self.handle, C.int(part))
}

// Returns current slide title (for presentations) or page title (for text documents)
func (self *Document) GetPartName(part int) string {
	c_part := C.int(part)
	c_part_name := C.get_document_part_name(self.handle, c_part)
	defer C.free(unsafe.Pointer(c_part_name))
	return C.GoString(c_part_name)
}

// GetSize returns width and height of document in twips (1 Twip = 1/1440th of an inch)
// You can convert twips to pixels by this formula: (width or height) * (1.0 / 1440.0) * DPI
func (self *Document) GetSize() (int, int) {
	width := C.long(0)
	heigth := C.long(0)
	C.get_document_size(self.handle, &width, &heigth)
	return int(width), int(heigth)
}

// Must be called before performing any rendering-related actions
func (self *Document) InitializeForRendering(arguments string) {
	c_arguments := C.CString(arguments)
	defer C.free(unsafe.Pointer(c_arguments))
	C.initialize_for_rendering(self.handle, c_arguments)
}

// Saves document at desired path in desired format with applied filter rules
// Actual (from libreoffice) error message can be read with Office.GetError
func (self *Document) SaveAs(path string, format string, filter string) error {
	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))
	c_format := C.CString(format)
	defer C.free(unsafe.Pointer(c_format))
	c_filter := C.CString(filter)
	defer C.free(unsafe.Pointer(c_filter))
	status := C.document_save(self.handle, c_path, c_format, c_filter)
	if status != 0 {
		return fmt.Errorf("Failed to save document")
	} else {
		return nil
	}
}
