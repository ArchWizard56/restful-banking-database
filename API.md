# API Documentation
## Basic Info
The API's authorization is done with JSON web tokens that expire every 5 minutes.
### Index (GET /)
Activated on `GET` to `/`, and simply returns "Placeholder"

Example:
```sh
$ curl localhost:8050
Placeholder
```
## User Account Routes
Routes that directly deal with authentication and account creation.
### Register (POST /register)
Activated on `POST` request to `/register` with a POST body containing JSON data with `username` and `password` parameters, and returns json containing a web token for interaction with the API.
It will return a `409 Conflict` error if the account already exists or a `400 Bad Request` if there aren't `username` or `password` parameters in the body.

Example:
```sh
$ curl -X POST localhost:8050/register -d '{"username":"foo", "password":"bar"}'
{"jwt":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImZvbyIsImV4cCI6MTU2MTMwODg1NCwiaWF0IjoxNTYxMzA4NTU0fQ.znXFS0gIMo0K7m5NJA4I1C9Fylzh3LpzwHR8zLutkbw"}
```
### Sign In (POST /signin)
Activated on `POST` to `/signin`. It accepts a POST body containing JSON with the `username` and `password` parameters of a valid user, and returns JSON containing a web token if the user account exists and the password is correct. It will return a `401 Unauthorized` error if the username or password is incorrect or a `400 Bad Request` if there aren't `username` or `password` parameters in the body.

Example:
```sh
$ curl -X POST localhost:8050/signin -d '{"username":"foo", "password":"bar"}'
{"jwt":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImZvbyIsImV4cCI6MTU2MTMwOTUxNywiaWF0IjoxNTYxMzA5MjE3fQ.88_T5HijoXF2etpaivu4YusCJ5Po3dEZ74QuwRG16FM"}
```
### Refresh (GET /refresh)
Activated on authorized `GET` to `/`. It requres an authorization header containing a valid token which should take the form: `Authorization: Bearer <token>`. If the token is valid, it will return another token that has the expiry date updated. It will return a `401 Unauthorized` if the token is invalid because it is expired or wrong.

Example:
```sh
$ curl localhost:8050/refresh -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImZvbyIsImV4cCI6MTU2MTMwOTUxNywiaWF0IjoxNTYxMzA5MjE3fQ.88_T5HijoXF2etpaivu4YusCJ5Po3dEZ74QuwRG16FM"
{"jwt":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImZvbyIsImV4cCI6MTU2MTMwOTgwOCwiaWF0IjoxNTYxMzA5NTA4fQ.5CmyWOAGMQTPLrkI8_c5ZVocJYzRVRo0WqBUtByfU_E"}
```
### Logout (POST /logout)
Activated on authorized `POST` to `/logout`. It requres an authorization header containing a valid token which should take the form: `Authorization: Bearer <token>`. If the token is valid, it will cause every token belonging to the user to become invalid, and the user will have to sign in again, and it will return a simple JSON body informing the user of the result. It will return a `401 Unauthorized` if the token is invalid.

Example:
```sh
$ curl -X POST localhost:8050/logout -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImZvbyIsImV4cCI6MTU2MTMxMDQwNSwiaWF0IjoxNTYxMzEwMTA1fQ.a7D0ri_9E1_TY7UNu697y4bXVe9czowMmOOHWXjJ2Ks"
{"result":"success"}
```
## Bank Account Routes
Routes that deal with accounts and their balances
### Accounts (GET /accounts)
Activated on authorized`GET` to `/accounts`. It requres an authorization header containing a valid token which should take the form: `Authorization: Bearer <token>`. If the token is valid, it will return a JSON array of all accounts and balances owned by the user. It will return a `401 Unauthorized` if the token is invalid.

Example:
```sh
$ curl localhost:8050/accounts -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImZvbyIsImV4cCI6MTU2MTMxMTI3NSwiaWF0IjoxNTYxMzEwOTc1fQ.ID3_vyvbYJvnNN4TNcURl4Deex2WixCYA9yG9H58zrc"
[{"Username":"foo","Number":"11198","CcBal":0,"DcBal":0,"ArBal":10}]
```

### Open Account (POST /openaccount)
Activated on authorized `POST` to `/openaccount.`It requres an authorization header containing a valid token which should take the form: `Authorization: Bearer <token>`. If the token is valid, it will create a sub account under the users name with a different account number. It will return JSON data about the created account. It will return a `401 Unauthorized` if the token is invalid.

Example:
```sh
$ curl -X POST localhost:8050/openaccount -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImZvbyIsImV4cCI6MTU2MTMxMTUzOCwiaWF0IjoxNTYxMzExMjM4fQ.5q4xtJ-KJACyG3-mEWPldO0kKyNkm-E3CP13b5x4a9o"
{"Username":"foo","Number":"25531","CcBal":0,"DcBal":0,"ArBal":0}
```
### Transfer (POST /transfer)
Activated on authorized `POST` to `/transfer`. It requres an authorization header containing a valid token which should take the form: `Authorization: Bearer <token>`. It also requires a transfer encoded in the body. The transfer takes the form of JSON with the following parameters: `fromaccount` with the value of a string tencoding the account number that the balance will be transfered from; `toaccount` with a string encoding the account receiving the transfer; `amount` with a non-negative integer representing the amount to transfer between the accounts; and `type` which encodes a string which represents which kind of balance to transfer which is either `"ArBal"`,`"CcBal"`, or `"DcBal"`. In sum, it should look like this:
```json
{
  "fromaccount": "<account number to transfer from>",
  "toaccount": "<account number to transfer to>",
  "amount": <transferamount>,
  "type": "<transfertype>"
}
```
It will return a simple JSON body informing the user of the result. It will return a `403 Forbidden` error if the user attempts to transfer from an unowned account, a `401 Unauthorized` if the token is invalid, a `400 Bad Request` if either account number is wrong, the transfer amount is too much, or a negative number, or if the type is wrong. It will return a `500 Internal Server Error` if the sql statements error. 
Example:
```sh
$ curl -X POST localhost:8050/transfer -d '{
  "fromaccount": "11198",
  "toaccount": "25531",
  "amount": 1,
  "type": "ArBal"
}' -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImZvbyIsImV4cCI6MTU2MTMyMzIyMSwiaWF0IjoxNTYxMzIyOTIxfQ.d0AMqhzwljnCntHSJB6H93qMnWFu1HdGMP1UciQTC5A"
{"result":"success"}
```
