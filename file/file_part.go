package file

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
)

const FilePartSize = 50 * 1024 * 1024 // 50MB

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
