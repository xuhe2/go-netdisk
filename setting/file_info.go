package setting

import (
	"encoding/json"
	"io"
)

type FileInfo struct {
	Key          []byte
	NumFileParts int
	FileParts    []string
}

func NewFileInfo(key []byte, numFileParts int, fileParts []string) *FileInfo {
	return &FileInfo{
		Key:          key,
		NumFileParts: numFileParts,
		FileParts:    fileParts,
	}
}

func (fi *FileInfo) Marshal() ([]byte, error) {
	// 二进制化
	return json.Marshal(*fi)
}

func (fi *FileInfo) Unmarshal(content []byte) error {
	// 从二进制中解析
	return json.Unmarshal(content, fi)
}

func (fi *FileInfo) WriteTo(w io.Writer) error {
	// 二进制化之后写入
	content, err := fi.Marshal()
	if err != nil {
		return err
	}
	_, err = w.Write(content)
	return err
}

func (fi *FileInfo) ReadFrom(r io.Reader) error {
	// 从二进制中读取
	buf := make([]byte, 2048)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf[:n], fi)
}
