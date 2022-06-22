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
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
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
func ParseOutput(file string) {
    out_file, err := os.Open(file)
    defer out_file.Close()
    
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
	iterationsFlag := "5000"
	matchIterFlag,_ := regexp.MatchString(iterationsFlag, line)
         
        if matchIterFlag {
	    lineField := strings.Fields(line)
	    peakValue, err := strconv.ParseFloat(lineField[5], 64)
	    if err != nil {
		log.Error(err)    
                continue lineparseloop
            }
	    averageValue, err := strconv.ParseFloat(lineField[6], 64)
	    if err != nil {
                log.Error(err)
		continue lineparseloop
            }
	    if (peakValue >= 70) && (peakValue <= 90) || (averageValue >= 70) && (averageValue <= 90) {
	        color.Yellow(line)
	    }	 
            if (peakValue < 70) || (averageValue < 70) {
		color.Red(line)
	    }

	}
    }	
}
