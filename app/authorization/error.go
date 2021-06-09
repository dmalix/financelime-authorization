package authorization

import "errors"

var ErrorBadParams = errors.New("BAD_PARAMETERS")                                           // one or more of the requestInput parameters are invalid
var ErrorBadParamEmail = errors.New("BAD_PARAM_EMAIL")                                      // the email param is not valid
var ErrorBadParamPassword = errors.New("BAD_PARAM_PASSWORD")                                // the password param is not valid
var ErrorBadParamInvite = errors.New("BAD_PARAM_INVITE")                                    // the InviteCode param is not valid
var ErrorBadParamLang = errors.New("BAD_PARAM_LANG")                                        // the Language param is not valid
var ErrorBadParamConfirmationKey = errors.New("BAD_PARAM_CONFIRMATION_KEY")                 // the ConfirmationKey param is not valid
var ErrorUserAlreadyExist = errors.New("USER_ALREADY_EXIST")                                // a user with the same email address already exists
var ErrorInviteNotFound = errors.New("INVITE_NOT_FOUND")                                    // the invite code does not exist or is expired
var ErrorInviteHasEnded = errors.New("INVITE_HAS_ENDED")                                    // the invite code has ended
var ErrorBadConfirmationKey = errors.New("BAD_CONFIRMATION_KEY")                            // the confirmation key not valid
var ErrorConfirmationKeyNotFound = errors.New("CONFIRMATION_KEY_NOT_FOUND")                 // the confirmation key hasn't found or expired
var ErrorConfirmationKeyAlreadyConfirmed = errors.New("CONFIRMATION_KEY_ALREADY_CONFIRMED") // the user email is already confirmed
var ErrorUserNotFound = errors.New("USER_NOT_FOUND")                                        // the user is not found
var ErrorBadRefreshToken = errors.New("BAD_REFRESH_TOKEN")                                  // failed to validate the Refresh Token (JWT)
var ErrorSessionNotFound = errors.New("SESSION_NOT_FOUND")                                  // the case (the session + hashedRefreshToken) does not exist
