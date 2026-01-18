# Chirp API - [Boot.dev](https://boot.dev) Project

## Introduction

**Chirp API** is a simple RESTful service for managing users, authentication, and short messages (“chirps”), built as part of the [Boot.dev](https://boot.dev) backend project. It demonstrates core backend concepts such as:

- User registration and authentication (JWT + refresh tokens)
- CRUD operations for Chirp messages
- Basic admin metrics and database reset endpoints
- Input validation and error handling

This project is ideal for learning modern TypeScript backend practices, including database integration, authentication, and API design.

## Installation

Follow these steps to get the project running locally:

1. Clone the repository and install

```bash
git clone https://github.com/lucasgjanot/go-http-server.git
cd go-http-server
go install .
```

2. Set environment variables

Create a .env.development file in the project root:

```bash
POSTGRES_PASSWORD=<password>
POSTGRES_HOST=<host>
POSTGRES_USER=<user>
POSTGRES_DB=<database_name>
POSTGRES_PORT=<database_port>
DB_URL=<database_connection_string>
JWT_SECRET=<secret>
POLKA_KEY=<ApiKey>
```

3. Init database and run migrations

```bash
make service-up
```

4. Start the server

```bash
go-http-server
```

The server should now be running at http://localhost:8080. You can test the endpoints using a tool like Postman or curl.

## API Documentation

### `User` resource

```json
{
  "id": "40ad210f-a739-4f6b-bb57-f401c06c7a35",
  "email": "username@example.com",
  "createdAt": "2025-10-27T19:36:41.669Z",
  "updatedAt": "2025-10-27T19:36:41.669Z",
  "isChirpyRed": false
}
```

#### GET /api/users

Returns an array of `Users`.

#### GET /api/users/:userId

Return `User` with `userId` if it exists, otherwise returns `HTTP Error 404`.

#### POST /api/users

Creates a new `User`.

**Request body:**

```json
{
  "email": "username@example.com",
  "password": "strongP4sw0rd!"
}
```

**Response body:** `User`

```json
{
  "id": "40ad210f-a739-4f6b-bb57-f401c06c7a35",
  "email": "username@example.com",
  "createdAt": "2025-10-28T22:57:00.627Z",
  "updatedAt": "2025-10-28T22:57:00.627Z",
  "isChirpyRed": false
}
```

#### PUT /api/user

##### Authentication

Requires `Authorization Header` with a valid `Bearer Token`.

```makefile
Authorization: Bearer <valid_token>
```

**Request body:**

```json
{
  "email": "newusername@example.com",
  "password": "NEWstrongP4sw0rd!"
}
```

**Response body:** `User`

```json
{
  "id": "40ad210f-a739-4f6b-bb57-f401c06c7a35",
  "email": "newusername@example.com",
  "createdAt": "2025-10-28T22:57:00.627Z",
  "updatedAt": "2025-10-28T22:57:00.627Z",
  "isChirpyRed": false
}
```

### `Auth` resource

#### POST /api/login

Authenticates an existing `User` and returns `User` with access and refresh tokens.

**Request body:**

```json
{
  "email": "newusername@example.com",
  "password": "NEWstrongP4sw0rd!"
}
```

**Response body:** `User`

```json
{
  "id": "40ad210f-a739-4f6b-bb57-f401c06c7a35",
  "email": "newusernamee@examle.com",
  "createdAt": "2025-10-27T12:46:21.707Z",
  "updatedAt": "2025-10-29T02:15:51.125Z",
  "isChirpyRed": false,
  "token": "<token>",
  "refreshToken": "<refresh_token>"
}
```

#### POST /api/refresh

Returns new `token` and new `refreshToken` for `User`

##### Authentication

Requires `Authorization Header` with a valid `Bearer Refresh Token`.

```makefile
Authorization: Bearer <valid_refresh_token>
```

**Response body:**

```json
{
  "token": "<new_token>",
  "refreshToken": "<new_refresh_token>"
}
```

#### POST /api/revoke

##### Authentication

Requires `Authorization Header` with a valid `Bearer Refresh Token`.

```makefile
Authorization: Bearer <valid_refresh_token>
```

**Response**: HTTP 204 on Success.

### `Chirp` resource

```json
{
  "id": "831e1728-ff0c-48b5-a4f4-57cc0266d219",
  "createdAt": "2025-10-28T23:44:44.949Z",
  "updatedAt": "2025-10-28T23:44:44.949Z",
  "body": "expemple chirp",
  "userId": "737b3766-fc7f-4040-b42e-32d6393eaee9"
}
```

#### GET /api/chirps

Returns a array of `Chirps`
You can optionally filter by author using the authorId query parameter.

| Parameter  | Type   | Required | Description                                                     |
| ---------- | ------ | -------- | --------------------------------------------------------------- |
| `authorId` | string | ❌ No    | Filters the chirps to only those created by the specified user. |

##### Example

```http
GET /api/chirps?authorId=737b3766-fc7f-4040-b42e-32d6393eaee9 HTTP/1.1 # Returns only this user chirps
GET /api/chirps HTTP/1.1  # Returns all chirps
```

#### GET /api/chirps/:chirpId

Return `Chirp` with `chirpId` if it exists, otherwise returns `HTTP Error 404`.

#### POST /api/chirps

Creates a new `Chirp`.

##### Authentication

Requires `Authorization Header` with a valid `Bearer Token`.

```makefile
Authorization: Bearer <valid_token>
```

**Request body:**

```json
{
  "body": "exemple"
}
```

**Response body:**

```json
{
  "id": "831e1728-ff0c-48b5-a4f4-57cc0266d219",
  "body": "exemple",
  "userId": "737b3766-fc7f-4040-b42e-32d6393eaee9"
}
```

#### DELETE /api/chirps/:chirpId

Delete `Chirp` with `chirpId` if exists, otherwise returns HTTP Error 404.

##### Authentication

Requires `Authorization Header` with a valid `Bearer Token`.

```makefile
Authorization: Bearer <valid_token>
```

### `Admin` Resource

#### GET /admin/metrics

Return server hits

#### POST /admin/reset

Delete all resources from database and set hits count to 0 (Only works at dev platform)

### `Readiness` Resource

#### GET /api/healthz

Returns a simple response indicating that the server is running and reachable.
