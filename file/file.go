package file

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/xuhe2/go-netdisk/setting"
)

type File struct {
	Name         string `json:"file_name"`
	NumFileParts int    `json:"num_file_parts"`
	Key          []byte `json:"key"`
	FileParts    []*FilePart
}

// 打开文件, 同时切分大文件为file parts.
func (f *File) Open(r io.Reader) error {
	f.NumFileParts = 0
	f.FileParts = make([]*FilePart, 0)

	content := make([]byte, FilePartSize)
	// read the file data
	for {
		n, err := r.Read(content)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		// create a new file part
		f.FileParts = append(f.FileParts, NewFilePart(f.Name+strconv.Itoa(f.NumFileParts)+".part", content[:n]))
		f.NumFileParts++
	}
	return nil
}

// 从路径加载文件
func (f *File) Load(path string) error {
	// 去除末尾的`/`
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	// 找到路径下面后缀是`.info`的所有文件
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	var infoFile *os.File
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(), ".info") {
			// load the file info
			infoFile, err = os.Open(path + "/" + file.Name())
			if err != nil {
				return err
			}
			defer infoFile.Close()
			break
		}
	}
	if infoFile == nil {
		return fmt.Errorf("file info not found")
	}
	// load file info from infoFile
	filePartsInfo := setting.FileInfo{}
	if _, err := filePartsInfo.ReadFrom(infoFile); err != nil {
		return err
	}
	// load file from file parts
	f.LoadInfo(filePartsInfo)
	// decrypt file
	if err := f.Decrypt(f.Key); err != nil {
		return err
	}
	log.Printf(string(f.FileParts[0].Data))
	return nil
}

// from file info to set file
func (f *File) LoadInfo(info setting.FileInfo) error {
	var err error = nil

	f.Name = info.Name
	f.NumFileParts = info.NumFileParts
	f.Key = info.Key
	f.FileParts = make([]*FilePart, f.NumFileParts)
	wg := sync.WaitGroup{}
	for i := 0; i < f.NumFileParts; i++ {
		f.FileParts[i] = NewFilePart(info.FileParts[i], nil)
		wg.Add(1)
		go func() {
			defer wg.Done()
			var file *os.File
			file, err = os.Open(f.FileParts[i].Name)
			if err != nil {
				log.Println(err)
				return
			}
			f.FileParts[i].LoadData(file)
		}()
	}
	wg.Wait()
	return err
}

func (f *File) Encrypt(key []byte) error {
	// set the key
	f.Key = key
	// encrypt each file part
	wg := sync.WaitGroup{}
	ok := true
	for i := 0; i < f.NumFileParts; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := f.FileParts[i].Encrypt(key); err != nil {
				ok = false
			}
		}(i)
	}
	wg.Wait()
	if !ok {
		return fmt.Errorf("encrypt file failed")
	}
	return nil
}

func (f *File) Decrypt(key []byte) error {
	wg := sync.WaitGroup{}
	ok := true
	for i := 0; i < f.NumFileParts; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := f.FileParts[i].Decrypt(key); err != nil {
				ok = false
			}
		}(i)
	}
	wg.Wait()
	if !ok {
		return fmt.Errorf("decrypt file failed")
	}
	return nil
}

func (f *File) SaveInfo() error {
	info := setting.FileInfo{
		Name:         f.Name,
		Key:          f.Key,
		NumFileParts: f.NumFileParts,
		FileParts:    make([]string, f.NumFileParts),
	}
	for i := 0; i < f.NumFileParts; i++ {
		info.FileParts[i] = f.FileParts[i].Name
	}
	// create a file to store the file info
	infoFile, err := os.Create(f.Name + ".info")
	if err != nil {
		return err
	}
	defer infoFile.Close()
	// write the file info to the file
	if _, err := info.WriteTo(infoFile); err != nil {
		return err
	}
	return nil
}

func (f *File) Save() error {
	wg := sync.WaitGroup{}
	ok := true
	for i := 0; i < f.NumFileParts; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := f.FileParts[i].Save(); err != nil {
				ok = false
			}
		}(i)
	}
	wg.Wait()
	if !ok {
		return fmt.Errorf("save file failed")
	}
	// create a file info json file to store the file parts info
	// this file will be used when decrypting the file
	if err := f.SaveInfo(); err != nil {
		return err
	}
	return nil
}
