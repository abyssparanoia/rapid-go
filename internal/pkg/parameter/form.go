package parameter

import (
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"

	"github.com/abyssparanoia/rapid-go/internal/pkg/util"
)

// GetForm ... get form value
func GetForm(r *http.Request, key string) string {
	return r.FormValue(key)
}

// GetFormByInt ... get form value by int
func GetFormByInt(ctx context.Context, r *http.Request, key string) (int, error) {
	str := r.FormValue(key)
	if str == "" {
		return 0, nil
	}
	num, err := strconv.Atoi(str)
	if err != nil {
		return num, err
	}
	return num, nil
}

// GetFormByIntOptional ... get form value by optional int
func GetFormByIntOptional(ctx context.Context, r *http.Request, key string) (*int, error) {
	str := r.FormValue(key)
	if str == "" {
		return nil, nil
	}
	num, err := strconv.Atoi(str)
	if err != nil {
		return nil, err
	}
	return &num, nil
}

// GetFormByInt64 ... get form value by int64
func GetFormByInt64(ctx context.Context, r *http.Request, key string) (int64, error) {
	str := r.FormValue(key)
	if str == "" {
		return 0, nil
	}
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return num, err
	}
	return num, nil
}

// GetFormByInt64Optional ... get form value by optional int64
func GetFormByInt64Optional(ctx context.Context, r *http.Request, key string) (*int64, error) {
	str := r.FormValue(key)
	if str == "" {
		return nil, nil
	}
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return nil, err
	}
	return &num, nil
}

// GetFormByFloat64 ... get form value by float64
func GetFormByFloat64(ctx context.Context, r *http.Request, key string) (float64, error) {
	str := r.FormValue(key)
	if str == "" {
		return 0, nil
	}
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return num, err
	}
	return num, nil
}

// GetFormByBool ... get form value by bool
func GetFormByBool(ctx context.Context, r *http.Request, key string) (bool, error) {
	str := r.FormValue(key)
	if str == "" {
		return false, nil
	}
	val, err := strconv.ParseBool(str)
	if err != nil {
		return val, err
	}
	return val, nil
}

// GetForms ... get form value
func GetForms(ctx context.Context, r *http.Request, dst interface{}) error {
	if reflect.TypeOf(dst).Kind() != reflect.Ptr {
		return errors.New("dst isn't a pointer")
	}

	paramType := reflect.TypeOf(dst).Elem()
	paramValue := reflect.ValueOf(dst).Elem()

	fieldCount := paramType.NumField()
	for i := 0; i < fieldCount; i++ {
		field := paramType.Field(i)

		formTag := paramType.Field(i).Tag.Get("form")
		if util.IsZero(formTag) {
			continue
		}

		fieldValue := paramValue.FieldByName(field.Name)
		if !fieldValue.CanSet() {
			return errors.New("fieldValue.CanSet")
		}
		switch field.Type.Kind() {
		case reflect.Int64:
			val, err := GetFormByInt64(ctx, r, formTag)
			if err != nil {
			}
			fieldValue.SetInt(val)
		case reflect.Int:
			val, err := GetFormByInt64(ctx, r, formTag)
			if err != nil {
			}
			fieldValue.SetInt(val)
		case reflect.Float64:
			val, err := GetFormByFloat64(ctx, r, formTag)
			if err != nil {
			}
			fieldValue.SetFloat(val)
		case reflect.String:
			val := GetForm(r, formTag)
			fieldValue.SetString(val)
		case reflect.Bool:
			val, err := GetFormByBool(ctx, r, formTag)
			if err != nil {
			}
			fieldValue.SetBool(val)
		}
	}
	return nil
}

// GetFormFile ... get form file
func GetFormFile(r *http.Request, key string) (multipart.File, *multipart.FileHeader, error) {
	return r.FormFile(key)
}
