//nolint:govet
package errors

import "github.com/abyssparanoia/goerr"

func NewBadRequestError(errCode string, msg string) *goerr.Error {
	return goerr.New(msg).
		WithCategory(ErrorCategoryBadRequest.String()).
		WithCode(errCode)
}

func NewUnauthorizedError(errCode string, msg string) *goerr.Error {
	return goerr.New(msg).
		WithCategory(ErrorCategoryUnauthorized.String()).
		WithCode(errCode)
}

func NewForbiddenError(errCode string, msg string) *goerr.Error {
	return goerr.New(msg).
		WithCategory(ErrorCategoryForbidden.String()).
		WithCode(errCode)
}

func NewNotFoundError(errCode string, msg string) *goerr.Error {
	return goerr.New(msg).
		WithCategory(ErrorCategoryNotFound.String()).
		WithCode(errCode)
}

func NewConflictError(errCode string, msg string) *goerr.Error {
	return goerr.New(msg).
		WithCategory(ErrorCategoryConflict.String()).
		WithCode(errCode)
}

func NewCanceledError(errCode string, msg string) *goerr.Error {
	return goerr.New(msg).
		WithCategory(ErrorCategoryCanceled.String()).
		WithCode(errCode)
}

func NewInternalError(errCode string, msg string) *goerr.Error {
	return goerr.New(msg).
		WithCategory(ErrorCategoryInternal.String()).
		WithCode(errCode)
}

func NewServiceAvailableError(errCode string, msg string) *goerr.Error {
	return goerr.New(msg).
		WithCategory(ErrorCategoryServiceAvailable.String()).
		WithCode(errCode)
}
