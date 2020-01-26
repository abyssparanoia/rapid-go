package parameter

import (
	"context"
	"net/http"
	"strconv"

	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
	"github.com/go-chi/chi"
)

// GetURL ... get url parameter
func GetURL(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

// GetURLByInt ... get url parameter by int
func GetURLByInt(ctx context.Context, r *http.Request, key string) (int, error) {
	str := chi.URLParam(r, key)
	if str == "" {
		return 0, nil
	}
	num, err := strconv.Atoi(str)
	if err != nil {
		log.Warningm(ctx, "strconv.Atoi", err)
		return 0, err
	}
	return num, nil
}

// GetURLByInt64 ... get url parameter by int64
func GetURLByInt64(ctx context.Context, r *http.Request, key string) (int64, error) {
	str := chi.URLParam(r, key)
	if str == "" {
		return 0, nil
	}
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Warningm(ctx, "strconv.ParseInt", err)
		return 0, err
	}
	return num, nil
}

// GetURLByFloat64 ... get url parameter by float64
func GetURLByFloat64(ctx context.Context, r *http.Request, key string) (float64, error) {
	str := chi.URLParam(r, key)
	if str == "" {
		return 0, nil
	}
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Warningm(ctx, "strconv.ParseFloat", err)
		return 0, err
	}
	return num, nil
}
