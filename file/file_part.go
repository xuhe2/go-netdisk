package file

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
)

const FilePartSize = 49 * 1024 * 1024 // 49MB

type FilePart struct {
	Name string `json:"file_name"`
	Data []byte `json:"file_data"`
}

func NewFilePart(name string, data []byte) *FilePart {
	return &FilePart{
		Name: name,
		Data: data,
	}
}

func (fp *FilePart) LoadData(r io.Reader) error {
	// read the reader and put all content into the file parts data
	// the file parts data may be nil
	n := FilePartSize
	// if r is a file, get the file size
	if f, ok := r.(*os.File); ok {
		fi, err := f.Stat()
		if err != nil {
			return err
		}
		n = int(fi.Size())
	}
	if fp.Data == nil {
		fp.Data = make([]byte, n)
	}
	_, err := io.ReadFull(r, fp.Data)
	return err
}

func (fp *FilePart) Encrypt(key []byte) error {
	// use AES to encrypt the file data
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	// encrypt the file data
	ciphertext := make([]byte, aes.BlockSize+len(fp.Data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}
	// encrypt the file data
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], fp.Data)
	// set the encrypted data
	fp.Data = ciphertext
	return nil
}

func (fp *FilePart) Decrypt(key []byte) error {
	// use AES to decrypt the file data
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	// decrypt the file data
	if len(fp.Data) < aes.BlockSize {
		return err
	}
	iv := fp.Data[:aes.BlockSize]
	chipertext := fp.Data[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(chipertext, chipertext)
	// set the decrypted data
	fp.Data = chipertext
	return nil
}

func (fp *FilePart) Save() error {
	file, err := os.Create(fp.Name)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(fp.Data)
	return err
}
