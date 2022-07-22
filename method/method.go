package method

import (
    "bufio"
    "io"
    "os"
    "fmt"
    "math/rand"
    "time"
    "bytes"
    "os/exec"
    log "github.com/sirupsen/logrus"
    "regexp"
    "strings"
    "strconv"
    "github.com/fatih/color"
)

/*
		NEW INFINIBAND METHOD
*/
func ConfCreate(sFlag string) string {
    var confPath string
    switch {
    case sFlag == "bdw":
        confPath = "./params-bdw.conf"
	conf, err := os.Create(confPath)
        if err != nil {
            log.Fatal(err)
        }
        defer conf.Close()
        _, err2 := conf.WriteString("SERVER=ib_write_bw  -s 65536 --report_gbits\nCLIENT=ib_write_bw  -s 65536 --report_gbits\n")
        if err2 != nil {
            log.Fatal(err2)
        }
    case sFlag == "lat":
        confPath = "./params-lat.conf"
        conf, err := os.Create(confPath)
        if err != nil {
            log.Fatal(err)
        }
        defer conf.Close()
        _, err2 := conf.WriteString("SERVER=ib_send_lat\nCLIENT=ib_send_lat\n")
        if err2 != nil {
            log.Fatal(err2)
        }
    case sFlag == "ipoIB":
        confPath = "./params-ipoIB.conf"
        conf, err := os.Create(confPath)
        if err != nil {
            log.Fatal(err)
        }
        defer conf.Close()
        _, err2 := conf.WriteString("SERVER=iperf3 -s -1\nCLIENT=sleep 1 ; iperf3 -c\n")
        if err2 != nil {
            log.Fatal(err2)
        }
    }
    return confPath
}


//ReadLine To Slice
func ReadLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
} 

//Write slice to File
func WriteLines(lines []string, path string) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    w := bufio.NewWriter(file)
    for _, line := range lines {
        fmt.Fprintln(w, line)
    }
    return w.Flush()
}

//Mix Slice
func ShuffleSlice(arr []string) []string {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	return arr
}

//Shell Command
const ShellToUse = "bash"

func Shellout(command string) (error, string, string) {
    var stdout bytes.Buffer
    var stderr bytes.Buffer
    cmd := exec.Command(ShellToUse, "-c", command)
    cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
    cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)
    err := cmd.Run()
    return err, stdout.String(), stderr.String()
}

func IsError(err error) bool {
    if err != nil {
        log.Error(err.Error())
    }

    return (err != nil)
}

//Delete file func
func DeleteFile(path string) {
    var err = os.Remove(path)
    if IsError(err) {
        return
    }
}

//Parser
func ParseOutput(file string, sFlag string, greenFlag string, yellowFlag string) {
    minFlag, err := strconv.ParseFloat(yellowFlag, 64)
    maxFlag, err := strconv.ParseFloat(greenFlag, 64)
    out_file, err := os.Open(file)
    defer out_file.Close()
    begin := time.Now()
    var iterationsFlag string

    if (sFlag == "bdw") {
	iterationsFlag = "5000"
    }
    if (sFlag == "lat") {
	iterationsFlag = "1000"
    }	    
    
    if err != nil {
        log.Fatal(err)
    }
    outReader := bufio.NewReader(out_file)
    lineparseloop:
    for {
        var buffer bytes.Buffer
        var l []byte
        var isPrefix bool
        for {
            l, isPrefix, err = outReader.ReadLine()
            buffer.Write(l)
            // If we've reached the end of the line, stop reading.
            if !isPrefix {
                break
            }
            // If we're just at the EOF, break
            if err != nil {
                break
            }
        }
        if err == io.EOF {
            break 
        }
        line := buffer.String()
	
	matchIterFlag,_ := regexp.MatchString(iterationsFlag, line)
         
        if matchIterFlag {
	    switch sFlag {
	    case "bdw":
	        fault_file, err := os.OpenFile("./bandwith_" + begin.Format(time.RFC850) + ".fault", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	        lineField := strings.Fields(line)
	        peakValue, err := strconv.ParseFloat(lineField[5], 64)
	        if err != nil {
		    log.Error(err)
		    fmt.Println("Something wrong with: " + lineField[0] + lineField[1] + lineField[2] + " interaction")
                    continue lineparseloop
                }
	        averageValue, err := strconv.ParseFloat(lineField[6], 64)
	        if err != nil {
                    log.Error(err)
		    fmt.Println("Something wrong with: " + lineField[0] + lineField[1] + lineField[2] + " interaction")
		    continue lineparseloop
                }
	        if (peakValue >= minFlag) && (peakValue <= maxFlag) || (averageValue >= minFlag) && (averageValue <= maxFlag) {
	            color.Yellow(line)
		    fault_file.WriteString("Warning: " + lineField[0] + " " + lineField[1]+ " " + lineField[2] + "    BW_peak[Gb/sec]: " + lineField[5] + "    -    BW_average[Gb/sec]: " + lineField[6] + "\n")
	        }	 
                if (peakValue < minFlag) || (averageValue < minFlag) {
		    color.Red(line)
		    fault_file.WriteString("Error:   " + lineField[0] + " " + lineField[1]+ " " + lineField[2] + "    BW_peak[Gb/sec]: " + lineField[5] + "    -    BW_average[Gb/sec]: " + lineField[6] + "\n")
	        }
	    case "lat":
                fault_file, err := os.OpenFile("./latency_" + begin.Format(time.RFC850) + ".fault", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	        lineField := strings.Fields(line)
                t_min, err := strconv.ParseFloat(lineField[5], 64)
                if err != nil {
                    log.Error(err)
                    fmt.Println("Something wrong with: " + lineField[0] + lineField[1] + lineField[2] + " interaction")
                    continue lineparseloop
                }
                t_avg, err := strconv.ParseFloat(lineField[8], 64)
                if err != nil {
                    log.Error(err)
                    fmt.Println("Something wrong with: " + lineField[0] + lineField[1] + lineField[2] + " interaction")
                    continue lineparseloop
                }
                if (t_min >= minFlag) && (t_min <= maxFlag) || (t_avg >= minFlag) && (t_avg <= maxFlag) {
                    color.Yellow(line)
                    fault_file.WriteString("Warning: " + lineField[0] + " " + lineField[1]+ " " + lineField[2] + "    t_min[usec]: " + lineField[5] + "    -    t_avg[usec]: " + lineField[8] + "\n")
                }
                if (t_min > maxFlag) || (t_avg > maxFlag) {
                    color.Red(line)
                    fault_file.WriteString("Error:   " + lineField[0] + " " + lineField[1]+ " " + lineField[2] + "    t_min[usec]: " + lineField[5] + "    -    t_avg[usec]: " + lineField[8] + "\n")
                }
	    }
	}
    }	
}
