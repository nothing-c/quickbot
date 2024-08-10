package main


import (
    "os"
    "io"
    "strings"
    "slices"
    "fmt"
)

var H map[string][]string //Hash
var U []string //Unique; possibly unnecessary, but I'll keep using it for now
var D [][]string //Data; [word1, word2] ad nausueam

func main() {
    H = make(map[string][]string) // Init hash
    I:="testnew.txt"
    F,e:=os.Open(I); if e!=nil{ panic(e) }
    C,e:=io.ReadAll(F); if e!=nil{ panic(e) }
    // structural regexes could come in handy here, or bufio.ScanWords
    for _,x := range strings.Split(string(C),"\n") {
        for i,y:= range strings.Split(x," ") {
            y=strings.ToLower(y) // Fix case
            if y == "" { continue } //skip unused
            if i != len(strings.Split(x," "))-1 { t:=strings.Split(x," ")[i:i+2]; D=append(D,t) }
            if slices.Contains(U, y) { //ditto
                continue
            } else {
                U=append(U,y)
                H[y]=[]string{"<abend>"}
            }
        }
    }

    for _,x := range D {
        //We can confidently assume that all of the words here will be in H
        H[strings.ToLower(x[0])]=append(H[strings.ToLower(x[0])],strings.ToLower(x[1]))
    }
//    fmt.Println(U)
    fmt.Println(H)
    fmt.Println(H["so"])
//    fmt.Println("~~~~"+U[1])
//    fmt.Println(D)
}
