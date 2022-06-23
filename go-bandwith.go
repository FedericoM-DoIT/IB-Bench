package main

import (
    "fmt"
    "os"
    "time"
    "flag"
    "infiniband_test/method"
    "path/filepath"
    log "github.com/sirupsen/logrus"
    _ "embed"
)

//go:embed pairs_not_random_ib_test_1.4
var f []byte 

func main() {
    
    
    //Flag Declare
    var pFlag string
    var cFlag bool
    var sFlag string
    
    //Flag Assignment
    flag.StringVar(&pFlag, "p", "./hosts.list", "Host list to test")
    flag.BoolVar(&cFlag, "c", false, "Clean All")
    flag.StringVar(&sFlag, "s", "bdw", "Test selection, accepted value: lat - latr - iperf - ipoIB. Default value bdw")

    flag.Parse()
    //Log File...............................................
    begin := time.Now()


    //Clean Folder
    if cFlag {
        files, err := filepath.Glob("./[Bb]andwith*")
        if err != nil {
            log.Panic(err)
        }
        for _, f := range files {
            if err := os.Remove(f); err != nil {
                log.Panic(err)
            }
        }
    os.Exit(0)
    }

    _ = os.WriteFile("./pairs_not_random_ib_test_1.4", f, 0755)

    //Read File
    lines, err := method.ReadLines(pFlag)
    if err != nil {
        log.Fatal("readLines: %s", err)
    }
    
    //Shuffle and print host list shuffled 
    shuffled := method.ShuffleSlice(lines)
    method.WriteLines(shuffled, "./hosts.shuffled")

    //Conf File Creation
    confPath := method.ConfCreate(sFlag)

    //Launch Bench Command
    //
    //
    err, out, errout := method.Shellout("./pairs_not_random_ib_test_1.4  -l hosts.shuffled --param params-bdw.conf -v -np -st 8 -P 200  --filter '65536'  --reverse")
    if err != nil {
        log.Printf("error: %v", err)
        log.Error(errout)
    }
    if len(errout) != 0 {
         fmt.Println("Ops, something went wrong, check logs.txt for Info")
    }
    out_file, err := os.OpenFile("./bandwith " + begin.Format(time.RFC850) + ".out", os.O_CREATE|os.O_WRONLY, 0666)
    _, err2 := out_file.WriteString(out)
    if err2 != nil {
        log.Fatal(err2)
    }

    //Parse output
    method.ParseOutput(out_file.Name())
    fmt.Println("\n*************************************")
    fmt.Println("\nCheck complete, refer to .out for complete output. Fault node will be listed in .fault")
    fmt.Println("\n*************************************")

    //Cleaning conf file
    method.DeleteFile(confPath)
    method.DeleteFile("./hosts.shuffled")
    //method.DeleteFile("./pairs_not_random_ib_test_1.4")
}

