# go-libreofficekit [![](https://godoc.org/github.com/docsbox/go-libreofficekit?status.svg)](https://godoc.org/github.com/docsbox/go-libreofficekit)
Cgo bindings to LibreOfficeKit

# Install
```bash 
# Latest version of LibreOffice (5.2) is required
$ sudo add-apt-repository ppa:libreoffice/ppa 
$ sudo apt-get update
$ sudo apt-get install libreoffice libreofficekit-dev
$ go get github.com/docsbox/go-libreofficekit
```

# Usage

This example demonstrates how to convert Microsoft Office document to PDF

```go
package main

import "github.com/docsbox/go-libreofficekit"

func main() {
    office, _ := libreofficekit.NewOffice("/path/to/libreoffice")
    
    document, _ := office.LoadDocument("kittens.docx")
    document.SaveAs("kittens.pdf", "pdf", "skipImages")

    document.Close()
    office.Close()
}

```

This example demonstrates how to get presentation slides titles

```go
package main

import "fmt"
import "github.com/docsbox/go-libreofficekit"

func main() {
    office, _ := libreofficekit.NewOffice("/path/to/libreoffice")
    
    document, _ := office.LoadDocument("kittens.pptx")
    slidesCount := document.GetParts()

    for i := 1; i < slidesCount; i++ {
        document.SetPart(i)
        currentPart = document.GetPart()
        fmt.Println("Current slide =", currentPart)
        currentPartName = document.GetPartName(i)
        fmt.Println("Current slide title =", currentPartName)
    }

    document.Close()
    office.Close()
}
```

Next example demonstrates how to use built-in LibreOffice rendering engine for creating page-by-page documents previews.

```go
package main

import (
    "os"
    "fmt"
    "unsafe"
    "image"
    "image/png"
)
import "github.com/docsbox/go-libreofficekit"

func main() {
    office, _ := libreofficekit.NewOffice("/path/to/libreoffice")
    document, _ := office.LoadDocument("kittens.docx")

    rectangles := document.GetPartPageRectangles()
    canvasWidth := libreofficekit.TwipsToPixels(rectangles[0].Dx(), 120)
    canvasHeight := libreofficekit.TwipsToPixels(rectangles[0].Dy(), 120)

    m := image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))

    for i, rectangle := range rectangles {
        document.PaintTile(unsafe.Pointer(&m.Pix[0]), canvasWidth, canvasHeight, rectangle.Min.X, rectangle.Min.Y, rectangle.Dx(), rectangle.Dy())
        libreofficekit.BGRA(m.Pix)
        out, _ := os.Create(fmt.Sprintf("page_%v.png", i))
        png.Encode(out, m)
        out.Close()
    }
}
```
