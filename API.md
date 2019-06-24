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
### Register (POST /register)
Activated on `POST` request to `/register` with a POST body containing JSON data with `username` and `password` parameters, and returns json containing a web token for interaction with the API.
It will return a `409 Conflict` error if the account already exists or a `400 Bad Request` if there aren't `username` or `password` parameters in the body.

Example:
```sh
$ curl -X POST localhost:8050/register -d '{"username":"foo", "password":"bar"}'
{"jwt":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImZvbyIsImV4cCI6MTU2MTMwODg1NCwiaWF0IjoxNTYxMzA4NTU0fQ.znXFS0gIMo0K7m5NJA4I1C9Fylzh3LpzwHR8zLutkbw"}
```
### Sign In (POST /signin)
Activated on `POST` to `/signin`. It accepts a POST body containing json with the `username` and `password` parameters of a valid user, and returns json containing a web token if the user account exists and the password is correct. It will return a `401 Unauthorized` error if the username or password is incorrect or a `400 Bad Request` if there aren't `username` or `password` parameters in the body.

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
Activated on authorized `POST` to `/logout`. It requres an authorization header containing a valid token which should take the form: `Authorization: Bearer <token>`. If the token is valid, it will cause every token belonging to the user to become invalid, and the user will have to sign in again, and it will return a simple json body informing the user of the result. It will return a `401 Unauthorized` if the token is invalid.

Example:
```sh
$ curl -X POST localhost:8050/logout -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImZvbyIsImV4cCI6MTU2MTMxMDQwNSwiaWF0IjoxNTYxMzEwMTA1fQ.a7D0ri_9E1_TY7UNu697y4bXVe9czowMmOOHWXjJ2Ks"
{"result":"success"}
```
