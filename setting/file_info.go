package setting

import (
	"encoding/json"
	"io"
)

type FileInfo struct {
	Name         string
	Key          []byte
	NumFileParts int
	FileParts    []string
}

func NewFileInfo(name string, key []byte, numFileParts int, fileParts []string) *FileInfo {
	return &FileInfo{
		Name:         name,
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

func (fi *FileInfo) WriteTo(w io.Writer) (int64, error) {
	// 二进制化之后写入
	content, err := fi.Marshal()
	if err != nil {
		return 0, err
	}
	_, err = w.Write(content)
	return int64(len(content)), err
}

func (fi *FileInfo) ReadFrom(r io.Reader) (int64, error) {
	// 从二进制中读取
	buf := make([]byte, 2048)
	n, err := r.Read(buf)
	if err != nil {
		return 0, err
	}
	return int64(n), fi.Unmarshal(buf[:n])
}
