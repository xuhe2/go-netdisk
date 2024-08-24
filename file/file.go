package file

import (
	"fmt"
	"io"
	"strconv"
	"sync"
)

type File struct {
	Name         string `json:"file_name"`
	NumFileParts int    `json:"num_file_parts"`
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
	return nil
}
