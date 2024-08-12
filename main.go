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
    "flag"
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

func unique(X map[string][]string) []string {
    var R []string
    for k,_ := range X { R = append(R,k) }
    return R
}

func run(S string, H map[string][]string) string {
    if !slices.Contains(U,S) { fmt.Println("foo"); fmt.Println(U); return "<abend>" }
    if len(H[S]) == 1 { return H[S][0] } // Otherwise rand.Int freaks
    n,e:=rand.Int(rand.Reader,big.NewInt(int64(len(H[S])-1))); if e!=nil{ panic(e) } // rand is very finicky
    return H[S][n.Int64()]
}

func help () {
    A := []string{"h","t","r","o","w","u","l"}
    fmt.Println("Quickbot: Markov bot training and usage in a small, fast package ")
    fmt.Println("Usage: qbot [-hul] [-t infile [-o outfile]] [-r [-w weightfile]]")
    for _,x := range A {
        fmt.Println("-" + x + ": " + flag.Lookup(x).Usage)
    }
    fmt.Println("If no options are provided, qbot will simply print this message and exit")
    os.Exit(0)
}

func train(I string, Up bool) {
    var E string // Escrow
    H = make(map[string][]string) // Init hash
    F,e:=os.Open(I); if e!=nil{ panic(e) }
    C,e:=io.ReadAll(F); if e!=nil{ panic(e) }
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
        if Up && slices.Contains(H[k], v) { continue } // Unweighted training
        H[k]=append(H[k],v)
    }
}

func interact(L int) {
    //Not doing the full thing yet
    Q := "in";
    fmt.Println(Q)
    for Q != "<abend>" {
        Q = run(Q,H)
        fmt.Println(Q)}
}

func main() {
    Ot:=flag.Bool("t",false,"Train only (outputting to \"weights.json\" if the output file is not specified)")
    Or:=flag.Bool("r",false,"Run only (using \"weights.json\" if the weight file is not specified)")
    Ow:=flag.String("w","","Specify weight file")
    Oo:=flag.String("o","","Specify output file")
    Ou:=flag.Bool("u",false,"Enable \"unweighted\" training; weights only contain distinct words, with no information about how often the occur")
    Ol:=flag.Int("l",0,"Specify minimum response length")
    Oh:=flag.Bool("h",false,"Display this message")
    flag.Parse()
    if *Oh==true { help() }
    if *Ot==true { if *Oo != "" { train(flag.Arg(0),*Ou); dumpweights(H,*Oo) } else { train(flag.Arg(0),*Ou); dumpweights(H,"weights.json") } 
        os.Exit(0)}
    if *Or==true { if *Ow != "" { H = pullweights(*Ow); U = unique(H); interact(*Ol) } else { H = pullweights("weightfile.json"); U = unique(H); interact(*Ol) } 
        os.Exit(0) }
    help()
}
