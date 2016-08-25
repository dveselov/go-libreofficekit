# go-libreofficekit
Cgo bindings to LibreOfficeKit

# Pre-requirements

You must have modern version of LibreOffice (e.g. greater than 4.3)  
Also, install a LibreOfficeKit headers package (`libreofficekit-dev` in Ubuntu).  

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
