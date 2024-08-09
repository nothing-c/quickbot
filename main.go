package main


import (
    "os"
    "io"
    "strings"
)

var H [string][[]string]

func main() {
    I:="testdata-clean.txt"
    F:=os.Open(I)
    C,e:=io.ReadAll(F); if e!=nil{ panic(e) }
    // structural regexes could come in handy here
    for _,x := range strings.Split(C,"\n") {
        
    }
}
