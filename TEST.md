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

## BACK tests

Important: Launch the server first abd open a terminal.

- To test the security headers

```shell
curl -s -D - -o /dev/null http://localhost:8080/test -v
```

Notes:
`-s`: stands for "silent" mode (curl will not show progress meter or error messages)
`-D -` or `--dump-header`: tells curl to dump the headers of the response to a specified file or to standard output. `-` after `-D` indicates that the headers should be written to standard output (stdout)
`-o /dev/null` or `--output`: specifies where to write the body of the response. In this case, `/dev/null` is used, which is a special file that discards all data written to it.

- To test the max header bytes size

```shell
curl -H "X-Request-Id: $(head -c 8096 < /dev/zero | tr '\0' 'a')" http://localhost:8080/test -v
```

Note: 8096 triggers the max size even if the server is configured for a maxsizebyte of `4096` bytes...

- to test the User Agent

```shell
curl -A "Mozilla/5.0 (X11; Linux x86_64; rv:60.0) Gecko/20100101 Firefox/81.0" -v
```

- To test the error middlware

```shell
curl -X GET http://localhost:8080/panic -v
#{"code":500,"message":"An unexpected error occurred. Please try again later."}
curl -X GET http://localhost:8080/not-found -v
#{"code":404,"message":"Not Found: The requested resource could not be found."}
curl -X GET http://localhost:8080/bad-request -v
#{"code":400,"message":"Bad Request: The request could not be understood or was missing required parameters."}
```

*work in progress*

- to test the static route (for files)

```shell
 curl http://localhost:8080/static/test.txt -v
```

- to test the db operations:

- be sure to create a .env (root of the project):

```shell
SERVER_ADDR=:8080
JWT_SECRET=666
```

- run the server

```shell
cd server
go run cmd/main.go
```

Add a new user `/api/register`

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

Login `/api/login`

```shell
curl -X POST http://localhost:8080/api/login -H "Content-Type: application/json" -d '{
  "email": "testuser@example.com",
  "password": "securepassword"
}'
```

Access a protected route `/api/protected` for create, update and delete operations

```shell
curl -X GET http://localhost:8080/api/protected -H "Authorization: Bearer <token>>"
```
