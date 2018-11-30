package util

import "unsafe"

// BytesToStr ... バイト列を文字列に変換する
func BytesToStr(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}
