package restv1

const (
	headerKeyContentType       = "content-type"
	headerValueApplicationJson = "application/json;charset=utf-8"
	headerValueTextPlain       = "text/plain;charset=utf-8"

	statusMessageNotFound            = "404 Not Found"
	statusMessageBadRequest          = "400 Bad Request"
	statusMessageUnauthorized        = "401 Unauthorized"
	statusMessageConflict            = "409 Conflict"
	statusMessageInternalServerError = "500 Internal Server Error"
)
