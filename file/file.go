package file

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
)

type File struct {
	Name string `json:"file_name"`
	Data []byte `json:"file_data"`
}

func (f *File) Open(name string) error {
	var err error
	f.Data, err = os.ReadFile(name)
	if err != nil {
		return err
	}
	f.Name = name
	return nil
}

func (f *File) Encrypt(key []byte) error {
	// use AES to encrypt the file data
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	// encrypt the file data
	ciphertext := make([]byte, aes.BlockSize+len(f.Data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}
	// encrypt the file data
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], f.Data)
	// set the encrypted data
	f.Data = ciphertext
	return nil
}

func (f *File) Decrypt(key []byte) error {
	// use AES to decrypt the file data
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	// decrypt the file data
	if len(f.Data) < aes.BlockSize {
		return err
	}
	iv := f.Data[:aes.BlockSize]
	chipertext := f.Data[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(chipertext, chipertext)
	// set the decrypted data
	f.Data = chipertext
	return nil
}
