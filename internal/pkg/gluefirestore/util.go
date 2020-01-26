package gluefirestore

import (
	"reflect"

	"cloud.google.com/go/firestore"
)

func setDocByDst(dst interface{}, ref *firestore.DocumentRef) {
	rv := reflect.Indirect(reflect.ValueOf(dst))
	rt := rv.Type()
	if rt.Kind() == reflect.Struct {
		for i := 0; i < rt.NumField(); i++ {
			f := rt.Field(i)
			tag := f.Tag.Get("gluefirestore")
			if tag == "id" && f.Type.Kind() == reflect.String {
				rv.Field(i).SetString(ref.ID)
				continue
			}
			if tag == "ref" && f.Type.Kind() == reflect.Ptr {
				rv.Field(i).Set(reflect.ValueOf(ref))
				continue
			}
		}
	}
}

func setDocByDsts(rv reflect.Value, rt reflect.Type, ref *firestore.DocumentRef) {
	if rt.Kind() == reflect.Struct {
		for i := 0; i < rt.NumField(); i++ {
			f := rt.Field(i)
			tag := f.Tag.Get("gluefirestore")
			if tag == "id" && f.Type.Kind() == reflect.String {
				rv.Elem().Field(i).SetString(ref.ID)
				continue
			}
			if tag == "ref" && f.Type.Kind() == reflect.Ptr {
				rv.Elem().Field(i).Set(reflect.ValueOf(ref))
				continue
			}
		}
	}
}
