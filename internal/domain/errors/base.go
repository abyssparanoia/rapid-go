//nolint:govet
package errors

import "github.com/abyssparanoia/goerr"

func NewBadRequestError(errCode string, msg string) *goerr.Error {
	return goerr.New("%s", msg).
		WithCategory(ErrorCategoryBadRequest.String()).
		WithCode(errCode)
}

func NewUnauthorizedError(errCode string, msg string) *goerr.Error {
	return goerr.New("%s", msg).
		WithCategory(ErrorCategoryUnauthorized.String()).
		WithCode(errCode)
}

func NewForbiddenError(errCode string, msg string) *goerr.Error {
	return goerr.New("%s", msg).
		WithCategory(ErrorCategoryForbidden.String()).
		WithCode(errCode)
}

func NewNotFoundError(errCode string, msg string) *goerr.Error {
	return goerr.New("%s", msg).
		WithCategory(ErrorCategoryNotFound.String()).
		WithCode(errCode)
}

func NewConflictError(errCode string, msg string) *goerr.Error {
	return goerr.New("%s", msg).
		WithCategory(ErrorCategoryConflict.String()).
		WithCode(errCode)
}

func NewCanceledError(errCode string, msg string) *goerr.Error {
	return goerr.New("%s", msg).
		WithCategory(ErrorCategoryCanceled.String()).
		WithCode(errCode)
}

func NewInternalError(errCode string, msg string) *goerr.Error {
	return goerr.New("%s", msg).
		WithCategory(ErrorCategoryInternal.String()).
		WithCode(errCode)
}

func NewServiceAvailableError(errCode string, msg string) *goerr.Error {
	return goerr.New("%s", msg).
		WithCategory(ErrorCategoryServiceAvailable.String()).
		WithCode(errCode)
}
