# Full-Stack Developer Test

Image gallery SPA with infinite scroll, hashtag filtering, Go API, and MySQL.

## Scope

- gallery page with masonry layout
- infinite scroll with configurable chunk size
- hashtag filter
- management page for image CRUD, tag CRUD, and image-tag assignment
- MySQL migrations and demo seed data
- OpenAPI spec and Swagger UI for the backend

## Stack

- Frontend: React, Vite, TypeScript, Tailwind CSS
- Backend: Go, Fiber
- Database: MySQL 8
- Docs: OpenAPI, Swagger UI
- Demo image source: placehold.co

Production architecture notes are documented in `docs/stack.md`.

## Project Layout

```text
backend/   Go API, migrations, Swagger, tests
frontend/  React SPA and management UI
docs/      assignment PDF and architecture notes
```

## Assignment Coverage

Implemented:

- gallery page showing images with tags
- incremental loading while scrolling
- variable image sizes in a masonry layout
- filtering by tag
- placeholder images and generated demo tags
- architecture notes for production stack and deployment approach

Also included:

- `/manage` page for image CRUD
- tag CRUD
- attach and detach tags from images

## Run Locally

### 1. Backend

Configure `backend/.env` from `backend/.env.example`, then run:

```powershell
cd backend
go run ./cmd/api
```

Default API URL:

```text
http://localhost:8080
```

Swagger:

```text
http://localhost:8080/swagger
```

### 2. Frontend

Configure `frontend/.env` from `frontend/.env.example`, then run:

```powershell
cd frontend
npm install
npm run dev
```

Default frontend URL:

```text
http://localhost:5173
```

## Run With Docker Compose

The repository now includes:

- `backend/Dockerfile`
- `frontend/Dockerfile`
- `frontend/nginx.conf`
- `docker-compose.yml`

Start the full stack:

```powershell
docker compose up --build
```

What this starts:

- `mysql` for the local database
- `migrate` to apply schema and seed demo data
- `backend` for the Go API
- `frontend` running on Nginx

Public URLs:

```text
Frontend: http://localhost:8080
Backend API: http://localhost:8081
Swagger: http://localhost:8081/swagger
```

The Docker Compose setup uses demo/mock data automatically:

- `000001_create_gallery_tables.up.sql` creates the schema
- `000002_seed_demo_data.up.sql` seeds the local MySQL database with demo images and tags

If you want a clean reset of the local database and seed data, remove the volume and start again:

```powershell
docker compose down -v
docker compose up --build
```

## Database

Migration and seed scripts are under `backend/migration/`.

Run from the repository root:

```powershell
powershell -ExecutionPolicy Bypass -File backend/scripts/migrate.ps1 up
```

Or use the app-backed migration command when you want `MOCKDATA=true/false` to control demo seeding:

```powershell
cd backend
go run ./cmd/api migrate up
```

Seed data currently includes:

- 100 images from placehold.co
- 20 readable tags
- multiple image-tag assignments

## Tests

Backend:

```powershell
cd backend
go test ./...
```

Frontend build check:

```powershell
cd frontend
npm run build
```

## Notes

- `docs/stack.md` contains the proposed production setup, including server sizing, Ubuntu, Nginx, Docker Compose, and GitHub Actions.
- `docs/database-schema.md` explains the current MySQL schema and the purpose of each field.
- `docs/deploy-railway.md` contains a Railway deployment guide for `frontend`, `backend`, and `mysql`.
- Dockerfiles, Docker Compose, and Nginx config are included in the repository.
- GitHub Actions workflow files are not included yet.
- `frontend/README.md` and `backend/README.md` contain app-specific details.
