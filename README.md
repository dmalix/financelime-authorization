## Financelime REST API service   

- [1. Environment variables](#1-0)
- [2. Identification at the device and request level](#2-0) 
- [3. Service version and state](#3-0)
- [4. Authorization](#4-0)
    - [4.1. Sign up](#4-1)
    
---
<a name="1-0"></a>
## 1. Environment variables

Name | Description
---|---
HTTP_SERVER_PORT |
AUTH_INVITE_CODE_REQUIRED |
DB_AUTH_MAIN_CONNECT_HOST |
DB_AUTH_MAIN_CONNECT_PORT |
DB_AUTH_MAIN_CONNECT_SSLMODE |
DB_AUTH_MAIN_CONNECT_DBNAME |
DB_AUTH_MAIN_CONNECT_USER |
DB_AUTH_MAIN_CONNECT_PASSWORD |
DB_AUTH_MAIN_MIGRATE_DROPFILE |
DB_AUTH_MAIN_MIGRATE_CREATEFILE |
DB_AUTH_MAIN_MIGRATE_INSERTFILE |
DB_AUTH_READ_CONNECT_HOST |
DB_AUTH_READ_CONNECT_PORT |
DB_AUTH_READ_CONNECT_SSLMODE |
DB_AUTH_READ_CONNECT_DBNAME |
DB_AUTH_READ_CONNECT_USER |
DB_AUTH_READ_CONNECT_PASSWORD |
DB_BLADE_CONNECT_HOST |
DB_BLADE_CONNECT_PORT |
DB_BLADE_CONNECT_SSLMODE |
DB_BLADE_CONNECT_DBNAME |
DB_BLADE_CONNECT_USER |
DB_BLADE_CONNECT_PASSWORD |
DB_BLADE_MIGRATE_DROPFILE |
DB_BLADE_MIGRATE_CREATEFILE |

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
## 3. Service version and state 
<a name="3-1"></a>
### 3.1. Get the current version of the REST API service 
- Method: `GET` 
- Endpoint:  `/dist`

The request will return a JSON object containing the following attributes:

Name        | Type   | Description
------------|--------|------------
version     | string | Service Version
build       | string | Service Build

#### cURL Example
```
curl -i -X GET \
-H "request-id: K7800-H7625-Z5852-N1693-K1972" \
"https://api.dev.financelime.com/dist"
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
   "version":"0.4.0",
   "build":"202007162345" 
}
```
<a name="3-2"></a>
### 3.2. Get the current state of a REST API service 
- Method: `GET` 
- Endpoint:  `/status`

The request will return a JSON object containing the following attributes:

Name        | Type   | Description
------------|--------|------------
api         | bool   | Status of API
db          | bool   | Status of DB

#### cURL Example
```
curl -i -X GET \
-H "request-id: K7800-H7625-Z5852-N1693-K1972" \
"https://api.dev.financelime.com/status"
```
#### Request Headers
```
content-type:application/json
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
   "api":true,
   "db":true
}

```
<a name="4-0"></a>
## 4. Authorization 
<a name="4-1"></a>
### 4.1. Sign up

- Method: `POST` 
- Endpoint:  `/authorization/signup`

#### Parameters:
Name        | Type   | Description
------------|--------|------------
email       | string | User Email
invite_code | string | Invite Code
language    | string | User Language 

#### Response:
The service sends a confirmation link to the specified email.
After confirmation, the service will send a password for authorization.

#### cURL Example
```
curl -i -X POST \
-H "content-type:application/json;charset=utf-8" \
-H "request-id:K7800-H7625-Z5852-N1693-K1972" \
-d '{"email":"testuser@financelime.com", "inviteCode":"testInviteCode", "language":"en"}' \
"https://api.dev.financelime.com/authorization/signup"
```
#### Request Headers
```
content-type: application/json;charset=utf-8
request-id: REQUEST_ID
```

##### Status Code

STATUS_CODE        | Description
-------------------|----------------
202 Accepted       | The API has successfully processed the request and the service sends a confirmation link to the specified email
400 Bad Request    | The API did not accept the request because headers or parameters did not pass validation (detail in API logs)
409 Conflict       | An error occurred while processing the request (details in the section `Domain Error Code`, in the next paragraph)

In the case of a 409 code, the request will return an `Error-Code` header containing the error code.

##### Domain Error Code:

When the API returns a `409 Conflict` error status, the header `domain-error-code` contains an additional code that can help to determine the cause.

DOMAIN_ERROR_CODE | Description
---|-------------
b1 | User already exists
b2 | The invite code does not exist or is expired
b3 | The limit for issuing the specified invite code has been exhausted

