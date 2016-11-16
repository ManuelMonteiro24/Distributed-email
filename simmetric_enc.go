package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func main() {
	originalText := "encrypt this"
	fmt.Println(originalText)

	key := make([]byte, 16)

	//this key needs to be saved and encrypted with dest public key
  _, err := rand.Read(key)
  Check(err)

	// encrypt value to base64
	cryptoText := encrypt(key, originalText)
	fmt.Println(cryptoText)

	// encrypt base64 crypto to original value
	text := decrypt(key, cryptoText)
	fmt.Printf(text)
}


// encrypt string to base64 crypto using AES
func encrypt(key []byte, text string) string {
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	Check(err)

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return Encode64(ciphertext)
}

// decrypt from base64 to decrypted string
func decrypt(key []byte, cryptoText string) string {
	ciphertext, _ := Decode64(cryptoText)

	block, err := aes.NewCipher(key)
	Check(err)

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s\n", ciphertext)
}
