## Financelime REST API service   

- [1. Environment variables](#1-0)
- [2. Identification at the device and request level](#2-0) 
- [3. Service version](#3-0)
    - [3.1. Get the current version of the REST API service](#3-1) 
- [4. Authorization](#4-0)
    - [4.1. Sign up](#4-1)
    - [4.2. Confirm user email](#4-2)
    - [4.5. Request Access token](#4-5)
    
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
## 2. Identification at the device and request level.
<a name="2-1"></a> 
### 2.1. Header 'request-id'.
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
### 3.1. Get the current version of the REST API service. 
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
"https://api.dev.financelime.com/system/dist"
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
- Endpoint:  `/authorization/signup`

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
`202 Accepted`    | The API has successfully processed the request and the service sends a confirmation link to the specified email
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
### 4.2. Confirm user email 

- Method: `GET` 
- Endpoint:  `/acue/LINK_KEY`

#### Параметры:
Name        | Type       | Description
------------|------------|------------
LINK_KEY    | string     | confirmation key

#### Ответ:
API returns HTML-page with a message (success or error).

#### cURL Example
```
curl -X GET \
"https://c/LINK_KEY"
```
#### Response Headers
```
status: STATUS_CODE
content-type: text/plain; charset=utf-8
```
##### Status Code

STATUS_CODE        | Description
-------------------|----------------
200 Ok             | API returns HTML-page with a message
404 Not Found      | Link not found

<a name="4-5"></a>
### 4.5. Request Access token 

- Method: `POST` 
- Endpoint:  `/authorization/oauth/token/request`

#### Параметры:

Name        | Type   | Description
------------|--------|------------
email       | string | User Email
password    | string | User Password
client_id   | string | User Client ID
device      | object | User Device Specification

** device object **

Name        | Type   
------------|--------
platform    | string 
height      | int     
width       | int    
language    | string 
timezone    | string 


#### Ответ
The request will return an object with three attributes:

Name            | Type     | Description
----------------|----------|------------
accessToken     | string   | Access Jwt-Token 
refreshToken    | string   | Refresh Jwt-Token 

#### cURL Example
```
curl -i -X POST \
-H "content-type:application/json;charset=utf-8" \
-H "request-id:K7800-H7625-Z5852-N1693-K1972" \
-d '{"email":"max@apivox.com","password":"12345", client_id: "PWA_v0.0.1", \
device":{"platform":"Linux x86_64", "height":1920, "width":1060, "language":"en-US", "timezone":"2"}}' \
"https://api.financelime.com/authorization/oauth/token/request"
```
#### Request Headers
```
content-type: application/json;charset=utf-8
request-id: REQUEST_ID
```
#### Response Headers
```
status: STATUS_CODE
content-type:application/json;charset=utf-8 (в случае успеха) 
content-type: text/plain; charset=utf-8 (в случае ошибки)
```
##### Status Code

Status Code        | Description
-------------------|----------------
200 OK             | API successfully processed the request, returned JWT tokens and sent a notification to the email account
401 Bad Request    | The API did not accept the request because headers or parameters did not pass validation (detail in API logs)
404 Not Found      | Account not found

#### Response Body
```
{
   "accessToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJGaW5hbmNlbGltZS5jb20iLCJzdWIiOiJBdXRob3JpemF0aW9uIiwicHVycG9zZSI6ImFjY2VzcyIsImlkIjoiMTZkN2RiNTM3MjQ3ZWFmMTEzZjRjOGFkNTllOWEyYTU4OWNlN2NhZjYxMzViY2Q3YmZlYzBiNTI1YWY0OGEyYSIsImlhdCI6MTU5NjgzNTM5OX0.d68bea3232f10c60483f838fff8d8c66661cb497b141213c9a006be2e7c9d723",
   "refreshToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJGaW5hbmNlbGltZS5jb20iLCJzdWIiOiJBdXRob3JpemF0aW9uIiwicHVycG9zZSI6InJlZnJlc2giLCJpZCI6IjE2ZDdkYjUzNzI0N2VhZjExM2Y0YzhhZDU5ZTlhMmE1ODljZTdjYWY2MTM1YmNkN2JmZWMwYjUyNWFmNDhhMmEiLCJpYXQiOjE1OTY4MzUzOTl9.b88345d361482865f1a7af533d41d66e922dcca4c76c2d4b1fcfa65616679471"
}
```

