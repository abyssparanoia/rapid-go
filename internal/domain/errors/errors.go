package errors

var (
	InternalErr                = NewInternalError("E100001", "An internal error has occurred")
	RequestInvalidArgumentErr  = NewBadRequestError("E100002", "Request argument is invalid")
	InvalidIDTokenErr          = NewUnauthorizedError("E100003", "Invalid ID token")
	RequireStaffSessionErr     = NewUnauthorizedError("E100004", "Require staff session")
	InvalidAdminRequestUserErr = NewUnauthorizedError("E100005", "Invalid admin request user")
	AssetInvalidErr            = NewBadRequestError("E100006", "Asset is invalid")
	AssetNotFoundErr           = NewNotFoundError("E100007", "Asset not found")
	CanceledErr                = NewCanceledError("E100008", "Canceled")

	// tenant error.
	TenantNotFoundErr = NewNotFoundError("E200101", "Tenant not found")

	// staff error.
	StaffNotFoundErr = NewNotFoundError("E200201", "Staff not found")
)
