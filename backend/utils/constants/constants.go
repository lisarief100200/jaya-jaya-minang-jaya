package constants

const (
	SUCCESS = "success"
	FAILED  = "failed"
	REJECT  = "reject"
)

const (
	PROD = "PROD"
	DEV  = "DEV"
)

const (
	//JSONType request api options type for json request
	JSONType = "application/json"

	//STREAMType request api options type for octet stream
	STREAMType = "application/octet-stream"

	//PDFType request api options type for type PDF
	PDFType = "application/pdf"

	//XMLType request api options type for xml request
	XMLType = "application/xml"

	//URLEncodedType request api options type for url encoded request
	URLEncodedType = "application/x-www-form-urlencoded"

	//FormDataType request api options type for type form data
	FormDataType = "multipart/form-data"

	//TextHtmlType request api options type for type html
	TextHtmlType = "text/html"
)

const (
	HeaderContentType   = "Content-Type"
	HeaderContentLength = "Content-Length"

	HeaderAuthorization = "Authorization"
	HeaderAuthenticate  = "WWW-Authenticate"
)

const (
	AuthTypeBearer = "Bearer"
	AuthTypeBasic  = "Basic"
)
