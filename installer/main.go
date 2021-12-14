package main

import (
	"errors"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"
)
var usingWind bool = false
var failed = false
func main(){
	const runOs = runtime.GOOS
	const runArch = runtime.GOARCH
	url := "https://github.com/KoperStudio/KoperManager/releases/download/release/koper_manager_"
	var saveTo string

	if runOs == "windows"{
		usingWind = true
		url += "windows_" + runArch + ".exe"
		saveTo = "C:\\Windows\\koper_manager.exe"
	} else if runOs == "darwin" {
		url += "darwin_amd64"
		log.Println("Using unix-like pathing")
		saveTo = "/usr/local/bin/koper_manager"
	} else {
		url += "linux" + "_" + runArch
		log.Println("Using unix-like pathing")
		saveTo = "/usr/bin/koper_manager"
	}
	DownloadFile(url, saveTo)
	if failed {
		log.Println("Installation failed! Press 'Enter' to continue")
	} else {
		log.Println("Successfully installed! Press 'Enter' to continue")
	}

	fmt.Scan()
}

func DownloadFile(url string, dest string) {
	file := path.Base(url)
	start := time.Now()
	err := os.RemoveAll(dest)
	if errors.Is(err, os.ErrExist) {
		log.Println("Found old version, deleting")
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln("Server respond with error", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36")
	resp, _ := http.DefaultClient.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln("Unable to close response body")
		}
	}(resp.Body)
	if resp == nil {
		log.Fatalln("Unable to access", url)
		return
	}

	f, _ := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0644)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil && !failed {
			failed = true
			log.Fatalln("Unable to close file")
		}
	}(f)

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading " + file,
	)

	_, err = io.Copy(io.MultiWriter(f, bar), resp.Body)

	if err != nil {
		if usingWind {
			failed = true
			log.Println("You need to run this program as administrator")
		}
		return
	}

	elapsed := time.Since(start)
	log.Printf("Download completed in %s", elapsed)
}