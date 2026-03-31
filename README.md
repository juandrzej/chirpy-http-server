# 🐦 Chirpy HTTP Server

A lightweight social media REST API built in **Go**, where users can post short messages called *chirps*. Built as a capstone project for the [Boot.dev](https://boot.dev) backend development course, Chirpy covers the full stack of backend fundamentals: routing, authentication, database persistence, and JWT-based authorization.

---

## What Is Chirpy?

Chirpy is a Twitter-like micro-posting platform exposed entirely through a JSON REST API. Users can register, log in, and post chirps — short messages capped at 140 characters. The server enforces content policy (automatically censoring a list of forbidden words), protects routes with JWT access tokens, and uses refresh tokens for seamless re-authentication.

---

## Why Should You Care?

This project is a great reference for anyone learning backend development in Go. It demonstrates:

- **Clean project structure** — handlers are split by resource (`handlers_chirps.go`, `handlers_users.go`, `handlers_refresh.go`)
- **Database-first design** — SQL queries generated with [sqlc](https://sqlc.dev/), backed by PostgreSQL
- **Proper auth flow** — bcrypt password hashing, short-lived JWTs (1 hour), and long-lived refresh tokens (60 days)
- **RESTful conventions** — correct HTTP methods, status codes, and JSON response shapes throughout

---

## Installation & Setup

### Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [PostgreSQL](https://www.postgresql.org/)
- [goose](https://github.com/pressly/goose) (for running migrations)

### 1. Clone the repository

```bash
git clone https://github.com/juandrzej/chirpy-http-server.git
cd chirpy-http-server
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Set up environment variables

Create a `.env` file in the project root:

```env
DB_URL=postgres://username:password@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
SECRET=your-super-secret-jwt-signing-key
```

| Variable   | Description                                                                 |
|------------|-----------------------------------------------------------------------------|
| `DB_URL`   | PostgreSQL connection string                                                |
| `PLATFORM` | Set to `dev` to enable the admin reset endpoint; any other value disables it |
| `SECRET`   | Secret key used to sign and verify JWTs                                     |

### 4. Run database migrations

```bash
goose -dir sql/schema postgres "$DB_URL" up
```

### 5. Run the server

```bash
go run .
```

The server starts on **port 8080**. You can verify it's running:

```bash
curl http://localhost:8080/api/healthz
# → OK
```

---

## API Documentation

All API endpoints are prefixed with `/api`. Responses and request bodies use JSON. Authenticated endpoints require a `Bearer` token in the `Authorization` header.

---

### Health Check

#### `GET /api/healthz`

Returns `200 OK` with a plain-text body confirming the server is alive. No authentication required.

---

### Users

#### `POST /api/users` — Register a new user

**Request body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response `201 Created`:**
```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "email": "user@example.com"
}
```

---

#### `PUT /api/users` — Update authenticated user's email and password

**Headers:** `Authorization: Bearer <access_token>`

**Request body:**
```json
{
  "email": "newemail@example.com",
  "password": "newpassword"
}
```

**Response `200 OK`:** Same shape as the create response (no token fields).

---

### Authentication

#### `POST /api/login` — Log in and receive tokens

**Request body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response `200 OK`:**
```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "email": "user@example.com",
  "token": "<jwt_access_token>",
  "refresh_token": "<refresh_token>"
}
```

The `token` is a JWT valid for **1 hour**. The `refresh_token` is valid for **60 days**.

---

#### `POST /api/refresh` — Get a new access token

Exchange a valid refresh token for a fresh JWT.

**Headers:** `Authorization: Bearer <refresh_token>`

**Response `200 OK`:**
```json
{
  "token": "<new_jwt_access_token>"
}
```

---

#### `POST /api/revoke` — Revoke a refresh token

Invalidates the provided refresh token immediately.

**Headers:** `Authorization: Bearer <refresh_token>`

**Response:** `204 No Content`

---

### Chirps

A chirp is limited to **140 characters**. Forbidden words (`kerfuffle`, `sharbert`, `fornax`) are automatically replaced with `****`.

---

#### `POST /api/chirps` — Create a chirp

**Headers:** `Authorization: Bearer <access_token>`

**Request body:**
```json
{
  "body": "Hello, Chirpy world!"
}
```

**Response `201 Created`:**
```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "body": "Hello, Chirpy world!",
  "user_id": "uuid"
}
```

---

#### `GET /api/chirps` — List all chirps

No authentication required. Returns all chirps sorted by creation time.

**Response `200 OK`:**
```json
[
  {
    "id": "uuid",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "body": "Hello, Chirpy world!",
    "user_id": "uuid"
  }
]
```

---

#### `GET /api/chirps/{chirpID}` — Get a single chirp by ID

No authentication required.

**Response `200 OK`:** Single chirp object (same shape as above).  
**Response `404 Not Found`:** If no chirp with that ID exists.

---

#### `DELETE /api/chirps/{chirpID}` — Delete a chirp

**Headers:** `Authorization: Bearer <access_token>`

Only the author of the chirp may delete it.

**Response:** `204 No Content`  
**Response `403 Forbidden`:** If the authenticated user does not own the chirp.

---

### Admin

These endpoints are intended for development and monitoring use only.

#### `GET /admin/metrics`

Returns an HTML page showing how many times the file server has been visited.

#### `POST /admin/reset`

Resets the visit counter and clears all users from the database. **Only available when `PLATFORM=dev`.**

---

## License

This project is licensed under the [GPL-3.0 License](LICENSE).
