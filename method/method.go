package method

import (
    "bufio"
    "os"
    "fmt"
    "math/rand"
    "time"
    "bytes"
    "os/exec"
    log "github.com/sirupsen/logrus"
)

/*
		NEW INFINIBAND METHOD
*/

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
