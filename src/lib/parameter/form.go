package parameter

import (
	"context"
	"net/http"
	"reflect"
	"strconv"

	"github.com/abyssparanoia/rapid-go-worker/src/lib/log"
	"github.com/abyssparanoia/rapid-go-worker/src/lib/util"
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
		log.Warningm(ctx, "strconv.Atoi", err)
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
		log.Warningm(ctx, "strconv.Atoi", err)
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
		log.Warningm(ctx, "strconv.ParseInt", err)
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
		log.Warningm(ctx, "strconv.ParseInt", err)
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
		log.Warningm(ctx, "strconv.ParseFloat", err)
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
		log.Warningm(ctx, "strconv.ParseInt", err)
		return val, err
	}
	return val, nil
}

// GetForms ... get form value
func GetForms(ctx context.Context, r *http.Request, dst interface{}) error {
	if reflect.TypeOf(dst).Kind() != reflect.Ptr {
		err := log.Errore(ctx, "dst isn't a pointer")
		return err
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
			err := log.Warningc(ctx, http.StatusBadRequest, "fieldValue.CanSet")
			return err
		}
		switch field.Type.Kind() {
		case reflect.Int64:
			val, err := GetFormByInt64(ctx, r, formTag)
			if err != nil {
				log.Debugm(ctx, "GetFormByInt64", err)
			}
			fieldValue.SetInt(val)
		case reflect.Int:
			val, err := GetFormByInt64(ctx, r, formTag)
			if err != nil {
				log.Debugm(ctx, "GetFormByInt64", err)
			}
			fieldValue.SetInt(val)
		case reflect.Float64:
			val, err := GetFormByFloat64(ctx, r, formTag)
			if err != nil {
				log.Debugm(ctx, "GetFormByFloat64", err)
			}
			fieldValue.SetFloat(val)
		case reflect.String:
			val := GetForm(r, formTag)
			fieldValue.SetString(val)
		case reflect.Bool:
			val, err := GetFormByBool(ctx, r, formTag)
			if err != nil {
				log.Debugm(ctx, "GetFormByBool", err)
			}
			fieldValue.SetBool(val)
		}
	}
	return nil
}
