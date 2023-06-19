package errors

import "github.com/pkg/errors"

const (
	// common error.
	InternalErr                InternalError     = "E100001"
	RequestInvalidArgumentErr  BadRequestError   = "E100002"
	InvalidIDTokenErr          UnauthorizedError = "E100003"
	RequireStaffSessionErr     UnauthorizedError = "E100004"
	InvalidAdminRequestUserErr UnauthorizedError = "E100005"

	// tenant error.
	TenantNotFoundErr NotFoundError = "E200101"

	// staff error.
	StaffNotFoundErr NotFoundError = "E200201"
)

var errorMessageMap = map[error]string{
	// common error
	InternalErr:                "An internal error has occurred",
	RequestInvalidArgumentErr:  "Request argument is invalid",
	InvalidIDTokenErr:          "Invalid ID token",
	RequireStaffSessionErr:     "Require staff session",
	InvalidAdminRequestUserErr: "Invalid admin request user",

	// tenant error
	TenantNotFoundErr: "Tenant not found",

	// staff error
	StaffNotFoundErr: "Staff not found",
}

func ExtractPlaneErrMessage(err error) (code string, message string) { //nolint: nonamedreturns
	switch errors.Cause(err) {
	// tenant error
	case TenantNotFoundErr:
		return TenantNotFoundErr.Error(), errorMessageMap[TenantNotFoundErr]

	// staff error
	case StaffNotFoundErr:
		return StaffNotFoundErr.Error(), errorMessageMap[StaffNotFoundErr]

	// common error
	case InvalidIDTokenErr:
		return InvalidIDTokenErr.Error(), errorMessageMap[InvalidIDTokenErr]
	case RequireStaffSessionErr:
		return RequireStaffSessionErr.Error(), errorMessageMap[RequireStaffSessionErr]
	case InvalidAdminRequestUserErr:
		return InvalidAdminRequestUserErr.Error(), errorMessageMap[InvalidAdminRequestUserErr]
	case RequestInvalidArgumentErr:
		return RequestInvalidArgumentErr.Error(), errorMessageMap[RequestInvalidArgumentErr]
	}
	return InternalErr.Error(), errorMessageMap[InternalErr]
}
