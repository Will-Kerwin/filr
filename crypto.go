package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func newEncryptionKey() []byte {
	keyBuf := make([]byte, 32)
	io.ReadFull(rand.Reader, keyBuf)
	return keyBuf
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

	var (
		buf    = make([]byte, 32*1024)
		stream = cipher.NewCTR(block, iv)
		// define total writer length starting with size of block i.e. 16 bytes
		nw = block.BlockSize()
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

	// write actual file
	var (
		buf    = make([]byte, 32*1024)
		stream = cipher.NewCTR(block, iv)
	)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf, buf[:n])
			if _, err := dst.Write(buf[:n]); err != nil {
				return 0, err
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
	}

	return 0, nil
}
