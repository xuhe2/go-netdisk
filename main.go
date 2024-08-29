package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/xuhe2/go-netdisk/file"
	"github.com/xuhe2/go-netdisk/setting"
)

var programSetting setting.ProgramSetting = setting.ProgramSetting{}

func main() {
	// read the args from command line
	key := flag.String("key", os.Getenv("KEY"), "key to unencrypt")
	flag.Parse()
	// read config.json
	settingFile, err := os.Open("config.json")
	if err == nil {
		defer settingFile.Close()
		programSetting.Parse(settingFile)
		log.Printf("programSetting: %v", programSetting)
	}
	// set key
	if *key != "" {
		programSetting.Key = *key
	}

	// get operation
	operation := os.Args[1]
	switch operation {
	case "push":
		push()
	case "pull":
		pull()
	default:
		log.Fatalf("operation %s not support", operation)
	}
}

func push() {
	// get data file name
	if len(os.Args) <= 1 {
		log.Fatalf("please input data file name")
	}
	dataFileName := os.Args[len(os.Args)-1]
	log.Printf("dataFileName: %s", dataFileName)

	reader, err := os.Open(dataFileName)
	if err != nil {
		log.Fatalf("open data file error: %v", err)
	}
	defer reader.Close()

	// open the file
	dataFile := file.File{Name: dataFileName}
	dataFile.Open(reader)

	if err := dataFile.Encrypt([]byte(programSetting.Key)); err != nil {
		log.Fatalf("encrypt file error: %v", err)
	}
	if err := dataFile.Save(); err != nil {
		log.Fatalf("save file error: %v", err)
	}
}

func pull() {
	// get data file name
	if len(os.Args) <= 1 {
		log.Fatalf("please input data file name")
	}
	path := os.Args[len(os.Args)-1]
	log.Printf("path: %s", path)

	dataFile := file.File{}

	// 如果是URL
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
	} else {
		//如果是文件路径
		if err := dataFile.Load(path); err != nil {
			log.Fatalf("load file error: %v", err)
		}
	}
}
