package errors

type ErrorCategory string

const (
	ErrorCategoryUnknown          ErrorCategory = "unknown"
	ErrorCategoryBadRequest       ErrorCategory = "bad_request"
	ErrorCategoryUnauthorized     ErrorCategory = "unauthorized"
	ErrorCategoryForbidden        ErrorCategory = "forbidden"
	ErrorCategoryNotFound         ErrorCategory = "not_found"
	ErrorCategoryConflict         ErrorCategory = "conflict"
	ErrorCategoryCanceled         ErrorCategory = "canceled"
	ErrorCategoryInternal         ErrorCategory = "internal"
	ErrorCategoryServiceAvailable ErrorCategory = "service_available"
)

func (c ErrorCategory) String() string {
	return string(c)
}
