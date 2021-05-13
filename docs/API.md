


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
| GET | /v1/u/{confirmationKey} | [confirm user email](#confirm-user-email) | Confirm User Email |
| POST | /v1/oauth/token | [create access token](#create-access-token) | Create Access Token (Domain Action: Log In) |
| POST | /v1/signup | [create new user](#create-new-user) | Create new user |
| GET | /v1/oauth/sessions | [get list active sessions](#get-list-active-sessions) | Get a list of active sessions |
| GET | /v1/version | [get version](#get-version) | Get the Service version |
| PUT | /v1/oauth/token | [refresh access token](#refresh-access-token) | Refresh Access Token (Domain Action: Renew authorization) |
| DELETE | /v1/oauth/sessions | [revoke refresh token](#revoke-refresh-token) | Revoke Refresh Token (Domain Action: Log Out) |
| POST | /v1/resetpassword | [user password reset](#user-password-reset) | Request to user password reset |
  


## Paths

### <span id="confirm-user-email"></span> Confirm User Email (*confirm_user_email*)

```
GET /v1/u/{confirmationKey}
```

API returns HTML-page with a message (success or error).

#### Produces
  * text/plain;charset=utf-8

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| confirmationKey | `path` | string | `string` |  | ✓ |  | Confirmation Key |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#confirm-user-email-200) | OK | Successful operation |  | [schema](#confirm-user-email-200-schema) |
| [404](#confirm-user-email-404) | Not Found | Not Found |  | [schema](#confirm-user-email-404-schema) |
| [500](#confirm-user-email-500) | Internal Server Error | Internal Server Error |  | [schema](#confirm-user-email-500-schema) |

#### Responses


##### <span id="confirm-user-email-200"></span> 200 - Successful operation
Status: OK

###### <span id="confirm-user-email-200-schema"></span> Schema

##### <span id="confirm-user-email-404"></span> 404 - Not Found
Status: Not Found

###### <span id="confirm-user-email-404-schema"></span> Schema
   
  

[AuthorizationAPICommonFailure](#authorization-api-common-failure)

##### <span id="confirm-user-email-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="confirm-user-email-500-schema"></span> Schema
   
  

[AuthorizationAPICommonFailure](#authorization-api-common-failure)

### <span id="create-access-token"></span> Create Access Token (Domain Action: Log In) (*create_access_token*)

```
POST /v1/oauth/token
```

Create Access Token

#### Consumes
  * application/json;charset=utf-8

#### Produces
  * application/json;charset=utf-8

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| apiCreateAccessTokenRequest | `body` | [AuthorizationAPICreateAccessTokenRequest](#authorization-api-create-access-token-request) | `models.AuthorizationAPICreateAccessTokenRequest` | | ✓ | | Data for creating a new token |

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
   
  

[AuthorizationAPIAccessTokenResponse](#authorization-api-access-token-response)

##### <span id="create-access-token-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="create-access-token-400-schema"></span> Schema
   
  

[AuthorizationAPICommonFailure](#authorization-api-common-failure)

##### <span id="create-access-token-404"></span> 404 - Not Found
Status: Not Found

###### <span id="create-access-token-404-schema"></span> Schema
   
  

[AuthorizationAPICommonFailure](#authorization-api-common-failure)

##### <span id="create-access-token-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="create-access-token-500-schema"></span> Schema
   
  

[AuthorizationAPICommonFailure](#authorization-api-common-failure)

### <span id="create-new-user"></span> Create new user (*create_new_user*)

```
POST /v1/signup
```

The service sends a confirmation link to the specified email. After confirmation, the service will send a password for authorization.

#### Consumes
  * application/json;charset=utf-8

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| apiSignUpRequest | `body` | [AuthorizationAPISignUpRequest](#authorization-api-sign-up-request) | `models.AuthorizationAPISignUpRequest` | | ✓ | | Data for creating a new user |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [204](#create-new-user-204) | No Content | Successful operation |  | [schema](#create-new-user-204-schema) |
| [400](#create-new-user-400) | Bad Request | Bad Request |  | [schema](#create-new-user-400-schema) |
| [404](#create-new-user-404) | Not Found | Not Found |  | [schema](#create-new-user-404-schema) |
| [409](#create-new-user-409) | Conflict | Conflict |  | [schema](#create-new-user-409-schema) |
| [500](#create-new-user-500) | Internal Server Error | Internal Server Error |  | [schema](#create-new-user-500-schema) |

#### Responses


##### <span id="create-new-user-204"></span> 204 - Successful operation
Status: No Content

###### <span id="create-new-user-204-schema"></span> Schema

##### <span id="create-new-user-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="create-new-user-400-schema"></span> Schema
   
  

[AuthorizationAPICommonFailure](#authorization-api-common-failure)

##### <span id="create-new-user-404"></span> 404 - Not Found
Status: Not Found

###### <span id="create-new-user-404-schema"></span> Schema
   
  

[AuthorizationAPICommonFailure](#authorization-api-common-failure)

##### <span id="create-new-user-409"></span> 409 - Conflict
Status: Conflict

###### <span id="create-new-user-409-schema"></span> Schema
   
  

[AuthorizationAPISignUpFailure409](#authorization-api-sign-up-failure409)

##### <span id="create-new-user-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="create-new-user-500-schema"></span> Schema
   
  

[AuthorizationAPICommonFailure](#authorization-api-common-failure)

### <span id="get-list-active-sessions"></span> Get a list of active sessions (*get_list_active_sessions*)

```
GET /v1/oauth/sessions
```

Get a list of active sessions

#### Produces
  * application/json;charset=utf-8

#### Security Requirements
  * authorization

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-list-active-sessions-200) | OK | Successful operation |  | [schema](#get-list-active-sessions-200-schema) |
| [401](#get-list-active-sessions-401) | Unauthorized | Unauthorized |  | [schema](#get-list-active-sessions-401-schema) |
| [404](#get-list-active-sessions-404) | Not Found | Not Found |  | [schema](#get-list-active-sessions-404-schema) |
| [500](#get-list-active-sessions-500) | Internal Server Error | Internal Server Error |  | [schema](#get-list-active-sessions-500-schema) |

#### Responses


##### <span id="get-list-active-sessions-200"></span> 200 - Successful operation
Status: OK

###### <span id="get-list-active-sessions-200-schema"></span> Schema
   
  

[][AuthorizationSession](#authorization-session)

##### <span id="get-list-active-sessions-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="get-list-active-sessions-401-schema"></span> Schema
   
  

[AuthorizationAPICommonFailure](#authorization-api-common-failure)

##### <span id="get-list-active-sessions-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-list-active-sessions-404-schema"></span> Schema
   
  

[AuthorizationAPICommonFailure](#authorization-api-common-failure)

##### <span id="get-list-active-sessions-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="get-list-active-sessions-500-schema"></span> Schema
   
  

[AuthorizationAPICommonFailure](#authorization-api-common-failure)

### <span id="get-version"></span> Get the Service version (*get_version*)

```
GET /v1/version
```

Get Version

#### Produces
  * application/json;charset=utf-8

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-version-200) | OK | Successful operation |  | [schema](#get-version-200-schema) |

#### Responses


##### <span id="get-version-200"></span> 200 - Successful operation
Status: OK

###### <span id="get-version-200-schema"></span> Schema
   
  

[SystemVersionResponse](#system-version-response)

### <span id="refresh-access-token"></span> Refresh Access Token (Domain Action: Renew authorization) (*refresh_access_token*)

```
PUT /v1/oauth/token
```

Refresh Access Token

#### Consumes
  * application/json;charset=utf-8

#### Produces
  * application/json;charset=utf-8

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| apiRefreshAccessTokenRequest | `body` | [AuthorizationAPIRefreshAccessTokenRequest](#authorization-api-refresh-access-token-request) | `models.AuthorizationAPIRefreshAccessTokenRequest` | | ✓ | | Data for refreshing the access token |

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
   
  

[AuthorizationAPIAccessTokenResponse](#authorization-api-access-token-response)

##### <span id="refresh-access-token-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="refresh-access-token-400-schema"></span> Schema
   
  

[AuthorizationAPICommonFailure](#authorization-api-common-failure)

##### <span id="refresh-access-token-404"></span> 404 - Not Found
Status: Not Found

###### <span id="refresh-access-token-404-schema"></span> Schema
   
  

[AuthorizationAPICommonFailure](#authorization-api-common-failure)

##### <span id="refresh-access-token-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="refresh-access-token-500-schema"></span> Schema
   
  

[AuthorizationAPICommonFailure](#authorization-api-common-failure)

### <span id="revoke-refresh-token"></span> Revoke Refresh Token (Domain Action: Log Out) (*revoke_refresh_token*)

```
DELETE /v1/oauth/sessions
```

This request revoke the Refresh Token associated with the specified session. Thus, when the Access Token expires, then it cannot be renewed. And only after that, the user will be log out. Be aware that this query is idempotent.

#### Consumes
  * application/json;charset=utf-8

#### Security Requirements
  * authorization

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| apiRevokeRefreshTokenRequest | `body` | [AuthorizationAPIRevokeRefreshTokenRequest](#authorization-api-revoke-refresh-token-request) | `models.AuthorizationAPIRevokeRefreshTokenRequest` | | ✓ | | Data for revoking the Refresh Token |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [204](#revoke-refresh-token-204) | No Content | Successful operation |  | [schema](#revoke-refresh-token-204-schema) |
| [400](#revoke-refresh-token-400) | Bad Request | Bad Request |  | [schema](#revoke-refresh-token-400-schema) |
| [401](#revoke-refresh-token-401) | Unauthorized | Unauthorized |  | [schema](#revoke-refresh-token-401-schema) |
| [404](#revoke-refresh-token-404) | Not Found | Not Found |  | [schema](#revoke-refresh-token-404-schema) |
| [500](#revoke-refresh-token-500) | Internal Server Error | Internal Server Error |  | [schema](#revoke-refresh-token-500-schema) |

#### Responses


##### <span id="revoke-refresh-token-204"></span> 204 - Successful operation
Status: No Content

###### <span id="revoke-refresh-token-204-schema"></span> Schema

##### <span id="revoke-refresh-token-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="revoke-refresh-token-400-schema"></span> Schema
   
  

[AuthorizationAPICommonFailure](#authorization-api-common-failure)

##### <span id="revoke-refresh-token-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="revoke-refresh-token-401-schema"></span> Schema
   
  

[AuthorizationAPICommonFailure](#authorization-api-common-failure)

##### <span id="revoke-refresh-token-404"></span> 404 - Not Found
Status: Not Found

###### <span id="revoke-refresh-token-404-schema"></span> Schema
   
  

[AuthorizationAPICommonFailure](#authorization-api-common-failure)

##### <span id="revoke-refresh-token-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="revoke-refresh-token-500-schema"></span> Schema
   
  

[AuthorizationAPICommonFailure](#authorization-api-common-failure)

### <span id="user-password-reset"></span> Request to user password reset (*user_password_reset*)

```
POST /v1/resetpassword
```

The service sends a confirmation link to the specified email. After confirmation, the service will send a new password for authorization.

#### Consumes
  * application/json;charset=utf-8

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| apiRequestUserPasswordResetRequest | `body` | [AuthorizationAPIRequestUserPasswordResetRequest](#authorization-api-request-user-password-reset-request) | `models.AuthorizationAPIRequestUserPasswordResetRequest` | | ✓ | | Data for resetting your password |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [204](#user-password-reset-204) | No Content | Successful operation |  | [schema](#user-password-reset-204-schema) |
| [400](#user-password-reset-400) | Bad Request | Bad Request |  | [schema](#user-password-reset-400-schema) |
| [404](#user-password-reset-404) | Not Found | Not Found |  | [schema](#user-password-reset-404-schema) |
| [500](#user-password-reset-500) | Internal Server Error | Internal Server Error |  | [schema](#user-password-reset-500-schema) |

#### Responses


##### <span id="user-password-reset-204"></span> 204 - Successful operation
Status: No Content

###### <span id="user-password-reset-204-schema"></span> Schema

##### <span id="user-password-reset-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="user-password-reset-400-schema"></span> Schema
   
  

[AuthorizationAPICommonFailure](#authorization-api-common-failure)

##### <span id="user-password-reset-404"></span> 404 - Not Found
Status: Not Found

###### <span id="user-password-reset-404-schema"></span> Schema
   
  

[AuthorizationAPICommonFailure](#authorization-api-common-failure)

##### <span id="user-password-reset-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="user-password-reset-500-schema"></span> Schema
   
  

[AuthorizationAPICommonFailure](#authorization-api-common-failure)

## Models

### <span id="authorization-api-access-token-response"></span> authorization.apiAccessTokenResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| accessToken | string| `string` |  | |  |  |
| refreshToken | string| `string` |  | |  |  |
| sessionID | string| `string` |  | |  |  |



### <span id="authorization-api-common-failure"></span> authorization.apiCommonFailure


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| code | integer| `int64` |  | |  | `404` |
| message | string| `string` |  | |  | `Not Found` |



### <span id="authorization-api-create-access-token-request"></span> authorization.apiCreateAccessTokenRequest


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| clientID | string| `string` | ✓ | | User Client ID | `PWA_v0.0.1` |
| device | [AuthorizationDevice](#authorization-device)| `AuthorizationDevice` | ✓ | |  |  |
| email | string| `string` | ✓ | | User Email | `test.user@financelime.com` |
| password | string| `string` | ✓ | | User Password | `qmhVXVC1%hVNa0Hcq` |



### <span id="authorization-api-refresh-access-token-request"></span> authorization.apiRefreshAccessTokenRequest


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| refreshToken | string| `string` | ✓ | |  | `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJmaW5hbmNlbGltZS5jb20iLCJzdWIiOiJBdXRob3JpemF0aW9uIiwicHVycG9zZSI6InJlZnJlc2giLCJzZXNzaW9uSUQiOiI2M2IyZjUyM2ZiZGQzMzFlZjQzM2U2NmU5NDZjYWQ2OTNkOTQ5MzdjMWUxNWNjMDI5YjFiNjE1YmExN2VjZWM3IiwidXNlckRhdGEiOiJpNHhTbDBmNXcrMGJSTE1KOURMVlp3NGJDZkpSSUFqQlhoN2VFZFlpNTV2L1QvVk1EK3RmNFNyK0NSV0ZnZEpoUkh2S1AyNnZGQ1AxZ05iOU4yejljMFRoYkRZNkFSdGt2WHkzMHJ3bTlDeXh0Vk1QdUUvRXh4UDdzaCs3MGVrbE5ObjdGS2xIIiwiaWF0IjoxNjIwNTIwNTg2fQ.4fd650daddded3a74c6fcfa28559d02c3ca6f32d55805b58ac88ccd302c5445e` |



### <span id="authorization-api-request-user-password-reset-request"></span> authorization.apiRequestUserPasswordResetRequest


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| email | string| `string` | ✓ | |  | `test.user@financelime.com` |



### <span id="authorization-api-revoke-refresh-token-request"></span> authorization.apiRevokeRefreshTokenRequest


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| sessionID | string| `string` | ✓ | |  | `f58f06a96b69083b7c4fb068faa6c8314af0636e44ecc710261abe1759b07755` |



### <span id="authorization-api-sign-up-failure409"></span> authorization.apiSignUpFailure409


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| code | integer| `int64` |  | |  |  |
| message | string| `string` |  | |  |  |



### <span id="authorization-api-sign-up-request"></span> authorization.apiSignUpRequest


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| email | string| `string` | ✓ | | User email | `test.user@financelime.com` |
| inviteCode | string| `string` | ✓ | | Invite code. Required depending on the setting of environment variable `AUTH_INVITE_CODE_REQUIRED`. The default is required. | `testInviteCode` |
| language | string| `string` | ✓ | | User language | `en` |



### <span id="authorization-device"></span> authorization.device


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| height | integer| `int64` |  | |  | `1920` |
| language | string| `string` |  | |  | `en-US` |
| platform | string| `string` |  | |  | `Linux x86_64` |
| timezone | string| `string` |  | |  | `2` |
| width | integer| `int64` |  | |  | `1060` |



### <span id="authorization-session"></span> authorization.session


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| platform | string| `string` |  | |  |  |
| sessionID | string| `string` |  | |  |  |
| updatedAt | string| `string` |  | |  |  |



### <span id="system-version-response"></span> system.versionResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| build | string| `string` |  | |  | `fc56bb1 [2021-05-07_11:12:09]` |
| number | string| `string` |  | |  | `v0.2.0-alpha` |


