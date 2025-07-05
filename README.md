# Information
**auth-service** - A pet-project written in Go that implements user registration and authentication.

The service provides functions such as:

* Issuance of JWTs tokens
* Hashing passwords
* Sign in with email, password
* Using PostgreSQL as a database


# Launch
1. Install PostgreSQL
2. Create environment variables with values like:

    * DB_HOST="localhost"
    * DB_PORT="5432"
    * DB_USER="test_user"
    * DB_PASSWORD="password"
    * DB_NAME="test_user"
    * DB_SSLMODE="disable"
    * SECRET="secret12345" 
3. Clone this repository
4. Build the auth-service binary: `make build`. You should see an output like this:
```
go build -C ./cmd -ldflags "-X main.Version=0.1.0 -X main.BuildTime=2025-07-05T12:46:06+0300" -o ./bin/auth
```
5. Execute the auth-service binary: `./cmd/bin/auth`

# Endpoints
**POST /register**

Register a new user with username, email and password
```
{
  "username": "Alex",  
  "email": "alex@example.com",
  "password": "12345678"
}
```
Response:
```
{
    "id": "e535b42e-7884-41e3-a18a-091dce9ef238",
    "username": "Alex",
    "email": "alex@example.com",
    "created_at": "2025-07-05T14:29:20.238934047+03:00",
    "updated_at": "2025-07-05T14:29:20.238934047+03:00"
}
```
**POST /login**

Login with email and password
```
{    
    "email": "alex@example.com",
    "password": "12345678"
}
```
Response:
```
{
    "token": "eyJhbGciOiJIUzI1...ZePZNHfBk",
    "user": {
        "id": "c5b520c4-cea6-4693-aee4-1e9ace519c84",
        "username": "Alex",
        "email": "alex@example.com",
        "created_at": "2025-06-01T12:14:26.465041+03:00",
        "updated_at": "2025-06-01T12:14:26.465041+03:00"
    }
}
```

**GET /validate**

Validate token
```
headers: {
  "Authorization" : "Bearer eyJhbGciOiJIUzI1...ZePZNHfBk"
}
```
Response:

```
{
    valid token
}
```