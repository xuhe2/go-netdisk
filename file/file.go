package file

import (
	"fmt"
	"io"
	"os"
	"strconv"
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
		f.FileParts = append(f.FileParts, NewFilePart(f.Name+strconv.Itoa(f.NumFileParts), content[:n]))
		f.NumFileParts++
	}
	return nil
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
