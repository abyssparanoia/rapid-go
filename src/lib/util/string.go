package util

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"unsafe"
)

// StrToMD5 ... generate hash string by MD5
func StrToMD5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// StrToSHA256 ... generate hash string by sha256
func StrToSHA256(str string) string {
	c := sha256.Sum256([]byte(str))
	return hex.EncodeToString(c[:])
}

// StrToBytes ... convert string to bytes
func StrToBytes(str string) []byte {
	return *(*[]byte)(unsafe.Pointer(&str))
}
