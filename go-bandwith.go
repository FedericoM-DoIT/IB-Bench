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

//go:embed extbin/pairs_not_random_ib_test_1.4
var f []byte 

func main() {
    
    
    //Flag Declare
    var pFlag string
    var cFlag bool
    var sFlag string
    var greenFlag string
    var yellowFlag string
    
    //Flag Assignment
    flag.StringVar(&pFlag, "p", "./hosts.list", "Host list to test.")
    flag.BoolVar(&cFlag, "c", false, "Clean all bandiwith* latency* file.")
    flag.StringVar(&sFlag, "s", "bdw", "Test selection, accepted value: <bdw> - <lat>.")
    flag.StringVar(&greenFlag, "max", "90", "Max Valor for Output") 
    flag.StringVar(&yellowFlag, "min", "70", "Min Valor for Output") 


    flag.Usage = func() {
    fmt.Fprintf(os.Stderr, "Infiniband Benchmark Tool.\n\nAuthor(s): Federico Mollo federico.mollo@doit-systems.it\n\nIB-Bench version 1.0.1\n\nIB-Bench comes with ABSOLUTELY NO WARRANTY.  This is free software, and you\nare welcome to redistribute it under certain conditions.  See the GNU\nGeneral Public Licence for details.\n\nIB-Bench is a performance testing tool.\n\n")
    fmt.Fprintf(os.Stderr, "Usage:    IB-Bench  [OPTION].\nIf not prompted -p or -s took Default VALUE\n\n")
    flag.PrintDefaults()
    }
     
    flag.Parse()
    //Log File...............................................
    begin := time.Now()


    //Clean Folder
    if cFlag {
        bandFile, _ := filepath.Glob("./bandwith*")
        latFile, _ := filepath.Glob("./latency*")
        for _, f := range bandFile {
            if err := os.Remove(f); err != nil {
                log.Panic(err)
            }
	}
        for _, f := range latFile {
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
    var out_file *os.File
    switch sFlag { 
    case "lat":
        err, out, errout := method.Shellout("./pairs_not_random_ib_test_1.4  -l hosts.shuffled --param params-lat.conf -v -np -st 3 -P 200  --filter 'iterations|2       1000' --reverse")
        if err != nil {
            log.Printf("error: %v", err)
            log.Error(errout)
        }
        if len(errout) != 0 {
             fmt.Println("Ops, something went wrong, check logs.txt for Info")
        }
        out_file, err = os.OpenFile("./latency_" + begin.Format(time.RFC850) + ".out", os.O_CREATE|os.O_WRONLY, 0666)
        _, err2 := out_file.WriteString(out)
        if err2 != nil {
            log.Fatal(err2)
        }
    case "bdw":
	    err, out, errout := method.Shellout("./pairs_not_random_ib_test_1.4  -l hosts.shuffled --param params-bdw.conf -v -np -st 8 -P 200  --filter 'iterations|65536'  --reverse")
        if err != nil {
            log.Printf("error: %v", err)
            log.Error(errout)
        }
        if len(errout) != 0 {
             fmt.Println("Ops, something went wrong, check logs.txt for Info")
        }
        out_file, err = os.OpenFile("./bandwith_" + begin.Format(time.RFC850) + ".out", os.O_CREATE|os.O_WRONLY, 0666)
        _, err2 := out_file.WriteString(out)
        if err2 != nil {
            log.Fatal(err2)
        }
    }    
    
    fmt.Println("\n*******************************************************")
    fmt.Println("\n               Printing Critical Value                 ")
    fmt.Println("\n*******************************************************")
    
    //Parse output
    method.ParseOutput(out_file.Name(), sFlag, greenFlag, yellowFlag)
    
    
    fmt.Println("\n*************************************")
    fmt.Println("\nCheck complete, refer to .out for complete output. Fault node will be listed in .fault")
    fmt.Println("\n*************************************")

    //Cleaning conf file
    method.DeleteFile(confPath)
    method.DeleteFile("./hosts.shuffled")
    method.DeleteFile("./pairs_not_random_ib_test_1.4")
}

