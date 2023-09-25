package constants

// Messages ID
const (
	WarnHealthSuccess    = "warning.healthSuccess"
	WarnUndefinedProcess = "warning.undefinedProcess"
	WarnUndefinedMethod  = "warning.undefinedMethod"
	WarnInternalError    = "warning.internalError"
	WarnInvalidSession   = "warning.invalidSession"
	WarnBadRequest       = "warning.badRequest"
	WarnNotFound         = "warning.notFound"
)

// Messages response
const (
	HealthCheck = "Health Check Success"

	BadRequest   = "Bad Request"
	EmptyRequest = "Empty Request"
	NotFound     = "Not Found"

	UndefinedProcess = "Undefined Process"
	UndefinedMethods = "Undefined Methods"

	InvalidSession = "Invalid Session"

	FailedToRetry  = "Failed to retry"
	FailedToCancel = "Failed to cancel"

	InternalServerError = "Internal server error, Please come back later"
)

/*

000: Success
100: Authentication
300: Data store
400: Invalid client
500: Utils
900: General error

*/

const (
	SuccessCode = "200"

	InvalidSessionCode = "100"

	FailedToRetryCode  = "300"
	FailedToCancelCode = "301"

	BadRequestCode   = "400"
	EmptyRequestCode = "401"
	NotFoundCode     = "404"

	UndefinedProcessCode    = "902"
	UndefinedMethodsCode    = "903"
	InternalServerErrorCode = "904"
	GeneralErrorCode        = "905"

	// error code for submit non-nego
	SystemInterferenceCode = "501"
	ExpiredSettlementCode  = "502"
	AfterHoursCode         = "503"
)

const (
	SuccessReqestId   = "SUCCESS for requestid: %v"
	ErrorForRequestId = " ERROR %v for requestId: %v"
)

const (
	// Admin
	LoginAdmin  = "[ADMIN][LoginAdmin]"
	GetItem     = "[ADMIN][GetItems]"
	CreateItem  = "[ADMIN][CreateItem]"
	UpdateItem  = "[ADMIN][UpdateItem]"
	DeleteItem  = "[ADMIN][DeleteItem]"
	LogoutAdmin = "[ADMIN][LogoutAdmin]"
)
