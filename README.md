# fdep

[![GoDoc](https://godoc.org/github.com/RangelReale/fdep?status.svg)](https://godoc.org/github.com/RangelReale/fdep)

Package for building relationships between proto files and extracting types, helping creating source code generators.

### install

    go get -u -v github.com/RangelReale/fdep

### usage

	package main

	import (
	    "fmt"
        "log"

        "github.com/RangelReale/fdep"
	)

	func main() {
	    parsedep := fdep.NewDep()
	    
        err := parsedep.AddPath("/protoc/include", fdep.DepType_Imported)
        if err != nil {
            log.Fatal(err)
        }
	
        err = parsedep.AddPath("/mysource/proto", fdep.DepType_Own)
        if err != nil {
            log.Fatal(err)
        }

        gftype, err := parsedep.GetType("google.protobuf.Empty")
        if err != nil {
            log.Fatal(err)
        }

        fmt.Printf("Empty type ALIAS=%s NAME=%s in FILE=%s\n", gftype.Alias, gftype.Name, gftype.DepFile.FilePath)

        ftype, err := parsedep.Files["application/user.proto"].GetType("User")
        if err != nil {
            log.Fatal(err)
        }

        fmt.Printf("User type ALIAS=%s NAME=%s in FILE=%s\n", ftype.Alias, ftype.Name, ftype.DepFile.FilePath)
	}

### author

Rangel Reale (rangelspam@gmail.com)
