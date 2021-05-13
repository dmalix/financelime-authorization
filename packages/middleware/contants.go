package middleware

const (
	statusMessageNotFound            = "404 Not Found"
	statusMessageBadRequest          = "400 Bad Request"
	statusMessageUnauthorized        = "401 Unauthorized"
	statusMessageConflict            = "409 Conflict"
	statusMessageInternalServerError = "500 Internal Server Error"
)

const (
	ContextKeyRemoteAddr       = "remoteAddr"
	ContextKeyRequestID        = "requestID"
	ContextKeyPublicSessionID  = "publicSessionID"
	ContextKeyEncryptedJWTData = "encryptedJWTData"
)
