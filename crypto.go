package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

func generateID() string {
	buf := make([]byte, 32)
	io.ReadFull(rand.Reader, buf)
	return hex.EncodeToString(buf)
}

// hashKey creates a one way encrypted hash key
// which can be used to store the key in file
func hashKey(key string) string {
	hash := md5.Sum([]byte(key))
	return hex.EncodeToString(hash[:])
}

func newEncryptionKey() []byte {
	keyBuf := make([]byte, 32)
	io.ReadFull(rand.Reader, keyBuf)
	return keyBuf
}

func copyStream(stream cipher.Stream, blockSize int, src io.Reader, dst io.Writer) (int, error) {
	var (
		buf = make([]byte, 32*1024)
		// define total writer length starting with size of block i.e. 16 bytes
		nw = blockSize
	)

	for {
		n, err := src.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf, buf[:n])
			// len of unencrypted
			nn, err := dst.Write(buf[:n])
			if err != nil {
				return 0, err
			}

			// add file lengnth (nn) to total length (nw)
			nw += nn
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
	}

	return nw, nil
}

func copyDecrypt(key []byte, src io.Reader, dst io.Writer) (int, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	// read the IV from the given io.Reader which in our case
	// should be the block.BlockSize() bytes we read.

	// create iv byte array
	iv := make([]byte, block.BlockSize())
	// read length into it
	if _, err := src.Read(iv); err != nil {
		return 0, err
	}

	stream := cipher.NewCTR(block, iv)
	return copyStream(stream, block.BlockSize(), src, dst)
}

func copyEncrypt(key []byte, src io.Reader, dst io.Writer) (int, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	// IV is important, must be in file for decryption.
	// Suggested prepend
	iv := make([]byte, block.BlockSize()) // 16 bytes
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return 0, err
	}

	// prepend iv to file
	if _, err := dst.Write(iv); err != nil {
		return 0, err
	}

	stream := cipher.NewCTR(block, iv)
	return copyStream(stream, block.BlockSize(), src, dst)
}
