package cmd

import "github.com/kmaasrud/doctor"

func Build() error {
    doc, err := doctor.NewDocument()
    if err != nil {
        return err
    }

    err = doc.Build()
    if err != nil {
        return err
    }
    return nil
}
