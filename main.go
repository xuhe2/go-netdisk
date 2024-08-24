package main

import (
	"flag"
	"log"
	"os"

	"github.com/xuhe2/go-netdisk/file"
	"github.com/xuhe2/go-netdisk/setting"
)

func main() {
	programSetting := setting.ProgramSetting{}
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

	// get data file name
	if len(os.Args) <= 1 {
		log.Fatalf("please input data file name")
	}
	dataFileName := os.Args[1]
	log.Printf("dataFileName: %s", dataFileName)

	dataFile := file.File{}
	dataFile.Open(dataFileName)

	if err := dataFile.Encrypt([]byte(programSetting.Key)); err != nil {
		log.Fatalf("encrypt error: %v", err)
	}
	log.Printf("dataFile content after encrypt: %v", string(dataFile.Data))

	if err := dataFile.Decrypt([]byte(programSetting.Key)); err != nil {
		log.Fatalf("decrypt error: %v", err)
	}
	log.Printf("dataFile content after decrypt: %v", string(dataFile.Data))
}
