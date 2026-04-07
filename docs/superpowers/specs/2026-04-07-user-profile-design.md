# User Profile Feature Design

## Overview

Add a user profile (personal homepage) feature to the existing User CRUD API. This includes JWT-based authentication (register/login), profile viewing/editing, and avatar upload. Profile fields extend the existing `users` table directly.

## Data Model

### `users` table — new columns

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `username` | TEXT | NOT NULL, UNIQUE | Login username |
| `password_hash` | TEXT | NOT NULL | bcrypt hash |
| `avatar` | TEXT | | Avatar file path |
| `bio` | TEXT | | Personal bio |
| `phone` | TEXT | | Contact phone |
| `city` | TEXT | | City |
| `website` | TEXT | | Personal website URL |

### Domain model

```go
type User struct {
    ID           int64     `json:"id"`
    Username     string    `json:"username"`
    Name         string    `json:"name"`
    Email        string    `json:"email"`
    Age          int       `json:"age"`
    Avatar       string    `json:"avatar,omitempty"`
    Bio          string    `json:"bio,omitempty"`
    Phone        string    `json:"phone,omitempty"`
    City         string    `json:"city,omitempty"`
    Website      string    `json:"website,omitempty"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

`password_hash` is excluded from JSON serialization.

## API

### Auth endpoints (public)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/auth/register` | Register: `{username, email, password, name}` |
| POST | `/auth/login` | Login: `{username, password}` → `{token, user}` |
| GET | `/auth/me` | Current user info (requires token) |

### Profile endpoints (authenticated)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/users/:id/profile` | View user profile |
| PUT | `/users/me/profile` | Edit own profile: `{name, bio, phone, city, website}` |
| POST | `/users/me/avatar` | Upload avatar (multipart form, field `avatar`) |

### Profile response example

```json
{
  "id": 1,
  "username": "zhangsan",
  "name": "张三",
  "email": "zhangsan@example.com",
  "age": 28,
  "avatar": "/avatars/1.jpg",
  "bio": "Go developer",
  "phone": "+86-138xxxx",
  "city": "Beijing",
  "website": "https://example.com",
  "created_at": "2026-04-07T10:00:00Z",
  "updated_at": "2026-04-07T12:00:00Z"
}
```

## Authentication

- **Password hashing**: `golang.org/x/crypto/bcrypt` (cost=12)
- **JWT**: `github.com/golang-jwt/jwt/v5`
  - Claims: `user_id` (int64), `username` (string), `exp` (unix timestamp)
  - Expiration: 24 hours
  - Secret: env `SKK_JWT_SECRET`, fallback `"dev-secret-change-me"` for development
- **Middleware**: `AuthMiddleware` extracts `Bearer <token>` from `Authorization` header, verifies JWT, sets `user_id` in `gin.Context`

## Avatar Upload

- Storage: local filesystem at `data/avatars/`
- File naming: `{user_id}.{ext}` (overwrite on re-upload)
- Allowed extensions: `.jpg`, `.jpeg`, `.png`, `.gif`, `.webp`
- Max size: 2 MB
- Served via Gin static route: `r.Static("/avatars", "data/avatars")`

## Code Structure

New and modified files:

```
internal/
├── middleware/
│   └── auth.go              # JWT auth middleware
├── domain/
│   └── user.go              # Extended User model
├── handler/
│   ├── user.go              # Modified: add profile endpoints
│   ├── auth.go              # New: register/login/me handlers
│   ├── auth_test.go         # New
│   └── user_test.go         # Updated
├── service/
│   ├── user.go              # Modified: add profile methods
│   ├── auth.go              # New: register/login/JWT logic
│   ├── auth_test.go         # New
│   └── user_test.go         # Updated
├── repository/
│   ├── user.go              # Modified: interface additions
│   └── sqlite/
│       ├── user.go          # Modified: DAO fields, new queries
│       └── user_test.go     # Updated
```

## Dependencies

- `github.com/golang-jwt/jwt/v5` — new, for JWT signing/verification
- `golang.org/x/crypto` — already indirect dependency, for bcrypt

## Error Handling

- Registration with duplicate username/email → 409 Conflict
- Invalid credentials → 401 Unauthorized
- Missing/expired token → 401 Unauthorized
- Profile not found → 404 Not Found
- Avatar too large / invalid type → 400 Bad Request

## Testing Strategy

- Unit tests for `service/auth.go` (register, login, token generation)
- Unit tests for `service/user.go` (profile CRUD)
- Handler tests using `httptest` for all new endpoints
- Repository tests against in-memory SQLite
