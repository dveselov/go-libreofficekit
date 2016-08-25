# go-libreofficekit
Cgo bindings to LibreOfficeKit

# Usage

```go
package main

import "github.com/docsbox/go-libreofficekit"

func main() {
    office := libreofficekit.NewOffice("/path/to/libreoffice")
    document := office.LoadDocument("kittens.docx")
    err := document.SaveAs("kittens.pdf", "pdf", "skipImages")
}

```
