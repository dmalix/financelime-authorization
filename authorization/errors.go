package authorization

const httpErrorTextNotFound = "Not Found"
const httpErrorTextInternalServerError = "Server Internal Error"

const domainErrorCodeBadParams = "BAD_PARAMETERS"      // one or more of the requestInput parameters are invalid
const domainErrorCodeBadParamEmail = "BAD_PARAM_EMAIL" // the Email param is not valid
const domainErrorCodeBadParamPassword = "BAD_PARAM_PASSWORD"
const domainErrorCodeBadParamInvite = "BAD_PARAM_INVITE"                                    // the InviteCode param is not valid
const domainErrorCodeBadParamLang = "BAD_PARAM_LANG"                                        // the Language param is not valid
const domainErrorCodeBadParamRemoteAddr = "BAD_PARAM_REMOTE_ADDR"                           // the RemoteAddr param is not valid
const domainErrorCodeBadParamConfirmationKey = "BAD_PARAM_CONFIRMATION_KEY"                 // the ConfirmationKey param is not valid
const domainErrorCodeUserAlreadyExist = "USER_ALREADY_EXIST"                                // the user with the email you specified already exists
const domainErrorCodeInviteNotFound = "INVITE_NOT_FOUND"                                    // the invite code does not exist or is expired
const domainErrorCodeInviteHasEnded = "INVITE_HAS_ENDED"                                    // the invite code has ended
const domainErrorCodeBadConfirmationKey = "BAD_CONFIRMATION_KEY"                            // the confirmation key not valid
const domainErrorCodeConfirmationKeyNotFound = "CONFIRMATION_KEY_NOT_FOUND"                 // the confirmation key hasn't found or expired
const domainErrorCodeConfirmationKeyAlreadyConfirmed = "CONFIRMATION_KEY_ALREADY_CONFIRMED" // the user email is already confirmed
const domainErrorCodeUserNotFound = "USER_NOT_FOUND"                                        // the user is not found                                       // User is not found
const domainErrorCodeBadRefreshToken = "BAD_REFRESH_TOKEN"                                  // failed to validate the Refresh Token (JWT)
const domainErrorSessionNotFound = "SESSION_NOT_FOUND"                                      // the case (the session + hashedRefreshToken) does not exist
