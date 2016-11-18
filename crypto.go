package main

import (
	"bytes"
	"compress/flate"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"

	"github.com/gtank/cryptopasta"
)

var ENCRYPTION_COMPRESSION = true
var ENCRYPTION_ENABLED = true

func decrypt(content, password []byte) ([]byte, error) {
	key := sha256.Sum256(password)
	decrypted, err := cryptopasta.Decrypt(content, &key)
	if ENCRYPTION_COMPRESSION {
		decrypted = decompressByte(decrypted)
	}
	return decrypted, err
}

func encrypt(content, password []byte) ([]byte, error) {
	if ENCRYPTION_COMPRESSION {
		content = compressByte(content)
	}
	key := sha256.Sum256([]byte(password))
	encrypted, err := cryptopasta.Encrypt(content, &key)
	return encrypted, err
}

func encryptString(content string, password string) (string, error) {
	if !ENCRYPTION_ENABLED {
		return content, nil
	}
	bEncrypted, err := encrypt([]byte(content), []byte(password))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bEncrypted), nil
}

func decryptString(encodedContent string, password string) (string, error) {
	if !ENCRYPTION_ENABLED {
		return encodedContent, nil
	}
	bEncrypted, err := hex.DecodeString(encodedContent)
	if err != nil {
		return "", err
	}
	bDecrypt, err := decrypt(bEncrypted, []byte(password))
	return string(bDecrypt), err
}

func encryptAndWrite(filename, content, password string) (err error) {
	e, err := encryptString(content, password)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, []byte(e), 0644)
}

func openAndDecrypt(filename, password string) (string, error) {
	contentBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return string(contentBytes), err
	}
	return decryptString(string(contentBytes), password)
}

// compressByte returns a compressed byte slice.
func compressByte(src []byte) []byte {
	compressedData := new(bytes.Buffer)
	compress(src, compressedData, 9)
	return compressedData.Bytes()
}

// decompressByte returns a decompressed byte slice.
func decompressByte(src []byte) []byte {
	compressedData := bytes.NewBuffer(src)
	deCompressedData := new(bytes.Buffer)
	decompress(compressedData, deCompressedData)
	return deCompressedData.Bytes()
}

// compress uses flate to compress a byte slice to a corresponding level
func compress(src []byte, dest io.Writer, level int) {
	compressor, _ := flate.NewWriter(dest, level)
	compressor.Write(src)
	compressor.Close()
}

// compress uses flate to decompress an io.Reader
func decompress(src io.Reader, dest io.Writer) {
	decompressor := flate.NewReader(src)
	io.Copy(dest, decompressor)
	decompressor.Close()
}
