


# Financelime Authorization
Financelime Authorization RESTful API service
  

## Informations

### Version

v0.2.0-alpha

### License

[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.html)

### Contact

API Support dmalix@financelime.com 

## Content negotiation

### URI Schemes
  * https

### Consumes
  * application/json

### Produces
  * application/json
  * text/plain

## Access control

### Security Schemes

#### authorization (header: authorization)



> **Type**: apikey

## All endpoints

###  operations

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| POST | /v1/oauth | [create access token](#create-access-token) | Create Access Token (Domain Action: Log In) |
| GET | /v1/sessions | [get list active sessions](#get-list-active-sessions) | Get a list of active sessions |
| GET | /v1/version | [get version](#get-version) | Get the Service version |
| PUT | /v1/oauth | [refresh access token](#refresh-access-token) | Refresh Access Token (Domain Action: Renew authorization) |
| PUT | /v1/user | [reset user password step1](#reset-user-password-step1) | Request to user password reset |
| GET | /p/{confirmationKey} | [reset user password step2](#reset-user-password-step2) | Confirm to user password reset |
| DELETE | /v1/session | [revoke refresh token](#revoke-refresh-token) | Revoke Refresh Token (Domain Action: Log Out) |
| POST | /v1/user | [signup step1](#signup-step1) | Create new user |
| GET | /u/{confirmationKey} | [signup step2](#signup-step2) | Confirm User Email |
  


## Paths

### <span id="create-access-token"></span> Create Access Token (Domain Action: Log In) (*create_access_token*)

```
POST /v1/oauth
```

Create Access Token

#### Consumes
  * application/json;charset=utf-8

#### Produces
  * application/json;charset=utf-8

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| request-id | `header` | string | `string` |  | ✓ |  | RequestID |
| model.CreateAccessTokenRequest | `body` | [ModelCreateAccessTokenRequest](#model-create-access-token-request) | `models.ModelCreateAccessTokenRequest` | | ✓ | | Data for creating a new token |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#create-access-token-200) | OK | Successful operation |  | [schema](#create-access-token-200-schema) |
| [400](#create-access-token-400) | Bad Request | Bad Request |  | [schema](#create-access-token-400-schema) |
| [404](#create-access-token-404) | Not Found | Not Found |  | [schema](#create-access-token-404-schema) |
| [500](#create-access-token-500) | Internal Server Error | Internal Server Error |  | [schema](#create-access-token-500-schema) |

#### Responses


##### <span id="create-access-token-200"></span> 200 - Successful operation
Status: OK

###### <span id="create-access-token-200-schema"></span> Schema
   
  

[ModelAccessTokenResponse](#model-access-token-response)

##### <span id="create-access-token-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="create-access-token-400-schema"></span> Schema
   
  

[ModelCreateAccessTokenFailure400](#model-create-access-token-failure400)

##### <span id="create-access-token-404"></span> 404 - Not Found
Status: Not Found

###### <span id="create-access-token-404-schema"></span> Schema
   
  

[ModelCreateAccessTokenFailure404](#model-create-access-token-failure404)

##### <span id="create-access-token-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="create-access-token-500-schema"></span> Schema
   
  

[ModelCommonFailure](#model-common-failure)

### <span id="get-list-active-sessions"></span> Get a list of active sessions (*get_list_active_sessions*)

```
GET /v1/sessions
```

Get a list of active sessions

#### Produces
  * application/json;charset=utf-8

#### Security Requirements
  * authorization

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| request-id | `header` | string | `string` |  | ✓ |  | RequestID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-list-active-sessions-200) | OK | Successful operation |  | [schema](#get-list-active-sessions-200-schema) |
| [500](#get-list-active-sessions-500) | Internal Server Error | Internal Server Error |  | [schema](#get-list-active-sessions-500-schema) |

#### Responses


##### <span id="get-list-active-sessions-200"></span> 200 - Successful operation
Status: OK

###### <span id="get-list-active-sessions-200-schema"></span> Schema
   
  

[][ModelSession](#model-session)

##### <span id="get-list-active-sessions-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="get-list-active-sessions-500-schema"></span> Schema
   
  

[ModelCommonFailure](#model-common-failure)

### <span id="get-version"></span> Get the Service version (*get_version*)

```
GET /v1/version
```

Get Version

#### Produces
  * application/json;charset=utf-8

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| request-id | `header` | string | `string` |  | ✓ |  | RequestID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-version-200) | OK | Successful operation |  | [schema](#get-version-200-schema) |

#### Responses


##### <span id="get-version-200"></span> 200 - Successful operation
Status: OK

###### <span id="get-version-200-schema"></span> Schema
   
  

[ModelVersionResponse](#model-version-response)

### <span id="refresh-access-token"></span> Refresh Access Token (Domain Action: Renew authorization) (*refresh_access_token*)

```
PUT /v1/oauth
```

Refresh Access Token

#### Consumes
  * application/json;charset=utf-8

#### Produces
  * application/json;charset=utf-8

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| request-id | `header` | string | `string` |  | ✓ |  | RequestID |
| model.RefreshAccessTokenRequest | `body` | [ModelRefreshAccessTokenRequest](#model-refresh-access-token-request) | `models.ModelRefreshAccessTokenRequest` | | ✓ | | Data for refreshing the access token |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#refresh-access-token-200) | OK | Successful operation |  | [schema](#refresh-access-token-200-schema) |
| [400](#refresh-access-token-400) | Bad Request | Bad Request |  | [schema](#refresh-access-token-400-schema) |
| [404](#refresh-access-token-404) | Not Found | Not Found |  | [schema](#refresh-access-token-404-schema) |
| [500](#refresh-access-token-500) | Internal Server Error | Internal Server Error |  | [schema](#refresh-access-token-500-schema) |

#### Responses


##### <span id="refresh-access-token-200"></span> 200 - Successful operation
Status: OK

###### <span id="refresh-access-token-200-schema"></span> Schema
   
  

[ModelAccessTokenResponse](#model-access-token-response)

##### <span id="refresh-access-token-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="refresh-access-token-400-schema"></span> Schema
   
  

[ModelRefreshAccessTokenFailure400](#model-refresh-access-token-failure400)

##### <span id="refresh-access-token-404"></span> 404 - Not Found
Status: Not Found

###### <span id="refresh-access-token-404-schema"></span> Schema
   
  

[ModelRefreshAccessTokenFailure404](#model-refresh-access-token-failure404)

##### <span id="refresh-access-token-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="refresh-access-token-500-schema"></span> Schema
   
  

[ModelCommonFailure](#model-common-failure)

### <span id="reset-user-password-step1"></span> Request to user password reset (*reset_user_password_step1*)

```
PUT /v1/user
```

The service sends a confirmation link to the specified email. After confirmation, the service will send a new password for authorization.

#### Consumes
  * application/json;charset=utf-8

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| request-id | `header` | string | `string` |  | ✓ |  | RequestID |
| model.ResetUserPasswordRequest | `body` | [ModelResetUserPasswordRequest](#model-reset-user-password-request) | `models.ModelResetUserPasswordRequest` | | ✓ | | Data for resetting your password |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [204](#reset-user-password-step1-204) | No Content | Successful operation |  | [schema](#reset-user-password-step1-204-schema) |
| [400](#reset-user-password-step1-400) | Bad Request | Bad Request |  | [schema](#reset-user-password-step1-400-schema) |
| [404](#reset-user-password-step1-404) | Not Found | Not Found |  | [schema](#reset-user-password-step1-404-schema) |
| [500](#reset-user-password-step1-500) | Internal Server Error | Internal Server Error |  | [schema](#reset-user-password-step1-500-schema) |

#### Responses


##### <span id="reset-user-password-step1-204"></span> 204 - Successful operation
Status: No Content

###### <span id="reset-user-password-step1-204-schema"></span> Schema

##### <span id="reset-user-password-step1-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="reset-user-password-step1-400-schema"></span> Schema
   
  

[ModelRequestUserPasswordResetFailure400](#model-request-user-password-reset-failure400)

##### <span id="reset-user-password-step1-404"></span> 404 - Not Found
Status: Not Found

###### <span id="reset-user-password-step1-404-schema"></span> Schema
   
  

[ModelRequestUserPasswordResetFailure404](#model-request-user-password-reset-failure404)

##### <span id="reset-user-password-step1-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="reset-user-password-step1-500-schema"></span> Schema
   
  

[ModelCommonFailure](#model-common-failure)

### <span id="reset-user-password-step2"></span> Confirm to user password reset (*reset_user_password_step2*)

```
GET /p/{confirmationKey}
```

API returns HTML-page with a message (success or error).

#### Produces
  * text/plain;charset=utf-8

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| confirmationKey | `path` | string | `string` |  | ✓ |  | Confirmation Key |
| rid | `query` | string | `string` |  | ✓ |  | RequestID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#reset-user-password-step2-200) | OK | Successful operation |  | [schema](#reset-user-password-step2-200-schema) |
| [404](#reset-user-password-step2-404) | Not Found | Not Found |  | [schema](#reset-user-password-step2-404-schema) |
| [500](#reset-user-password-step2-500) | Internal Server Error | Internal Server Error |  | [schema](#reset-user-password-step2-500-schema) |

#### Responses


##### <span id="reset-user-password-step2-200"></span> 200 - Successful operation
Status: OK

###### <span id="reset-user-password-step2-200-schema"></span> Schema

##### <span id="reset-user-password-step2-404"></span> 404 - Not Found
Status: Not Found

###### <span id="reset-user-password-step2-404-schema"></span> Schema
   
  

[ModelCommonFailure](#model-common-failure)

##### <span id="reset-user-password-step2-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="reset-user-password-step2-500-schema"></span> Schema
   
  

[ModelCommonFailure](#model-common-failure)

### <span id="revoke-refresh-token"></span> Revoke Refresh Token (Domain Action: Log Out) (*revoke_refresh_token*)

```
DELETE /v1/session
```

This request revoke the Refresh Token associated with the specified session. Thus, when the Access Token expires, then it cannot be renewed. And only after that, the user will be log out. Be aware that this query is idempotent.

#### Consumes
  * application/json;charset=utf-8

#### Security Requirements
  * authorization

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| request-id | `header` | string | `string` |  | ✓ |  | RequestID |
| model.RevokeRefreshTokenRequest | `body` | [ModelRevokeRefreshTokenRequest](#model-revoke-refresh-token-request) | `models.ModelRevokeRefreshTokenRequest` | | ✓ | | Data for revoking the Refresh Token |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [204](#revoke-refresh-token-204) | No Content | Successful operation |  | [schema](#revoke-refresh-token-204-schema) |
| [400](#revoke-refresh-token-400) | Bad Request | Bad Request |  | [schema](#revoke-refresh-token-400-schema) |
| [500](#revoke-refresh-token-500) | Internal Server Error | Internal Server Error |  | [schema](#revoke-refresh-token-500-schema) |

#### Responses


##### <span id="revoke-refresh-token-204"></span> 204 - Successful operation
Status: No Content

###### <span id="revoke-refresh-token-204-schema"></span> Schema

##### <span id="revoke-refresh-token-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="revoke-refresh-token-400-schema"></span> Schema
   
  

[ModelRevokeRefreshTokenFailure400](#model-revoke-refresh-token-failure400)

##### <span id="revoke-refresh-token-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="revoke-refresh-token-500-schema"></span> Schema
   
  

[ModelCommonFailure](#model-common-failure)

### <span id="signup-step1"></span> Create new user (*signup_step1*)

```
POST /v1/user
```

The service sends a confirmation link to the specified email. After confirmation, the service will send a password for authorization.

#### Consumes
  * application/json;charset=utf-8

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| request-id | `header` | string | `string` |  | ✓ |  | RequestID |
| model.SignUpRequest | `body` | [ModelSignUpRequest](#model-sign-up-request) | `models.ModelSignUpRequest` | | ✓ | | Data for creating a new user |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [204](#signup-step1-204) | No Content | Successful operation |  | [schema](#signup-step1-204-schema) |
| [400](#signup-step1-400) | Bad Request | Bad Request |  | [schema](#signup-step1-400-schema) |
| [404](#signup-step1-404) | Not Found | Not Found |  | [schema](#signup-step1-404-schema) |
| [409](#signup-step1-409) | Conflict | Conflict |  | [schema](#signup-step1-409-schema) |
| [500](#signup-step1-500) | Internal Server Error | Internal Server Error |  | [schema](#signup-step1-500-schema) |

#### Responses


##### <span id="signup-step1-204"></span> 204 - Successful operation
Status: No Content

###### <span id="signup-step1-204-schema"></span> Schema

##### <span id="signup-step1-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="signup-step1-400-schema"></span> Schema
   
  

[ModelSignUpFailure400](#model-sign-up-failure400)

##### <span id="signup-step1-404"></span> 404 - Not Found
Status: Not Found

###### <span id="signup-step1-404-schema"></span> Schema
   
  

[ModelCommonFailure](#model-common-failure)

##### <span id="signup-step1-409"></span> 409 - Conflict
Status: Conflict

###### <span id="signup-step1-409-schema"></span> Schema
   
  

[ModelSignUpFailure409](#model-sign-up-failure409)

##### <span id="signup-step1-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="signup-step1-500-schema"></span> Schema
   
  

[ModelCommonFailure](#model-common-failure)

### <span id="signup-step2"></span> Confirm User Email (*signup_step2*)

```
GET /u/{confirmationKey}
```

API returns HTML-page with a message (success or error).

#### Produces
  * text/plain;charset=utf-8

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| confirmationKey | `path` | string | `string` |  | ✓ |  | Confirmation Key |
| rid | `query` | string | `string` |  | ✓ |  | RequestID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#signup-step2-200) | OK | Successful operation |  | [schema](#signup-step2-200-schema) |
| [404](#signup-step2-404) | Not Found | Not Found |  | [schema](#signup-step2-404-schema) |
| [500](#signup-step2-500) | Internal Server Error | Internal Server Error |  | [schema](#signup-step2-500-schema) |

#### Responses


##### <span id="signup-step2-200"></span> 200 - Successful operation
Status: OK

###### <span id="signup-step2-200-schema"></span> Schema

##### <span id="signup-step2-404"></span> 404 - Not Found
Status: Not Found

###### <span id="signup-step2-404-schema"></span> Schema
   
  

[ModelCommonFailure](#model-common-failure)

##### <span id="signup-step2-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="signup-step2-500-schema"></span> Schema
   
  

[ModelCommonFailure](#model-common-failure)

## Models

### <span id="model-access-token-response"></span> model.AccessTokenResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| accessToken | string| `string` |  | |  |  |
| refreshToken | string| `string` |  | |  |  |
| sessionID | string| `string` |  | |  |  |



### <span id="model-common-failure"></span> model.CommonFailure


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| code | integer| `int64` |  | |  | `404` |
| message | string| `string` |  | |  | `404 Not Found` |



### <span id="model-create-access-token-failure400"></span> model.CreateAccessTokenFailure400


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| code | integer| `int64` |  | |  | `400` |
| message | string| `string` |  | |  | `BAD_PARAMETERS` |



### <span id="model-create-access-token-failure404"></span> model.CreateAccessTokenFailure404


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| code | integer| `int64` |  | |  | `404` |
| message | string| `string` |  | |  | `USER_NOT_FOUND` |



### <span id="model-create-access-token-request"></span> model.CreateAccessTokenRequest


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| clientID | string| `string` | ✓ | | User Client ID | `PWA_v0.0.1` |
| device | [ModelDevice](#model-device)| `ModelDevice` | ✓ | |  |  |
| email | string| `string` | ✓ | | User Email | `test.user@financelime.com` |
| password | string| `string` | ✓ | | User Password | `qmhVXVC1%hVNa0Hcq` |



### <span id="model-device"></span> model.Device


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| height | integer| `int64` |  | |  | `1920` |
| language | string| `string` |  | |  | `en-US` |
| platform | string| `string` |  | |  | `Linux x86_64` |
| timezone | string| `string` |  | |  | `2` |
| width | integer| `int64` |  | |  | `1060` |



### <span id="model-refresh-access-token-failure400"></span> model.RefreshAccessTokenFailure400


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| code | integer| `int64` |  | |  | `400` |
| message | string| `string` |  | |  | `BAD_PARAMETERS` |



### <span id="model-refresh-access-token-failure404"></span> model.RefreshAccessTokenFailure404


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| code | integer| `int64` |  | |  | `404` |
| message | string| `string` |  | |  | `USER_NOT_FOUND` |



### <span id="model-refresh-access-token-request"></span> model.RefreshAccessTokenRequest


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| refreshToken | string| `string` | ✓ | |  | `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJmaW5hbmNlbGltZS5jb20iLCJzdWIiOiJBdXRob3JpemF0aW9uIiwicHVycG9zZSI6InJlZnJlc2giLCJzZXNzaW9uSUQiOiI2M2IyZjUyM2ZiZGQzMzFlZjQzM2U2NmU5NDZjYWQ2OTNkOTQ5MzdjMWUxNWNjMDI5YjFiNjE1YmExN2VjZWM3IiwidXNlckRhdGEiOiJpNHhTbDBmNXcrMGJSTE1KOURMVlp3NGJDZkpSSUFqQlhoN2VFZFlpNTV2L1QvVk1EK3RmNFNyK0NSV0ZnZEpoUkh2S1AyNnZGQ1AxZ05iOU4yejljMFRoYkRZNkFSdGt2WHkzMHJ3bTlDeXh0Vk1QdUUvRXh4UDdzaCs3MGVrbE5ObjdGS2xIIiwiaWF0IjoxNjIwNTIwNTg2fQ.4fd650daddded3a74c6fcfa28559d02c3ca6f32d55805b58ac88ccd302c5445e` |



### <span id="model-request-user-password-reset-failure400"></span> model.RequestUserPasswordResetFailure400


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| code | integer| `int64` |  | |  | `400` |
| message | string| `string` |  | |  | `BAD_PARAMETERS` |



### <span id="model-request-user-password-reset-failure404"></span> model.RequestUserPasswordResetFailure404


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| code | integer| `int64` |  | |  | `404` |
| message | string| `string` |  | |  | `USER_NOT_FOUND` |



### <span id="model-reset-user-password-request"></span> model.ResetUserPasswordRequest


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| email | string| `string` | ✓ | |  | `test.user@financelime.com` |



### <span id="model-revoke-refresh-token-failure400"></span> model.RevokeRefreshTokenFailure400


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| code | integer| `int64` |  | |  | `400` |
| message | string| `string` |  | |  | `BAD_PARAMETERS` |



### <span id="model-revoke-refresh-token-request"></span> model.RevokeRefreshTokenRequest


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| sessionID | string| `string` | ✓ | |  | `f58f06a96b69083b7c4fb068faa6c8314af0636e44ecc710261abe1759b07755` |



### <span id="model-session"></span> model.Session


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| platform | string| `string` |  | |  |  |
| sessionID | string| `string` |  | |  |  |
| updatedAt | string| `string` |  | |  |  |



### <span id="model-sign-up-failure400"></span> model.SignUpFailure400


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| code | integer| `int64` |  | |  | `400` |
| message | string| `string` |  | |  | `BAD_PARAMETERS` |



### <span id="model-sign-up-failure409"></span> model.SignUpFailure409


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| code | integer| `int64` |  | |  | `409` |
| message | string| `string` |  | |  | `USER_ALREADY_EXIST` |



### <span id="model-sign-up-request"></span> model.SignUpRequest


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| email | string| `string` | ✓ | | User email | `test.user@financelime.com` |
| inviteCode | string| `string` | ✓ | | Invite code. Required depending on the setting of environment variable `AUTH_INVITE_CODE_REQUIRED`. The default is required. | `testInviteCode` |
| language | string| `string` | ✓ | | User language | `en` |



### <span id="model-version-response"></span> model.VersionResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| build | string| `string` |  | |  | `fc56bb1 [2021-05-07_11:12:09]` |
| number | string| `string` |  | |  | `v0.2.0-alpha` |


