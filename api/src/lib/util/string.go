package util

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"unsafe"
)

// StrToMD5 ... 文字列のハッシュ(MD5)を取得する
func StrToMD5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// StrToSHA256 ... 文字列のハッシュ(SHA256)を取得する
func StrToSHA256(str string) string {
	c := sha256.Sum256([]byte(str))
	return hex.EncodeToString(c[:])
}

// StrToBytes ... 文字列をバイト列に変換する
func StrToBytes(str string) []byte {
	return *(*[]byte)(unsafe.Pointer(&str))
}
