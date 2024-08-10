package main

import (
    "os"
    "io"
    "strings"
    "slices"
    "fmt"
    "crypto/rand"
    "math/big"
    "encoding/json"
)

var H map[string][]string //Hash, output to JSON
var U []string //Unique; possibly unnecessary, but I'll keep using it for now
var D [][]string //Data; [word1, word2] ad nausueam

func dumpweights(X map[string][]string, O string) {
    F,e:=os.OpenFile(O, os.O_RDWR|os.O_CREATE, 0644); if e!=nil { panic(e) }
    E:=json.NewEncoder(F); E.Encode(X)
    e=F.Close(); if e!=nil { panic(e) }
}

func pullweights(I string) map[string][]string {
    var R map[string][]string
    F,e:=os.Open(I); if e!=nil { panic(e) }
    D:=json.NewDecoder(F)
    e=D.Decode(&R); if e!=nil { panic(e) }
    e=F.Close(); if e!=nil { panic(e) }
    return R
}

func run(S string) string {
    if !slices.Contains(U,S) { return "<abend>" }
    if len(H[S]) == 1 { return H[S][0] } // Otherwise rand.Int freaks
    n,e:=rand.Int(rand.Reader,big.NewInt(int64(len(H[S])-1))); if e!=nil{ panic(e) } // rand is very finicky
    return H[S][n.Int64()]
}

func main() {
    H = make(map[string][]string) // Init hash
    I:="testnew.txt"
    F,e:=os.Open(I); if e!=nil{ panic(e) }
    C,e:=io.ReadAll(F); if e!=nil{ panic(e) }
    var E string // Escrow
    // structural regexes could come in handy here, or bufio.ScanWords, but I /do/ want to buffer by lines
    for _,x := range strings.Split(string(C),"\n") {
        for i,y:= range strings.Split(x," ") {
            y=strings.ToLower(y) // Fix case
            if y == "" { continue } // Skip unused
            if i != len(strings.Split(x," ")) { // Build D
                if i != 0 { D=append(D,[]string{E,y}) } // Don't overlap sentences
                E = y 
            }
            if slices.Contains(U, y) { // Skip if already found
                continue
            } else {
                U=append(U,y)
                H[y]=[]string{"<abend>"}
            }
        }
    }
    for _,x := range D {
        k:=strings.ToLower(x[0]) // Key
        v:=strings.ToLower(x[1]) // Value
        //We can confidently assume that all of the words here will be in H
        if slices.Contains(H[k], v) { continue } // I /could/ take this out and have some natural weighting occur; maybe a cmdline option?
        H[k]=append(H[k],v)
    }
    Q := "so"
    fmt.Println(Q)
    for Q != "<abend>" {
        Q = run(Q)
        fmt.Println(Q)
    }
    // Example: dumpweights(H,"weightfile.json")
    // Example: fmt.Println(pullweights("weightfile.json"))
}
