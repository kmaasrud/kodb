package doctor

import (
    "fmt"
)

var logLevel int = 0

func main() {
    doc, err := NewDocument()
    if err != nil {
        fmt.Println(err.Error())
    }
    fmt.Printf("%s\n", doc)
}
