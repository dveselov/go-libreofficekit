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

// TwipsToPixels converts given twips to pixels with given dpi
func TwipsToPixels(twips int, dpi int) int {
	return int(float32(twips) / 1440.0 * float32(dpi))
}

// PixelsToTwips is like TwipsToPixels, but to another way
func PixelsToTwips(pixels int, dpi int) int {
	return int((float32(pixels) / float32(dpi)) * 1440.0)
}

// BGRA converts BGRA array of pixels to RGBA
// https://github.com/golang/exp/blob/master/shiny/driver/internal/swizzle/swizzle_common.go#L13
func BGRA(p []uint8) {
	for i := 0; i < len(p); i += 4 {
		p[i+0], p[i+2] = p[i+2], p[i+0]
	}
}

type Office struct {
	handle *C.struct__LibreOfficeKit
	Mutex  *sync.Mutex
}

// NewOffice returns new Office or error if LibreOfficeKit fails to load
// required libs (actually, when libreofficekit-dev package isn't installed or path is invalid)
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

// Close destroys C LibreOfficeKit instance
func (office *Office) Close() {
	C.destroy_office(office.handle)
}

// GetError returns last happened error message in human-readable format
func (office *Office) GetError() string {
	message := C.get_error(office.handle)
	return C.GoString(message)
}

// LoadDocument return Document or error, if LibreOffice fails to open document at provided path.
// Actual error message can be retrieved by office.GetError method
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

// Types of tile color mode
const (
	RGBATilemode = iota
	BGRATilemode
)

const (
	SetGraphicSelectionStart = iota
	SetGraphicSelectionEnd
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

// GetTileMode returns tile mode of document, currently only RGBA or BGRA (5.2).
// You can compare returned int with RGBATilemode / BGRATilemode.
func (document *Document) GetTileMode() int {
	return int(C.get_tile_mode(document.handle))
}

func (document *Document) GetTextSelection(mimetype string) string {
	cMimetype := C.CString(mimetype)
	defer C.free(unsafe.Pointer(cMimetype))
	return C.GoString(C.get_text_selection(document.handle, cMimetype))
}

func (document *Document) SetTextSelection(sType int, x int, y int) {
	C.set_text_selection(
		document.handle,
		C.int(sType),
		C.int(x),
		C.int(y),
	)
}

func (document *Document) ResetTextSelection() {
	C.reset_selection(document.handle)
}

// GetPartPageRectangles array of image.Rectangle, with actually TextDocument page
// rectangles. Useful, when rendering text document page-by-page.
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
		x0, y0 := intPoints[0], intPoints[1]
		x1, y1 := x0+intPoints[2], y0+intPoints[3]
		rectangles = append(rectangles, image.Rect(x0, y0, x1, y1))
	}
	return rectangles
}

// PaintTile renders tile to given buf (which size must be a `4 * canvasWidth * canvasHeight`).
// In practice buf must be a pointer to image.Image.Pix array's first element, e.g. unsafe.Pointer(&image.Pix[0])
func (document *Document) PaintTile(buf unsafe.Pointer, canvasWidth int, canvasHeight int, tilePosX int, tilePosY int, tileWidth int, tileHeight int) {
	C.paint_tile(
		document.handle,
		(*C.uchar)(buf),
		(C.int)(canvasWidth),
		(C.int)(canvasHeight),
		(C.int)(tilePosX),
		(C.int)(tilePosY),
		(C.int)(tileWidth),
		(C.int)(tileHeight),
	)
}
