# Agnos Backend Assignment

Backend API for staff authentication and patient lookup, running with `Go + Gin + PostgreSQL + Nginx`.

## Quick Start (for reviewers)

### 1) Start all services

```bash
docker compose up --build -d
```

Services:
- API via Nginx: `http://localhost:8080`
- PostgreSQL: `localhost:5432`

### 2) (Optional) Reset DB and reseed

Use this when you want a clean database from `database/init.sql`.

```bash
docker compose down -v
docker compose up --build -d
```

## API Test with cURL

### 1) Login (seed user)

```bash
curl -X POST http://localhost:8080/staff/login \
  -H "Content-Type: application/json" \
  -d '{"username":"staff_a","password":"password123"}'
```

Expected: returns JWT token.

### 2) Create new staff

```bash
curl -X POST http://localhost:8080/staff/create \
  -H "Content-Type: application/json" \
  -d '{"username":"staff_c","password":"password1234","hospital":"Hospital A"}'
```

Expected: `{"message":"staff created successfully"}`

### 3) Get patient by ID (national id or passport id)

```bash
curl -X GET http://localhost:8080/patient/P88888888
```

### 4) Search patient (requires Bearer token)

Replace `<TOKEN>` with token from login.

```bash
curl -G http://localhost:8080/patient/search \
  -H "Authorization: Bearer <TOKEN>" \
  --data-urlencode "passport_id=P88888888" \
  --data-urlencode "page=1" \
  --data-urlencode "limit=10"
```

Expected: returns filtered list and pagination object (`page`, `perPage`, `total`).

## Notes

- Seed users from `database/init.sql`:
  - `staff_a / password123`
  - `staff_b / password123`
- `GET /patient/:id` is public route
- `GET /patient/search` requires JWT