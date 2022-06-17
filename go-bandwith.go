package main

import (
    "fmt"
    log "github.com/sirupsen/logrus"
    "os"
    "time"
    "flag"
    "infiniband_test/method"
    "path/filepath"
)

func main() {
    
    //Flag Declare
    var dFlag bool
    var pFlag string
    var cFlag bool
    
    //Flag Assignment
    flag.BoolVar(&dFlag, "d", false, "DEBUG MODE")
    flag.StringVar(&pFlag, "p", "./hosts.list", "Host list to test")
    flag.BoolVar(&cFlag, "c", false, "Clean All")

    flag.Parse()
    
    //Log File...............................................
    begin := time.Now()
    if dFlag {
        log.SetLevel(log.DebugLevel)
        fmt.Println("Debug Level active, see logs.txt for details")
    } else {
        log.SetLevel(log.ErrorLevel)
    }

    log_file, err := os.OpenFile("./Bandwith " + begin.Format(time.RFC850) + ".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatal(err)
    }
    log.SetOutput(log_file)
    log.Info("BEGIN")

    //Clean Folder
    if cFlag {
        files, err := filepath.Glob("./Bandwith*")
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

    //Read File
    lines, err := method.ReadLines(pFlag)
    if err != nil {
        log.Fatal("readLines: %s", err)
    }
    
    //Shuffle and print host list shuffled 
    shuffled := method.ShuffleSlice(lines)
    method.WriteLines(shuffled, "./hosts.shuffled")

    err, out, errout := method.Shellout("cat ./hosts.shuffled")
    if err != nil {
        log.Printf("error: %v", err)
        log.Error(errout)
    }
    if len(errout) != 0 {
         fmt.Println("Ops, something went wrong, check logs.txt for Info")
    }
    fmt.Println(out)
}

