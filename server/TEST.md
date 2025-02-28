# TESTS

## Usage

Create a .env at the root of the project

```shell
SERVER_ADDR=:8080
JWT_SECRET=666
```

Run the server

```shell
cd server
go run cmd/main.go
```

## BACK

to test the error middlware

Launch the server first abd open a terminal. 

```shell
curl -X GET http://localhost:8080/panic
#{"code":500,"message":"An unexpected error occurred. Please try again later."}
curl -X GET http://localhost:8080/not-found
#{"code":404,"message":"Not Found: The requested resource could not be found."}
curl -X GET http://localhost:8080/bad-request
#{"code":400,"message":"Bad Request: The request could not be understood or was missing required parameters."}
```

*work in progress*

to test the db operations:

- create a JWT token:
```shell
export JWT_SECRET="your secret"
```
- run the server
```shell
cd server
go run cmd/main.go
```

Add a new user

```shell
curl -X POST http://localhost:8080/api/register -H "Content-Type: application/json" -d '{
  "nickname": "testuser",
  "age": 23,
  "gender": "male",
  "first_name": "Test",
  "last_name": "User",
  "email": "testuser@example.com",
  "password": "securepassword",
  "date_register": "'$(date -u +"%Y-%m-%dT%H:%M:%SZ")'"
}'
```

Login

```shell
curl -X POST http://localhost:8080/api/login -H "Content-Type: application/json" -d '{
  "email": "testuser@example.com",
  "password": "securepassword"
}'
```

Access a protected route

```shell
curl -X GET http://localhost:8080/api/protected -H "Authorization: Bearer <token>>"
```
