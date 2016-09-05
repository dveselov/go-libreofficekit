package libreofficekit

/*
#cgo CFLAGS: -I ./ -D LOK_USE_UNSTABLE_API
#cgo LDFLAGS: -ldl
#include <lokbridge.h>
*/
import "C"
import (
	"fmt"
	"image"
	"strconv"
	"strings"
	"sync"
	"unsafe"
)

// TwipsToPixels converts given twips to pixels with given dpi & zoom
func TwipsToPixels(twips int, dpi int) int {
	return int(float32(twips) / 1440.0 * float32(dpi))
}

func PixelsToTwips(pixels int, dpi int) int {
	return int((float32(pixels) / float32(dpi)) * 1440.0)
}

type Office struct {
	handle *C.struct__LibreOfficeKit
	mutex  *sync.Mutex
}

func NewOffice(path string) (*Office, error) {
	office := new(Office)

	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	lokit := C.lok_init(cPath)
	if lokit == nil {
		return nil, fmt.Errorf("Failed to initialize LibreOfficeKit with path: '%s'", path)
	}

	office.handle = lokit
	office.mutex = &sync.Mutex{}

	return office, nil

}

func (office *Office) Close() {
	C.destroy_office(office.handle)
}

func (office *Office) GetError() string {
	message := C.get_error(office.handle)
	return C.GoString(message)
}

func (office *Office) LoadDocument(path string) (*Document, error) {
	document := new(Document)
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	handle := C.document_load(office.handle, cPath)
	if handle == nil {
		return nil, fmt.Errorf("Failed to load document")
	}
	document.handle = handle
	return document, nil
}

func (office *Office) GetFilters() string {
	filters := C.get_filter_types(office.handle)
	defer C.free(unsafe.Pointer(filters))
	return C.GoString(filters)
}

// Types of documents returned by Document.GetType function
const (
	TextDocument = iota
	SpreadsheetDocument
	PresentationDocument
	DrawingDocument
	OtherDocument
)

const (
	RGBATilemode = iota
	BGRATilemode
)

type Document struct {
	handle *C.struct__LibreOfficeKitDocument
}

// Close destroys document
func (document *Document) Close() {
	C.destroy_document(document.handle)
}

// GetType returns type of loaded document
func (document *Document) GetType() int {
	return int(C.get_document_type(document.handle))
}

// GetParts returns count of slides (for presentations) or pages (for text documents)
func (document *Document) GetParts() int {
	return int(C.get_document_parts(document.handle))
}

// GetPart returns current part of document, e.g.
// if document was just loaded it's current part will be 0
func (document *Document) GetPart() int {
	return int(C.get_document_part(document.handle))
}

// SetPart updates current part of document
func (document *Document) SetPart(part int) {
	C.set_document_part(document.handle, C.int(part))
}

// GetPartName returns current slide title (for presentations) or page title (for text documents)
func (document *Document) GetPartName(part int) string {
	cPart := C.int(part)
	cPartName := C.get_document_part_name(document.handle, cPart)
	defer C.free(unsafe.Pointer(cPartName))
	return C.GoString(cPartName)
}

// GetSize returns width and height of document in twips (1 Twip = 1/1440th of an inch)
// You can convert twips to pixels by this formula: (width or height) * (1.0 / 1440.0) * DPI
func (document *Document) GetSize() (int, int) {
	width := C.long(0)
	heigth := C.long(0)
	C.get_document_size(document.handle, &width, &heigth)
	return int(width), int(heigth)
}

// InitializeForRendering must be called before performing any rendering-related actions
func (document *Document) InitializeForRendering(arguments string) {
	cArguments := C.CString(arguments)
	defer C.free(unsafe.Pointer(cArguments))
	C.initialize_for_rendering(document.handle, cArguments)
}

// SaveAs saves document at desired path in desired format with applied filter rules
// Actual (from libreoffice) error message can be read with Office.GetError
func (document *Document) SaveAs(path string, format string, filter string) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	cFormat := C.CString(format)
	defer C.free(unsafe.Pointer(cFormat))
	cFilter := C.CString(filter)
	defer C.free(unsafe.Pointer(cFilter))
	status := C.document_save(document.handle, cPath, cFormat, cFilter)
	if status != 0 {
		return fmt.Errorf("Failed to save document")
	}
	return nil
}

// CreateView return id if newly created view
func (document *Document) CreateView() int {
	return int(C.create_view(document.handle))
}

// GetView returns current document view id
func (document *Document) GetView() int {
	return int(C.get_view(document.handle))
}

// GetViews returns total number of views in document
func (document *Document) GetViews() int {
	return int(C.get_views(document.handle))
}

func (document *Document) GetTileMode() int {
	return int(C.get_tile_mode(document.handle))
}

func (document *Document) GetPartPageRectangles() []image.Rectangle {
	var rectangles []image.Rectangle
	rawRectangles := C.GoString(C.get_part_page_rectangles(document.handle))
	pageRectangles := strings.Split(rawRectangles, ";")
	for _, points := range pageRectangles {
		var intPoints []int
		strPoints := strings.Split(points, ",")
		for _, point := range strPoints {
			point = strings.Trim(point, " ")
			i, _ := strconv.Atoi(point)
			intPoints = append(intPoints, i)
		}
		x0, y0, x1, y1 := intPoints[0], intPoints[1], intPoints[2], intPoints[3]
		rectangles = append(rectangles, image.Rect(x0, y0, x1, y1))
	}
	return rectangles
}

// PaintTile renders tile to given buf (which size must be a `4 * canvasWidth * canvasHeight`)
func (document *Document) PaintTile(buf []C.uchar, canvasWidth int, canvasHeight int, tilePosX int, tilePosY int, tileWidth int, tileHeight int) {
	C.paint_tile(
		document.handle,
		(*C.uchar)(unsafe.Pointer(&buf[0])),
		(C.int)(canvasWidth),
		(C.int)(canvasHeight),
		(C.int)(tilePosX),
		(C.int)(tilePosY),
		(C.int)(tileWidth),
		(C.int)(tileHeight),
	)
}
