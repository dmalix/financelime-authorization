# Financelime Authorization service
## RESTfull API Version 1  

- [1. Environment variables](#1-0)
- [2. Identification at the device and request level](#2-0) 
- [3. System commands](#3-0)
    - [3.1. Get the current version of the service](#3-1) 
- [4. Authorization](#4-0)
    - [4.1. Sign Up](#4-1)
    - [4.2. Confirm User Email](#4-2)
    - ...
    - [4.5. Request Access Token (Domain Action: Log In)](#4-5)
    - [4.6. Refresh Access Token (Domain Action: Renew Authorization)](#4-6)
    - [4.7. Revoke Refresh Token (Domain Action: Log Out)](#4-7)
    - [4.8. Get a list of Active Sessions](#4-8)
    
---
<a name="1-0"></a>
## 1. Environment variables

Name | Description
---|---
`DOMAIN_APP` |
`DOMAIN_API` |
`HTTP_SERVER_PORT` |
`AUTH_INVITE_CODE_REQUIRED` |
`DB_AUTH_MAIN_CONNECT_HOST` |
`DB_AUTH_MAIN_CONNECT_PORT` |
`DB_AUTH_MAIN_CONNECT_SSLMODE` |
`DB_AUTH_MAIN_CONNECT_DBNAME` |
`DB_AUTH_MAIN_CONNECT_USER` |
`DB_AUTH_MAIN_CONNECT_PASSWORD` |
`DB_AUTH_MAIN_MIGRATE_DROPFILE` |
`DB_AUTH_MAIN_MIGRATE_CREATEFILE` |
`DB_AUTH_MAIN_MIGRATE_INSERTFILE` |
`DB_AUTH_READ_CONNECT_HOST` |
`DB_AUTH_READ_CONNECT_PORT` |
`DB_AUTH_READ_CONNECT_SSLMODE` |
`DB_AUTH_READ_CONNECT_DBNAME` |
`DB_AUTH_READ_CONNECT_USER` |
`DB_AUTH_READ_CONNECT_PASSWORD` |
`DB_BLADE_CONNECT_HOST` |
`DB_BLADE_CONNECT_PORT` |
`DB_BLADE_CONNECT_SSLMODE` |
`DB_BLADE_CONNECT_DBNAME` |
`DB_BLADE_CONNECT_USER` |
`DB_BLADE_CONNECT_PASSWORD` |
`DB_BLADE_MIGRATE_DROPFILE` |
`DB_BLADE_MIGRATE_CREATEFILE` |
`SMTP_USER` |
`SMTP_PASSWORD` |
`SMTP_HOST` |
`SMTP_PORT` |
`MAIL_MESSAGE_FROM_NAME` |
`MAIL_MESSAGE_FROM_ADDR` |
`LANGUAGE_CONTENT_FILE` |
`CRYPTO_SALT` |
`JWT_SECRET_KEY` |
`JWT_SIGNING_ALGORITHM` |
`JWT_ISSUER` |
`JWT_SUBJECT` |
`JWT_ACCESS_TOKEN_LIFETIME` |
`JWT_REFRESH_TOKEN_LIFETIME` |

<a name="2-0"></a>
## 2. Identification at the device and request level
<a name="2-1"></a> 
### 2.1. Header 'request-id'
API requests must contain a `request-id` header. Its content must be generated immediately before the request, according to a certain algorithm. API service validates the header before processing some requests. The `request-id` header allows you to identify the request and, together with the` Authorization` header, quickly track the chain of events, for example, filter events in the logs by device and request.  
If a request requires a `request-id` header but is missing or failed validation, the request is rejected and a `400 Bad Request` response is returned. Since the algorithm changes slightly with each version, the `request-id` header effectively filters out inappropriate requests from bots, being a CAPTCHA for them.
 
Checks can be disabled in the API configuration for all or some requests.

#### Algorithm for generating the 'request-id' header:

```
        +------------------------- | 4 random character sets 
        |     +------------------- |          pattern:                                
        |     |     +------------- |            LNNNN                          
        |     |     |     +------- |       (1 letter + 4 numbers)
check   |     |     |     |
 sum    1     2     3     4
----- ----- ----- ----- -----
LXXXX-LNNNN-LNNNN-LNNNN-LNNNN  <-- pattern
W7000-T6755-T7700-P4010-W6778  <-- example of generated value  
=----                   =----
^                       |
|                       |
+-----------------------+                                      

LXXXX | L - The letter that is always equal to a letter from set 4
      | XXXX - This is the sum of all digits from sets 1,2,3 and 4
      |        The number of characters is always 4 (missing ones get zeros)
```  
<a name="3-0"></a>
## 3. System commands 
<a name="3-1"></a>
### 3.1. Get the current version of the service 
- Method: `GET` 
- Endpoint:  `/v1/dist`

The request will return a JSON object containing the following attributes:

Name        | Type   | Description
------------|--------|------------
`version`   | string | Service Version
`build`     | string | Service Build

#### cURL Example
```
curl -i -X GET \
-H "request-id: K7800-H7625-Z5852-N1693-K1972" \
"https://api.auth.financelime.com/v1/dist"
```
#### Request Headers
```
request-id: REQUEST_ID
```
#### Response Headers
```
status: 200
content-type:application/json;charset=utf-8
```
#### Response Body
```json
{
  "version":"v0.1.5-alpha",
  "build":"373dcab [2020-11-05_19:23:56]"
}
```
<a name="4-0"></a>
## 4. Authorization 
<a name="4-1"></a>
### 4.1. Sign up

- Method: `POST` 
- Endpoint:  `/v1/signup`

#### Parameters:
Name        | Type   | Description
------------|--------|------------
`email`       | `string` | User Email
`invite_code` | `string` | Invite Code
`language`    | `string` | User Language 

#### Response:
The service sends a confirmation link to the specified email.
After confirmation, the service will send a password for authorization.

#### cURL Example
```
curl -i -X POST \
-H "content-type:application/json;charset=utf-8" \
-H "request-id:K7800-H7625-Z5852-N1693-K1972" \
-d '{"email":"testuser@financelime.com", "inviteCode":"testInviteCode", "language":"en"}' \
"https://api.auth.financelime.com/v1/signup"
```
#### Request Headers
```
content-type: application/json;charset=utf-8
request-id: REQUEST_ID
```

##### Status Code

STATUS_CODE        | Description
-------------------|----------------
`204 Not Content`  | The API has successfully processed the request and the service sends a confirmation link to the specified email
`400 Bad Request` | The API did not accept the request because headers or parameters did not pass validation (detail in API logs)
`409 Conflict`    | An error occurred while processing the request (details in the section `Domain Error Code`, in the next paragraph)

##### Domain Error Code:

When the API returns a `409 Conflict` error status, the header `domain-error-code` contains an additional human-readable code that can help to determine the cause.

DOMAIN_ERROR_CODE | Description
---|-------------
`USER_ALREADY_EXIST` | User already exists
`INVITE_NOT_EXIST_EXPIRED` | The invite code does not exist or is expired
`INVITE_LIMIT` | The limit for issuing the specified invite code has been exhausted

<a name="4-2"></a>
### 4.2. Confirm User Email 

- Method: `GET` 
- Endpoint:  `/u/LINK_KEY`

#### Параметры:
Name        | Type       | Description
------------|------------|------------
`LINK_KEY`    | string     | confirmation key

#### Ответ:
API returns HTML-page with a message (success or error).

#### cURL Example
```
curl -X GET \
"https://u/LINK_KEY"
```
#### Response Headers
```
status: STATUS_CODE
content-type: text/plain; charset=utf-8
```
##### Status Code

STATUS_CODE        | Description
-------------------|----------------
`200 Ok`           | API returns HTML-page with a message
`404 Not Found`    | Link not found

<a name="4-5"></a>
### 4.5. Request Access Token (Domain Action: Log In)

- Method: `POST` 
- Endpoint:  `/v1/oauth/token`

#### Параметры:

Name        | Type   | Description
------------|--------|------------
`email`     | string | User Email
`password`  | string | User Password
`client_id` | string | User Client ID
`device`    | object | User Device Specification

**device object**

Name        | Type   
------------|--------
`platform`    | string 
`height`      | int     
`width`       | int    
`language`    | string 
`timezone`    | string 


#### Ответ
The request will return an object with three attributes:

Name            | Type     | Description
----------------|----------|------------
`sessionID`     | string   | Public Session ID
`accessToken`   | string   | Access Jwt-Token 
`refreshToken`  | string   | Refresh Jwt-Token 

#### cURL Example
```
curl -i -X POST \
-H "content-type:application/json;charset=utf-8" \
-H "request-id:K7800-H7625-Z5852-N1693-K1972" \
-d '{"email":"test.user@financelime.com","password":"12345", "client_id": "PWA_v0.0.1", \
"device":{"platform":"Linux x86_64", "height":1920, "width":1060, "language":"en-US", "timezone":"2"}}' \
"https://api.auth.financelime.com/v1/oauth/token"
```
#### Request Headers
```
content-type: application/json;charset=utf-8
request-id: REQUEST_ID
```
#### Response Headers
```
status: STATUS_CODE
content-type:application/json;charset=utf-8 (in case of success) 
content-type: text/plain; charset=utf-8 (in case of error)
```
##### Status Code

Status Code        | Description
-------------------|----------------
`200 OK`           | API successfully processed the request, returned JWT tokens and sent a notification to the email account
`401 Bad Request`  | The API did not accept the request because headers or parameters did not pass validation (detail in API logs)
`404 Not Found`    | Account not found

#### Response Body
```
{
   "sessionID":"46512f3a52cf8e6ddbe8ef34573f7b954086f8714b6a96f1df57234df8ee3735",
   "accessToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJGaW5hbmNlbGltZS5jb20iLCJzdWIiOiJBdXRob3JpemF0aW9uIiwicHVycG9zZSI6ImFjY2VzcyIsImlkIjoiMTZkN2RiNTM3MjQ3ZWFmMTEzZjRjOGFkNTllOWEyYTU4OWNlN2NhZjYxMzViY2Q3YmZlYzBiNTI1YWY0OGEyYSIsImlhdCI6MTU5NjgzNTM5OX0.d68bea3232f10c60483f838fff8d8c66661cb497b141213c9a006be2e7c9d723",
   "refreshToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJGaW5hbmNlbGltZS5jb20iLCJzdWIiOiJBdXRob3JpemF0aW9uIiwicHVycG9zZSI6InJlZnJlc2giLCJpZCI6IjE2ZDdkYjUzNzI0N2VhZjExM2Y0YzhhZDU5ZTlhMmE1ODljZTdjYWY2MTM1YmNkN2JmZWMwYjUyNWFmNDhhMmEiLCJpYXQiOjE1OTY4MzUzOTl9.b88345d361482865f1a7af533d41d66e922dcca4c76c2d4b1fcfa65616679471"
}
```

<a name="4-6"></a>
### 4.6. Refresh Access Token (Domain Action: Renew Authorization)

- Method: `PUT` 
- Endpoint:  `/v1/oauth/token`

#### Параметры:

Name            | Type   | Description
----------------|--------|------------
refreshToken    | string | Refresh Token

#### Ответ
The request will return an object with three attributes:

Name            | Type     | Description
----------------|----------|------------
sessionID       | string   | Session ID
accessToken     | string   | Access Token 
refreshToken    | string   | Refresh Token 

#### cURL Example
```
curl -i -X PUT \
-H "content-type:application/json;charset=utf-8" \
-H "request-id:K7800-H7625-Z5852-N1693-K1972" \
-d '{"refreshToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJGaW5hbmNlbGltZS5jb20iLCJzdWIiOiJBdXRob3JpemF0aW9uIiwicHVycG9zZSI6InJlZnJlc2giLCJpZCI6IjE2ZDdkYjUzNzI0N2VhZjExM2Y0YzhhZDU5ZTlhMmE1ODljZTdjYWY2MTM1YmNkN2JmZWMwYjUyNWFmNDhhMmEiLCJpYXQiOjE1OTY4MzUzOTl9.b88345d361482865f1a7af533d41d66e922dcca4c76c2d4b1fcfa65616679471"}' \
"https://api.auth.financelime.com/v1/oauth/token"
```
#### Request Headers
```
content-type: application/json;charset=utf-8
request-id: REQUEST_ID
```
#### Response Headers
```
status: STATUS_CODE
content-type:application/json;charset=utf-8 (in case of success) 
content-type: text/plain; charset=utf-8 (in case of error)
```
##### Status Code

Status Code        | Description
-------------------|----------------
200 OK             | API successfully processed the request, returned JWT tokens
401 Bad Request    | The API did not accept the request because headers or parameters did not pass validation (detail in API logs)
404 Not Found      | Invalid or expired `Refresh Token` or an account not found or deleted  

#### Response Body

```
{
   "sessionID":"16d7db537247eaf113f4c8ad59e9a2a589ce7caf6135bcd7bfec0b525af48a2a",
   "accessToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJGaW5hbmNlbGltZS5jb20iLCJzdWIiOiJBdXRob3JpemF0aW9uIiwicHVycG9zZSI6ImFjY2VzcyIsImlkIjoiMTZkN2RiNTM3MjQ3ZWFmMTEzZjRjOGFkNTllOWEyYTU4OWNlN2NhZjYxMzViY2Q3YmZlYzBiNTI1YWY0OGEyYSIsImlhdCI6MTU5NjgzNTM5OX0.d68bea3232f10c60483f838fff8d8c66661cb497b141213c9a006be2e7c9d723",
   "refreshToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJGaW5hbmNlbGltZS5jb20iLCJzdWIiOiJBdXRob3JpemF0aW9uIiwicHVycG9zZSI6InJlZnJlc2giLCJpZCI6IjE2ZDdkYjUzNzI0N2VhZjExM2Y0YzhhZDU5ZTlhMmE1ODljZTdjYWY2MTM1YmNkN2JmZWMwYjUyNWFmNDhhMmEiLCJpYXQiOjE1OTY4MzUzOTl9.b88345d361482865f1a7af533d41d66e922dcca4c76c2d4b1fcfa65616679471"
}
```

<a name="4-7"></a>
### 4.7. Revoke Refresh Token (Domain Action: Log Out)

This request revoke the `Refresh Token` associated with the specified session. Thus, when the `Access Token` expires, it cannot be renewed. And only after that, the user will be log out. Be aware that this query is idempotent.

- Method: `DELETE` 
- Endpoint:  `/v1/oauth/sessions`

#### Параметры:

Name         | Type   | Description
-------------|--------|------------
`sessionID`    | `string` | Here you need to specify the `sessionID`, the `refresh token` of which you want to cancel. If no `sessionID` is specified, then the current session `refresh token` will be canceled. 

#### cURL Example
```
curl -i -X DELETE \
-H "content-type:application/json;charset=utf-8" \
-H "request-id:K7800-H7625-Z5852-N1693-K1972" \
-H "authorization:bearer ACCESS_TOKEN" \
-d '{"sessionID":"870bd06be766720b7348f6baf946355b71d23401978f7199b8437f52377f62e1"}' \
"https://api.auth.financelime.com/v1/oauth/sessions"
```
#### Request Headers
```
content-type: application/json;charset=utf-8
request-id: REQUEST_ID
authorization: ACCESS_TOKEN
```
#### Response Headers
```
status: STATUS_CODE
content-type: text/plain; charset=utf-8
```
##### Status Code

Status Code        | Description
-------------------|----------------
204 No Content     | API successfully processed the request
401 Bad Request    | The API did not accept the request because headers or parameters did not pass validation (detail in API logs)

<a name="4-8"></a>
### 4.8. Get a list of active sessions

- Method: `GET` 
- Endpoint:  /v1/oauth/sessions

#### cURL Example
```
curl -i -X GET \
-H "request-id:K7800-H7625-Z5852-N1693-K1972" \
-H "authorization:bearer ACCESS_TOKEN" \
"https://api.financelime.com/v1/oauth/sessions"
```
#### Request Headers
```
request-id: REQUEST_ID
authorization: ACCESS_TOKEN
```
#### Response Headers
```
status: 200
content-type: application/json;charset=utf-8
```
##### Status Code

Status Code        | Description
-------------------|----------------
200 OK             | API successfully processed the request
400 Bad Request    | The API did not accept the request because headers did not pass validation (detail in API logs)
403 Forbidden      | API did not accept request because authorization is required

#### Response Body

```json
[
  {
    "sessionID":"8d8e1a32b7c349ee306eb3ec2b82e1fd37a97eba0903c66f061e5fdc774067f0",
    "updatedAt":"2020-08-23T17:59:46.558594Z",
    "platform":"Linux x86_64"
  },
  { 
    "sessionID":"828d96230888e9aad01f874212a5a8abf3a74ec6bbd44272bb112d0727418c38",
    "updatedAt":"2020-08-23T18:01:16.2457Z",
    "platform":"iPhone"
  },
  {
    "sessionID":"904eebcbcb7a9911b934d5ba0ac248247d66e6e7fb96dbaeff1ebf09073b7c4c",
    "updatedAt":"2020-08-23T18:01:32.233222Z",
    "platform":"Android"
  }
]
```