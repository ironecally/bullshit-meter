package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var bullshitScout *regexp.Regexp
var fileCount float32
var pendingWork float32
var verbose bool
var ignoreList []string

func init() {
	ignoreList = []string{
		"bullshit-meter",
		"vendor",
		"swagger",
		"files",
	}

}

func main() {
	//flag
	directoryFlag := flag.String("dir", "", "directory to be checked, if its empty, it will check current directory")
	verbosePtr := flag.Bool("v", false, "show verbose")
	flag.Parse()

	if verbosePtr != nil {
		verbose = *verbosePtr
	}

	//regexp
	r, _ := regexp.Compile(`(\/\/\s*TODO)|(\/\/\s*FIXME)`)
	bullshitScout = r

	//begin
	fmt.Println("ಠ_ಠ Smells like bullshit..")
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if directoryFlag != nil && *directoryFlag != "" {
		if strings.HasPrefix(*directoryFlag, "/") {
			dir = *directoryFlag
		} else {
			dir = dir + "/" + *directoryFlag
		}
	}
	if verbose {
		log.Println("scanning:", dir)
	}

	readDir(dir)

	if pendingWork > 0.0 {
		fmt.Printf("TODO/FIX: %.0f\nTotal File:%.0f\n", pendingWork, fileCount)
		fmt.Printf("ಠ_ಠ Your repository is %.2f %% Bullshit\n", pendingWork/fileCount*100)
	} else {
		fmt.Printf("( ͒˃⌂˂ ͒) false alarm!\n")
	}
}

func isIgnored(a string) bool {
	for _, b := range ignoreList {
		if b == a {
			return true
		}
	}
	return false
}

func readDir(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, val := range files {

		if isIgnored(val.Name()) {
			continue
		}
		if !strings.HasPrefix(val.Name(), ".") && !val.IsDir() {
			fileCount++
			if getTODO(dir, val) {
				pendingWork++
			}
		} else if !strings.HasPrefix(val.Name(), ".") && val.IsDir() {

			readDir(dir + "/" + val.Name())
		}
	}
}

func getTODO(dir string, filedata os.FileInfo) bool {
	if verbose {
		fmt.Printf("checking file:%s/%s\n", dir, filedata.Name())
	}

	f, err := os.Open(dir + "/" + filedata.Name())
	if err != nil {
		log.Println("fail to read file", filedata.Name(), err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		if bullshitScout.MatchString(string(scanner.Text())) {
			return true
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("fail to read line", filedata.Name(), err)
	}
	return false
}
