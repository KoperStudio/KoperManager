package main

import (
	"log"
	"runtime"
)


func main(){
	const runOs = runtime.GOOS
	const runArch = runtime.GOARCH
	log.Println("Detected OS:", runOs, "ARCH:", runArch)
	DownloadFile("https://google.com/", "aboba.txt")
}