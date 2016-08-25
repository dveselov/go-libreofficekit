package libreofficekit

// #cgo CFLAGS: -I /usr/include/LibreOfficeKit/
// #cgo LDFLAGS: -ldl
// #include <stdlib.h>
// #include "LibreOfficeKit/LibreOfficeKitInit.h"
// #include "LibreOfficeKit/LibreOfficeKit.h"
/*
typedef void (*voidFunc) ();
typedef int (*intFunc) ();
void destroy_bridge(voidFunc f, void* handle) {
	return f(handle);
};
LibreOfficeKitDocument* document_load_bridge(voidFunc f,
											 LibreOfficeKit* pThis,
											 const char* pURL) {
	f(pThis, pURL);
};
int document_save_bridge(intFunc f,
						LibreOfficeKitDocument* pThis,
						const char* pUrl,
						const char* pFormat,
						const char* pFilterOptions) {
	return f(pThis, pUrl, pFormat, pFilterOptions);
};
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type Office struct {
	handle *C.struct__LibreOfficeKit
}

func NewOffice(libreofficePath string) (*Office, error) {
	office := new(Office)

	path := C.CString(libreofficePath)
	defer C.free(unsafe.Pointer(path))

	lokit := C.lok_init(path)
	if lokit == nil {
		return nil, fmt.Errorf("Failed to initialize LibreOfficeKit with libreofficePath: '%s'", libreofficePath)
	}

	office.handle = lokit

	return office, nil

}

func (self *Office) Close() {
	selfPointer := unsafe.Pointer(self.handle)
	C.destroy_bridge(self.handle.pClass.destroy, selfPointer)
}

func (self *Office) LoadDocument(path string) *Document {
	document := new(Document)
	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))
	handle := C.document_load_bridge(self.handle.pClass.documentLoad, self.handle, c_path)
	document.handle = handle
	return document
}

type Document struct {
	handle *C.struct__LibreOfficeKitDocument
}

func (self *Document) Close() {
	selfPointer := unsafe.Pointer(self.handle)
	C.destroy_bridge(self.handle.pClass.destroy, selfPointer)
}

func (self *Document) SaveAs(path string, format string, filter string) error {
	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))
	c_format := C.CString(format)
	defer C.free(unsafe.Pointer(c_format))
	c_filter := C.CString(filter)
	defer C.free(unsafe.Pointer(c_filter))
	status := C.document_save_bridge(self.handle.pClass.saveAs, self.handle, c_path, c_format, c_filter)
	if status != 0 {
		return fmt.Errorf("Failed to save document")
	} else {
		return nil
	}
}
