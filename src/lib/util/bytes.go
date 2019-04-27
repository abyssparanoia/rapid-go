package util

import "unsafe"

// BytesToStr ... convert byte to string
func BytesToStr(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}
