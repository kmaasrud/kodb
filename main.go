package main

import (
    "fmt"

    "github.com/kmaasrud/doctor/core"
)

func main() {
    doc, err := core.NewDocument()
    if err != nil {
        fmt.Println(err.Error())
    }
    fmt.Printf("%s\n", doc)
}
