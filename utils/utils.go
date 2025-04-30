package utils

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func CompressData(data []byte) ([]byte, error) {

	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	w.Close()

	return b.Bytes(), nil
}

// UnCompressData decompresses a zlib-compressed byte slice
func UnCompressData(data []byte) (string, error) {
	b := bytes.NewReader(data)

	r, err := zlib.NewReader(b)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var out bytes.Buffer
	_, err = io.Copy(&out, r)
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func GenerateHash(data []byte) (string, error) {
	h := sha1.New()

	// Write the input data to the hash
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	// Return the hex-encoded hash
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
