## Release v0.3.0-beta

- Added new environment variables: JWT_ACCESS_ENCRYPT_DATA and JWT_REFRESH_ENCRYPT_DATA.
- Now the content of the JWT token in the Claim Data field is encrypted using the JWT_ACCESS_ENCRYPT_DATA and JWT_REFRESH_ENCRYPT_DATA keys.
- Added new environment variables: JWT_ACCESS_AUDIENCE and JWT_REFRESH_AUDIENCE.
- Updated package github.com/dmalix/middleware to v0.1.13
- Updated package github.com/dmalix/jwt to v0.1.8
- Updated package github.com/dmalix/secretdata to v0.1.4

## Release v0.2.0-alpha

- Implemented the generation of the OpenAPI spec. 
- Added context to some functions.
- Changed the graceful shutdown.
- Implemented Zap-logger.
- Added a confirmation for a user password reset (a step 2).

## Pre-Release v0.1.5-alpha

- Get the current version of the REST API service
- Sign Up 
- Confirm User Email
- Request Access Token
- Get a list of Active Sessions
- Refresh Access Token
- Revoke Refresh Token
- Request User Password reset
- Add error labels.
- Make some improvements to function and variable names.
- Add human-readable errors (domain error code).
- Fix minor inaccuracies in SQL-queries.
- Create a new "utils/email" package and implement it into the authorization package.
- Make minor improvements to the project structure.
- Add a sending email function with a link to verify the email address.
- Small changes in the user model.

## Pre-Release v0.1.4-alpha

- Transfer some functions from the router to utilities.

## Pre-Release v0.1.3-alpha

- Add method handling to the router.

## Pre-Release v0.1.2-alpha

- Make some router improvements.

## Pre-Release v0.1.1-alpha

- Add checking the `request-id` header.

## Pre-Release v0.1.0-alpha

Initial Release

- Added a Signup method.
