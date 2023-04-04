package errors

import "github.com/pkg/errors"

const (
	// common error
	InternalErr                InternalError     = "E100001"
	UnauthorizedErr            UnauthorizedError = "E100002"
	RequestInvalidArgumentErr  BadRequestError   = "E100003"
	NotFoundErr                NotFoundError     = "E100004"
	ConflictErr                ConflictError     = "E100005"
	InvalidAdminRequestUserErr UnauthorizedError = "E100006"

	// tenant error
	TenantNotFoundErr NotFoundError = "E200101"

	// staff error
	StaffNotFoundErr NotFoundError = "E200201"
)

var errorMessageMap = map[error]string{
	// common error
	InternalErr:                "An internal error has occurred",
	UnauthorizedErr:            "Unauthroized",
	RequestInvalidArgumentErr:  "Request argument is invalid",
	NotFoundErr:                "Not found",
	ConflictErr:                "Already Exist",
	InvalidAdminRequestUserErr: "Invalid admin request user",

	// tenant error
	TenantNotFoundErr: "Tenant not found",

	// staff error
	StaffNotFoundErr: "Staff not found",
}

func ExtractPlaneErrMessage(err error) (code string, message string) {
	switch errors.Cause(err) {

	// tenant error
	case TenantNotFoundErr:
		return TenantNotFoundErr.Error(), errorMessageMap[TenantNotFoundErr]

	// staff error
	case StaffNotFoundErr:
		return StaffNotFoundErr.Error(), errorMessageMap[StaffNotFoundErr]

	// common error
	case InvalidAdminRequestUserErr:
		return InvalidAdminRequestUserErr.Error(), errorMessageMap[InvalidAdminRequestUserErr]
	case UnauthorizedErr:
		return UnauthorizedErr.Error(), errorMessageMap[UnauthorizedErr]
	case RequestInvalidArgumentErr:
		return RequestInvalidArgumentErr.Error(), errorMessageMap[RequestInvalidArgumentErr]
	case NotFoundErr:
		return NotFoundErr.Error(), errorMessageMap[NotFoundErr]
	case ConflictErr:
		return ConflictErr.Error(), errorMessageMap[ConflictErr]
	}
	return InternalErr.Error(), errorMessageMap[InternalErr]
}
