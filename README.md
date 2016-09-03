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

Next example demonstrates how to use built-in LibreOffice rendering (result will looks like [this](http://i.imgur.com/GozPaFc.png)).

```go
package main

import (
    "bufio"
    "os"
    "image"
    "image/png"
)
import "github.com/docsbox/go-libreofficekit"

func main() {
    office, _ := libreofficekit.NewOffice("/path/to/libreoffice")
    
    document, _ := office.LoadDocument("kittens.docx")
    // Get document width & height, in twips (1 twip = 1/1440 of inch)
    width, height := document.GetSize()
    // Convert document width/height to pixels with DPI = 100
    canvasWidth := libreofficekit.TwipsToPixels(width, 100)
    canvasHeight := libreofficekit.TwipsToPixels(height, 100)
    buf := make([]_Ctype_uchar, 4*canvasWidth*canvasHeight)
    document.PaintTile(buf, canvasWidth, canvasHeight, 0, 0, width, height)
    m := image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))
    pixels := make([]uint8, len(buf))
    for i := 0; i < len(buf); i++ {
        pixels[i] = (uint8)(buf[i])
    }
    m.Pix = pixels
    out, _ := os.Create("output.png")
    defer out.Close()
    outBuf := bufio.NewWriter(out)
    png.Encode(out, m)
    document.Close()
    office.Close()
}
```
